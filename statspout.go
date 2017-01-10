package main

import (
	"log"
	"os"
	"os/signal"
	"time"
	"flag"
	"errors"

	"github.com/mijara/statspout/repo"
	"github.com/mijara/statspout/backend"
)

var (
	// seconds between each stat, in seconds. Minimum is 1 second.
	interval = flag.Int(
		"interval",
		5,
		"Interval between each stats query.")

	// which repository to use.
	repository = flag.String(
		"repository",
		"stdout",
		"One of: stdout, mongodb, prometheus, influxdb, rest.")

	// Which mode use for the connection.
	mode = flag.String(
		"mode",
		"socket",
		"Mode to create the client: socket, http, tls.")

	modeSocketPath = flag.String(
		"socket.path",
		"/var/run/docker.sock",
		"Unix socket to connect to Docker.")

	modeHTTPAddress = flag.String(
		"http.address",
		"localhost:4243",
		"Docker API Address.")

	modeTLSAddress = flag.String(
		"tls.address",
		"localhost:4243",
		"Docker API Address.")

	modeTLSCert = flag.String(
		"tls.cert",
		"",
		"TLS Certificate.")

	modeTLSKey = flag.String(
		"tls.key",
		"",
		"TLS Key.")

	modeTLSCA = flag.String(
		"tls.ca",
		"",
		"TLS CA.")

	// specific maps of options.
	influxDBOpts   = repo.CreateInfluxDBOpts()
	mongoDBOpts    = repo.CreateMongoOpts()
	prometheusOpts = repo.CreatePrometheusOpts()
	restOpts       = repo.CreateRestOpts()
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
	flag.Parse()

	if *interval < 1 {
		log.Fatal("Interval cannot be less than 1.")
	}

	// start the Docker Endpoint.
	endpoint, err := createClientFromFlags()
	if err != nil {
		log.Fatal(err)
	}

	// start the Repo.
	repository, err := createRepositoryFromFlags()
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

func createRepositoryFromFlags() (repo.Interface, error) {
	switch *repository {
	case "stdout":
		return repo.NewStdout(), nil
	case "mongodb":
		return repo.NewMongo(
			*mongoDBOpts["address"],
			*mongoDBOpts["database"],
			*mongoDBOpts["collection"],
		)
	case "prometheus":
		return repo.NewPrometheus(
			*prometheusOpts["address"],
		)
	case "influxdb":
		return repo.NewInfluxDB(
			*influxDBOpts["address"],
			*influxDBOpts["database"],
		)
	case "rest":
		return repo.NewRest(
			*restOpts["address"],
			*restOpts["path"],
		)
	}

	return nil, errors.New("Unknown repository: " + *repository)
}

func createClientFromFlags() (*backend.Endpoint, error) {
	switch *mode {
	case "socket":
		return backend.NewUnixEndpoint(*modeSocketPath)
	case "http":
		return backend.NewHTTPEndpoint(*modeHTTPAddress)
	case "tls":
		return backend.NewTLSEndpoint(*modeTLSAddress, *modeTLSCert, *modeTLSKey, *modeTLSCA)
	}

	return nil, errors.New("Unknown mode: " + *mode)
}
