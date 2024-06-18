package middleware

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"time"
)

func (m *Middleware) LogRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			m.L.WithFields(
				m.L.TraceWrap(err),
			).Error(err)
		}
		r.Body.Close()
		r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

		defer func() {
			byteDump, err := httputil.DumpRequest(r, false)

			if err != nil {
				m.L.WithFields(
					m.L.TraceWrap(err),
				).Error(err)
			} else {
				msg := string(byteDump)

				if len(bodyBytes) > 0 {
					msg += fmt.Sprintf("\n Body: %s", string(bodyBytes))
				}

				msg += fmt.Sprintf("\n Duration: %v", time.Since(t))
				m.L.Info(msg)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
