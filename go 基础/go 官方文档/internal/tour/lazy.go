// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tour

import (
	"net/http"
	"sync"
)

// RegisterLazyHandlers defers loading the tour until its first request.
// All tour pages, lessons, scripts, and static assets use the /tour/ prefix.
func RegisterLazyHandlers(mux *http.ServeMux) {
	var once sync.Once
	var err error
	tourMux := http.NewServeMux()
	mux.HandleFunc("/tour/", func(w http.ResponseWriter, r *http.Request) {
		once.Do(func() {
			err = RegisterHandlers(tourMux)
		})
		if err != nil {
			http.Error(w, "loading tour: "+err.Error(), http.StatusInternalServerError)
			return
		}
		tourMux.ServeHTTP(w, r)
	})
}
