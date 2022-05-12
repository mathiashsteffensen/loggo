package loggo

import (
	"context"
	"github.com/fatih/color"
	"google.golang.org/grpc"
	"net/http"
	"os"
	"time"
)

func measureDurationGRPC(handler grpc.UnaryHandler, ctx context.Context, req interface{}) (time.Duration, interface{}, error) {
	start := time.Now()

	resp, err := handler(ctx, req)

	diff := time.Since(start)

	return diff, resp, err
}

func measureDurationHTTP(handler http.Handler, w http.ResponseWriter, req *http.Request) time.Duration {
	start := time.Now()

	handler.ServeHTTP(w, req)

	diff := time.Since(start)

	return diff
}

var (
	blue     = color.New(color.FgHiBlue).SprintfFunc()
	cyan     = color.New(color.FgHiCyan).SprintfFunc()
	red      = color.New(color.FgRed).SprintfFunc()
	yellow   = color.New(color.FgYellow).SprintfFunc()
	greenBg  = color.New(color.BgGreen).SprintfFunc()
	redBg    = color.New(color.BgRed).SprintfFunc()
	yellowBg = color.New(color.BgYellow).SprintfFunc()
	bold     = color.New(color.Bold).SprintfFunc()
)

func env(key, defaultValue string) string {
	if value := os.Getenv(key); value == "" {
		return defaultValue
	} else {
		return value
	}
}
