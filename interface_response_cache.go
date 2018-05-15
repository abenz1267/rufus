package rufus

import "net/http"

type responseCache interface {
	Check() func(next http.Handler) http.Handler
}
