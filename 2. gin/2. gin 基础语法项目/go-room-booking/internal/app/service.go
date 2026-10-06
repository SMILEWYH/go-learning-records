package app

import (
	"errors"
	"example.com/go-room-booking/internal/filedata"
	"example.com/go-room-booking/internal/web"
	"fmt"
	"os"
	"slices"
	"strings"
	"sync"
	"time"
	"unicode/utf8"
)

// 房间作为本地固定配置，避免把学习重点转为后台管理系统。
var rooms = []Room{{ID: "room-a", Name: "青竹会议室", Capacity: 4}, {ID: "room-b", Name: "云杉会议室", Capacity: 10}, {ID: "room-c", Name: "远程会议间", Capacity: 2}}

func Rooms() []Room { return slices.Clone(rooms) }
func roomByID(id string) (Room, error) {
	for _, r := range rooms {
		if r.ID == id {
			return r, nil
		}
	}
	return Room{}, web.Fail(404, "会议室不存在")
}

type Saver interface{ Save(State) error }
type disk struct{ path string }

func (d disk) Save(state State) error { return filedata.Save(d.path, state) }

type Service struct {
	mu    sync.Mutex
	state State
	saver Saver
	now   func() time.Time
}

func Open(path string) (*Service, error) {
	state := State{}
	err := filedata.Load(path, &state)
	if errors.Is(err, os.ErrNotExist) {
		state = State{Version: 1, NextID: 1, Bookings: map[string]Booking{}}
	} else if err != nil {
		return nil, fmt.Errorf("加载预约: %w", err)
	}
	if err := validateState(state); err != nil {
		return nil, err
	}
	return &Service{state: state, saver: disk{path}, now: time.Now}, nil
}
func clone(state State) State {
	next := state
	next.Bookings = make(map[string]Booking, len(state.Bookings))
	for k, v := range state.Bookings {
		next.Bookings[k] = v
	}
	return next
}
func (s *Service) commit(next State) error {
	if err := s.saver.Save(next); err != nil {
		return fmt.Errorf("保存预约: %w", err)
	}
	s.state = next
	return nil
}
func validText(s string, limit int) bool {
	return strings.TrimSpace(s) != "" && utf8.ValidString(s) && utf8.RuneCountInString(s) <= limit
}

// 时间必须带时区。统一转 UTC 后比较；区间采用 [start,end)，首尾相接可以预约。
func parseInterval(startRaw, endRaw string) (time.Time, time.Time, error) {
	start, e1 := time.Parse(time.RFC3339, startRaw)
	end, e2 := time.Parse(time.RFC3339, endRaw)
	if e1 != nil || e2 != nil {
		return time.Time{}, time.Time{}, web.Fail(400, "start/end 必须为带时区的 RFC3339 时间")
	}
	start = start.UTC()
	end = end.UTC()
	if start.Year() < 2000 || end.Year() > 2100 || start.Unix()%900 != 0 || end.Unix()%900 != 0 || start.Nanosecond() != 0 || end.Nanosecond() != 0 {
		return time.Time{}, time.Time{}, web.Fail(400, "时间范围为 2000..2100 年，且须对齐 UTC 的 15 分钟刻度")
	}
	duration := end.Sub(start)
	if duration < 15*time.Minute || duration > 4*time.Hour {
		return time.Time{}, time.Time{}, web.Fail(400, "预约时长必须为 15 分钟到 4 小时")
	}
	return start, end, nil
}
func validateInput(in BookingInput) (BookingInput, error) {
	in.Organizer = strings.TrimSpace(in.Organizer)
	in.Title = strings.TrimSpace(in.Title)
	room, err := roomByID(in.RoomID)
	if err != nil {
		return BookingInput{}, err
	}
	if !validText(in.Organizer, 40) || !validText(in.Title, 100) || in.Attendees < 1 || in.Attendees > room.Capacity {
		return BookingInput{}, web.Fail(400, "组织者 1..40 字符，标题 1..100 字符，人数不能超过会议室容量")
	}
	start, end, err := parseInterval(in.Start, in.End)
	if err != nil {
		return BookingInput{}, err
	}
	in.Start = start.Format(time.RFC3339)
	in.End = end.Format(time.RFC3339)
	return in, nil
}
func validateState(state State) error {
	bad := errors.New("预约数据内容或版本不合法；请检查数据文件")
	if state.Version != 1 || state.NextID < 1 || state.NextID > 10_000_001 || state.Bookings == nil {
		return bad
	}
	for id, b := range state.Bookings {
		var n int
		if _, err := fmt.Sscanf(id, "booking-%d", &n); err != nil || n < 1 || n >= state.NextID || id != fmt.Sprintf("booking-%06d", n) || b.ID != id {
			return bad
		}
		clean, err := validateInput(b.BookingInput)
		if err != nil || clean != b.BookingInput {
			return bad
		}
		if b.Status != Active && b.Status != Canceled {
			return bad
		}
		if _, err := time.Parse(time.RFC3339, b.CreatedAt); err != nil {
			return bad
		}
	}
	return nil
}
func overlaps(aStart, aEnd, bStart, bEnd time.Time) bool {
	return aStart.Before(bEnd) && bStart.Before(aEnd)
}
func future(start, now time.Time) error {
	if !start.After(now) || start.After(now.Add(90*24*time.Hour)) {
		return web.Fail(400, "开始时间必须晚于当前时间，且在未来 90 天内")
	}
	return nil
}
func (s *Service) Create(input BookingInput) (Booking, error) {
	in, err := validateInput(input)
	if err != nil {
		return Booking{}, err
	}
	start, end, _ := parseInterval(in.Start, in.End)
	s.mu.Lock()
	defer s.mu.Unlock()
	now := s.now()
	if err := future(start, now); err != nil {
		return Booking{}, err
	}
	for _, existing := range s.state.Bookings {
		if existing.Status != Active {
			continue
		}
		a, b, _ := parseInterval(existing.Start, existing.End)
		if overlaps(start, end, a, b) {
			if existing.RoomID == in.RoomID {
				return Booking{}, web.Fail(409, "这个会议室在所选时间已被预约")
			}
			if existing.Organizer == in.Organizer {
				return Booking{}, web.Fail(409, "同一组织者在所选时间已有其他会议")
			}
		}
	}
	if s.state.NextID > 10_000_000 {
		return Booking{}, web.Fail(409, "预约编号达到教学版上限")
	}
	next := clone(s.state)
	booking := Booking{ID: fmt.Sprintf("booking-%06d", next.NextID), BookingInput: in, Status: Active, CreatedAt: now.UTC().Format(time.RFC3339)}
	next.NextID++
	next.Bookings[booking.ID] = booking
	if err := s.commit(next); err != nil {
		return Booking{}, err
	}
	return booking, nil
}
func (s *Service) Get(id string) (Booking, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	b, ok := s.state.Bookings[id]
	if !ok {
		return Booking{}, web.Fail(404, "预约不存在")
	}
	return b, nil
}
func (s *Service) List(roomID, status string) ([]Booking, error) {
	if roomID != "" {
		if _, err := roomByID(roomID); err != nil {
			return nil, err
		}
	}
	if status != "" && status != string(Active) && status != string(Canceled) {
		return nil, web.Fail(400, "未知预约状态")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	out := []Booking{}
	for _, b := range s.state.Bookings {
		if (roomID == "" || b.RoomID == roomID) && (status == "" || string(b.Status) == status) {
			out = append(out, b)
		}
	}
	slices.SortFunc(out, func(a, b Booking) int {
		if n := strings.Compare(a.Start, b.Start); n != 0 {
			return n
		}
		return strings.Compare(a.ID, b.ID)
	})
	return out, nil
}
func (s *Service) Cancel(id string) (Booking, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	b, ok := s.state.Bookings[id]
	if !ok {
		return Booking{}, web.Fail(404, "预约不存在")
	}
	if b.Status == Canceled {
		return b, nil
	}
	start, _, _ := parseInterval(b.Start, b.End)
	if !start.After(s.now()) {
		return Booking{}, web.Fail(409, "已开始或已结束的预约不能取消")
	}
	b.Status = Canceled
	next := clone(s.state)
	next.Bookings[id] = b
	if err := s.commit(next); err != nil {
		return Booking{}, err
	}
	return b, nil
}

// Available 是即时查询结果；真正创建时仍要在锁内重新检查，不能凭查询结果承诺占位。
func (s *Service) Available(roomID, startRaw, endRaw string) (bool, error) {
	if _, err := roomByID(roomID); err != nil {
		return false, err
	}
	start, end, err := parseInterval(startRaw, endRaw)
	if err != nil {
		return false, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := future(start, s.now()); err != nil {
		return false, err
	}
	for _, b := range s.state.Bookings {
		if b.RoomID == roomID && b.Status == Active {
			a, z, _ := parseInterval(b.Start, b.End)
			if overlaps(start, end, a, z) {
				return false, nil
			}
		}
	}
	return true, nil
}
