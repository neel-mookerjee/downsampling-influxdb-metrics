package main

import (
	"github.com/rcrowley/go-metrics"
	"github.com/vrischmann/go-metrics-influxdb"
	"time"
)

type Metrics struct {
	Running                metrics.Gauge
	PointsWrittenCount     metrics.Counter
	ZeroPointsWrittenCount metrics.Counter
	ErrorCount             metrics.Counter
	Duration               metrics.Gauge
}

func (m *Metrics) startMetrics(r metrics.Registry, config Config, jobId string) {
	go influxdb.InfluxDBWithTags(r, 5e9,
		config.MetricsConfig.Host,
		config.MetricsConfig.Database,
		config.MetricsConfig.Username,
		config.MetricsConfig.Password,
		map[string]string{"job": jobId, "env": config.Environment, "app": "downsampling.historical.job"},
	)
}

func NewMetrics(config Config, jobId string) *Metrics {
	var m Metrics
	r := metrics.NewRegistry()
	m.Duration = metrics.NewGauge()
	r.Register("downsampling.historical.duration", m.Duration)
	m.Running = metrics.NewGauge()
	r.Register("downsampling.historical.running", m.Running)
	m.PointsWrittenCount = metrics.NewCounter()
	r.Register("downsampling.historical.pointsWrittenCount", m.PointsWrittenCount)
	m.ZeroPointsWrittenCount = metrics.NewCounter()
	r.Register("downsampling.historical.zeroPointsWrittenCount", m.ZeroPointsWrittenCount)
	m.ErrorCount = metrics.NewCounter()
	r.Register("downsampling.historical.errorCount", m.ErrorCount)

	metrics.RegisterDebugGCStats(r)
	go metrics.CaptureDebugGCStats(r, 5e9)

	metrics.RegisterRuntimeMemStats(r)
	go metrics.CaptureRuntimeMemStats(r, 5e9)

	m.startMetrics(r, config, jobId)
	return &m
}

func (m *Metrics) ResetErrorCount() {
	m.ErrorCount.Clear()
}

func (m *Metrics) SetDuration(started int64) {
	m.Duration.Update(time.Now().Unix() - started)
}

func (m *Metrics) ResetMetrics(restarted bool) {
	if restarted == true {
		return
	}
	m.Running.Update(0)
	m.PointsWrittenCount.Clear()
	m.ZeroPointsWrittenCount.Clear()
	m.Duration.Update(0)
}

func (m *Metrics) PointsWrittenByJob(points int64) {
	m.PointsWrittenCount.Inc(points)
	if points == 0 {
		m.IncrementZeroPointsWrittenCount()
	}
}

func (m *Metrics) IncrementErrorCount() {
	m.ErrorCount.Inc(1)
}

func (m *Metrics) IncrementZeroPointsWrittenCount() {
	m.ZeroPointsWrittenCount.Inc(1)
}

func (m *Metrics) JobRunning(running int64) {
	m.Running.Update(running)
}
