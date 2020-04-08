package main

import (
	"encoding/json"
	"io/ioutil"
)

type FileRepo struct {
	Configs *Config
}

func NewFileRepo(config *Config) *FileRepo {
	db := &FileRepo{}
	return db
}

func (d *FileRepo) GetItemByQueryId(id string) (QueryObject, error) {
	return d.getItem()
}

func (d *FileRepo) PutItem(Query QueryObject) (string, error) {
	queryJson, err := json.Marshal(Query)
	if err != nil {
		return "Failure", err
	}

	err = ioutil.WriteFile("./query.json", queryJson, 0644)
	if err != nil {
		return "Failure", err
	}

	status := "Success"
	return status, nil
}

func (d *FileRepo) getItem() (QueryObject, error) {
	query := QueryObject{}
	queryJson, err := ioutil.ReadFile("./query.json")
	if err != nil {
		return query, err
	}

	err = json.Unmarshal(queryJson, &query)
	if err != nil {
		return query, err
	}
	return query, nil
}
