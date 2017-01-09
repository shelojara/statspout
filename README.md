Statspout
=========

Service that retrieves stats from Docker Containers and sends them to some repository (a.k.a DB).

Supported Repositories:

- Stdout [stdout] (for testing)
- MongoDB [mongodb] (using https://github.com/go-mgo/mgo)
- Prometheus [prometheus] (as a scapre source, using https://github.com/prometheus/client_golang)
- InfluxDB [influxdb] (using https://github.com/influxdata/influxdb/tree/master/client)


## Usage

As a CLI, run the following on a console:

```
./statspout -interval=<interval> -repository=<repository> {options}
```

Where:
- `interval`: seconds between each stat, in seconds. Minimum is 1 second. Default is `5`.
- `repository`: which repository to use (they're listed in the Supported Repositories list, in square brackets)
  each repository will bound different options. Default is `stdout`.

## Specific Options

### MongoDB
- `mongo.address`: Address of the MongoDB Endpoint. Default: `localhost:27017`

### Prometheus
- `prometheus.address`: Address on which the Prometheus HTTP Server will publish metrics. Default: `:8080`

### InfluxDB
- `influxdb.address`: Address of the InfluxDB Endpoint. Default: `http://localhost:8086`
- `influxdb.database`: Database to store data. Default: `statspout`
