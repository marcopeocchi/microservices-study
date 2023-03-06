package workers

import (
	"context"
	conversionpb "fuu/v/gen/go/conversion/v1"
	config "fuu/v/pkg/config"
	"fuu/v/pkg/instrumentation"
	"sync"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	FormatAvif string = "avif"
	FormatWebP string = "webp"
)

func Converter(workingDir string, images []string, format string, logger *zap.SugaredLogger) {

	// if os.IsExist(err) {
	// 	return
	// }

	// if err != nil && !os.IsExist(err) {
	// 	logger.Errorw(
	// 		"error while creating conversion directory",
	// 		"error", err,
	// 	)
	// }

	available := []*grpc.ClientConn{}

	for _, node := range config.Instance().ImageProcessors {
		conn, err := getGRPCCLient(node)
		if err == nil {
			available = append(available, conn)
		}
	}

	wg := &sync.WaitGroup{}
	wg.Add(len(available))

	for i, chunk := range partition(images, len(images)/len(available)) {
		client := conversionpb.NewConversionServiceClient(available[i%len(available)])

		go func(part []string) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			start := time.Now()
			logger.Infow(
				"requested images conversion",
				"path", workingDir,
				"count", len(images),
				"format", format,
				"cores", 8,
			)

			res, err := client.Run(ctx, &conversionpb.RunRequest{
				Job: &conversionpb.ConversionJob{
					Path:   workingDir,
					Format: format,
					Files:  part,
				},
			})
			if err != nil {
				logger.Errorw("failed to run conversion", "error", err)
			}
			stop := time.Since(start)
			logger.Infow(
				"completed images conversion",
				"path", workingDir,
				"count", len(part),
				"format", format,
				"elapsed", stop,
				"remote", res,
			)
			instrumentation.TimePerOpGuage.Set(float64(stop / 1_000_000))
			wg.Done()
		}(chunk)
	}

	wg.Wait()

	for _, conn := range available {
		defer conn.Close()
	}
}

func getGRPCCLient(addr string) (*grpc.ClientConn, error) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	return grpc.Dial(addr, opts...)
}

func partition(arr []string, chunkSize int) (temp [][]string) {
	temp = [][]string{}
	for i := 0; i < len(arr); i += chunkSize {
		if i == len(arr)-1 {
			temp = append(temp, arr[i:])
			return
		}
		if i >= len(arr) {
			return
		}
		temp = append(temp, arr[i:i+chunkSize])
	}
	return
}
