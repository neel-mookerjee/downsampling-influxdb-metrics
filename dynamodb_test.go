package main

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"testing"
)

type DynamodbTestSuite struct {
	Dynamodb DbInterface
}

func NewDynamodbTestSuite(getItemError error, putItemError error) *DynamodbTestSuite {
	return &DynamodbTestSuite{&Dynamodb{&Config{}, &mockDynamoDBClient{getItemError: getItemError, putItemError: putItemError}}}
}

// Define a mock struct to be used in your unit tests of functions.
type mockDynamoDBClient struct {
	dynamodbiface.DynamoDBAPI
	getItemError error
	putItemError error
}

func (m *mockDynamoDBClient) GetItem(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
	var av map[string]*dynamodb.AttributeValue
	av, _ = dynamodbattribute.ConvertToMap(QueryObject{
		BeginSample:     1472713200,
		CompletedSample: 1505415556,
		CreatedAt:       1505411956,
		EndSample:       1505415556,
		JobId:           "SEARCHBYQUERY",
		Note:            "Created",
		Query:           "SELECT MEAN(\"Capacity_Value\") AS \"mean_Capacity_Value\" INTO \"drproduct\".\"ds\".\"testTEST\" FROM \"drproduct\".\"autogen\".\"cassandraCache\" WHERE time >= $ds_start_ts AND time <= $ds_end_ts GROUP BY time(600s),\"scope\"",
		QueryId:         "query1",
		StartedAt:       1505412427,
		State:           "Complete",
		UpdatedAt:       1505413794,
		Database:        "db",
	})
	op := &dynamodb.GetItemOutput{Item: av}
	return op, m.getItemError
}

func (m *mockDynamoDBClient) PutItem(*dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	return &dynamodb.PutItemOutput{}, m.putItemError
}

func Test_Dynamodb_GetItemByQueryId(t *testing.T) {
	tc := NewDynamodbTestSuite(nil, nil)
	Query, _ := tc.Dynamodb.GetItemByQueryId("query1")
	if Query.JobId != "SEARCHBYQUERY" {
		t.Error(fmt.Sprintf("%s expected to be %s but found %s", "JobId", "SEARCHBYQUERY", Query.JobId))
	}
}

func Test_Dynamodb_GetItemByQueryId_Error(t *testing.T) {
	tc := NewDynamodbTestSuite(errors.New("error"), nil)
	_, err := tc.Dynamodb.GetItemByQueryId("query1")
	if err == nil {
		t.Error(fmt.Sprintf("%s expected but not found", "Error"))
	}
}

func Test_Dynamodb_PutItem(t *testing.T) {
	tc := NewDynamodbTestSuite(nil, nil)
	status, err := tc.Dynamodb.PutItem(QueryObject{
		BeginSample:     1472713200,
		CompletedSample: 1505415556,
		CreatedAt:       1505411956,
		EndSample:       1505415556,
		JobId:           "job1",
		Note:            "Created",
		Query:           "SELECT MEAN(\"Capacity_Value\") AS \"mean_Capacity_Value\" INTO \"drproduct\".\"ds\".\"testTEST\" FROM \"drproduct\".\"autogen\".\"cassandraCache\" WHERE time >= $ds_start_ts AND time <= $ds_end_ts GROUP BY time(600s),\"scope\"",
		QueryId:         "SEARCHBYJOB",
		StartedAt:       1505412427,
		State:           "Complete",
		UpdatedAt:       1505413794,
		Database:        "db",
	})
	if status != "Success" || err != nil {
		t.Error(fmt.Sprintf("%s was not expected but found, %v", "Error", err))
	}
}

func Test_Dynamodb_PutItem_Error(t *testing.T) {
	tc := NewDynamodbTestSuite(nil, errors.New("error"))
	status, err := tc.Dynamodb.PutItem(QueryObject{
		BeginSample:     1472713200,
		CompletedSample: 1505415556,
		CreatedAt:       1505411956,
		EndSample:       1505415556,
		JobId:           "job1",
		Note:            "Created",
		Query:           "SELECT MEAN(\"Capacity_Value\") AS \"mean_Capacity_Value\" INTO \"drproduct\".\"ds\".\"testTEST\" FROM \"drproduct\".\"autogen\".\"cassandraCache\" WHERE time >= $ds_start_ts AND time <= $ds_end_ts GROUP BY time(600s),\"scope\"",
		QueryId:         "SEARCHBYJOB",
		StartedAt:       1505412427,
		State:           "Complete",
		UpdatedAt:       1505413794,
		Database:        "db",
	})
	if status == "Success" || err == nil {
		t.Error(fmt.Sprintf("%s expected but not found", "Error"))
	}
}
