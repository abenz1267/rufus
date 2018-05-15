package rufus

import "net/http"

func (m *Middleware) setContentType() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Header.Get("Accept") {
			case "application/json":
				w.Header().Add("Content-Type", "application/json; charset=utf-8")
			default:
				w.Header().Add("Content-Type", "text/html; charset=utf-8")
			}
			next.ServeHTTP(w, r)
		})
	}
}
