package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

type reponseWriter struct {
	http.ResponseWriter
	statusCode   int
	bytesWritten int
}

func Logger(next http.Handler) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wr := &reponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(wr, r)

		duration := time.Since(start)

		slog.Info("HTTP request",
			"method", r.Method,
			"path", r.URL.Path,
			"query", r.URL.RawQuery,
			"status", wr.statusCode,
			"duration", duration.String(),
			"user_agent", r.UserAgent(),
			"ip", r.RemoteAddr,
			"time", start.Format(time.RFC3339),
		)
	}
	return http.HandlerFunc(f)
}
