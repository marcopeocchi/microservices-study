package main

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"fuu/v/cmd/internal"
	"fuu/v/internal/domain"
	"fuu/v/internal/gallery"
	"fuu/v/internal/listing"
	"fuu/v/internal/user"
	"fuu/v/pkg/config"
	"fuu/v/pkg/instrumentation"
	"fuu/v/pkg/middlewares"
	"fuu/v/pkg/workers"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	//go:embed frontend/dist
	reactApp embed.FS
	cfg      = config.Instance()
)

func main() {
	errChan := run()

	if err := <-errChan; err != nil {
		log.Fatalln(err)
	}
}

func run() <-chan error {
	zapConfig := zap.NewProductionConfig()

	zapConfig.OutputPaths = []string{"stdout"}
	if cfg.LogPath != "" {
		zapConfig.OutputPaths = []string{cfg.LogPath, "stdout"}
	}

	logger, _ := zapConfig.Build()

	percevalConn, err := getGrpcClient("localhost:10099")
	if err != nil {
		panic(err)
	}

	defer logger.Sync()

	logger.Info("starting fuu")

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPass,
		DB:       0,
	})

	err = initCacheDir()
	if err != nil {
		panic(err)
	}

	db, err := initDatabase()
	if err != nil {
		panic(err)
	}

	instrumentation.InitTracing()

	thumbnailer := workers.Thumbnailer{
		BaseDir:    cfg.WorkingDir,
		ImgHeight:  cfg.ThumbnailHeight,
		ImgQuality: cfg.ThumbnailQuality,
		CacheDir:   cfg.CacheDir,
		Database:   db,
		Logger:     logger.Sugar(),
		Conn:       percevalConn,
	}

	fileWatcher := workers.FileWatcher{
		WorkingDir: cfg.WorkingDir,
		OnFileCreated: func(event string) {
			thumbnailer.Generate()
		},
		OnFileDeleted: func(event string) {
			thumbnailer.Remove(event)
		},
		Logger: logger.Sugar(),
	}
	fileWatcher.Setup()
	// fileWatcher.Poll()

	// Discourage the execution of this program as super user.
	// Unless in executed docker because of obvious reasons.
	uid := os.Getuid()
	if uid == 0 {
		logger.Warn("you're running this program as root (UID 0)")
	}

	rmq, err := internal.NewRabbitMQ(cfg.RabbitMQEnpoint)
	if err != nil {
		panic(err)
	}

	// ********** MAIN COMPONENTS GOROUTINES **********

	// HTTP Server
	server := initServer(ServerConfig{
		app:   &reactApp,
		port:  cfg.Port,
		db:    db,
		rdb:   rdb,
		rmq:   rmq,
		sugar: logger.Sugar(),
		conn:  percevalConn,
	})

	// Thumbnailer worker
	go func() {
		start := time.Now()
		thumbnailer.Generate()
		stop := time.Since(start)

		logger.Info("thumbnailer", zap.Duration("time", stop))
	}()

	// Ionotify filewatcher worker
	go func() {
		defer func() {
			fileWatcher.Close()
		}()

		fileWatcher.Start()
	}()

	// collect metrics for prometheus
	go instrumentation.CollectMetrics(db)

	// gracefully shutdown
	errChan := make(chan error, 1)

	ctx, stop := signal.NotifyContext(context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	go func() {
		<-ctx.Done()

		logger.Info("shutdown signal received")
		ctxTimeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		defer func() {
			logger.Sync()
			percevalConn.Close()
			fileWatcher.Close()
			rdb.Close()
			rmq.Close()
			stop()
			cancel()
			close(errChan)
		}()

		server.SetKeepAlivesEnabled(false)

		if err := server.Shutdown(ctxTimeout); err != nil {
			errChan <- err
		}

		logger.Info("shutdown completed")
	}()

	go func() {
		logger.Sugar().Infow("listening and serving", "port", cfg.Port)

		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errChan <- err
		}
	}()

	return errChan
}

type ServerConfig struct {
	port  int
	app   *embed.FS
	db    *gorm.DB
	rdb   *redis.Client
	rmq   *internal.RabbitMQ
	sugar *zap.SugaredLogger
	conn  *grpc.ClientConn
}

func initServer(sc ServerConfig) *http.Server {
	// Dependency injection containers
	userContainer := user.Container(sc.db, sc.sugar)
	listingContainer := listing.Container(sc.db, sc.rdb, sc.conn, sc.sugar)
	galleryContainer := gallery.Container(sc.rdb, sc.sugar, sc.rmq.Channel, cfg.WorkingDir)

	// Logging middleware
	lmdw := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sc.sugar.Desugar().Info(
				r.Method,
				zap.Time("time", time.Now()),
				zap.String("url", r.URL.String()),
			)
			next.ServeHTTP(w, r)
		})
	}

	// Router
	r := mux.NewRouter()
	r.Use(lmdw)
	r.Use(otelmux.Middleware("fuu"))

	// User group router
	ur := r.PathPrefix("/user").Subrouter()
	ur.HandleFunc("/login", userContainer.Login())
	ur.HandleFunc("/logout", userContainer.Logout())
	ur.HandleFunc("/signup", userContainer.SingUp())

	ur.Use(middlewares.CORS)

	// Overlay functionalites router
	or := r.PathPrefix("/overlay").Subrouter()
	or.HandleFunc("/list", listingContainer.ListAllDirectories())
	or.HandleFunc("/gallery", galleryContainer.DirectoryContent())
	or.Use(middlewares.CORS)
	or.Use(middlewares.Authenticated)

	// Static resources related router
	sr := r.PathPrefix("/static").Subrouter()
	sr.PathPrefix("/").Handler(http.StripPrefix("/static", http.FileServer(http.Dir(cfg.WorkingDir))))
	sr.Use(middlewares.Neuter)
	sr.Use(middlewares.Authenticated)

	// Thumbnails related router
	tr := r.PathPrefix("/thumbs").Subrouter()
	tr.PathPrefix("/").Handler(http.StripPrefix("/thumbs", http.FileServer(http.Dir(cfg.CacheDir))))
	tr.Use(middlewares.Neuter)
	tr.Use(middlewares.ServeThumbnail)
	tr.Use(middlewares.Authenticated)

	// Prometheus
	r.Handle("/metrics", promhttp.Handler())

	// Frontend
	build, _ := fs.Sub(*sc.app, "frontend/dist")

	sh := middlewares.SpaHandler{
		Entrypoint: "index.html",
		Filesystem: &build,
	}

	sh.AddRoute("/login")
	sh.AddRoute("/gallery")
	sh.AddRoute("/help")

	r.PathPrefix("/").Handler(sh.Handler())
	r.Use(middlewares.CORS)

	return &http.Server{
		Addr:         fmt.Sprintf(":%d", sc.port),
		Handler:      r,
		WriteTimeout: 2 * time.Second,
		ReadTimeout:  2 * time.Second,
	}
}

func initDatabase() (*gorm.DB, error) {
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

	db.AutoMigrate(&domain.Directory{})
	db.AutoMigrate(&domain.User{})

	p, err := bcrypt.GenerateFromPassword(
		[]byte(config.Instance().Masterpass),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return nil, err
	}

	db.Create(&domain.User{
		Username: "admin",
		Password: string(p),
		Role:     int(domain.Admin),
	})

	return db, nil
}

func initCacheDir() error {
	var cacheDir string
	homeDir, err := os.UserHomeDir()

	if err == nil {
		cacheDir = filepath.Join(homeDir, ".cache", "fuu")
		os.MkdirAll(cacheDir, os.ModePerm)
	} else {
		cacheDir = "/cache"
		_, err := os.Stat(cacheDir)
		if err != nil {
			return err
		}
	}

	cfg.CacheDir = cacheDir
	return nil
}

func getGrpcClient(addr string) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*2500)
	defer cancel()

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()),
	}

	return grpc.DialContext(ctx, addr, opts...)
}
