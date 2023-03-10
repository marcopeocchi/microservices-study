package workers

// import (
// 	"context"
// 	conversionpb "fuu/v/gen/go/conversion/v1"
// 	"fuu/v/pkg/config"
// 	"sync"
// 	"time"

// 	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
// 	"go.uber.org/zap"
// 	"google.golang.org/grpc"
// 	"google.golang.org/grpc/credentials/insecure"
// )

// const (
// 	FormatAvif string = "avif"
// 	FormatWebP string = "webp"
// )

// func Converter(ctx context.Context, workingDir string, images []string, format string, logger *zap.SugaredLogger) {
// 	available := []*grpc.ClientConn{}

// 	for _, node := range config.Instance().ImageProcessors {
// 		conn, err := getGRPCCLient(node)
// 		if err == nil {
// 			available = append(available, conn)
// 		}
// 	}

// 	if len(available) == 0 {
// 		logger.Warnw("no workers available", "found", len(available))
// 		return
// 	}

// 	partitions := partition(images, len(images)/len(available))
// 	wg := &sync.WaitGroup{}
// 	wg.Add(len(partitions))

// 	for i, chunk := range partitions {
// 		client := conversionpb.NewConversionServiceClient(available[i%len(available)])

// 		go func(part []string) {
// 			ctx, cancel := context.WithCancel(context.Background())
// 			defer cancel()

// 			start := time.Now()
// 			logger.Infow(
// 				"requested images conversion",
// 				"path", workingDir,
// 				"count", len(images),
// 				"format", format,
// 				"cores", 8,
// 			)

// 			res, err := client.Run(ctx, &conversionpb.RunRequest{
// 				Job: &conversionpb.ConversionJob{
// 					Path:   workingDir,
// 					Format: format,
// 					Files:  part,
// 				},
// 			})
// 			if err != nil {
// 				logger.Errorw("failed to run conversion", "error", err)
// 			}
// 			stop := time.Since(start)
// 			logger.Infow(
// 				"completed images conversion",
// 				"path", workingDir,
// 				"count", len(part),
// 				"format", format,
// 				"elapsed", stop,
// 				"remote", res,
// 			)
// 			wg.Done()
// 		}(chunk)
// 	}

// 	wg.Wait()

// 	for _, conn := range available {
// 		defer conn.Close()
// 	}
// }

// func getGRPCCLient(addr string) (*grpc.ClientConn, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*2500)
// 	defer cancel()

// 	// creds, err := crypto.LoadTLSCreds()
// 	// if err != nil {
// 	// 	return nil, err
// 	// }

// 	opts := []grpc.DialOption{
// 		grpc.WithTransportCredentials(insecure.NewCredentials()),
// 		grpc.WithBlock(),
// 		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
// 		grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()),
// 	}

// 	return grpc.DialContext(ctx, addr, opts...)
// }

// func partition(arr []string, chunkSize int) (temp [][]string) {
// 	temp = [][]string{}
// 	for i := 0; i < len(arr); i += chunkSize {
// 		if i == len(arr)-1 {
// 			temp = append(temp, arr[i:])
// 			return
// 		}
// 		if i >= len(arr) {
// 			return
// 		}
// 		temp = append(temp, arr[i:i+chunkSize])
// 	}
// 	return
// }
