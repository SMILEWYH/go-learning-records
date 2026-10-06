// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"runtime"
	"strings"
	"testing"
	"testing/fstest"
)

func TestLocalPackageDocumentation(t *testing.T) {
	previous := *localDocs
	*localDocs = true
	t.Cleanup(func() { *localDocs = previous })
	mux := http.NewServeMux()
	if _, err := newSite(mux, "", os.DirFS("../../_content"), os.DirFS(runtime.GOROOT())); err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{"/pkg/", "/pkg/fmt/", "/cmd/go/"} {
		t.Run(path, func(t *testing.T) {
			r := httptest.NewRecorder()
			mux.ServeHTTP(r, httptest.NewRequest(http.MethodGet, path, nil))
			if r.Code != http.StatusOK || r.Header().Get("Location") != "" {
				t.Fatalf("local documentation redirected or failed: status=%d location=%q", r.Code, r.Header().Get("Location"))
			}
			if !strings.Contains(r.Body.String(), "<!DOCTYPE html>") {
				t.Fatal("local documentation did not render HTML")
			}
			if path != "/pkg/" && !strings.Contains(r.Body.String(), "ReferenceTranslationNotice") {
				t.Fatal("translated snapshot or version notice missing")
			}
		})
	}
	for _, path := range []string{"/pkg/fmt/?lang=en", "/pkg/fmt/?GOOS=windows&GOARCH=amd64", "/cmd/go/?lang=en"} {
		t.Run(path, func(t *testing.T) {
			r := httptest.NewRecorder()
			mux.ServeHTTP(r, httptest.NewRequest(http.MethodGet, path, nil))
			if r.Code != http.StatusOK || strings.Contains(r.Body.String(), "ReferenceTranslationNotice") {
				t.Fatal("explicit English or platform request did not use live documentation")
			}
		})
	}
	r := httptest.NewRecorder()
	mux.ServeHTTP(r, httptest.NewRequest(http.MethodGet, "/pkg/golang.org/x/tools/gopls/", nil))
	if r.Code != http.StatusTemporaryRedirect || r.Header().Get("Location") != "https://pkg.go.dev/golang.org/x/tools/gopls" {
		t.Fatalf("third-party package link failed: status=%d location=%q", r.Code, r.Header().Get("Location"))
	}
}

func TestTranslatedReferenceSelection(t *testing.T) {
	content := fstest.MapFS{
		"pkg/fmt/index.html": {Data: []byte("translated")},
		"cmd/go/index.html":  {Data: []byte("translated")},
	}
	marker := func(text string) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(text))
		})
	}
	h := translatedReference(content, marker("translated"), marker("live"))
	for _, tc := range []struct{ url, want string }{
		{"/pkg/fmt/", "translated"},
		{"/pkg/fmt", "translated"},
		{"/cmd/go/", "translated"},
		{"/pkg/fmt/?lang=zh", "translated"},
		{"/pkg/fmt/?lang=en", "live"},
		{"/pkg/fmt/?GOOS=windows&GOARCH=amd64", "live"},
		{"/pkg/fmt/?lang=zh&GOOS=windows", "live"},
		{"/pkg/fmt/?m=all", "live"},
		{"/pkg/bytes/", "live"},
		{"/pkg/", "live"},
		{"/pkg/golang.org/x/tools/gopls/", "live"},
	} {
		t.Run(tc.url, func(t *testing.T) {
			r := httptest.NewRecorder()
			h.ServeHTTP(r, httptest.NewRequest(http.MethodGet, tc.url, nil))
			if got := r.Body.String(); got != tc.want {
				t.Fatalf("got %q, want %q", got, tc.want)
			}
		})
	}
}
