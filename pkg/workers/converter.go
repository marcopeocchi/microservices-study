package workers

import (
	"fmt"
	config "fuu/v/pkg/config"
	"fuu/v/pkg/utils"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"
	"time"
)

var (
	pipeline = make(chan int, maxParallelizationGrade())
	quality  = 80
)

const (
	FormatAvif string = "avif"
	FormatWebP string = "webp"
)

func Converter(workingDir string, images []string, format string) {
	err := os.Mkdir(filepath.Join(workingDir, format), os.ModePerm)

	if os.IsExist(err) {
		log.Println(workingDir, "already contains optimized elements")
		log.Println(err.Error())
		return
	}

	start := time.Now()
	log.Println("Requested", workingDir, format, "conversion")

	wg := new(sync.WaitGroup)
	wg.Add(len(images))

	for _, image := range images {
		pipeline <- 1
		go func(img string) {
			if utils.IsImagePath(img) {
				cmd := exec.Command(
					"convert", filepath.Join(workingDir, img),
					"-format", format,
					"-quality", strconv.Itoa(quality),
					filepath.Join(workingDir, format, fmt.Sprint(img, ".", format)),
				)
				cmd.Start()
				cmd.Wait()
			}
			<-pipeline
			wg.Done()
		}(image)
	}

	wg.Wait()

	stop := time.Since(start)
	log.Println("Completed", workingDir, format, "conversion in", stop)
}

func maxParallelizationGrade() int {
	cores := runtime.NumCPU()
	format := config.Instance().ImageOptimizationFormat
	if cores == 1 {
		return 1
	}
	if cores <= 2 && format == FormatAvif {
		return 1
	}
	if cores <= 2 && format == FormatWebP {
		return 2
	}
	if cores > 2 && format == FormatAvif {
		return 1
	}
	return cores
}
