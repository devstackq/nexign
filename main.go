package main

import (
	"net/http"

	"github.com/devstackq/nexign/config"
	"github.com/devstackq/nexign/handler"
	"go.uber.org/zap"
)

func main() {
	cfg := config.New()
	cfg.Load()

	lg, err := zap.NewProduction()
	if err != nil {
		lg.Fatal(err.Error())
	}

	h := handler.New(lg, cfg)

	srv := http.Server{
		Addr:           cfg.Port,
		Handler:        h.Route(),
		MaxHeaderBytes: 1 << 20,
	}

	lg.Info("starting web server on port: ")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		lg.Fatal(err.Error())
	}
}
