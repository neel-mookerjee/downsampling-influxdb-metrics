package main

import "fmt"
import "testing"

type FileRepoTestSuite struct {
	FileRepo FileRepo
}

func NewFileRepoTestSuite() *FileRepoTestSuite {
	return &FileRepoTestSuite{FileRepo: FileRepo{}}
}

func Test_TaskManager_ReadWrite(t *testing.T) {
	ts := NewFileRepoTestSuite()
	ds := QueryObject{
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
	}
	ts.FileRepo.PutItem(ds)
	ds2, _ := ts.FileRepo.GetItemByQueryId("")
	if ds2 != ds {
		t.Error(fmt.Sprintf("%s expected to be %v but found %v", "object", ds, ds2))
	}
}
