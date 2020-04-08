package main

import "fmt"
import "testing"

func Test_QueryObject_ToString(t *testing.T) {
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
	str, _ := query.ToString()
	strExp := `{"queryId":"MARKPROGRESS","query":"NA","database":"db","targetMeasurement":"","job":"job1","beginSample":1000,"endSample":8000,"startedAt":0,"completedSample":0,"createdAt":0,"updatedAt":1234567890,"note":"NA","state":"NA"}`
	if str != strExp {
		t.Error(fmt.Sprintf("%s expected to be %d but found %d", "query.ToString()", strExp, str))
	}
}
