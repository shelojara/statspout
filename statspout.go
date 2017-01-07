package main

import (
	"fmt"

	"github.com/fsouza/go-dockerclient"
	"time"
)


func queryContainer(client *docker.Client, container *docker.APIContainers, done <-chan bool) {
	statsC := make(chan *docker.Stats)
	errC := make(chan error, 1)

	go func() {
		errC <- client.Stats(docker.StatsOptions{ID: container.ID, Stats: statsC, Stream: true, Done: done})
		close(errC)
	}()

	for {
		stats, ok := <- statsC
		if !ok {
			break
		}

		fmt.Println(stats.CPUStats.CPUUsage.TotalUsage)
	}
}


func main() {
	client, err := docker.NewClient("unix:///var/run/docker.sock")
	if err != nil {
		panic(err)
	}

	containers, err := client.ListContainers(docker.ListContainersOptions{})
	if err != nil {
		panic(err)
	}

	done := make(chan bool)

	for _, container := range containers {
		go queryContainer(client, &container, done)
	}

	time.Sleep(5 * time.Second)

	done <- true
}
