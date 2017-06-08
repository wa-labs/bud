package service

import (
	"fmt"
	"net/http"
	"time"

	"github.com/wa-labs/bud/pkg/service"
)

// Logger ...
type Logger struct {
	Logger service.Logger
}

// RequestLog ...
func (l *Logger) RequestLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		scheme := "http"
		if r.TLS != nil {
			scheme = "https"
		}
		d := fmt.Sprintf("%s://%s%s %s\" ", scheme, r.Host, r.RequestURI, r.Proto)

		begin := time.Now()
		defer func() {
			l.Logger.Prod(
				"method", r.Method,
				"request", d,
				"addr", r.RemoteAddr,
				"took", time.Since(begin),
			)

		}()
		next.ServeHTTP(w, r)
	})
}
