package utils

import (
	"os"
	"runtime"
	"strconv"
)

const (
	FormatAvif string = "avif"
	FormatWebP string = "webp"
)

func MaxParallelizationGrade() int {
	cores := runtime.NumCPU()

	if os.Getenv("CPUS") != "" {
		i, _ := strconv.Atoi(os.Getenv("CPUS"))
		cores = i
	}

	format := os.Getenv("PROCESSING_FORMAT")
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
