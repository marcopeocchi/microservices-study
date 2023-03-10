package instrumentation

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	ThumbnailsConverted = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "thumbnails_converted_count",
		Help: "Number of thumbnails converted",
	})
)
