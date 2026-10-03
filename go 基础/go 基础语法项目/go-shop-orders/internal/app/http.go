package app

import (
	"example.com/go-shop-orders/internal/web"
	"net/http"
)

func Handler(s *Service, version string) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		web.Reply(w, 200, map[string]string{"status": "ok", "version": version})
	})
	mux.HandleFunc("GET /products", func(w http.ResponseWriter, r *http.Request) {
		result, err := web.Paginate(r, s.Products())
		if err != nil {
			web.ReplyError(w, err)
			return
		}
		web.Reply(w, 200, result)
	})
	mux.HandleFunc("POST /products", func(w http.ResponseWriter, r *http.Request) {
		var in Product
		if err := web.Decode(w, r, &in); err != nil {
			web.ReplyError(w, err)
			return
		}
		p, err := s.AddProduct(in)
		if err != nil {
			web.ReplyError(w, err)
			return
		}
		web.Reply(w, 201, p)
	})
	mux.HandleFunc("POST /products/{id}/restock", func(w http.ResponseWriter, r *http.Request) {
		var in struct {
			Quantity int `json:"quantity"`
		}
		if err := web.Decode(w, r, &in); err != nil {
			web.ReplyError(w, err)
			return
		}
		p, err := s.Restock(r.PathValue("id"), in.Quantity)
		if err != nil {
			web.ReplyError(w, err)
			return
		}
		web.Reply(w, 200, p)
	})
	mux.HandleFunc("GET /orders", func(w http.ResponseWriter, r *http.Request) {
		items, err := s.Orders(r.URL.Query().Get("status"))
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
	mux.HandleFunc("GET /orders/{id}", func(w http.ResponseWriter, r *http.Request) {
		o, err := s.Get(r.PathValue("id"))
		if err != nil {
			web.ReplyError(w, err)
			return
		}
		web.Reply(w, 200, o)
	})
	mux.HandleFunc("POST /orders", func(w http.ResponseWriter, r *http.Request) {
		var in CreateOrder
		if err := web.Decode(w, r, &in); err != nil {
			web.ReplyError(w, err)
			return
		}
		o, err := s.Create(in)
		if err != nil {
			web.ReplyError(w, err)
			return
		}
		web.Reply(w, 201, o)
	})
	mux.HandleFunc("PATCH /orders/{id}/status", func(w http.ResponseWriter, r *http.Request) {
		var in struct {
			Status OrderStatus `json:"status"`
		}
		if err := web.Decode(w, r, &in); err != nil {
			web.ReplyError(w, err)
			return
		}
		o, err := s.Transition(r.PathValue("id"), in.Status)
		if err != nil {
			web.ReplyError(w, err)
			return
		}
		web.Reply(w, 200, o)
	})
	return mux
}
