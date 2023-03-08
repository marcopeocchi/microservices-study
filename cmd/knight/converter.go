package main

import (
	"fuu/v/cmd/knight/instrumentation"
	"fuu/v/cmd/knight/utils"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"go.uber.org/zap"
)

const quality = "80"

var pipeline = make(chan int8, utils.MaxParallelizationGrade())

func convert(path, format string, logger *zap.SugaredLogger) error {
	_, err := os.Stat(path)
	if err != nil {
		return err
	}

	file := filepath.Base(path)
	directory := filepath.Dir(path)
	outputDirectory := filepath.Join(directory, format)

	_ = os.Mkdir(outputDirectory, os.ModePerm)

	logger.Infow(
		"requested images conversion",
		"path", path,
		"format", format,
		"cores", runtime.NumCPU(),
	)

	pipeline <- 1

	// if strings.HasSuffix(filepath.Ext(file), format) {
	// }

	if utils.IsImagePath(file) {
		out := file[:len(file)-len(filepath.Ext(file))]
		outfile := filepath.Join(outputDirectory, out+"."+format)

		_, err := os.Stat(outfile)
		if err == nil {
			<-pipeline
			return nil
		}

		start := time.Now()

		cmd := exec.Command(
			"convert", filepath.Join(path),
			"-format", format,
			"-quality", quality,
			outfile,
		)
		cmd.Start()

		logger.Infow(
			"processing",
			"image", path,
			"format", format,
			"time", time.Now(),
		)

		cmd.Wait()
		stop := time.Since(start)
		instrumentation.TimePerOpGuage.Set(float64(stop) / 1_000_000)
		instrumentation.OpsCounter.Add(1)

		logger.Infow(
			"completed images conversion",
			"path", path,
			"format", format,
			"elapsed", stop,
		)
	}
	<-pipeline

	return nil
}
