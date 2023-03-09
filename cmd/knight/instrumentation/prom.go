package instrumentation

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	ItemsGuage = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "direcories_managed_guage",
		Help: "Number of directories managed by the backend",
	})
	TimePerOpGuage = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "time_per_conversion_ms",
		Help: "Latest time of completed conversion",
	})
	OpsCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "images_processed_counter",
		Help: "Number of image processed with imagemagick",
	})
	HardlinkedCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "images_hardlined_counter",
		Help: "Number of image hardlinked to convertion folder",
	})
)
