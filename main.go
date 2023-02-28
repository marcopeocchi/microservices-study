package main

import (
	"embed"
	"fmt"
	"fuu/v/pkg"
	config "fuu/v/pkg/config"
	"fuu/v/pkg/domain"
	"log"
	"os"
	"path/filepath"

	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

//go:embed frontend/dist
var reactApp embed.FS

func main() {
	cfg := config.Instance()

	var cacheDir string
	homeDir, err := os.UserHomeDir()

	if err == nil {
		cacheDir = filepath.Join(homeDir, ".cache", "fuu")
		os.MkdirAll(cacheDir, os.ModePerm)
	} else {
		cacheDir = "/cache"
		_, err := os.Stat(cacheDir)
		if err != nil {
			log.Fatalln("Cannot find a valid cache directory")
		}
	}

	cfg.CacheDir = cacheDir

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.MysqlUser,
		cfg.MysqlPass,
		cfg.MysqlAddr,
		cfg.MysqlPort,
		cfg.MysqlDBName,
	)

	var db *gorm.DB

	if cfg.UseMySQL {
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	} else {
		dbPath := filepath.Join(cfg.CacheDir, "fuu.db")
		db, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	}

	if err != nil {
		log.Panicln(err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPass,
		DB:       0,
	})

	initDatabase(db)
	pkg.RunBlocking(db, rdb, &reactApp)
}

func initDatabase(db *gorm.DB) {
	db.AutoMigrate(&domain.Directory{})
	db.AutoMigrate(&domain.User{})

	p, err := bcrypt.GenerateFromPassword(
		[]byte(config.Instance().Masterpass),
		bcrypt.DefaultCost,
	)
	if err != nil {
		log.Fatalln(err)
	}

	db.Create(&domain.User{
		Username: "admin",
		Password: string(p),
		Role:     int(domain.Admin),
	})
}
