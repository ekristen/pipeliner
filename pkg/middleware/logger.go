package middleware

import (
	"bufio"
	"context"
	"errors"
	"net"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
)

type timer interface {
	Now() time.Time
	Since(time.Time) time.Duration
}

// realClock save request times
type realClock struct{}

func (rc *realClock) Now() time.Time {
	return time.Now()
}

func (rc *realClock) Since(t time.Time) time.Duration {
	return time.Since(t)
}

// LogOptions logging middleware options
type LogOptions struct {
	Formatter      logrus.Formatter
	EnableStarting bool
}

// LoggingMiddleware is a middleware handler that logs the request as it goes in and the response as it goes out.
type LoggingMiddleware struct {
	logger *logrus.Entry
	clock  timer
	user   string
	store  *sessions.CookieStore
}

// NewLogger returns a new *LoggingMiddleware, yay!
func NewLogger(log *logrus.Entry, store *sessions.CookieStore) *LoggingMiddleware {
	return &LoggingMiddleware{
		logger: log,
		clock:  &realClock{},
		store:  store,
	}
}

// realIP get the real IP from http request
func realIP(req *http.Request) string {
	ra := req.RemoteAddr
	if ip := req.Header.Get("X-Forwarded-For"); ip != "" {
		ra = strings.Split(ip, ", ")[0]
	} else if ip := req.Header.Get("X-Real-IP"); ip != "" {
		ra = ip
	} else {
		ra, _, _ = net.SplitHostPort(ra)
	}
	return ra
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
	hijacker   http.Hijacker
}

func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	hijacker, _ := w.(http.Hijacker)
	return &loggingResponseWriter{w, http.StatusOK, hijacker}
}

func (lw *loggingResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if lw.hijacker == nil {
		return nil, nil, errors.New("http.Hijacker not implemented by underlying http.ResponseWriter")
	}
	return lw.hijacker.Hijack()
}

func (lw *loggingResponseWriter) WriteHeader(code int) {
	lw.statusCode = code
	lw.ResponseWriter.WriteHeader(code)
}

func (lw *loggingResponseWriter) Write(b []byte) (int, error) {
	return lw.ResponseWriter.Write(b)
}

// Middleware implement mux middleware interface
func (m *LoggingMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		entry := m.logger
		start := m.clock.Now()

		ctx := r.Context()

		if reqID := r.Header.Get("X-Request-Id"); reqID != "" {
			ctx = context.WithValue(ctx, "requestId", reqID)
			entry = entry.WithField("requestId", reqID)
		} else {
			reqID := uuid.New().String()
			ctx = context.WithValue(ctx, "requestId", reqID)
			entry = entry.WithField("requestId", reqID)
		}

		r = r.WithContext(ctx)

		if remoteAddr := realIP(r); remoteAddr != "" {
			entry = entry.WithField("remoteAddr", remoteAddr)
		}

		re := regexp.MustCompile("(Authorization|Token)=([A-Za-z0-9-_=]+.[A-Za-z0-9-_=]+.?[A-Za-z0-9-_.+/=]*)")
		sRequestURI := string(re.ReplaceAll([]byte(r.RequestURI), []byte("$1=[MASKED]")))

		entry.WithFields(logrus.Fields{
			"request": sRequestURI,
			"method":  r.Method,
		}).Debug("started handling request")

		lw := newLoggingResponseWriter(w)
		next.ServeHTTP(lw, r)

		latency := m.clock.Since(start)

		entry.WithFields(logrus.Fields{
			"status":  lw.statusCode,
			"took":    latency,
			"request": sRequestURI,
		}).Info("handled: request")
	})
}
