package service

import (
	"fmt"
	"net/http"
	"time"
)

// RequestLog ...
func RequestLog(l Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			scheme := "http"
			if r.TLS != nil {
				scheme = "https"
			}
			d := fmt.Sprintf("%s://%s%s %s\" ", scheme, r.Host, r.RequestURI, r.Proto)

			begin := time.Now()
			defer func() {
				l.Prod(
					"method", r.Method,
					"request", d,
					"addr", r.RemoteAddr,
					"took", time.Since(begin),
				)

			}()
			next.ServeHTTP(w, r)
		})
	}
}
