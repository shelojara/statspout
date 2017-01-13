package opts

import (
	"flag"
	"errors"

	"github.com/mijara/statspout/repo"
	"github.com/mijara/statspout/backend"
	"github.com/mijara/statspout/common"
)

// Structure to hold different options given by the client.
type options struct {
	Interval   int    // Seconds between each stats query.
	Repository string // Which repository to use.

	Mode struct {
		Name string // Client mode name

		Socket struct {
			Path string // Unix socket to connect Docker
		}

		HTTP struct {
			Address string // Docker API address
		}

		TLS struct {
			Address string // Docker API address
			Cert    string // TLS certificate
			Key     string // TLS key
			CA      string // TLS CA
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

	flag.StringVar(&i.Repository,
		"repository",
		"stdout",
		"One of: stdout, mongodb, prometheus, influxdb, rest.")

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

	flag.StringVar(&i.Mode.TLS.Address,
		"tls.address",
		"localhost:4243",
		"Docker API Address.")

	flag.StringVar(&i.Mode.TLS.Cert,
		"tls.cert",
		"",
		"TLS Certificate.")

	flag.StringVar(&i.Mode.TLS.Key,
		"tls.key",
		"",
		"TLS Key.")

	flag.StringVar(&i.Mode.TLS.CA,
		"tls.ca",
		"",
		"TLS CA.")

	/*
	common.CreateInfluxDBOpts(&i.Influx)
	common.CreateMongoOpts(&i.Mongo)
	common.CreatePrometheusOpts(&i.Prometheus)
	common.CreateRestOpts(&i.Rest)
	*/

	// flag.Parse()

	return i
}

func (*options) Parse() {
	flag.Parse()
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
func CreateClientFromFlags() (*backend.Endpoint, error) {
	switch GetOpts().Mode.Name {
	case "socket":
		return backend.NewUnixEndpoint(GetOpts().Mode.Socket.Path)
	case "http":
		return backend.NewHTTPEndpoint(GetOpts().Mode.HTTP.Address)
	case "tls":
		return backend.NewTLSEndpoint(GetOpts().Mode.TLS.Address, GetOpts().Mode.TLS.Cert, GetOpts().Mode.TLS.Key, GetOpts().Mode.TLS.CA)
	}

	return nil, errors.New("Unknown mode: " + GetOpts().Mode.Name)
}
