package main

import (
	"errors"
	"fmt"
	"github.com/rcrowley/go-metrics"
	"testing"
)

type TaskManagerTestSuite struct {
	TaskManager TaskManager
}

func NewTaskManagerTestSuite(queryObject QueryObject, errorBeginning error, erroDownsample error, errorEnd error) *TaskManagerTestSuite {
	return &TaskManagerTestSuite{TaskManager{Configs: &Config{JobConfig: &JobConfig{SleepIdle: 1}}, Metrics: getFakeMetrics(), Downsampler: &FakeDownsampler{erroDownsample}, DbService: FakeDbService2{queryObject, errorBeginning, errorEnd}}}
}

type FakeDbService2 struct {
	queryObject    QueryObject
	errorBeginning error
	errorEnd       error
}

func (f FakeDbService2) GetAssignment(jobId string) (QueryObject, error) {
	return f.queryObject, nil
}

func (f FakeDbService2) MarkBeginning(query *QueryObject) error {
	return f.errorBeginning
}

func (f FakeDbService2) MarkProgress(query QueryObject) error {
	return nil
}

func (f FakeDbService2) MarkEnd(query *QueryObject) error {
	return f.errorEnd
}

type FakeDownsampler struct {
	error error
}

func (f FakeDownsampler) Downsample(Query *QueryObject) error {
	return f.error
}

func getFakeMetrics() *Metrics {
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

func Test_TaskManager_Iterate_NoStatus(t *testing.T) {
	ts := NewTaskManagerTestSuite(QueryObject{
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
	}, errors.New("error"), nil, nil)
	restart := ts.TaskManager.Iterate("qeury", false)
	if restart {
		t.Error(fmt.Sprintf("%s expected to be %t but found %t", "restart", false, true))
	}
}

func Test_TaskManager_Iterate_Error1(t *testing.T) {
	ts := NewTaskManagerTestSuite(QueryObject{
		BeginSample:     1000,
		CompletedSample: 0,
		CreatedAt:       0,
		EndSample:       8000,
		JobId:           "job1",
		Note:            "NA",
		Query:           "NA",
		QueryId:         "MARKPROGRESS",
		StartedAt:       0,
		State:           "Running",
		UpdatedAt:       1234567890,
		Database:        "db",
	}, errors.New("error"), nil, nil)
	restart := ts.TaskManager.Iterate("qeury", false)
	if !restart {
		t.Error(fmt.Sprintf("%s expected to be %t but found %t", "restart", true, false))
	}
}

func Test_TaskManager_Iterate_Error2(t *testing.T) {
	ts := NewTaskManagerTestSuite(QueryObject{
		BeginSample:     1000,
		CompletedSample: 0,
		CreatedAt:       0,
		EndSample:       8000,
		JobId:           "job1",
		Note:            "NA",
		Query:           "NA",
		QueryId:         "MARKPROGRESS",
		StartedAt:       0,
		State:           "Running",
		UpdatedAt:       1234567890,
		Database:        "db",
	}, nil, nil, errors.New("error"))
	restart := ts.TaskManager.Iterate("qeury", false)
	if !restart {
		t.Error(fmt.Sprintf("%s expected to be %t but found %t", "restart", true, false))
	}
}

func Test_TaskManager_Iterate_Error3(t *testing.T) {
	ts := NewTaskManagerTestSuite(QueryObject{
		BeginSample:     1000,
		CompletedSample: 0,
		CreatedAt:       0,
		EndSample:       8000,
		JobId:           "job1",
		Note:            "NA",
		Query:           "NA",
		QueryId:         "MARKPROGRESS",
		StartedAt:       0,
		State:           "Running",
		UpdatedAt:       1234567890,
		Database:        "db",
	}, nil, errors.New("error"), nil)
	restart := ts.TaskManager.Iterate("qeury", false)
	if !restart {
		t.Error(fmt.Sprintf("%s expected to be %t but found %t", "restart", true, false))
	}
}

func Test_TaskManager_Iterate(t *testing.T) {
	ts := NewTaskManagerTestSuite(QueryObject{
		BeginSample:     1000,
		CompletedSample: 0,
		CreatedAt:       0,
		EndSample:       8000,
		JobId:           "job1",
		Note:            "NA",
		Query:           "NA",
		QueryId:         "MARKPROGRESS",
		StartedAt:       0,
		State:           "Running",
		UpdatedAt:       1234567890,
		Database:        "db",
	}, nil, nil, nil)
	restart := ts.TaskManager.Iterate("qeury", false)
	if restart {
		t.Error(fmt.Sprintf("%s expected to be %t but found %t", "restart", false, true))
	}
}
