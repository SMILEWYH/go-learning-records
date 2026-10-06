package app

import (
	"errors"
	"example.com/go-shop-orders/internal/web"
	"fmt"
	"path/filepath"
	"sync"
	"testing"
)

type memorySaver struct {
	err   error
	saves int
}

func (m *memorySaver) Save(State) error { m.saves++; return m.err }
func service() (*Service, *memorySaver) {
	m := &memorySaver{}
	return &Service{state: initial(), saver: m}, m
}
func input(key string, quantity int) CreateOrder {
	return CreateOrder{RequestID: key, Customer: "小林", Items: []ItemInput{{ProductID: "p-1001", Quantity: quantity}}}
}
func assertStatus(t *testing.T, err error, status int) {
	t.Helper()
	var e *web.Error
	if !errors.As(err, &e) || e.Status != status {
		t.Fatalf("error=%v want status %d", err, status)
	}
}
func TestOrderLifecycle(t *testing.T) {
	s, m := service()
	o, err := s.Create(input("req-1", 2))
	if err != nil || o.Total != 3980 || s.Products()[0].Stock != 18 {
		t.Fatal(o, err)
	}
	same, err := s.Create(input("req-1", 2))
	if err != nil || same.ID != o.ID || m.saves != 1 {
		t.Fatal("repeat deducted stock", err)
	}
	_, err = s.Create(input("req-1", 3))
	assertStatus(t, err, 409)
	o.Items[0].Quantity = 999
	got, err := s.Get(o.ID)
	if err != nil || got.Items[0].Quantity != 2 {
		t.Fatal("mutable state escaped lock")
	}
	if _, err = s.Transition(o.ID, Canceled); err != nil {
		t.Fatal(err)
	}
	if _, err = s.Transition(o.ID, Canceled); err != nil {
		t.Fatal(err)
	}
	if s.Products()[0].Stock != 20 || m.saves != 2 {
		t.Fatal("cancel restored stock twice")
	}
	_, err = s.Transition(o.ID, Shipped)
	assertStatus(t, err, 409)
	second, err := s.Create(input("req-2", 1))
	if err != nil {
		t.Fatal(err)
	}
	if _, err = s.Transition(second.ID, Shipped); err != nil {
		t.Fatal(err)
	}
	_, err = s.Transition(second.ID, Canceled)
	assertStatus(t, err, 409)
}
func TestValidationAndAllOrNothing(t *testing.T) {
	for _, q := range []int{0, -1, 1001} {
		s, _ := service()
		_, err := s.Create(input("req", q))
		assertStatus(t, err, 400)
	}
	s, m := service()
	in := input("req", 2)
	in.Items = append(in.Items, ItemInput{"p-1002", 11})
	_, err := s.Create(in)
	assertStatus(t, err, 409)
	if s.Products()[0].Stock != 20 || m.saves != 0 {
		t.Fatal("partial stock deduction")
	}
	in.Items[1] = in.Items[0]
	_, err = s.Create(in)
	assertStatus(t, err, 400)
	in.Items[1] = ItemInput{"missing", 1}
	_, err = s.Create(in)
	assertStatus(t, err, 404)
	if _, err := s.Get("missing"); err == nil {
		t.Fatal("missing order accepted")
	}
}
func TestDiskFailureRollsBack(t *testing.T) {
	s, m := service()
	m.err = errors.New("disk full")
	if _, err := s.Create(input("req", 2)); !errors.Is(err, m.err) {
		t.Fatal(err)
	}
	if s.Products()[0].Stock != 20 || s.state.NextID != 1 || len(s.state.Orders) != 0 {
		t.Fatal("failed save changed memory")
	}
	m.err = nil
	o, err := s.Create(input("req", 2))
	if err != nil {
		t.Fatal(err)
	}
	m.err = errors.New("disk full")
	if _, err = s.Transition(o.ID, Canceled); err == nil {
		t.Fatal("expected failure")
	}
	got, _ := s.Get(o.ID)
	if got.Status != Confirmed || s.Products()[0].Stock != 18 {
		t.Fatal("failed cancel changed state")
	}
}
func TestConcurrentStockAndRequestDeduplication(t *testing.T) {
	s, _ := service()
	var wg sync.WaitGroup
	results := make(chan error, 40)
	for i := 0; i < 40; i++ {
		wg.Add(1)
		go func(i int) { defer wg.Done(); _, err := s.Create(input(fmt.Sprintf("req-%d", i), 1)); results <- err }(i)
	}
	wg.Wait()
	close(results)
	success := 0
	for err := range results {
		if err == nil {
			success++
		} else {
			assertStatus(t, err, 409)
		}
	}
	if success != 20 || s.Products()[0].Stock != 0 {
		t.Fatal(success, s.Products())
	}
	s, _ = service()
	results = make(chan error, 20)
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() { defer wg.Done(); _, err := s.Create(input("same-request", 1)); results <- err }()
	}
	wg.Wait()
	close(results)
	for err := range results {
		if err != nil {
			t.Fatal(err)
		}
	}
	if len(s.state.Orders) != 1 || s.Products()[0].Stock != 19 {
		t.Fatal("same request created multiple orders")
	}
}
func TestPersistence(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.json")
	s, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	o, err := s.Create(input("persist", 2))
	if err != nil {
		t.Fatal(err)
	}
	reopened, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	got, err := reopened.Get(o.ID)
	if err != nil || got.Total != 3980 || reopened.Products()[0].Stock != 18 {
		t.Fatal(got, err)
	}
	if _, err = reopened.Create(input("persist", 2)); err != nil || reopened.Products()[0].Stock != 18 {
		t.Fatal("dedup lost after restart", err)
	}
}
func TestRestockAndProducts(t *testing.T) {
	s, _ := service()
	_, err := s.AddProduct(Product{ID: "new", Name: "新商品", Price: 1, Stock: 0})
	if err != nil {
		t.Fatal(err)
	}
	_, err = s.AddProduct(Product{ID: "new", Name: "新商品", Price: 1})
	assertStatus(t, err, 409)
	p, err := s.Restock("new", 3)
	if err != nil || p.Stock != 3 {
		t.Fatal(p, err)
	}
	_, err = s.Restock("new", 0)
	assertStatus(t, err, 400)
	_, err = s.Restock("missing", 1)
	assertStatus(t, err, 404)
}
