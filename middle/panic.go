package middle

import (
	"log"
	"net/http"
)

func (m *Middle) Panic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Println("panic:", err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Ошибка сервера"))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
