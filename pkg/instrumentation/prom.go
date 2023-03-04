package instrumentation

import (
	"fuu/v/pkg/domain"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"gorm.io/gorm"
)

var (
	ItemsGuage = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "direcories_managed_guage",
		Help: "Number of directories managed by the backend",
	})
	OpsCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "images_processed_counter",
		Help: "Number of image processed with imagemagick",
	})
	CacheHitCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "cache_hit_counter",
		Help: "Number of cache hits for listing function",
	})
	CacheMissCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "cache_miss_counter",
		Help: "Number of cache miss for listing function",
	})
)

func CollectMetrics(db *gorm.DB) {
	go func() {
		for {
			var count int64
			db.Model(&domain.Directory{}).Count(&count)

			ItemsGuage.Set(float64(count))

			time.Sleep(time.Second * 2)
		}
	}()
}
