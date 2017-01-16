package statspout

import (
	"os"
	"log"
	"os/signal"

	"github.com/mijara/statspout/opts"
	"time"
	"github.com/mijara/statspout/backend"
	"fmt"
)

func loop(client *backend.Client, containers []string) {
	ticker := time.NewTicker(time.Duration(opts.GetOpts().Interval) * time.Second)

	closeC := make(chan os.Signal, 1)
	signal.Notify(closeC, os.Interrupt, os.Kill)

	// initial loop.
	for i := 0; i < len(containers); i++ {
		client.Query(containers[i])
	}

	for {
		select {
		case <-closeC:
			fmt.Println(" * Stopping: closing Goroutines and Clients. Please wait...")
			ticker.Stop()
			return
		case <-ticker.C:
			// query containers.
			for i := 0; i < len(containers); i++ {
				client.Query(containers[i])
			}
		}
	}
}

func Start(cfg *opts.Config) {
	opts.GetOpts().Parse()

	if opts.GetOpts().Interval < 1 {
		log.Fatal("Interval cannot be less than 1.")
	}

	// start the Repo.
	repository, err := opts.CreateRepositoryFromFlags(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer repository.Close()

	// start the Docker Endpoint.
	client, err := opts.CreateClientFromFlags(repository)
	if err != nil {
		log.Fatal(err)
	}

	// get containers.
	containers, err := client.GetContainers()
	if err != nil {
		log.Fatal(err)
	}

	// loop indefinitely until interrupt is received.
	loop(client, containers)

	// close all connections and goroutines.
	client.Close()
}
