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
	endpoint, err := backend.NewEndpointUnix("unix:///var/run/docker.sock")
	if err != nil {
		log.Fatal(err)
	}

	containers, err := backend.GetContainers(endpoint)
	if err != nil {
		log.Fatal(err)
	}

	output := repo.Stdout{}

	// query containers and store done channels to stop each goroutine.
	var doneChannels []chan bool
	for i := 0; i < len(containers); i++ {
		doneChannels = append(doneChannels, backend.Query(endpoint, &containers[i], output))
	}

	// graceful Ctrl-C quit.
	gracefulQuitInterrupt(doneChannels)
}
