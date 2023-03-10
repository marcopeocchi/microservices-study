package main

import (
	"context"
	"fmt"
	"fuu/v/cmd/perceval/config"
	"fuu/v/cmd/perceval/instrumentation"
	model "fuu/v/cmd/perceval/model"
	thumbnailspb "fuu/v/gen/go/grpc/thumbnails/v1"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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
	errChan := make(chan error, 1)

	logger, _ := zap.NewProduction()
	sugar := logger.Sugar()

	lis, err := net.Listen("tcp", ":10099")
	if err != nil {
		sugar.Fatalw("failed to listen", "error", err)
	}

	db, err := initDatabase()
	if err != nil {
		errChan <- err
	}

	prune(db, sugar)

	instrumentation.InitTracing()

	grpcSrv := grpc.NewServer(
		grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
		grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor()),
	)

	thumbnailspb.RegisterThumbnailServiceServer(grpcSrv, &ThumbnailsService{
		db:     db,
		Logger: sugar,
	})

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	httpSrv := &http.Server{
		Addr:         ":9899",
		Handler:      mux,
		WriteTimeout: 2 * time.Second,
		ReadTimeout:  2 * time.Second,
	}

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
			lis.Close()
			stop()
			cancel()
			close(errChan)
		}()

		if err := httpSrv.Shutdown(ctxTimeout); err != nil {
			errChan <- err
		}

		grpcSrv.GracefulStop()

		logger.Info("shutdown completed")
	}()

	go func() {
		logger.Info("http metrics server listening and serving")

		if err := httpSrv.ListenAndServe(); err != nil {
			errChan <- err
		}
	}()

	go func() {
		logger.Info("gRPC server listening and serving")

		if err := grpcSrv.Serve(lis); err != nil {
			errChan <- err
		}
	}()

	return errChan, nil
}

func initDatabase() (*gorm.DB, error) {
	cfg := config.Instance()

	var db *gorm.DB
	var err error

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.MysqlUser,
		cfg.MysqlPass,
		cfg.MysqlAddr,
		cfg.MysqlPort,
		cfg.MysqlDBName,
	)

	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&model.Thumbnail{})

	return db, nil
}
