package main

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

type DbServiceTestSuite struct {
	DbService DbServiceInterface
}

func NewDbServiceTestCase(config Config) *DbServiceTestSuite {
	return &DbServiceTestSuite{DbService: &DbService{&FakeDynamodb{}, config}}
}

// Define a mock struct to be used in your unit tests of functions.
type FakeDynamodb struct {
}

func (f *FakeDynamodb) GetItemByQueryId(id string) (QueryObject, error) {
	return QueryObject{}, nil
}

func (f *FakeDynamodb) GetItemByJobId(id string) (QueryObject, error) {
	return QueryObject{}, nil
}

func (f *FakeDynamodb) PutItem(Query QueryObject) (string, error) {
	retn := "Success"
	err := validate(Query)
	if err != nil {
		retn = "Failure"
	}
	return retn, err
}

func validate(query QueryObject) error {
	if query.QueryId == "MARKBEGINNING" {
		if query.StartedAt < time.Now().Unix()-2 {
			return errors.New(fmt.Sprintf("%s expected to be %d but found %d", "StartedAt", time.Now().Unix(), query.StartedAt))
		}
		if query.State != "Running" {
			return errors.New(fmt.Sprintf("%s expected to be %s but found %s", "State", "Running", query.State))
		}
	} else if query.QueryId == "MARKPROGRESS" {
		if query.UpdatedAt != 1234567890 {
			return errors.New(fmt.Sprintf("%s expected to be %d but found %d", "UpdatedAt", time.Now().Unix(), query.UpdatedAt))
		}
	} else if query.QueryId == "MARKEND" {
		if query.UpdatedAt < time.Now().Unix()-2 {
			return errors.New(fmt.Sprintf("%s expected to be %d but found %d", "UpdatedAt", time.Now().Unix(), query.UpdatedAt))
		}
		if query.State != "Complete" {
			return errors.New(fmt.Sprintf("%s expected to be %s but found %s", "State", "Complete", query.State))
		}
		if query.Note != "Execution completed" {
			return errors.New(fmt.Sprintf("%s expected to be %s but found %s", "Note", "\"Execution completed\"", query.Note))
		}
		if query.JobId != "Unassigned" {
			return errors.New(fmt.Sprintf("%s expected to be %s but found %s", "JobId", "Unassigned", query.JobId))
		}
	} else {
		return errors.New("invalid test")
	}
	return nil
}

func Test_DbService_MarkBeginning(t *testing.T) {
	tc := NewDbServiceTestCase(Config{JobConfig: &JobConfig{JobName: "test"}})
	query := &QueryObject{
		BeginSample:     0,
		CompletedSample: 0,
		CreatedAt:       0,
		EndSample:       0,
		JobId:           "job1",
		Note:            "NA",
		Query:           "NA",
		QueryId:         "MARKBEGINNING",
		StartedAt:       0,
		State:           "NA",
		UpdatedAt:       0,
		Database:        "db",
	}

	err := tc.DbService.MarkBeginning(query)
	if err != nil {
		t.Error(fmt.Sprintf("%s expected to be %s but found - %s", "MarkBeginning", "complete", err))
	}
}

func Test_DbService_MarkProgress(t *testing.T) {
	tc := NewDbServiceTestCase(Config{JobConfig: &JobConfig{}})
	query := &QueryObject{
		BeginSample:     0,
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
	err := tc.DbService.MarkProgress(*query)
	if err != nil {
		t.Error(fmt.Sprintf("%s expected to be %s but found - %s", "MarkProgress", "complete", err))
	}
}

func Test_DbService_MarkEnd(t *testing.T) {
	tc := NewDbServiceTestCase(Config{JobConfig: &JobConfig{}})
	query := &QueryObject{
		BeginSample:     0,
		CompletedSample: 0,
		CreatedAt:       0,
		EndSample:       0,
		JobId:           "Unassigned",
		Note:            "Execution completed",
		Query:           "NA",
		QueryId:         "MARKEND",
		StartedAt:       0,
		State:           "Complete",
		UpdatedAt:       time.Now().Unix(),
		Database:        "db",
	}
	err := tc.DbService.MarkEnd(query)
	if err != nil {
		t.Error(fmt.Sprintf("%s expected to be %s but found - %s", "MarkEnd", "complete", err))
	}
}
