package main

import (
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/mijara/statspout/repo"
	"github.com/mijara/statspout/backend"
	"flag"
)

var (
	// which repository to use.
	repository = flag.String(
		"repository",
		"stdout",
		"One of: stdout, mongodb, prometheus, influxdb")

	// seconds between each stat, in seconds. Minimum is 1 second.
	interval = flag.Int(
		"interval",
		5,
		"Interval between each stats query")

	// specific maps of options.
	influxDBOpts   = repo.CreateInfluxDBOpts()
	mongoDBOpts    = repo.CreateMongoOpts()
	prometheusOpts = repo.CreatePrometheusOpts()
)

func gracefulQuitInterrupt(doneChannels []chan bool) {
	// graceful Ctrl-C quit.
	closeC := make(chan os.Signal, 1)
	signal.Notify(closeC, os.Interrupt)

	stop := false

	go func() {
		for _ = range closeC {
			for _, done := range doneChannels {
				done <- true
			}

			stop = true
			return
		}
	}()

	for {
		if stop {
			break
		}

		time.Sleep(1 * time.Second)
	}
}

func main() {
	flag.Parse()

	// start the Docker Endpoint.
	endpoint, err := backend.NewEndpointUnix("unix:///var/run/docker.sock")
	if err != nil {
		log.Fatal(err)
	}

	// start the Repo.
	repository, err := getRepositoryObject()
	if err != nil {
		log.Fatal(err)
	}

	defer repository.Close()

	// get containers.
	containers, err := backend.GetContainers(endpoint)
	if err != nil {
		log.Fatal(err)
	}

	// query containers and store done channels to stop each goroutine.
	var doneChannels []chan bool
	for i := 0; i < len(containers); i++ {
		doneChannels = append(doneChannels, backend.Query(endpoint, &containers[i], repository,
			time.Duration(*interval)*time.Second))
	}

	// graceful Ctrl-C quit.
	gracefulQuitInterrupt(doneChannels)
}

func getRepositoryObject() (repo.Interface, error) {
	var r repo.Interface
	var err error

	switch *repository {
	case "stdout":
		r = repo.NewStdout()
	case "mongodb":
		r, err = repo.NewMongo(
			*mongoDBOpts["address"],
			*mongoDBOpts["database"],
			*mongoDBOpts["collection"],
		)
	case "prometheus":
		r, err = repo.NewPrometheus(
			*prometheusOpts["address"],
		)
	case "influxdb":
		r, err = repo.NewInfluxDB(
			*influxDBOpts["address"],
			*influxDBOpts["database"],
		)
	}

	if err != nil {
		return nil, err
	}

	return r, nil
}
