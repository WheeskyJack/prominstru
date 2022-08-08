package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/version"

	"github.com/mackerelio/go-osstat/cpu"
	"github.com/mackerelio/go-osstat/memory"
)

var collector = "query_exporter"
var subSystem = "sys_stat"
var cpuMetricName = "cpu"
var memMetricName = "memory"
var cpuLabelNames = []string{"cpu"}
var memLabelNames = []string{"mem"}

func main() {
	var bind string
	flag.StringVar(&bind, "bind", "0.0.0.0:9104", "bind")
	flag.Parse()
	go recordMetrics()
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(bind, nil))
}

// SysStatCollector defines a struct for collector that contains pointers
// to prometheus descriptors for each metric we wish to expose.
type SysStatCollector struct {
	cpuMetric *prometheus.Desc
	memMetric *prometheus.Desc
}

func NewSysStatCollector(namespace string) *SysStatCollector {
	return &SysStatCollector{
		cpuMetric: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subSystem, cpuMetricName),
			"cpu stat",
			cpuLabelNames, nil,
		),
		memMetric: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subSystem, memMetricName),
			"mem stat",
			memLabelNames, nil,
		),
	}
}

func (c *SysStatCollector) Describe(ch chan<- *prometheus.Desc) {
	// Update this section with the each metric created for a given collector
	ch <- c.cpuMetric
	ch <- c.memMetric
}

func (c *SysStatCollector) Collect(ch chan<- prometheus.Metric) {
	// get latest val of each metric
	cr := getCpuMetric()
	mr := getMemMetric()

	// Write latest value for each metric in the prometheus metric channel.
	m1 := prometheus.MustNewConstMetric(c.cpuMetric, prometheus.GaugeValue, cr.val, cr.labelVals...)
	m2 := prometheus.MustNewConstMetric(c.memMetric, prometheus.GaugeValue, mr.val, mr.labelVals...)
	ch <- m1
	ch <- m2
}

func init() {
	prometheus.Register(version.NewCollector(collector))
	prometheus.Register(NewSysStatCollector(collector))
}

type MetRes struct {
	val       float64
	labelVals []string
}

func getCpuMetric() MetRes {
	var metricValue float64
	now, err := cpu.Get()
	if err == nil {
		metricValue = (float64(now.Total-now.Idle) / float64(now.Total)) * 100
	}

	// CPU Metric labels
	cpuLabelVals := []string{"total"}
	return MetRes{
		val:       metricValue,
		labelVals: cpuLabelVals,
	}
}

func getMemMetric() MetRes {
	var metricValue float64
	now, err := memory.Get()
	if err == nil {
		metricValue = float64(now.Used)
	}
	// Mem Metric labels
	memLabelVals := []string{"total"}
	return MetRes{
		val:       metricValue,
		labelVals: memLabelVals,
	}
}

// https://github.com/prometheus/client_golang/blob/main/prometheus/examples_test.go

func recordMetrics() {
	var taskCounterVec = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: collector,
			Name:      "completed_tasks_by_id",
			Subsystem: "worker_pool",
			Help:      "Total number of tasks completed.",
		},
		[]string{"worker_id", "status"},
	)
	prometheus.Register(taskCounterVec)
	var wo = map[string]string{ // id to status map
		"read":  "OK",
		"write": "NOK",
	}
	for n, s := range wo {
		go worker(taskCounterVec, n, s)
	}
}

func worker(c *prometheus.CounterVec, name, status string) {
	wc := c.WithLabelValues(name, status)
	for {
		time.Sleep(500 * time.Millisecond)
		wc.Inc()
	}
}
