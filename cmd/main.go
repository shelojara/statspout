package main

import (
	"github.com/mijara/statspout"
	"github.com/mijara/statspout/common"
	"github.com/mijara/statspout/opts"
)

func main() {
	cfg := opts.NewConfig()

	cfg.AddRepository(&common.Stdout{}, nil)

	cfg.AddRepository(&common.Rest{}, common.CreateRestOpts())

	cfg.AddRepository(&common.Prometheus{}, common.CreatePrometheusOpts())
	cfg.AddRepository(&common.InfluxDB{}, common.CreateInfluxDBOpts())
	cfg.AddRepository(&common.Mongo{}, common.CreateMongoOpts())

	statspout.Start(cfg)
}
