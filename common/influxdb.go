package common

import (
	"flag"

	"github.com/mijara/statspout/repo"
	"github.com/influxdata/influxdb/client/v2"
	"github.com/mijara/statspout/stats"
)

type InfluxDB struct {
	client   client.Client
	database string
}

type InfluxOpts struct {
	Address  string
	Database string
}

// Creates a new InfluxDB repository.
func NewInfluxDB(opts *InfluxOpts) (*InfluxDB, error) {
	c, err := client.NewHTTPClient(client.HTTPConfig{Addr: opts.Address})
	if err != nil {
		return nil, err
	}

	return &InfluxDB{
		database: opts.Database,
		client:   c,
	}, nil
}

func (*InfluxDB) Create(v interface{}) (repo.Interface, error) {
	return NewInfluxDB(v.(*InfluxOpts))
}

func (influx *InfluxDB) Push(s *stats.Stats) error {
	if err := influx.pushResource(s, "cpu_usage", s.CpuPercent); err != nil {
		return err
	}

	if err := influx.pushResource(s, "mem_usage", s.MemoryPercent); err != nil {
		return err
	}

	if err := influx.pushResource(s, "tx_bytes", s.TxBytesTotal); err != nil {
		return err
	}

	if err := influx.pushResource(s, "rx_bytes", s.RxBytesTotal); err != nil {
		return err
	}

	if err := influx.pushResource(s, "blkio_read", s.BlockIOBytesRead); err != nil {
		return err
	}

	if err := influx.pushResource(s, "blkio_write", s.BlockIOBytesWrite); err != nil {
		return err
	}

	return nil
}

func CreateInfluxDBOpts() *InfluxOpts {
	o := &InfluxOpts{}

	flag.StringVar(&o.Address,
		"influxdb.address",
		"http://localhost:8086",
		"Address of the InfluxDB Endpoint")

	flag.StringVar(&o.Database,
		"influxdb.database",
		"statspout",
		"Database to store data")

	return o
}

func (*InfluxDB) Name() string {
	return "influxdb"
}

func (influx *InfluxDB) Close() {
	influx.client.Close()
}

func (influx *InfluxDB) Clear(name string) {
	// not used.
}

// Pushes certain a single value to the database, using the resource as the name and
// the name of the container as a tag.
func (influx *InfluxDB) pushResource(s *stats.Stats, resource string, value interface{}) error {
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  influx.database,
		Precision: "s",
	})
	if err != nil {
		return err
	}

	tags := map[string]string{"container": s.Name}
	fields := map[string]interface{}{"value": value}

	pt, err := client.NewPoint(resource, tags, fields, s.Timestamp)
	if err != nil {
		return err
	}

	bp.AddPoint(pt)

	err = influx.client.Write(bp)
	if err != nil {
		return err
	}

	return nil
}
