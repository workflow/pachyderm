package worker

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// PrometheusPort is the port the aggregated metrics are served on for scraping
	PrometheusPort = int32(9090)

	datumCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "pachyderm",
			Subsystem: "worker",
			Name:      "datum_count",
			Help:      "Number of datums processed by pipeline ID and state (started|errored|finished)",
		},
		[]string{
			"pipeline",
			"job",
			"state",
		},
	)

	bucketFactor  = 2.0
	bucketCount   = 20 // Which makes the max bucket 2^20 seconds or ~12 days in size
	datumProcTime = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "pachyderm",
			Subsystem: "worker",
			Name:      "datum_proc_time",
			Help:      "Time running user code",
			Buckets:   prometheus.ExponentialBuckets(1.0, bucketFactor, bucketCount),
		},
		[]string{
			"pipeline",
			"job",
			"state", // Since both finished and errored datums can have proc times
		},
	)

	datumDownloadTime = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "pachyderm",
			Subsystem: "worker",
			Name:      "datum_download_time",
			Help:      "Time to download input data",
			Buckets:   prometheus.ExponentialBuckets(1.0, bucketFactor, bucketCount),
		},
		[]string{
			"pipeline",
			"job",
		},
	)

	datumUploadTime = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "pachyderm",
			Subsystem: "worker",
			Name:      "datum_upload_time",
			Help:      "Time to upload output data",
			Buckets:   prometheus.ExponentialBuckets(1.0, bucketFactor, bucketCount),
		},
		[]string{
			"pipeline",
			"job",
		},
	)

	datumDownloadSize = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "pachyderm",
			Subsystem: "worker",
			Name:      "datum_download_size",
			Help:      "Size of downloaded input data",
			Buckets:   prometheus.ExponentialBuckets(1.0, bucketFactor, bucketCount),
		},
		[]string{
			"pipeline",
			"job",
		},
	)

	datumUploadSize = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "pachyderm",
			Subsystem: "worker",
			Name:      "datum_upload_size",
			Help:      "Size of uploaded output data",
			Buckets:   prometheus.ExponentialBuckets(1.0, bucketFactor, bucketCount),
		},
		[]string{
			"pipeline",
			"job",
		},
	)
)

func initPrometheus() {
	metrics := []prometheus.Collector{
		datumCount,
		datumProcTime,
		datumDownloadTime,
		datumUploadTime,
		datumDownloadSize,
		datumUploadSize,
	}
	for _, metric := range metrics {
		if err := prometheus.Register(metric); err != nil {
			fmt.Printf("error registering prometheus metric: %v\n", err)
		}
	}
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		if err := http.ListenAndServe(fmt.Sprintf(":%v", PrometheusPort), nil); err != nil {
			fmt.Printf("error serving prometheus metrics: %v\n", err)
		}
	}()
}
