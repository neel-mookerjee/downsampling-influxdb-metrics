package main

import (
	"github.com/influxdata/influxdb/client/v2"
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"time"
)

type InfluxdbInterface interface {
	Execute(Query QueryObject, finishAtNew int64) (int, error)
}

type Influxdb struct {
	Config              *InfluxDbConfig
	InfluxdbReadClient  client.Client
	InfluxdbWriteClient client.Client
}

func NewInfluxdb(config *InfluxDbConfig) (*Influxdb, error) {
	i := Influxdb{}
	i.Config = config
	clntW, err := i.getWriteClient()
	if err != nil {
		return &i, err
	}
	clntR, err := i.getReadClient()
	if err != nil {
		return &i, err
	}

	i.InfluxdbWriteClient = clntW
	i.InfluxdbReadClient = clntR

	return &i, nil
}

func (i *Influxdb) getReadClient() (client.Client, error) {
	// Create a new influxdb HTTPClient
	return i.getClient(i.Config.Source)
}

func (i *Influxdb) getWriteClient() (client.Client, error) {
	// Create a new influxdb HTTPClient
	return i.getClient(i.Config.Destination)
}

func (i *Influxdb) getClient(config *InfluxDbClientConfig) (client.Client, error) {
	// Create a new influxdb HTTPClient
	return client.NewHTTPClient(client.HTTPConfig{
		Addr:     config.Url,
		Username: config.User,
		Password: config.Password,
	})
}

func (i *Influxdb) Execute(Query QueryObject, finishAtNew int64) (int, error) {
	Query_to_use := i.buildQuery(Query, finishAtNew)
	//Query_to_use := `SELECT COUNT("duration") AS "count_duration",MIN("duration") AS "min_duration",MAX("duration") AS "max_duration",PERCENTILE("duration",90) AS "p90_duration" FROM "omni"."autogen"."nsa_duration" WHERE time >= 1505886000000000000 AND time <= 1505886600000000000 GROUP BY time(600s),"az","code","consumer","entity","group","host","region"`
	log.Printf("Query: %s\n", Query_to_use)

	res, err := i.query(Query_to_use)
	if err != nil {
		return 0, err
	}

	// the output object is received
	//written, err := strconv.Atoi(string(res[0].Series[0].Values[0][1].(json.Number)))
	written, err := i.write(Query, res)

	if err != nil {
		return written, err
	}

	return written, nil
}

func (i *Influxdb) write(Query QueryObject, res []client.Result) (int, error) {
	// Create a new point batch
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:        Query.Database,
		Precision:       "ns",
		RetentionPolicy: "histdownsample2",
	})
	if err != nil {
		return 0, err
	}

	// add points to batch
	for _, row := range res[0].Series {
		//fmt.Println("row:", row)

		for _, values := range row.Values {
			//fmt.Println("columns:", row.Columns)
			//fmt.Println("values:", values)

			fields := map[string]interface{}{}

			for i, col := range row.Columns {
				//fmt.Println("value:", values[i])
				fields[col] = values[i]
			}

			time, err := time.Parse(time.RFC3339, values[0].(string))
			if err != nil {
				return 0, err
			}

			//fmt.Println("fields:", fields)
			//fmt.Println("tags:", row.Tags)

			pt, err := client.NewPoint(Query.TargetMeasurement, row.Tags, fields, time)
			if err != nil {
				return 0, err
			}
			//fmt.Println("point:", pt)

			bp.AddPoint(pt)
		}
	}

	//clnt, err := i.getWriteClient()
	clnt := i.InfluxdbWriteClient
	defer clnt.Close()
	// Write the batch
	if err := clnt.Write(bp); err != nil {
		return 0, err
	}

	return len(bp.Points()), nil
}

func (i *Influxdb) buildQuery(Query QueryObject, finishAtNew int64) string {
	Query_to_use := strings.Replace(Query.Query, "$ds_end_ts", strconv.FormatInt(finishAtNew*1000000000, 10), 1)
	Query_to_use = strings.Replace(Query_to_use, "$ds_start_ts", strconv.FormatInt(Query.CompletedSample*1000000000, 10), 1)

	return Query_to_use
}

func (i *Influxdb) query(cmd string) (res []client.Result, err error) {
	q := client.Query{
		Command: cmd,
	}

	//clnt, err := i.getReadClient()
	//defer clnt.Close()

	clnt := i.InfluxdbReadClient

	if response, err := clnt.Query(q); err == nil {
		if response.Error() != nil {
			return res, response.Error()
		}
		res = response.Results
	} else {
		return res, err
	}

	//log.Println(res)
	return res, nil
}
