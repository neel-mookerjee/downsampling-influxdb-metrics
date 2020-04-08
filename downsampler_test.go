package main

import (
	"fmt"
	"github.com/rcrowley/go-metrics"
	"testing"
)

type DownsampleTestSuite struct {
	FakeDbService DbServiceInterface
	FakeInfluxdb  InfluxdbInterface
}

func NewDownsampleTestSuite() *DownsampleTestSuite {
	return &DownsampleTestSuite{FakeInfluxdb: FakeInfluxdb{}, FakeDbService: FakeDbService{}}
}

func (ts *DownsampleTestSuite) getFakeMetrics(tc int) *Metrics {
	var m Metrics
	m.Duration = metrics.NewGauge()
	metrics.Register("duration", m.Duration)
	m.Running = metrics.NewGauge()
	metrics.Register("runningCount", m.Running)
	m.PointsWrittenCount = metrics.NewCounter()
	metrics.Register("pointsWrittenCount", m.PointsWrittenCount)
	m.ZeroPointsWrittenCount = metrics.NewCounter()
	metrics.Register("zeroPointsWrittenCount", m.ZeroPointsWrittenCount)
	m.ErrorCount = metrics.NewCounter()
	metrics.Register("errorCount", m.ErrorCount)
	return &m
}

type FakeDbService struct {
}

type FakeInfluxdb struct {
}

func (f FakeDbService) GetAssignment(jobId string) (QueryObject, error) {
	return QueryObject{}, nil
}

func (f FakeDbService) MarkBeginning(query *QueryObject) error {
	return nil
}

func (f FakeDbService) MarkProgress(query QueryObject) error {
	return nil
}

func (f FakeDbService) MarkEnd(query *QueryObject) error {
	return nil
}

func (f FakeInfluxdb) Execute(Query QueryObject, finishAtNew int64) (int, error) {
	return 100, nil
}

func Test_Downsampler_Downsample(t *testing.T) {
	ts := NewDownsampleTestSuite()
	ds := &Downsampler{Metrics: &Metrics{}, Configs: &Config{}, DbService: ts.FakeDbService, Influxdb: ts.FakeInfluxdb}
	ds.Downsample(&QueryObject{})
}

func Test_Downsampler_Downsample_CompletedSample(t *testing.T) {
	ts := NewDownsampleTestSuite()
	ds := &Downsampler{Metrics: ts.getFakeMetrics(1), Configs: &Config{}, DbService: ts.FakeDbService, Influxdb: ts.FakeInfluxdb}
	query := &QueryObject{
		BeginSample:     1000,
		CompletedSample: 0,
		CreatedAt:       0,
		EndSample:       0,
		JobId:           "job1",
		Note:            "NA",
		Query:           "NA",
		QueryId:         "MARKPROGRESS",
		StartedAt:       0,
		State:           "NA",
		UpdatedAt:       1234567890,
		Database:        "db",
	}
	ds.Downsample(query)
	if query.CompletedSample != 1000 {
		t.Error(fmt.Sprintf("%s expected to be %d but found %d", "CompletedSample", 1000, query.CompletedSample))
	}
}

func Test_Downsampler_Downsample_CompletedSample2(t *testing.T) {
	ts := NewDownsampleTestSuite()
	ds := &Downsampler{Metrics: ts.getFakeMetrics(2), Configs: &Config{"test", "FILE", "test", &InfluxDbConfig{}, &JobConfig{"test", 10, 10, 3600, 24, 100}, &MetricsConfig{}}, DbService: ts.FakeDbService, Influxdb: ts.FakeInfluxdb}
	query := &QueryObject{
		BeginSample:     1000,
		CompletedSample: 0,
		CreatedAt:       0,
		EndSample:       8000,
		JobId:           "job1",
		Note:            "NA",
		Query:           "NA",
		QueryId:         "MARKPROGRESS",
		StartedAt:       0,
		State:           "NA",
		UpdatedAt:       1234567890,
		Database:        "db",
	}
	ds.Downsample(query)
	if query.CompletedSample != 8000 {
		t.Error(fmt.Sprintf("%s expected to be %d but found %d", "CompletedSample", 8000, query.CompletedSample))
	}
}
