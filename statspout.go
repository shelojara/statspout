package main

import (
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/mijara/statspout/repo"
	"github.com/mijara/statspout/backend"
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
	// start the Docker Endpoint.
	endpoint, err := backend.NewEndpointUnix("unix:///var/run/docker.sock")
	if err != nil {
		log.Fatal(err)
	}

	// start the Repo.
	repository, err := repo.NewInfluxDB("http://localhost:8086", "allstats")
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
		doneChannels = append(doneChannels, backend.Query(endpoint, &containers[i], repository, 5 * time.Second))
	}

	// graceful Ctrl-C quit.
	gracefulQuitInterrupt(doneChannels)
}
