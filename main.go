package main

import (
	"embed"
	"fuu/v/pkg"
	"fuu/v/pkg/common"
	config "fuu/v/pkg/config"
	"fuu/v/pkg/domain"
	"log"
	"os"
	"path/filepath"

	"golang.org/x/crypto/bcrypt"
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
	dbPath := filepath.Join(cfg.CacheDir, "fuu.db")

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Panicln(err)
	}

	initDatabase(db)
	pkg.RunBlocking(db, &reactApp)
}

func initDatabase(db *gorm.DB) {
	db.AutoMigrate(&domain.Directory{})
	db.AutoMigrate(&domain.User{})

	p, err := bcrypt.GenerateFromPassword(
		[]byte(config.Instance().Masterpass),
		common.BCRYPT_ROUNDS,
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
