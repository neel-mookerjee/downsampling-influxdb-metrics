package main

import (
	log "github.com/sirupsen/logrus"
	"time"
)

type TaskManagerInterface interface {
	Start(job_id string)
}

type TaskManager struct {
	Metrics     *Metrics
	Configs     *Config
	DbService   DbServiceInterface
	Downsampler DownsamplerInterface
}

func NewTaskManager(metrics *Metrics, config *Config) (*TaskManager, error) {
	ds, err := NewDownsampler(metrics, config)
	if err != nil {
		return nil, err
	}

	return &TaskManager{Metrics: metrics, Configs: config, DbService: NewDbService(config), Downsampler: ds}, nil
}

func (it *TaskManager) Start(queryId string) {
	it.Recursive(queryId, false)
}

func (it *TaskManager) Recursive(queryId string, restarted bool) {
	restart := it.Iterate(queryId, restarted)
	if restart {
		it.Recursive(queryId, restart)
	}
}

func (it *TaskManager) Iterate(queryId string, restarted bool) bool {

	it.Metrics.ResetMetrics(restarted)

	// choose Query
	log.Printf("Checking for pending downsampling task for Query Id: %s", queryId)
	QueryObj, err := it.DbService.GetAssignment(queryId)
	log.Printf("Query Object from DB: %v", QueryObj)
	if err != nil {
		it.doWhenError(err)
		restarted = true
		return restarted
	}

	// if received Query then execute
	if it.VerifyQueryState(QueryObj) == true {
		it.Metrics.JobRunning(1)
		it.Metrics.SetDuration(QueryObj.StartedAt)

		// mark first run
		if restarted == false {
			err = it.DbService.MarkBeginning(&QueryObj)
			if err != nil {
				it.doWhenError(err)
				restarted = true
				return restarted
			}
		}

		// traverse
		err = it.downsample(&QueryObj)
		if err != nil {
			it.doWhenError(err)
			restarted = true
			return restarted
		}

		// mark end
		err = it.DbService.MarkEnd(&QueryObj)
		if err != nil {
			it.doWhenError(err)
			restarted = true
			return restarted
		}

	} else {
		log.Println("No pending task found")
	}

	// reset stuff
	it.Metrics.ResetErrorCount()
	it.sleep()
	restarted = false
	return restarted
}

func (it *TaskManager) VerifyQueryState(query QueryObject) bool {
	if query.State == "Running" || query.State == "Ready" {
		return true
	}
	return false
}

func (it *TaskManager) downsample(QueryObj *QueryObject) error {
	log.Printf("Query execution start:: Query id: %s\n", QueryObj.QueryId)
	log.Println(QueryObj.ToString())

	err := it.Downsampler.Downsample(QueryObj)
	if err != nil {
		return err
	}

	log.Printf("Query execution end:: Query id: %s\n\n", QueryObj.QueryId)
	return nil
}

func (it *TaskManager) doWhenError(err error) {
	it.Metrics.IncrementErrorCount()
	log.Println("An error has occurred. Exiting the loop. Details as follows – \n", err)
	time.Sleep(time.Duration(it.Configs.JobConfig.RetryErrorInterval) * time.Second)
}

func (it *TaskManager) sleep() {
	it.Metrics.JobRunning(0)
	time.Sleep(time.Duration(it.Configs.JobConfig.SleepIdle) * time.Second)
}
