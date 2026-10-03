package app

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestHTTPValidationAndRealExchange(t *testing.T) {
	s, _ := service()
	handler := Handler(s, "test")
	server := httptest.NewServer(handler)
	defer server.Close()
	client := &http.Client{Timeout: 3 * time.Second}
	defer client.CloseIdleConnections()
	for _, tc := range []struct {
		method, path, body, media string
		status                    int
	}{
		{"GET", "/healthz", "", "", 200},
		{"POST", "/bookings", `{"room_id":"room-a","organizer":"小林","title":"评审","attendees":3,"start":"2026-09-28T09:00:00+08:00","end":"2026-09-28T10:00:00+08:00"}`, "application/json", 201},
		{"POST", "/bookings", `null`, "application/json", 400},
		{"POST", "/bookings", `{"room_id":"room-a","organizer":"小林","title":"评审","attendees":3,"start":"2026-09-28T09:00:00+08:00","end":"2026-09-28T10:00:00+08:00"}`, "text/plain", 415},
		{"POST", "/bookings", `{"room_id":"room-a","organizer":"小林","title":"评审","attendees":3,"start":"2026-09-28T09:00:00+08:00","end":"2026-09-28T10:00:00+08:00"}` + `{}`, "application/json", 400},
		{"GET", "/bookings?page=0", "", "", 400},
		{"GET", "/bookings?status=unknown", "", "", 400},
		{"GET", "/bookings/missing", "", "", 404},
		{"DELETE", "/bookings", "", "", 405},
		{"GET", "/missing", "", "", 404},
	} {
		req, err := http.NewRequest(tc.method, server.URL+tc.path, strings.NewReader(tc.body))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", tc.media)
		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		data, readErr := io.ReadAll(resp.Body)
		closeErr := resp.Body.Close()
		if readErr != nil || closeErr != nil {
			t.Fatal(readErr, closeErr)
		}
		if resp.StatusCode != tc.status {
			t.Fatalf("%s %s: %d want %d body=%s", tc.method, tc.path, resp.StatusCode, tc.status, data)
		}
		if tc.status == 201 {
			var created struct {
				ID string `json:"id"`
			}
			if err := json.Unmarshal(data, &created); err != nil || created.ID == "" {
				t.Fatal(string(data), err)
			}
		}
	}
}
