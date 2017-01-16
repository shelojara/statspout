package backend

import (
	"net/http/httputil"
	"errors"
	"fmt"
	"net"
	"net/http"
	"bufio"
	"encoding/json"
	"os"
	"time"
	"io/ioutil"

	"github.com/mijara/statspout/repo"
	"github.com/mijara/statspout/stats"
)

const (
	DAEMONS     = 10
	STATS_QUERY = "/containers/%s/stats?stream=0"
)

type Backend struct {
	service *Service
	clients chan *httputil.ClientConn
	repo    repo.Interface
	exit    bool
}

// Work to process by daemons.
type Workload struct {
	connection *httputil.ClientConn // connection on which the request is going to be made.
	name       string               // name of the docker container to request.
}

type CpuUsage struct {
	Total  uint64   `json:"total_usage"`
	PerCpu []uint64 `json:"percpu_usage"`
}

type CpuStats struct {
	Usage          CpuUsage `json:"cpu_usage"`
	SystemCpuUsage uint64   `json:"system_cpu_usage"`
}

type MemoryStats struct {
	Usage uint64 `json:"usage"`
	Limit uint64 `json:"limit"`
}

type ContainerStats struct {
	Cpu    CpuStats `json:"cpu_stats"`
	PreCpu CpuStats `json:"precpu_stats"`

	Memory MemoryStats `json:"memory_stats"`

	Read time.Time `json:"read"`
}

type Container struct {
	Names []string `json:"Names"`
}

func New(repo repo.Interface, http bool, address string) (*Backend, error) {
	cli := &Backend{
		repo: repo,
	}

	cli.service = NewService(DAEMONS, cli.process, cli.onError)
	cli.clients = make(chan *httputil.ClientConn, DAEMONS)

	for i := 0; i < DAEMONS; i++ {
		var (
			conn net.Conn
			err  error
		)

		if http {
			conn, err = net.Dial("tcp", address)
		} else {
			conn, err = net.Dial("unix", address)
		}

		if err != nil {
			return nil, err
		}

		cli.clients <- httputil.NewClientConn(conn, nil)
	}

	return cli, nil
}

func (cli *Backend) Query(name string) {
	conn := <-cli.clients

	cli.service.Send(Workload{
		connection: conn,
		name:       name,
	})

	cli.clients <- conn
}

func (cli *Backend) GetContainers() ([]string, error) {
	conn, err := net.Dial("unix", "/var/run/docker.sock")
	if err != nil {
		return nil, err
	}

	c := httputil.NewClientConn(conn, nil)

	req, err := http.NewRequest("GET", "/containers/json", nil)
	if err != nil {
		return nil, err
	}

	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var containers []Container
	json.Unmarshal(body, &containers)

	names := make([]string, len(containers))

	for i, container := range containers {
		name := container.Names[0]
		names[i] = name[1:]
	}

	return names, nil
}

func (cli *Backend) Close() {
	cli.exit = true

	cli.service.Close()

	for i := 0; i < DAEMONS; i++ {
		conn := <-cli.clients
		conn.Close()
	}
}

func (cli *Backend) process(v interface{}) error {
	// client wants to exit, ignore workload.
	if cli.exit {
		return nil
	}

	wl, ok := v.(Workload)
	if !ok {
		return errors.New(fmt.Sprintf("This is not a workload %T", v))
	}

	req, err := http.NewRequest("GET", fmt.Sprintf(STATS_QUERY, wl.name), nil)
	if err != nil {
		return err
	}

	res, err := wl.connection.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	reader := bufio.NewReader(res.Body)
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			break
		}

		container := &ContainerStats{}
		json.Unmarshal(line, container)

		if container.Read == (time.Time{}) {
			return errors.New("Not a valid container: " + wl.name)
		}

		cli.repo.Push(&stats.Stats{
			MemoryPercent: calcMemoryPercent(container),
			CpuPercent:    calcCpuPercent(container),
			MemoryUsage:   container.Memory.Usage,
			Timestamp:     container.Read,
			Name:          wl.name,
		})
	}

	return nil
}

func (cli *Backend) onError(err error) {
	fmt.Fprintln(os.Stderr, err)
}
