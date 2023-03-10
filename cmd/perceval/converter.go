package main

import (
	"context"
	"fuu/v/cmd/perceval/config"
	"fuu/v/cmd/perceval/model"
	"fuu/v/cmd/perceval/utils"
	"fuu/v/internal/domain"
	"os"
	"path/filepath"
	"runtime"

	"github.com/bits-and-blooms/bloom/v3"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	outputPath = config.Instance().CacheDir
	pipeline   = make(chan int8, 1)
)

func convert(path, folder, format string, db *gorm.DB, logger *zap.SugaredLogger) error {
	_, err := os.Stat(path)
	if err != nil {
		return err
	}

	_ = os.Mkdir(outputPath, os.ModePerm)

	pipeline <- 1

	if utils.IsImagePath(path) {
		uuid := uuid.New()
		outfile := filepath.Join(outputPath, uuid.String())

		cmd := utils.GetCmd(path, outfile, format)
		cmd.Start()

		logger.Infow(
			"generating thumbnail",
			"path", path,
			"format", format,
			"cores", runtime.NumCPU(),
		)

		cmd.Wait()

		db.FirstOrCreate(&model.Thumbnail{
			Thumbnail: uuid.String(),
			Path:      outfile,
			Folder:    folder,
		})

		logger.Infow(
			"generated thumbnail",
			"id", uuid.String(),
			"path", path,
			"format", format,
			"cores", runtime.NumCPU(),
		)
	}

	<-pipeline
	return nil
}

func deleteFile(ctx context.Context, path string, db *gorm.DB, logger *zap.SugaredLogger) (string, error) {
	logger.Infow("deleting thumbnail", "path", path)
	id, toDelete, err := getByPath(ctx, path, db, logger)

	if id != "" {
		db.Delete(&model.Thumbnail{}, "path = ? OR thumbnail = ?", path, id)
		err := os.Remove(toDelete)
		return id, err
	}

	return "", err
}

func getByPath(ctx context.Context, path string, db *gorm.DB, logger *zap.SugaredLogger) (string, string, error) {
	logger.Infow("requesting thumbnail", "path", path)

	res := model.Thumbnail{}
	err := db.WithContext(ctx).First(&res, "path = ?", path).Error

	return res.Thumbnail, res.Path, err
}

func getManyByPath(ctx context.Context, paths []string, db *gorm.DB, logger *zap.SugaredLogger) (*[]model.Thumbnail, error) {
	logger.Infoln("requesting thumbnails")

	res := new([]model.Thumbnail)

	err := db.WithContext(ctx).Where(paths).Find(res).Error

	return res, err
}

func prune(db *gorm.DB, logger *zap.SugaredLogger) {
	all := &[]model.Thumbnail{}
	db.Find(all)

	filter := bloom.NewWithEstimates(uint(len(*all)), 0.01)

	logger.Infoln("started database prune")
	count := 0

	for _, entry := range *all {
		_, err := os.Stat(entry.Path)
		if os.IsNotExist(err) {
			db.Where("path = ?", entry.Path).Delete(&domain.Directory{})
			count++
		}
		if err == nil {
			filter.AddString(entry.Thumbnail)
		}
	}

	files, _ := os.ReadDir(config.Instance().CacheDir)
	for _, file := range files {
		if !filter.TestString(file.Name()) && filepath.Ext(file.Name()) != ".db" {
			toRemove := filepath.Join(config.Instance().CacheDir, file.Name())
			logger.Infow("deleting dead enrty", "file", toRemove)
			os.Remove(toRemove)
		}
	}

	logger.Infow("finished database prune", "count", count)
}
