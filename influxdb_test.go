package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

type InfluxDbTestSuite struct {
}

func NewInfluxDbTestSuite() *InfluxDbTestSuite {
	return &InfluxDbTestSuite{}
}

type MockInfluxDbServerResponse struct {
	ResponseCode   int
	ResponseString string
}

var MockInfluxDbServerResponses = []MockInfluxDbServerResponse{
	{
		200,
		`{
		"results": [
			{
				"statement_id": 0,
				"series": [
					{
						"name": "cpu_load_short",
						"columns": [
							"time",
							"value"
						],
						"values": [
							[
								"2015-01-29T21:55:43.702900257Z",
								500
							],
							[
								"2015-01-29T21:55:43.702900257Z",
								0.55
							],
							[
								"2015-06-11T20:46:02Z",
								0.64
							]
						]
					}
				]
			}
		]
	}`,
	},
	{
		403,
		"",
	},
}

func (tc *MockInfluxDbServerResponse) influxHttpHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(tc.ResponseCode)
	w.Write([]byte(tc.ResponseString))
}

func (tc *InfluxDbTestSuite) buildQueryObject() QueryObject {
	return QueryObject{
		BeginSample:     1472713200,
		CompletedSample: 1505415556,
		CreatedAt:       1505411956,
		EndSample:       1505415556,
		JobId:           "job1",
		Note:            "Created",
		Query:           "SELECT MEAN(\"Capacity_Value\") AS \"mean_Capacity_Value\" INTO \"drproduct\".\"ds\".\"testTEST\" FROM \"drproduct\".\"autogen\".\"cassandraCache\" WHERE time >= $ds_start_ts AND time <= $ds_end_ts GROUP BY time(600s),\"scope\"",
		QueryId:         "id",
		StartedAt:       1505412427,
		State:           "Complete",
		UpdatedAt:       1505413794,
		Database:        "db",
	}
}

func Test_Influxdb_Execute(t *testing.T) {
	tc := NewInfluxDbTestSuite()
	influxSpy := httptest.NewServer(http.HandlerFunc(MockInfluxDbServerResponses[0].influxHttpHandler))
	defer influxSpy.Close()
	client, err := NewInfluxdb(&InfluxDbConfig{&InfluxDbClientConfig{Url: influxSpy.URL, User: "", Password: ""}, &InfluxDbClientConfig{Url: influxSpy.URL, User: "", Password: ""}})
	if err != nil {
		t.Error(err)
	}
	res, err := client.Execute(tc.buildQueryObject(), 1505415556)
	if err != nil {
		t.Error(err)
	}
	if res != 3 {
		t.Error(fmt.Sprintf("%s expected to be %d but found %d", "points.written", 3, res))
	}
}

func Test_Influxdb_Execute_Error(t *testing.T) {
	tc := NewInfluxDbTestSuite()
	influxSpy := httptest.NewServer(http.HandlerFunc(MockInfluxDbServerResponses[1].influxHttpHandler))
	defer influxSpy.Close()
	client, err := NewInfluxdb(&InfluxDbConfig{&InfluxDbClientConfig{Url: influxSpy.URL, User: "", Password: ""}, &InfluxDbClientConfig{Url: influxSpy.URL, User: "", Password: ""}})
	if err != nil {
		t.Error(err)
	}
	_, err = client.Execute(tc.buildQueryObject(), 1505415556)
	if err == nil {
		t.Error(fmt.Sprintf("%s expected but not found none", "error"))
	}
}

func Test_Influxdb_buildQuery(t *testing.T) {
	tc := NewInfluxDbTestSuite()
	influxSpy := httptest.NewServer(http.HandlerFunc(MockInfluxDbServerResponses[1].influxHttpHandler))
	influxdb, err := NewInfluxdb(&InfluxDbConfig{&InfluxDbClientConfig{Url: influxSpy.URL, User: "", Password: ""}, &InfluxDbClientConfig{Url: influxSpy.URL, User: "", Password: ""}})
	if err != nil {
		t.Error(err)
	}
	q := influxdb.buildQuery(tc.buildQueryObject(), 1505415557)
	tobe := "SELECT MEAN(\"Capacity_Value\") AS \"mean_Capacity_Value\" INTO \"drproduct\".\"ds\".\"testTEST\" FROM \"drproduct\".\"autogen\".\"cassandraCache\" WHERE time >= 1505415556000000000 AND time <= 1505415557000000000 GROUP BY time(600s),\"scope\""
	if q != tobe {
		t.Error(fmt.Sprintf("%s expected to be %d but found %d", "query", tobe, q))
	}
}
