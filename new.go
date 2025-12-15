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

var grpcLogger = New().WithPrefix("gRPC: ")

func NewGRPC() grpc.UnaryServerInterceptor {
	logger := grpcLogger

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
	logger := grpcLogger

	log := func(status string, info *grpc.StreamServerInfo, duration time.Duration) {
		logger.Infof("%s | STATUS: %s | Completed in %s", info.FullMethod, status, duration)
	}

	return func(server interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		duration, err := measureDurationStreamGRPC(handler, server, ss)

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

func NewHTTP(next http.Handler, logger *Logger) http.Handler {
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
	case "test":
		return &Logger{level: QuietLevel}
	case "production":
		return &Logger{level: InfoLevel}
	default:
		return &Logger{level: DebugLevel}
	}
}

func New() *Logger {
	l := logpkg.New(os.Stdout, "", 0)

	logger := newLogger()

	logger.l = l

	return logger
}
