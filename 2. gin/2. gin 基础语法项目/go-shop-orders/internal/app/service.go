package app

import (
	"errors"
	"example.com/go-shop-orders/internal/filedata"
	"example.com/go-shop-orders/internal/web"
	"fmt"
	"os"
	"slices"
	"strings"
	"sync"
	"time"
	"unicode/utf8"
)

// 小接口定义在使用方，测试可以注入“保存失败”，无需真实损坏磁盘（第 13、24 章）。
type Saver interface{ Save(State) error }
type disk struct{ path string }

func (d disk) Save(s State) error { return filedata.Save(d.path, s) }

type Service struct {
	mu    sync.Mutex
	state State
	saver Saver
}

func initial() State {
	return State{Version: 1, NextID: 1, Products: map[string]Product{
		"p-1001": {ID: "p-1001", Name: "办公笔记本", Price: 1990, Stock: 20},
		"p-1002": {ID: "p-1002", Name: "桌面收纳盒", Price: 3500, Stock: 10},
	}, Orders: map[string]Order{}}
}
func Open(path string) (*Service, error) {
	var state State
	err := filedata.Load(path, &state)
	if errors.Is(err, os.ErrNotExist) {
		state = initial()
	} else if err != nil {
		return nil, fmt.Errorf("加载订单数据: %w", err)
	}
	if err := validateState(state); err != nil {
		return nil, err
	}
	return &Service{state: state, saver: disk{path}}, nil
}
func validText(s string, max int) bool {
	return strings.TrimSpace(s) != "" && utf8.ValidString(s) && utf8.RuneCountInString(s) <= max
}
func validID(s string) bool {
	if len(s) < 1 || len(s) > 64 {
		return false
	}
	for _, c := range s {
		if !(c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9' || c == '-' || c == '_') {
			return false
		}
	}
	return true
}
func validateState(s State) error {
	bad := errors.New("数据文件内容不合法或版本不支持；请检查文件，程序不会自动覆盖")
	if s.Version != 1 || s.NextID < 1 || s.NextID > 10_000_001 || s.Products == nil || s.Orders == nil {
		return bad
	}
	for id, p := range s.Products {
		if id != p.ID || !validID(id) || !validText(p.Name, 80) || p.Price < 1 || p.Price > 100_000_000 || p.Stock < 0 || p.Stock > 1_000_000 {
			return bad
		}
	}
	requests := map[string]bool{}
	for id, o := range s.Orders {
		if id != o.ID || !validText(o.Customer, 80) || !validID(o.RequestID) || requests[o.RequestID] || len(o.Items) < 1 || len(o.Items) > 100 {
			return bad
		}
		requests[o.RequestID] = true
		var n int
		if _, err := fmt.Sscanf(id, "ord-%d", &n); err != nil || n < 1 || n >= s.NextID || id != fmt.Sprintf("ord-%06d", n) {
			return bad
		}
		if _, err := time.Parse(time.RFC3339, o.CreatedAt); err != nil {
			return bad
		}
		if o.Status != Confirmed && o.Status != Shipped && o.Status != Canceled {
			return bad
		}
		var total Cents
		seen := map[string]bool{}
		for _, item := range o.Items {
			if _, ok := s.Products[item.ProductID]; !ok || seen[item.ProductID] || item.Quantity < 1 || item.Quantity > 1000 || item.UnitPrice < 1 || item.UnitPrice > 100_000_000 || !validText(item.Name, 80) {
				return bad
			}
			seen[item.ProductID] = true
			total += item.UnitPrice * Cents(item.Quantity)
		}
		if total != o.Total {
			return bad
		}
	}
	return nil
}
func copyOrder(o Order) Order { o.Items = slices.Clone(o.Items); return o }
func clone(s State) State {
	next := s
	next.Products = make(map[string]Product, len(s.Products))
	next.Orders = make(map[string]Order, len(s.Orders))
	for k, v := range s.Products {
		next.Products[k] = v
	}
	for k, v := range s.Orders {
		next.Orders[k] = copyOrder(v)
	}
	return next
}

// 为了让小型文件版保持顺序，锁覆盖“检查、复制、保存、替换内存”。
// 锁内磁盘写入会限制吞吐量，这是明确的教学取舍；不能移到锁外后忽略并发覆盖。
func (s *Service) commit(next State) error {
	if err := s.saver.Save(next); err != nil {
		return fmt.Errorf("保存订单数据: %w", err)
	}
	s.state = next
	return nil
}
func (s *Service) Products() []Product {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]Product, 0, len(s.state.Products))
	for _, p := range s.state.Products {
		out = append(out, p)
	}
	slices.SortFunc(out, func(a, b Product) int { return strings.Compare(a.ID, b.ID) })
	return out
}
func (s *Service) AddProduct(p Product) (Product, error) {
	p.Name = strings.TrimSpace(p.Name)
	if !validID(p.ID) || !validText(p.Name, 80) || p.Price < 1 || p.Price > 100_000_000 || p.Stock < 0 || p.Stock > 1_000_000 {
		return Product{}, web.Fail(400, "id 使用 1..64 位字母数字-_，name 为 1..80 字符，价格 1..100000000 分，库存 0..1000000")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.state.Products[p.ID]; ok {
		return Product{}, web.Fail(409, "商品编号已存在")
	}
	next := clone(s.state)
	next.Products[p.ID] = p
	if err := s.commit(next); err != nil {
		return Product{}, err
	}
	return p, nil
}
func (s *Service) Restock(id string, count int) (Product, error) {
	if count < 1 || count > 100_000 {
		return Product{}, web.Fail(400, "补货数量必须为 1..100000")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	p, ok := s.state.Products[id]
	if !ok {
		return Product{}, web.Fail(404, "商品不存在")
	}
	// 保留已确认订单可能取消归还的空间，使取消始终能完整恢复库存。
	reserved := 0
	for _, o := range s.state.Orders {
		if o.Status == Confirmed {
			for _, item := range o.Items {
				if item.ProductID == id {
					reserved += item.Quantity
				}
			}
		}
	}
	if p.Stock+reserved+count > 1_000_000 {
		return Product{}, web.Fail(409, "可售库存加待履约数量不能超过 1000000")
	}
	p.Stock += count
	next := clone(s.state)
	next.Products[id] = p
	if err := s.commit(next); err != nil {
		return Product{}, err
	}
	return p, nil
}
func (s *Service) Create(in CreateOrder) (Order, error) {
	in.Customer = strings.TrimSpace(in.Customer)
	if !validID(in.RequestID) || !validText(in.Customer, 80) || len(in.Items) < 1 || len(in.Items) > 100 {
		return Order{}, web.Fail(400, "request_id 格式不合法、customer 为空或 items 数量不在 1..100")
	}
	// 复制后排序，以便同一请求改变行顺序仍能识别；不修改调用方的切片。
	in.Items = slices.Clone(in.Items)
	slices.SortFunc(in.Items, func(a, b ItemInput) int { return strings.Compare(a.ProductID, b.ProductID) })
	for i, item := range in.Items {
		if !validID(item.ProductID) || item.Quantity < 1 || item.Quantity > 1000 || (i > 0 && in.Items[i-1].ProductID == item.ProductID) {
			return Order{}, web.Fail(400, "商品不能重复，每项数量必须为 1..1000")
		}
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, o := range s.state.Orders {
		if o.RequestID == in.RequestID {
			same := o.Customer == in.Customer && len(o.Items) == len(in.Items)
			if same {
				for i, item := range in.Items {
					if o.Items[i].ProductID != item.ProductID || o.Items[i].Quantity != item.Quantity {
						same = false
						break
					}
				}
			}
			if !same {
				return Order{}, web.Fail(409, "同一 request_id 不能提交不同内容")
			}
			return copyOrder(o), nil
		}
	}
	if s.state.NextID > 10_000_000 {
		return Order{}, web.Fail(409, "订单编号达到教学版上限")
	}
	next := clone(s.state)
	o := Order{ID: fmt.Sprintf("ord-%06d", next.NextID), RequestID: in.RequestID, Customer: in.Customer, Items: []OrderItem{}, Status: Confirmed, CreatedAt: time.Now().UTC().Format(time.RFC3339)}
	for _, item := range in.Items {
		p, ok := next.Products[item.ProductID]
		if !ok {
			return Order{}, web.Fail(404, "商品不存在："+item.ProductID)
		}
		if p.Stock < item.Quantity {
			return Order{}, web.Fail(409, "库存不足："+item.ProductID)
		}
		p.Stock -= item.Quantity
		next.Products[p.ID] = p
		o.Items = append(o.Items, OrderItem{ProductID: p.ID, Name: p.Name, UnitPrice: p.Price, Quantity: item.Quantity})
		o.Total += p.Price * Cents(item.Quantity)
	}
	next.NextID++
	next.Orders[o.ID] = o
	if err := s.commit(next); err != nil {
		return Order{}, err
	}
	return copyOrder(o), nil
}
func (s *Service) Get(id string) (Order, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	o, ok := s.state.Orders[id]
	if !ok {
		return Order{}, web.Fail(404, "订单不存在")
	}
	return copyOrder(o), nil
}
func (s *Service) Orders(status string) ([]Order, error) {
	if status != "" && status != string(Confirmed) && status != string(Shipped) && status != string(Canceled) {
		return nil, web.Fail(400, "未知订单状态")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	out := []Order{}
	for _, o := range s.state.Orders {
		if status == "" || string(o.Status) == status {
			out = append(out, copyOrder(o))
		}
	}
	slices.SortFunc(out, func(a, b Order) int { return strings.Compare(a.ID, b.ID) })
	return out, nil
}
func (s *Service) Transition(id string, target OrderStatus) (Order, error) {
	if target != Shipped && target != Canceled {
		return Order{}, web.Fail(400, "目标状态只能是 shipped 或 canceled")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	o, ok := s.state.Orders[id]
	if !ok {
		return Order{}, web.Fail(404, "订单不存在")
	}
	if o.Status == target {
		return copyOrder(o), nil
	} // 重复点击不重复退库存。
	if o.Status != Confirmed {
		return Order{}, web.Fail(409, "只有 confirmed 订单可以发货或取消")
	}
	next := clone(s.state)
	if target == Canceled {
		for _, item := range o.Items {
			p := next.Products[item.ProductID]
			p.Stock += item.Quantity
			next.Products[p.ID] = p
		}
	}
	o.Status = target
	next.Orders[id] = o
	if err := s.commit(next); err != nil {
		return Order{}, err
	}
	return copyOrder(o), nil
}
