package main

import (
	"context"
	"fuu/v/cmd/internal"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

func main() {
	errChan, err := run()
	if err != nil {
		panic(err)
	}
	if err := <-errChan; err != nil {
		log.Fatalf("Error while running: %s", err)
	}
}

func run() (<-chan error, error) {
	zapConfig := zap.NewProductionConfig()
	zapConfig.OutputPaths = []string{"stdout"}

	if os.Getenv("LOG_PATH") != "" {
		zapConfig.OutputPaths = []string{os.Getenv("LOG_PATH"), "stdout"}
	}
	logger, _ := zapConfig.Build()

	rmq, err := internal.NewRabbitMQ(os.Getenv("RMQ_ENDPOINT"))
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	httpServer := &http.Server{
		Addr:         ":9898",
		Handler:      mux,
		WriteTimeout: 2 * time.Second,
		ReadTimeout:  2 * time.Second,
	}

	srv := &Server{
		logger: logger.Sugar(),
		srv:    httpServer,
		rmq:    rmq,
		done:   make(chan struct{}),
	}

	errChan := make(chan error, 1)

	ctx, stop := signal.NotifyContext(context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	go func() {
		<-ctx.Done()

		logger.Info("Shutdown signal received")

		ctxTimeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		defer func() {
			logger.Sync()
			rmq.Close()
			stop()
			cancel()
			close(errChan)
		}()

		if err := srv.Shutdown(ctxTimeout); err != nil { //nolint: contextcheck
			errChan <- err
		}

		logger.Info("Shutdown completed")
	}()

	go func() {
		logger.Info("Listening and serving")

		if err := srv.ListenAndServe(); err != nil {
			errChan <- err
		}
	}()

	return errChan, nil
}
