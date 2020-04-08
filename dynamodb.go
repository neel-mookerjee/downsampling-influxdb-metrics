package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

type Dynamodb struct {
	Configs *Config
	Svc     dynamodbiface.DynamoDBAPI
}

func (d *Dynamodb) getApiService() dynamodbiface.DynamoDBAPI {
	config := &aws.Config{
		Region: aws.String("us-west-2"),
	}
	return dynamodb.New(session.New(), config)
}

func NewDynamodb(config *Config) *Dynamodb {
	db := &Dynamodb{}
	db.Svc = db.getApiService()
	db.Configs = config
	return db
}

func (d *Dynamodb) GetItemByQueryId(id string) (QueryObject, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(d.formTableName("metrics_downsampling_state")),
		Key: map[string]*dynamodb.AttributeValue{
			"queryId": {
				S: aws.String(id),
			},
		},
	}
	var query QueryObject
	result, err := d.Svc.GetItem(input)
	if err != nil {
		return query, err
	}
	err = dynamodbattribute.UnmarshalMap(result.Item, &query)
	if err != nil {
		return query, err
	}
	return query, nil
}

func (d *Dynamodb) PutItem(Query QueryObject) (string, error) {
	q, err := dynamodbattribute.MarshalMap(Query)
	if err != nil {
		return "Failure", err
	}

	_, err = d.Svc.PutItem(
		&dynamodb.PutItemInput{
			TableName: aws.String(d.formTableName("metrics_downsampling_state")),
			Item:      q,
		})
	if err != nil {
		return "Failure", err
	}

	status := "Success"
	return status, nil
}

func (d *Dynamodb) formTableName(table string) string {
	return fmt.Sprintf("%s%s", d.Configs.DbTablePrefix, table)
}
