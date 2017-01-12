package common

import (
	"flag"

	"github.com/influxdata/influxdb/client/v2"
	"github.com/mijara/statspout/data"
)

type InfluxDB struct {
	client   client.Client
	database string
}

type InfluxOpts struct {
	Address  string
	Database string
}

func NewInfluxDB(opts InfluxOpts) (*InfluxDB, error) {
	c, err := client.NewHTTPClient(client.HTTPConfig{Addr: opts.Address})
	if err != nil {
		return nil, err
	}

	return &InfluxDB{
		database: opts.Database,
		client:   c,
	}, nil
}

func (influx *InfluxDB) Push(stats *statspout.Stats) error {
	influx.pushResource(stats, "cpu_usage", stats.CpuPercent)
	influx.pushResource(stats, "mem_usage", stats.MemoryPercent)
	return nil
}

func (influx *InfluxDB) pushResource(stats *statspout.Stats, resource string, value float64) error {
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  influx.database,
		Precision: "s",
	})
	if err != nil {
		return err
	}

	tags := map[string]string{
		"container": stats.Name,
	}

	fields := map[string]interface{}{
		"value":   value,
	}

	pt, err := client.NewPoint(resource, tags, fields, stats.Timestamp)
	bp.AddPoint(pt)

	err = influx.client.Write(bp)
	if err != nil {
		return err
	}

	return nil
}

func (influx *InfluxDB) Close() {
	influx.client.Close()
}

func CreateInfluxDBOpts(opts *InfluxOpts) {
	flag.StringVar(&opts.Address,
		"influxdb.address",
		"http://localhost:8086",
		"Address of the InfluxDB Endpoint")

	flag.StringVar(&opts.Database,
		"influxdb.database",
		"statspout",
		"Database to store data")
}
