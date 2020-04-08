package main

import (
	log "github.com/sirupsen/logrus"
	"time"
)

type DownsamplerInterface interface {
	Downsample(Query *QueryObject) error
}

type Downsampler struct {
	Metrics   *Metrics
	Configs   *Config
	DbService DbServiceInterface
	Influxdb  InfluxdbInterface
}

func NewDownsampler(metrics *Metrics, config *Config) (*Downsampler, error) {
	influxdb, err := NewInfluxdb(config.influxConfig)
	if err != nil {
		return nil, err
	}
	return &Downsampler{Metrics: metrics, Configs: config, DbService: NewDbService(config), Influxdb: influxdb}, nil
}

func (ds *Downsampler) Downsample(Query *QueryObject) error {
	var counter int = 1
	// initiate
	if Query.CompletedSample < Query.BeginSample {
		Query.CompletedSample = Query.BeginSample
	}

	// loop until EndSample is not reached
	for Query.EndSample > Query.CompletedSample {
		log.Printf("Loop:: %d\n", counter)
		ds.Metrics.SetDuration(Query.StartedAt)

		completed := Query.CompletedSample + int64(ds.Configs.JobConfig.SampleWindow)
		// adding up sampling window might cross the end sample
		if completed > Query.EndSample {
			completed = Query.EndSample
		}

		// call influxdb
		written, err := ds.Influxdb.Execute(*Query, completed)
		if err != nil {
			// update the progress in case of error
			err2 := ds.DbService.MarkProgress(*Query)
			if err2 != nil {
				return err2
			}
			return err
		}

		// update variables
		log.Printf("Points written: %0d\n", written)
		ds.Metrics.PointsWrittenByJob(int64(written))

		Query.CompletedSample = completed
		Query.UpdatedAt = time.Now().Unix()
		Query.Note = "Updated with progress timestamp"

		// update only at the offset or when it's complete
		if counter%ds.Configs.JobConfig.ProgressUpdateOffsetWindows == 0 || completed == Query.EndSample {
			err := ds.DbService.MarkProgress(*Query)
			if err != nil {
				return err
			}
		}

		counter++
		ds.rest()
	}

	return nil
}

func (ds *Downsampler) rest() {
	time.Sleep(time.Duration(ds.Configs.JobConfig.RestBetweenWritesMs) * time.Millisecond)
}
