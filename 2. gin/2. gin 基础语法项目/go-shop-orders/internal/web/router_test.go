package web

import (
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRouterErrorsAndRecovery(t *testing.T) {
	router := NewRouter()
	router.GET("/panic", func(c *gin.Context) { panic("private detail") })
	router.GET("/ok", func(c *gin.Context) { c.JSON(200, gin.H{"ok": true}) })
	for _, tc := range []struct {
		method, path string
		status       int
	}{
		{"GET", "/missing", 404}, {"POST", "/ok", 405}, {"GET", "/panic", 500}, {"GET", "/ok", 200},
	} {
		w := httptest.NewRecorder()
		router.ServeHTTP(w, httptest.NewRequest(tc.method, tc.path, nil))
		if w.Code != tc.status || !json.Valid(w.Body.Bytes()) || !strings.HasPrefix(w.Header().Get("Content-Type"), "application/json") {
			t.Fatalf("%s %s: %d %s", tc.method, tc.path, w.Code, w.Body.String())
		}
		if strings.Contains(w.Body.String(), "private detail") {
			t.Fatal("panic details leaked")
		}
		if tc.status == 405 && !strings.Contains(w.Header().Get("Allow"), "GET") {
			t.Fatal("missing Allow header")
		}
	}
}
