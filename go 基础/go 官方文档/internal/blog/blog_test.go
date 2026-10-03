// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package blog

import (
	"fmt"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"testing/fstest"

	"golang.org/x/website/internal/web"
)

type countingFS struct {
	fs.FS
	opens     atomic.Int64
	blogOpens atomic.Int64
}

func (f *countingFS) Open(name string) (fs.File, error) {
	f.opens.Add(1)
	// Error rendering can probe blog subdirectories for layout templates;
	// count only directory listings and posts as feed generation work.
	if name == "blog" || strings.HasPrefix(name, "blog/") && (strings.HasSuffix(name, ".md") || strings.HasSuffix(name, ".html")) {
		f.blogOpens.Add(1)
	}
	return f.FS.Open(name)
}

func feedTestFS() fstest.MapFS {
	return fstest.MapFS{
		"blogfeed.tmpl": {Data: []byte(`{{.Content}}`)},
		"site.tmpl":     {Data: []byte(`{{template "layout" .}}`)},
		"error.tmpl":    {Data: []byte(`{{define "layout"}}{{.error}}{{end}}`)},
		"blog/example.md": {Data: []byte(`---
title: 中文示例
date: 2026-09-29
by:
- Go Team
summary: A sample post.
layout: none
---
A **sample** post.
`)},
	}
}

func TestRegisterLazyFeeds(t *testing.T) {
	// Record the public responses and I/O of one eager initialization. Lazy
	// initialization must preserve Atom aliases, JSON, and JSONP behavior.
	eagerFS := &countingFS{FS: feedTestFS()}
	eager := http.NewServeMux()
	if err := RegisterFeeds(eager, "", web.NewSite(eagerFS)); err != nil {
		t.Fatal(err)
	}
	targets := []string{
		"/blog/feed.atom",
		"/blog/feeds/posts/default",
		"/blog/.json",
		"/blog/.json?jsonp=callback.feed",
		"/blog/.json?jsonp=invalid!",
	}
	want := make([]*httptest.ResponseRecorder, len(targets))
	for i, target := range targets {
		want[i] = httptest.NewRecorder()
		eager.ServeHTTP(want[i], httptest.NewRequest("GET", target, nil))
		if want[i].Code != http.StatusOK || !strings.Contains(want[i].Body.String(), "中文示例") {
			t.Fatalf("invalid fixture response for %s: %d %s", target, want[i].Code, want[i].Body)
		}
	}

	lazyFS := &countingFS{FS: feedTestFS()}
	lazy := http.NewServeMux()
	RegisterLazyFeeds(lazy, "", web.NewSite(lazyFS))
	if got := lazyFS.opens.Load(); got != 0 {
		t.Fatalf("registration opened %d files; want no filesystem access", got)
	}

	// Simultaneous first requests must share a single initialization. Repeating
	// the requests also checks that later requests use the cached feeds.
	for round := 0; round < 2; round++ {
		var wg sync.WaitGroup
		errs := make(chan error, len(targets))
		for i, target := range targets {
			wg.Go(func() {
				got := httptest.NewRecorder()
				lazy.ServeHTTP(got, httptest.NewRequest("GET", target, nil))
				if got.Code != want[i].Code || got.Header().Get("Content-Type") != want[i].Header().Get("Content-Type") || got.Body.String() != want[i].Body.String() {
					errs <- fmt.Errorf("GET %s: got %d %q %s; want %d %q %s", target, got.Code, got.Header().Get("Content-Type"), got.Body, want[i].Code, want[i].Header().Get("Content-Type"), want[i].Body)
				}
			})
		}
		wg.Wait()
		close(errs)
		for err := range errs {
			t.Error(err)
		}
		if got, want := lazyFS.opens.Load(), eagerFS.opens.Load(); got != want {
			t.Errorf("round %d opened %d files; one initialization opens %d", round, got, want)
		}
	}
}

func TestRegisterLazyFeedsError(t *testing.T) {
	files := feedTestFS()
	files["blog/example.md"].Data = []byte("---\ntitle: Missing author\ndate: 2026-09-29\n---\nPost.\n")
	fsys := &countingFS{FS: files}
	mux := http.NewServeMux()
	RegisterLazyFeeds(mux, "", web.NewSite(fsys))
	if got := fsys.opens.Load(); got != 0 {
		t.Fatalf("registration opened %d files; want no filesystem access", got)
	}
	var firstReads int64
	for i, target := range []string{"/blog/feed.atom", "/blog/.json?jsonp=callback", "/blog/feeds/posts/default"} {
		got := httptest.NewRecorder()
		mux.ServeHTTP(got, httptest.NewRequest("GET", target, nil))
		if got.Code != http.StatusInternalServerError || !strings.Contains(got.Body.String(), "no author specified") {
			t.Errorf("GET %s: got %d %s; want the rendered feed error with status 500", target, got.Code, got.Body)
		}
		if i == 0 {
			firstReads = fsys.blogOpens.Load()
			if firstReads == 0 {
				t.Fatal("first request did not read blog files")
			}
		} else if got := fsys.blogOpens.Load(); got != firstReads {
			t.Errorf("GET %s reopened blog files after initialization failed: %d opens, want %d", target, got, firstReads)
		}
	}
}
