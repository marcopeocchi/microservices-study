package main

import (
	"embed"
	"fuu/v/pkg"
	"log"
	"os"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

//go:embed frontend/dist
var reactApp embed.FS

func main() {
	r := pkg.ConfigReader{}
	cfg := r.Load()

	var cacheDir string
	homeDir, err := os.UserHomeDir()

	if err == nil {
		cacheDir = filepath.Join(homeDir, ".cache", "fuu")
		os.MkdirAll(filepath.Join(cacheDir), os.ModePerm)
	} else {
		cacheDir = "/cache"
		_, err := os.Stat(cacheDir)
		if err != nil {
			log.Fatalln("Cannot find a valid cache directory")
		}
	}

	cfg.CacheDir = cacheDir

	dbPath := filepath.Join(cacheDir, "fuu_thumbs.db")

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Panicln(err)
	}

	pkg.RunBlocking(cfg, db, &reactApp)
}
