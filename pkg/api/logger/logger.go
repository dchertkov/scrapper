package logger

import (
	"net/http"
	"time"

	"github.com/dchertkov/scrapper/pkg/api/context"

	"github.com/sirupsen/logrus"
)

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		log := logrus.WithFields(logrus.Fields{
			"method":  r.Method,
			"request": r.RequestURI,
			"remote":  r.RemoteAddr,
		})

		ctx := context.ToLog(r.Context(), log)
		next.ServeHTTP(w, r.WithContext(ctx))

		log.WithFields(logrus.Fields{
			"latency": time.Since(start),
		}).Info("http request")
	})
}
