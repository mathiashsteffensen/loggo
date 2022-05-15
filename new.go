package loggo

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	logpkg "log"
	"net/http"
	"os"
	"time"
)

func NewGRPC() grpc.UnaryServerInterceptor {
	logger := NewWithPrefix("gRPC: ")

	log := func(status string, info *grpc.UnaryServerInfo, duration time.Duration) {
		logger.Infof("%s | STATUS: %s | Completed in %s", info.FullMethod, status, duration)
	}

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		duration, resp, err := measureDurationGRPC(handler, ctx, req)

		if err != nil {
			log(redBg(bold(" ERROR ")), info, duration)
		}

		log(greenBg(bold(" OK ")), info, duration)

		return
	}
}

func NewGRPCStream() grpc.StreamServerInterceptor {
	logger := NewWithPrefix("gRPC: ")

	log := func(status string, info *grpc.StreamServerInfo, duration time.Duration) {
		logger.Infof("%s | STATUS: %s | Completed in %s", info.FullMethod, status, duration)
	}

	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		duration, err := measureDurationStreamGRPC(handler, srv, ss)

		if err != nil {
			log(redBg(bold(" ERROR ")), info, duration)
		}

		log(greenBg(bold(" OK ")), info, duration)

		return err
	}
}

type ResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *ResponseWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func NewHTTP(next http.Handler, logger ILogger) http.Handler {
	log := func(status, processedIn string, req *http.Request) {
		logger.Infof("%s %s | %s | %s", req.Method, req.URL.Path, status, processedIn)
	}

	return http.HandlerFunc(
		func(w http.ResponseWriter, req *http.Request) {
			rw := &ResponseWriter{w, 200}

			duration := measureDurationHTTP(next, rw, req)

			processedIn := fmt.Sprintf("Processed request in %s", duration)

			switch status := rw.statusCode; {
			case status < 400:
				log(greenBg(bold(" %v ", status)), processedIn, req)
			case 500 > status && status >= 400:
				log(yellowBg(bold(" %v ", status)), processedIn, req)
			case status >= 500:
				log(redBg(bold(" %v ", status)), processedIn, req)
			}
		},
	)
}

func newLogger() *Logger {
	switch env("GO_ENV", "development") {
	case "development":
		return &Logger{level: DebugLevel}
	case "test":
		return &Logger{level: QuietLevel}
	case "production":
		return &Logger{level: InfoLevel}
	default:
		panic("GO_ENV must be one of [production, development, test]")
	}
}

func NewVanilla() ILogger {
	l := logpkg.New(os.Stdout, "", 0)

	logger := newLogger()

	logger.l = l

	return logger
}

func NewWithPrefix(prefix string) ILogger {
	l := logpkg.New(os.Stdout, prefix, 0)

	logger := newLogger()

	logger.l = l

	return logger
}
