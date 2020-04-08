package main

import (
	log "github.com/sirupsen/logrus"
	"time"
)

type DbServiceInterface interface {
	GetAssignment(queryId string) (QueryObject, error)
	MarkBeginning(query *QueryObject) error
	MarkProgress(query QueryObject) error
	MarkEnd(query *QueryObject) error
}

type DbService struct {
	db     DbInterface
	config Config
}

func NewDbService(Configs *Config) *DbService {
	if Configs.RepoType == "FILE" {
		return &DbService{NewFileRepo(Configs), *Configs}
	} else {
		return &DbService{NewDynamodb(Configs), *Configs}
	}
}

func (r *DbService) GetAssignment(queryId string) (QueryObject, error) {
	Query, err := r.db.GetItemByQueryId(queryId)
	if err != nil {
		return Query, err
	}

	return Query, nil
}

func (r *DbService) MarkBeginning(query *QueryObject) error {
	log.Printf("Marking beginning of query execution for query id: %s\n", query.QueryId)
	query.StartedAt = time.Now().Unix()
	query.State = "Running"
	query.JobId = r.config.JobConfig.JobName
	status, err := r.db.PutItem(*query)
	log.Printf("Marking status: %s\n", status)
	if err != nil {
		return err
	}

	return nil
}

func (r *DbService) MarkProgress(query QueryObject) error {
	log.Printf("Updating progress: %d for Query id: %s\n", query.CompletedSample, query.QueryId)
	status, err := r.db.PutItem(query)
	log.Printf("Update progress status: %s\n", status)
	return err
}

func (r *DbService) MarkEnd(query *QueryObject) error {
	log.Printf("Marking end of query execution for query id: %s\n", query.QueryId)
	query.UpdatedAt = time.Now().Unix()
	query.State = "Complete"
	query.Note = "Execution completed"
	query.JobId = "Unassigned"
	status, err := r.db.PutItem(*query)
	log.Printf("Marking status: %s\n", status)
	if err != nil {
		return err
	}

	return nil
}
