package main

import (
	"os"
	"os/signal"
	"time"
	"log"

	"github.com/mijara/statspout/opts"
	"github.com/mijara/statspout/backend"
)

func gracefulQuitInterrupt(doneChannels []chan bool) {
	// graceful Ctrl-C quit.
	closeC := make(chan os.Signal, 1)
	signal.Notify(closeC, os.Interrupt, os.Kill)

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
	if opts.GetOpts().Interval < 1 {
		log.Fatal("Interval cannot be less than 1.")
	}

	// start the Docker Endpoint.
	endpoint, err := opts.CreateClientFromFlags()
	if err != nil {
		log.Fatal(err)
	}

	// start the Repo.
	repository, err := opts.CreateRepositoryFromFlags()
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
			time.Duration(opts.GetOpts().Interval)*time.Second))
	}

	// graceful Ctrl-C quit.
	gracefulQuitInterrupt(doneChannels)
}
