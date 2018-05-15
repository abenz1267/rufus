package rufus

import (
	"net/http"
	"strings"
)

func (m *Middleware) redirectWithoutWWW() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.Contains(r.Host, "www.") {
				var b strings.Builder
				b.WriteString("https://")
				b.WriteString(strings.TrimPrefix(r.Host, "www."))
				b.WriteString(r.RequestURI)

				http.Redirect(w, r, b.String(), http.StatusMovedPermanently)
			} else {
				next.ServeHTTP(w, r)
			}
		})
	}
}
