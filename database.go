package main

type DbInterface interface {
	GetItemByQueryId(id string) (QueryObject, error)
	PutItem(Query QueryObject) (string, error)
}
