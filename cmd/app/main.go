package main

import (
	"context"
	"github.com/adamluzsi/frameless/pkg/logger"
	"github.com/adamluzsi/frameless/pkg/tasker"
	"net/http"
	"os"
)

func Main(ctx context.Context) error {
	return tasker.Main(ctx, tasker.HTTPServerTask(&http.Server{
		Addr: "0.0.0.0:8080",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello, world!"))
		}),
	}))
}

func main() {
	if err := Main(context.Background()); err != nil {
		logger.Fatal(nil, "main encountered an error", logger.ErrField(err))
		os.Exit(1)
	}
}
