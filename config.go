package main

import (
	"os"
	"strconv"
)

type InfluxDbConfig struct {
	Source      *InfluxDbClientConfig
	Destination *InfluxDbClientConfig
}

type InfluxDbClientConfig struct {
	Url      string
	User     string
	Password string
}

type JobConfig struct {
	// seconds by default
	JobName                     string
	SleepIdle                   int64
	RetryErrorInterval          int64
	SampleWindow                int64
	ProgressUpdateOffsetWindows int // = ProgressUpdateOffsetWindows * SampleWindow s
	RestBetweenWritesMs         int64
}

type MetricsConfig struct {
	Host     string
	Database string
	Username string
	Password string
}

type Config struct {
	Environment   string
	RepoType      string
	DbTablePrefix string
	influxConfig  *InfluxDbConfig
	JobConfig     *JobConfig
	MetricsConfig *MetricsConfig
}

func NewConfig(jobName string) (*Config, error) {
	SleepIdle, err := strconv.ParseInt(os.Getenv("JOB_SLEEP_IDLE_S"), 10, 64)
	var c *Config
	if err != nil {
		return c, err
	}

	RetryErrorInterval, err := strconv.ParseInt(os.Getenv("JOB_ERROR_RETRY_INTERVAL_S"), 10, 64)
	if err != nil {
		return c, err
	}

	SampleWindow, err := strconv.ParseInt(os.Getenv("JOB_METRICS_SAMPLE_WINDOW_S"), 10, 64)
	if err != nil {
		return c, err
	}

	ProgressUpdateOffsetWindows, err := strconv.ParseInt(os.Getenv("JOB_PROGRESS_UPDATE_OFFSET_WINDOWS"), 10, 32)
	if err != nil {
		return c, err
	}

	RestBetweenWrites, err := strconv.ParseInt(os.Getenv("JOB_REST_BETWEEN_WRITES_MS"), 10, 64)
	if err != nil {
		return c, err
	}

	c = &Config{
		os.Getenv("ENVIRONMENT"),
		os.Getenv("REPO_TYPE"),
		os.Getenv("DB_TABLE_PREFIX"),
		&InfluxDbConfig{
			&InfluxDbClientConfig{
				Url:      os.Getenv("INFLUXDB_SRC_URL"),
				User:     os.Getenv("INFLUXDB_SRC_USERNAME"),
				Password: os.Getenv("INFLUXDB_SRC_PASSWORD"),
			},
			&InfluxDbClientConfig{
				Url:      os.Getenv("INFLUXDB_DEST_URL"),
				User:     os.Getenv("INFLUXDB_DEST_USERNAME"),
				Password: os.Getenv("INFLUXDB_DEST_PASSWORD"),
			},
		},
		&JobConfig{
			JobName:                     jobName,
			SleepIdle:                   SleepIdle,
			RetryErrorInterval:          RetryErrorInterval,
			SampleWindow:                SampleWindow,
			ProgressUpdateOffsetWindows: int(ProgressUpdateOffsetWindows),
			RestBetweenWritesMs:         RestBetweenWrites,
		},
		&MetricsConfig{
			Host:     os.Getenv("METRICS_HOST"),
			Database: os.Getenv("METRICS_DATABASE"),
			Username: os.Getenv("METRICS_USERNAME"),
			Password: os.Getenv("METRICS_PASSWORD"),
		},
	}

	return c, nil
}
