Statspout
=========

Service that retrieves stats from Docker Containers and sends them to some repository (a.k.a DB).

Supported Repositories:

- Stdout `stdout` (for testing)
- MongoDB `mongodb` (using https://github.com/go-mgo/mgo)
- Prometheus `prometheus` (as a scapre source, using https://github.com/prometheus/client_golang)
- InfluxDB `influxdb` (using https://github.com/influxdata/influxdb/tree/master/client)
- RestAPI `rest`


## Usage

As a CLI, run the following on a console:

```
./statspout [-mode=<mode>] [-interval=<interval>] [-repository=<repository>] {options}
```

If no option is given, the program will run on the default Docker Socket, with an interval of 5 seconds and `stdout` as
repository (this is done so you can quickly check what this tool does, without setting a DB or service).


### Top Level Opts:
- `mode`: mode to create the client: `socket`, `http`, `tls`. Default: `socket`
- `interval`: seconds between each stat, in seconds. Minimum is 1 second. Default is `5`.
- `repository`: which repository to use (they're listed in the Supported Repositories list, in special font)
                each repository will bound different options. Default is `stdout`.


### Mode Options

#### Socket

- `socket.path`: unix socket to connect to Docker. Default: `/var/run/docker.sock`

#### HTTP

- `http.address`: Docker API address. Default: `localhost:4243`

#### TLS

- `tls.address`: Docker API address. Default: `localhost:4243`
- `tls.cert`: TLS certificate.
- `tls.key`: TLS key.
- `tls.ca`: TLS CA.


### Specific Repository Options


#### MongoDB
- `mongo.address`: Address of the MongoDB Endpoint. Default: `localhost:27017`
- `mongo.database`: Database for the collection. Default: `statspout`
- `mongo.collection`: Collection for the stats. Default: `stats`


#### Prometheus
- `prometheus.address`: Address on which the Prometheus HTTP Server will publish metrics. Default: `:8080`


#### InfluxDB
- `influxdb.address`: Address of the InfluxDB Endpoint. Default: `http://localhost:8086`
- `influxdb.database`: Database to store data. Default: `statspout`


#### Rest
- `rest.address`: Address on which the Rest HTTP Server will publish data. Default: `:8080`
- `rest.path`: Path on which data is served. Default: `/stats`

## Run as a Docker Container

The container version is available at https://hub.docker.com/r/mijara/statspout/

Run with:

```
docker run -v /var/run/docker.sock:/var/run/docker.sock -p 8080:8080 mijara/statspout [options...]
```
