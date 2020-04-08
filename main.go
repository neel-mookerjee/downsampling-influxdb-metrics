package main

import (
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
)

func main() {
	// set logger
	log.SetFormatter(&log.JSONFormatter{})

	var queryId string
	// get query id
	if len(os.Args) < 2 {
		handleError(AppError{"Query id not supplied as a parameter. "})
	}
	log.Print("Param: " + os.Args[1])
	queryId = strings.Split(os.Args[1], "-")[1]
	if queryId == "" {
		handleError(AppError{"Query id is empty. "})
	}

	jobName := "historic-downsample-" + queryId
	log.Printf("Job: %s", jobName)

	// init
	configs, metrics, err := initialize(jobName)
	handleError(err)

	iterator, err := NewTaskManager(metrics, configs)
	handleError(err)

	iterator.Start(queryId)
}

func initialize(jobName string) (*Config, *Metrics, error) {
	// Config
	c, err := NewConfig(jobName)

	// start metrics
	m := NewMetrics(*c, jobName)

	return c, m, err
}

func handleError(err error) {
	if err != nil {
		log.Fatal("An error has occurred. Exiting the loop. Details as follows – \n", err)
	}
}
