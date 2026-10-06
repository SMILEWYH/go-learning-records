package app

import (
	"errors"
	"example.com/go-room-booking/internal/web"
	"fmt"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

type memorySaver struct {
	err   error
	saves int
}

func (m *memorySaver) Save(State) error { m.saves++; return m.err }
func fixedNow() time.Time               { return time.Date(2026, 9, 27, 0, 0, 0, 0, time.UTC) }
func service() (*Service, *memorySaver) {
	m := &memorySaver{}
	return &Service{state: State{Version: 1, NextID: 1, Bookings: map[string]Booking{}}, saver: m, now: fixedNow}, m
}
func input() BookingInput {
	return BookingInput{RoomID: "room-a", Organizer: "小林", Title: "需求评审", Attendees: 4, Start: "2026-09-28T09:00:00+08:00", End: "2026-09-28T10:00:00+08:00"}
}
func assertStatus(t *testing.T, err error, status int) {
	t.Helper()
	var e *web.Error
	if !errors.As(err, &e) || e.Status != status {
		t.Fatalf("error=%v want status=%d", err, status)
	}
}
func TestBookingLifecycleAndTimezone(t *testing.T) {
	s, m := service()
	b, err := s.Create(input())
	if err != nil || b.Start != "2026-09-28T01:00:00Z" {
		t.Fatal(b, err)
	}
	same := input()
	same.Start = "2026-09-28T01:00:00Z"
	same.End = "2026-09-28T02:00:00Z"
	same.Organizer = "小周"
	_, err = s.Create(same)
	assertStatus(t, err, 409)
	adjacent := input()
	adjacent.Start = "2026-09-28T10:00:00+08:00"
	adjacent.End = "2026-09-28T10:30:00+08:00"
	if _, err := s.Create(adjacent); err != nil {
		t.Fatal("adjacent should pass", err)
	}
	if _, err := s.Cancel(b.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Cancel(b.ID); err != nil {
		t.Fatal(err)
	}
	if m.saves != 3 {
		t.Fatal("repeat cancel persisted twice")
	}
	if available, err := s.Available(b.RoomID, b.Start, b.End); err != nil || !available {
		t.Fatal(available, err)
	}
	if _, err = s.Create(same); err != nil {
		t.Fatal("canceled slot not freed", err)
	}
}
func TestIntervalCases(t *testing.T) {
	for _, tc := range []struct {
		name, start, end string
		valid            bool
	}{
		{"minimum", "2026-09-28T09:00:00+08:00", "2026-09-28T09:15:00+08:00", true},
		{"max", "2026-09-28T09:00:00+08:00", "2026-09-28T13:00:00+08:00", true},
		{"long", "2026-09-28T09:00:00+08:00", "2026-09-28T13:15:00+08:00", false},
		{"backward", "2026-09-28T10:00:00+08:00", "2026-09-28T09:00:00+08:00", false},
		{"unaligned", "2026-09-28T09:01:00+08:00", "2026-09-28T10:00:00+08:00", false},
		{"missing zone", "2026-09-28T09:00:00", "2026-09-28T10:00:00", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, _, err := parseInterval(tc.start, tc.end)
			if (err == nil) != tc.valid {
				t.Fatal(err)
			}
		})
	}
}
func TestBusinessValidation(t *testing.T) {
	s, _ := service()
	for _, mutate := range []func(*BookingInput){func(in *BookingInput) { in.Attendees = 5 }, func(in *BookingInput) { in.Title = " " }, func(in *BookingInput) { in.Start = "2026-09-26T09:00:00+08:00"; in.End = "2026-09-26T10:00:00+08:00" }, func(in *BookingInput) { in.Start = "2027-03-01T09:00:00+08:00"; in.End = "2027-03-01T10:00:00+08:00" }} {
		in := input()
		mutate(&in)
		_, err := s.Create(in)
		assertStatus(t, err, 400)
	}
	b, err := s.Create(input())
	if err != nil {
		t.Fatal(err)
	}
	other := input()
	other.RoomID = "room-b"
	_, err = s.Create(other)
	assertStatus(t, err, 409)
	s.now = func() time.Time { return time.Date(2026, 9, 28, 1, 0, 0, 0, time.UTC) }
	_, err = s.Cancel(b.ID)
	assertStatus(t, err, 409)
}
func TestConcurrentBooking(t *testing.T) {
	s, _ := service()
	results := make(chan error, 20)
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			in := input()
			in.Organizer = fmt.Sprintf("user-%d", i)
			_, err := s.Create(in)
			results <- err
		}(i)
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
	if success != 1 {
		t.Fatal("overlapping bookings accepted", success)
	}
}
func TestSaveFailureAndReload(t *testing.T) {
	s, m := service()
	m.err = errors.New("disk full")
	if _, err := s.Create(input()); !errors.Is(err, m.err) {
		t.Fatal(err)
	}
	if len(s.state.Bookings) != 0 || s.state.NextID != 1 {
		t.Fatal("failed save changed state")
	}
	m.err = nil
	b, err := s.Create(input())
	if err != nil {
		t.Fatal(err)
	}
	m.err = errors.New("disk full")
	if _, err = s.Cancel(b.ID); err == nil {
		t.Fatal("failure expected")
	}
	got, _ := s.Get(b.ID)
	if got.Status != Active {
		t.Fatal("cancel changed memory")
	}
	path := filepath.Join(t.TempDir(), "state.json")
	s, err = Open(path)
	if err != nil {
		t.Fatal(err)
	}
	s.now = fixedNow
	b, err = s.Create(input())
	if err != nil {
		t.Fatal(err)
	}
	reopened, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	got, err = reopened.Get(b.ID)
	if err != nil || got != b {
		t.Fatal(got, err)
	}
}
