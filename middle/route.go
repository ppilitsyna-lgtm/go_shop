package middle

import (
	"log"
	"net/http"
)

func (m *Middle) Route(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s query=%v", r.Method, r.URL.Path, r.URL.Query())
		next(w, r)
	})
}
