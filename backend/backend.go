// This package is intended to provide a level of abstraction for container data collection, making it
// easier to switch the Docker Client library.
package backend

import (
	"github.com/mijara/statspout/data"
	"github.com/mijara/statspout/repo"
	"github.com/fsouza/go-dockerclient"
	"time"
)

type Endpoint struct {
	client *docker.Client
}

// Creates an Endpoint holding data connection for the current Client.
func NewEndpointUnix(sockPath string) (*Endpoint, error) {
	client, err := docker.NewClient(sockPath)
	if err != nil {
		return nil, err
	}

	return &Endpoint{
		client: client,
	}, nil
}

// Retrieves and returns all containers data for future reference.
func GetContainers(endpoint *Endpoint) ([]statspout.Container, error) {
	containers, err := endpoint.client.ListContainers(docker.ListContainersOptions{})
	if err != nil {
		return nil, err
	}

	var result []statspout.Container
	for _, container := range containers {
		result = append(result, statspout.Container{
			ID:    container.ID,
			Names: container.Names,
			Image: container.Image,
		})
	}

	return result, nil
}

// Queries all containers using the specific Endpoint Client implementation. Each container
// is queried in a different Goroutine to improve performance.
func Query(endpoint *Endpoint, container *statspout.Container, repo repo.Interface) (chan bool) {
	done := make(chan bool)
	go queryContainer(endpoint.client, container, repo, done)
	return done
}

func queryContainer(client *docker.Client, container *statspout.Container, repo repo.Interface, done chan bool) {
	statsC := make(chan *docker.Stats)
	errC := make(chan error, 1)

	go func() {
		errC <- client.Stats(docker.StatsOptions{ID: container.ID, Stats: statsC, Stream: true, Done: done})
		close(errC)
	}()

	containerStats := statspout.Stats{
		ID: container.ID,
		Name: container.Names[0][1:],
	}

	// receive stats from container, ignore stats that are received in between ticker times.
	// TODO: this may be not the best approach, but we have to test.
	ticker := time.NewTicker(5 * time.Second)

	stats, ok := <-statsC
	if !ok {
		return
	}

	pushStats(&containerStats, repo, stats)

	for {
		stats, ok := <-statsC
		if !ok {
			return
		}

		select {
		case <-ticker.C:
			pushStats(&containerStats, repo, stats)
		default:
		//
		}
	}
}

func pushStats(containerStats *statspout.Stats, repo repo.Interface, stats *docker.Stats) {
	containerStats.Timestamp = stats.Read
	containerStats.CpuPercent = calcCpuPercent(stats)
	containerStats.MemoryPercent = calcMemoryPercent(stats)
	containerStats.MemoryUsage = stats.MemoryStats.Usage

	repo.Push(containerStats)
}

// taken from: https://github.com/portainer/portainer/blob/develop/app/components/stats/statsController.js#L177-L193
func calcCpuPercent(stats *docker.Stats) float64 {
	cpuPercent := 0.0

	cpuDelta := float64(stats.CPUStats.CPUUsage.TotalUsage - stats.PreCPUStats.CPUUsage.TotalUsage)
	systemDelta := float64(stats.CPUStats.SystemCPUUsage - stats.PreCPUStats.SystemCPUUsage)

	if systemDelta > 0.0 && cpuDelta > 0.0 {
		cpuPercent = (cpuDelta / systemDelta) * float64(len(stats.CPUStats.CPUUsage.PercpuUsage)) * 100.0
	}

	return cpuPercent
}

func calcMemoryPercent(stats *docker.Stats) float64 {
	return float64(stats.MemoryStats.Usage) * 100.0 / float64(stats.MemoryStats.Limit)
}
