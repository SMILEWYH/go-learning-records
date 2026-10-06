package app

import (
	"example.com/go-room-booking/internal/web"
	"net/http"
)

func Handler(s *Service, version string) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		web.Reply(w, 200, map[string]string{"status": "ok", "version": version})
	})
	mux.HandleFunc("GET /rooms", func(w http.ResponseWriter, r *http.Request) { web.Reply(w, 200, Rooms()) })
	mux.HandleFunc("GET /rooms/{id}/availability", func(w http.ResponseWriter, r *http.Request) {
		ok, err := s.Available(r.PathValue("id"), r.URL.Query().Get("start"), r.URL.Query().Get("end"))
		if err != nil {
			web.ReplyError(w, err)
			return
		}
		web.Reply(w, 200, map[string]bool{"available": ok})
	})
	mux.HandleFunc("POST /bookings", func(w http.ResponseWriter, r *http.Request) {
		var in BookingInput
		if err := web.Decode(w, r, &in); err != nil {
			web.ReplyError(w, err)
			return
		}
		b, err := s.Create(in)
		if err != nil {
			web.ReplyError(w, err)
			return
		}
		web.Reply(w, 201, b)
	})
	mux.HandleFunc("GET /bookings", func(w http.ResponseWriter, r *http.Request) {
		items, err := s.List(r.URL.Query().Get("room_id"), r.URL.Query().Get("status"))
		if err != nil {
			web.ReplyError(w, err)
			return
		}
		page, err := web.Paginate(r, items)
		if err != nil {
			web.ReplyError(w, err)
			return
		}
		web.Reply(w, 200, page)
	})
	mux.HandleFunc("GET /bookings/{id}", func(w http.ResponseWriter, r *http.Request) {
		b, err := s.Get(r.PathValue("id"))
		if err != nil {
			web.ReplyError(w, err)
			return
		}
		web.Reply(w, 200, b)
	})
	mux.HandleFunc("POST /bookings/{id}/cancel", func(w http.ResponseWriter, r *http.Request) {
		b, err := s.Cancel(r.PathValue("id"))
		if err != nil {
			web.ReplyError(w, err)
			return
		}
		web.Reply(w, 200, b)
	})
	return mux
}
