package rufus

import (
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/rs/zerolog/hlog"
)

func (m *Middleware) logRequests() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			nw := httptest.NewRecorder()
			now := time.Now()
			next.ServeHTTP(nw, r)

			hlog.FromRequest(r).Info().
				Str("URL", r.RequestURI).
				Int("status", nw.Code).
				Dur("in", time.Since(now)).
				Msg("")

			w.Header().Set("Content-Type", nw.Header().Get("Content-Type"))

			w.Write(nw.Body.Bytes())
		})
	}
}
