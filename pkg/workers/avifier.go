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
	"sync"
	"time"
)

var (
	pipeline = make(chan int, runtime.NumCPU())
	quality  = 75
)

func Avifier(workingDir string, images []string) {
	err := os.Mkdir(filepath.Join(workingDir, "avif"), os.ModePerm)

	if os.IsExist(err) {
		log.Println(workingDir, "already contains avif elements")
		log.Println(err.Error())
		return
	}

	start := time.Now()
	log.Println("Requested", workingDir, "AVIF conversion")

	wg := new(sync.WaitGroup)
	wg.Add(len(images))

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
			wg.Done()
		}(image)
	}

	wg.Wait()

	stop := time.Since(start)
	log.Println("Completed", workingDir, "AVIF conversion in", stop)
}
