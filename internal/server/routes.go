// Copyright © 2017 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package server

import (
	"expvar"
	"net/http"

	"github.com/maier/go-appstats"
)

func (s *Server) router(w http.ResponseWriter, r *http.Request) {
	_ = appstats.IncrementInt("requests_total")

	s.logger.Debug().Str("method", r.Method).Str("url", r.URL.String()).Msg("request")

	switch r.Method {
	case "GET":
		switch {
		case pluginPathRx.MatchString(r.URL.Path): // run plugin(s)
			s.run(w, r)
		case inventoryPathRx.MatchString(r.URL.Path): // plugin inventory
			s.inventory(w)
		case statsPathRx.MatchString(r.URL.Path): // app stats
			expvar.Handler().ServeHTTP(w, r)
		case promPathRx.MatchString(r.URL.Path): // output prom format...
			s.promOutput(w)
		default:
			_ = appstats.IncrementInt("requests_bad")
			s.logger.Warn().Str("method", r.Method).Str("url", r.URL.String()).Msg("not found")
			http.NotFound(w, r)
		}
	case "POST":
		fallthrough
	case "PUT":
		switch {
		case writePathRx.MatchString(r.URL.Path):
			s.write(w, r)
		case promPathRx.MatchString(r.URL.Path):
			s.promReceiver(w, r)
		default:
			_ = appstats.IncrementInt("requests_bad")
			s.logger.Warn().Str("method", r.Method).Str("url", r.URL.String()).Msg("not found")
			http.NotFound(w, r)
		}
	default:
		_ = appstats.IncrementInt("requests_bad")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
