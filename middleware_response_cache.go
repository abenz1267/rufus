package rufus

import (
	"net/http"
	"net/http/httptest"
	"sync"
)

// Cache is a response cache middleware
type Cache struct {
	liveCache, strippedCache, jsonCache sync.Map
}

// Check is used to perform the caching
func (c *Cache) Check() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "GET" {
				switch r.Header.Get("Accept") {
				case "application/json":
					c.checkCache(&c.jsonCache, w, r, next)
				case "text/html-raw":
					c.checkCache(&c.strippedCache, w, r, next)
				default:
					c.checkCache(&c.liveCache, w, r, next)
				}
			} else {
				next.ServeHTTP(w, r)
			}
		})
	}
}

// Invalidate is used to delete a cache entry
func (c *Cache) Invalidate(path string) bool {
	if _, ok := c.liveCache.Load(path); ok {
		c.liveCache.Delete(path)
		return true
	}

	if _, ok := c.strippedCache.Load(path); ok {
		c.strippedCache.Delete(path)
		return true
	}

	if _, ok := c.jsonCache.Load(path); ok {
		c.jsonCache.Delete(path)
		return true
	}

	return false
}

func (c *Cache) checkCache(cacheMap *sync.Map, w http.ResponseWriter, r *http.Request, next http.Handler) {
	if val, ok := cacheMap.Load(r.RequestURI); ok {
		w.Write(val.([]byte))
	} else {
		nw := httptest.NewRecorder()
		next.ServeHTTP(nw, r)

		go cacheMap.Store(r.RequestURI, nw.Body.Bytes())
		w.Write(nw.Body.Bytes())
	}
}
