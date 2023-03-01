package workers

import (
	"fmt"
	"fuu/v/pkg/utils"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"time"
)

var (
	maxWorkers = make(chan int, 1)
	pipeline   = make(chan int, runtime.NumCPU())
	quality    = 75
)

func Avifier(workingDir string, images []string) {
	err := os.Mkdir(filepath.Join(workingDir, "avif"), os.ModePerm)
	if os.IsExist(err) {
		log.Println(workingDir, "already contains avif elements")
		log.Println(err.Error())
		return
	}

	start := time.Now()

	maxWorkers <- 1
	log.Println("Requested", workingDir, "AVIF conversion")

	for _, image := range images {
		pipeline <- 1
		go func(img string) {
			if utils.IsImagePath(img) {
				cmd := exec.Command(
					"convert", filepath.Join(workingDir, img),
					"-format", "avif",
					"-quality", strconv.Itoa(quality),
					filepath.Join(workingDir, "avif", fmt.Sprint(img, ".avif")),
				)
				cmd.Start()
				cmd.Wait()
			}
			<-pipeline
		}(image)
	}

	stop := time.Since(start)

	<-maxWorkers
	log.Println("Completed", workingDir, "AVIF conversion in", stop)
}
