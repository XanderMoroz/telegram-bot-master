package middleware

import (
	"net/http"
)

func (m *Middleware) RecoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				if convertedError, ok := err.(error); ok {
					m.L.WithFields(
						m.L.TraceWrap(convertedError),
					).Error(convertedError)
				}
				http.Error(w, "", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
