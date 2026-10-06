// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"io/fs"
	"net/http"
	"path"
	"strings"
)

// translatedReference serves checked-in documentation snapshots when available.
// Explicit language, platform, and display-mode requests use the live GOROOT
// documentation so the snapshot cannot misrepresent another platform or mode.
func translatedReference(content fs.FS, translated, live http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		if query.Get("lang") == "zh" {
			query.Del("lang")
		}
		if len(query) == 0 {
			name := strings.TrimPrefix(path.Clean(r.URL.Path), "/") + "/index.html"
			if info, err := fs.Stat(content, name); err == nil && !info.IsDir() {
				translated.ServeHTTP(w, r)
				return
			}
		}
		live.ServeHTTP(w, r)
	})
}
