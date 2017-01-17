package opts

import (
	"strings"
	"errors"
	"flag"

	"github.com/mijara/statspout/repo"
	"github.com/mijara/statspout/backend"
	"github.com/mijara/statspout/common"
)

// Structure to hold different options given by the client.
type options struct {
	Interval   int      // Seconds between each stats query.
	Repository string   // Which repository to use.
	Daemons    int      // Number of daemons to handle requests.
	Ignore     []string // Container names to ignore, as an array.

	ignoreBuff string // Container names to ignore, separated by comma.

	Mode struct {
		Name string // Client mode name

		Socket struct {
			Path string // Unix socket to connect Docker
		}

		HTTP struct {
			Address string // Docker API address
		}
	}

	Influx     common.InfluxOpts     // Influx specific options
	Mongo      common.MongoOpts      // Mongo specific options.
	Rest       common.RestOpts       // Rest specific options.
	Prometheus common.PrometheusOpts // Prometheus specific options.
}

// Single instance of this package.
var i *options

// Gets options struct instance.
func GetOpts() *options {
	if i != nil {
		return i
	}

	i = &options{}

	flag.IntVar(&i.Interval,
		"interval",
		5,
		"Interval between each stats query.")

	flag.IntVar(&i.Daemons,
		"daemons",
		10,
		"Number of daemons to handle requests.")

	flag.StringVar(&i.Repository,
		"repository",
		"stdout",
		"One of: stdout, mongodb, prometheus, influxdb, rest.")

	flag.StringVar(&i.ignoreBuff,
		"ignore",
		"",
		"Repository names to ignore, separated by comma.")

	flag.StringVar(&i.Mode.Name,
		"mode",
		"socket",
		"Mode to create the client: socket, http, tls.")

	flag.StringVar(&i.Mode.Socket.Path,
		"socket.path",
		"/var/run/docker.sock",
		"Unix socket to connect to Docker.")

	flag.StringVar(&i.Mode.HTTP.Address,
		"http.address",
		"localhost:4243",
		"Docker API Address.")

	return i
}

func (*options) Parse() {
	flag.Parse()

	names := strings.Split(i.ignoreBuff, ",")
	i.Ignore = make([]string, 0)

	for _, name := range names {
		if name != "" {
			i.Ignore = append(i.Ignore, name)
		}
	}
}

// Creates the repository from the options given by the client.
func CreateRepositoryFromFlags(cfg *Config) (repo.Interface, error) {
	for name, b := range cfg.Repositories {
		if name == GetOpts().Repository {
			return b.Repository.Create(b.Options)
		}
	}

	return nil, errors.New("Unknown repository: " + i.Repository)
}

// Creates the client from the options given by the client.
func CreateClientFromFlags(repo repo.Interface) (*backend.Client, error) {
	switch GetOpts().Mode.Name {
	case "socket":
		return backend.New(repo, false, GetOpts().Mode.Socket.Path, GetOpts().Daemons)
	case "http":
		return backend.New(repo, true, GetOpts().Mode.HTTP.Address, GetOpts().Daemons)
	}

	return nil, errors.New("Unknown mode: " + GetOpts().Mode.Name)
}
