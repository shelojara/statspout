package repo

import (
	"github.com/influxdata/influxdb/client/v2"

	"github.com/mijara/statspout/data"
	"flag"
)

type InfluxDB struct {
	client client.Client
	database string
}

func NewInfluxDB(address string, database string) (*InfluxDB, error) {
	c, err := client.NewHTTPClient(client.HTTPConfig{Addr: address})
	if err != nil {
		return nil, err
	}

	return &InfluxDB{
		database: database,
		client: c,
	}, nil
}

func (influx *InfluxDB) Push(stats *statspout.Stats) error {
	influx.pushResource(stats, "cpu-usage", stats.CpuPercent)
	influx.pushResource(stats, "mem-usage", stats.MemoryPercent)

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
		"resource": resource,
	}

	fields := map[string]interface{}{
		"value":   value,
	}

	pt, err := client.NewPoint(stats.Name, tags, fields, stats.Timestamp)
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

func CreateInfluxDBFlagsMap() map[string]*string {
	return map[string]*string {
		"address": flag.String(
			"influxdb.address",
			"http://localhost:8086",
			"Address of the InfluxDB Endpoint"),

		"database": flag.String(
			"influxdb.database",
			"dockerstats",
			"Database to store data"),
	}
}
