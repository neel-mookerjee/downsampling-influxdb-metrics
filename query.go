package main

import (
	"encoding/json"
)

type QueryObject struct {
	QueryId           string `json:"queryId"`
	Query             string `json:"query"`
	Database          string `json:"database"`
	TargetMeasurement string `json:"targetMeasurement"`
	JobId             string `json:"job"`
	BeginSample       int64  `json:"beginSample"`
	EndSample         int64  `json:"endSample"`
	StartedAt         int64  `json:"startedAt"`
	CompletedSample   int64  `json:"completedSample"`
	CreatedAt         int64  `json:"createdAt"`
	UpdatedAt         int64  `json:"updatedAt"`
	Note              string `json:"note"`
	State             string `json:"state"`
	ready             bool
}

func (q QueryObject) ToString() (string, error) {
	return toJsonStr(q)
}

func toJsonStr(q interface{}) (string, error) {
	bytes, err := json.Marshal(q)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
