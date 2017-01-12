package opts

import (
	"flag"
	"errors"

	"github.com/mijara/statspout/repo"
	"github.com/mijara/statspout/backend"
	"github.com/mijara/statspout/common"
)

type Opts struct {
	Interval   int
	Repository string

	Mode struct {
		Name string

		Socket struct {
			Path string
		}

		HTTP struct {
			Address string
		}

		TLS struct {
			Address string
			Cert    string
			Key     string
			CA      string
		}
	}

	Influx     common.InfluxOpts
	Mongo      common.MongoOpts
	Rest       common.RestOpts
	Prometheus common.PrometheusOpts
}

var (
	i *Opts
)

func GetOpts() *Opts {
	if i != nil {
		return i
	}

	i = &Opts{}

	// seconds between each stat, in seconds. Minimum is 1 second.
	flag.IntVar(&i.Interval,
		"interval",
		5,
		"Interval between each stats query.")

	// which repository to use.
	flag.StringVar(&i.Repository,
		"repository",
		"stdout",
		"One of: stdout, mongodb, prometheus, influxdb, rest.")

	// Which mode use for the connection.
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

	// specific maps of options.
	common.CreateInfluxDBOpts(&i.Influx)
	common.CreateMongoOpts(&i.Mongo)
	common.CreatePrometheusOpts(&i.Prometheus)
	common.CreateRestOpts(&i.Rest)

	flag.Parse()

	return i
}

func CreateRepositoryFromFlags() (repo.Interface, error) {
	switch GetOpts().Repository {
	case "stdout":
		return common.NewStdout(), nil
	case "mongodb":
		return common.NewMongo(i.Mongo)
	case "prometheus":
		return common.NewPrometheus(i.Prometheus)
	case "influxdb":
		return common.NewInfluxDB(i.Influx)
	case "rest":
		return common.NewRest(i.Rest)
	}

	return nil, errors.New("Unknown repository: " + i.Repository)
}

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
