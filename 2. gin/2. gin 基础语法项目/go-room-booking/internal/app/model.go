package app

type Status string

const (
	Active   Status = "active"
	Canceled Status = "canceled"
)

type Room struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Capacity int    `json:"capacity"`
}
type BookingInput struct {
	RoomID    string `json:"room_id"`
	Organizer string `json:"organizer"`
	Title     string `json:"title"`
	Attendees int    `json:"attendees"`
	Start     string `json:"start"`
	End       string `json:"end"`
}
type Booking struct {
	ID string `json:"id"`
	BookingInput
	Status    Status `json:"status"`
	CreatedAt string `json:"created_at"`
}
type State struct {
	Version  int                `json:"version"`
	NextID   int                `json:"next_id"`
	Bookings map[string]Booking `json:"bookings"`
}
