package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

// SlogMiddleware logs every request using the provided slogger
func SlogMiddleware(next http.Handler, l *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		timeStart := time.Now()
		respWriter := NewWriter(w, false)
		// serve the request
		next.ServeHTTP(respWriter, r)
		// get the duration
		timeDiff := time.Since(timeStart)
		log(l, r, respWriter.StatusCode(), timeDiff)
	})
}

func log(l *slog.Logger, r *http.Request, statusCode int, dur time.Duration) {

	if IsStatusError(statusCode) {

		// TODO read the error string from the http error response
		l.Error("",
			slog.String("method", r.Method),
			slog.String("url", r.RequestURI),
			slog.Duration("req-dur", dur),
			slog.Int("response-code", statusCode),
			slog.String("ip", userIp(r)),
			slog.String("req-id", r.Header.Get("Request-Id")),
			slog.String("err-message", "TODO"),
		)
	} else {
		l.Info("",
			slog.String("method", r.Method),
			slog.String("url", r.RequestURI),
			slog.Duration("req-dur", dur),
			slog.Int("response-code", statusCode),
			slog.String("ip", userIp(r)),
			slog.String("req-id", r.Header.Get("Request-Id")),
		)
	}
}

func userIp(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}
	return IPAddress
}
