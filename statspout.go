package statspout

import (
	"os"
	"os/signal"
	"runtime"
	"time"

	"github.com/mijara/statspout/backend"
	"github.com/mijara/statspout/log"
	"github.com/mijara/statspout/opts"
)

func loop(client *backend.Client, containers map[string]backend.Container) {
	ticker := time.NewTicker(time.Duration(opts.GetOpts().Interval) * time.Second)

	closeC := make(chan os.Signal, 1)
	signal.Notify(closeC, os.Interrupt, os.Kill)

	// initial loop.
	for name := range containers {
		if !contains(opts.GetOpts().Ignore, name) {
			client.Query(containers[name])
		}
	}

	for {
		select {
		case <-closeC:
			log.Info.Printf("Stopping: closing Goroutines and Clients. Please wait...")
			ticker.Stop()
			return
		case <-ticker.C:
			// query containers.
			for name := range containers {
				if !contains(opts.GetOpts().Ignore, name) {
					client.Query(containers[name])
				}
			}
		}
	}
}

func inspect() {
	ticker := time.NewTicker(10 * time.Second)

	for range ticker.C {
		n := runtime.NumGoroutine()
		log.Debug.Printf("Goroutines: %d", n)
	}
}

func Start(cfg *opts.Config) {
	opts.GetOpts().Parse()

	if opts.GetOpts().Interval < 1 {
		log.Error.Fatal("Interval cannot be less than 1.")
	}

	// start the Repo.
	repository, err := opts.CreateRepositoryFromFlags(cfg)
	if err != nil {
		log.Error.Fatal(err)
	}
	defer repository.Close()

	// start the Docker Endpoint.
	client, err := opts.CreateClientFromFlags(repository)
	if err != nil {
		log.Error.Fatal(err)
	}

	// get containers.
	containers, err := client.GetContainers()
	if err != nil {
		log.Error.Fatal(err)
	}

	// small goroutine inspector.
	go inspect()

	log.Info.Printf("Statspout started: %d daemons, %d interval, %s mode, %s repo",
		opts.GetOpts().Daemons,
		opts.GetOpts().Interval,
		opts.GetOpts().Mode.Name,
		opts.GetOpts().Repository)

	client.StartMonitor(containers)

	// loop indefinitely until interrupt is received.
	loop(client, containers)

	// close all connections and goroutines.
	client.Close()

	// force exit
	// TODO: this is needed because EventsMonitor is not able to exit gracefully when reading the HTTP
	// TODO: stream Docker API, I don't know how to fix it at this time.
	os.Exit(0)
}
