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
		{"POST", "/orders", `{"request_id":"http-1","customer":"小林","items":[{"product_id":"p-1001","quantity":2}]}`, "application/json", 201},
		{"POST", "/orders", `null`, "application/json", 400},
		{"POST", "/orders", `{"request_id":"http-1","customer":"小林","items":[{"product_id":"p-1001","quantity":2}]}`, "text/plain", 415},
		{"POST", "/orders", `{"request_id":"http-1","customer":"小林","items":[{"product_id":"p-1001","quantity":2}]}` + `{}`, "application/json", 400},
		{"GET", "/orders?page=0", "", "", 400},
		{"GET", "/orders?status=unknown", "", "", 400},
		{"GET", "/orders/missing", "", "", 404},
		{"DELETE", "/orders", "", "", 405},
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
