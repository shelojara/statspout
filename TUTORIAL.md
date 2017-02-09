
## Creating your own Repository

This application is shaped as a framework, with the main file being the piece of code that makes the default "moves".
However, you can completely discard that main file and compile your own with solely the `statspout` framework.

First, take a look at the main file:

```go
// https://github.com/mijara/statspout/blob/master/cmd/main.go

package main

import (
	"github.com/mijara/statspout"
	"github.com/mijara/statspout/common"
	"github.com/mijara/statspout/opts"
)

func main() {
	cfg := opts.NewConfig()

	cfg.AddRepository(&common.Stdout{}, nil)

	cfg.AddRepository(&common.Rest{}, common.CreateRestOpts())

	cfg.AddRepository(&common.Prometheus{}, common.CreatePrometheusOpts())
	cfg.AddRepository(&common.InfluxDB{}, common.CreateInfluxDBOpts())
	cfg.AddRepository(&common.Mongo{}, common.CreateMongoOpts())

	statspout.Start(cfg)
}
```   

In this file we add every repository we want to be shipped with the application, it is a simple piece of code that
creates a configuration (`cfg`) for the application, adds repositories to said config and then fires up the application
passing the `cfg` to `statspout.Start`.

At this point you should be able to remove any module you don't want to be used.

How about creating your repository? A repository is usable by the framework as long as it implements the following
interface:

```go
// Repository Interface, used to define any repository.
// A repository is any service that can have data pushed to it.
// The close method is provided in case the service needs it.
type Interface interface {
	// Creates a new instance of the repository.
	// v holds the set of options for this new instance.
	Create(v interface{}) (Interface, error)

	// Push container stats to this service.
	// The repository should return an error if it's not capable of pushing the stats.
	Push(stats *stats.Stats) error

	// Close the service.
	Close()

	// Canonical name of this repository, used to identify it in the command line flags.
	Name() string
}

```

For example, the most simple repository, Stdout:

```go
package common

import (
	"fmt"

	"github.com/mijara/statspout/repo"
	"github.com/mijara/statspout/stats"
)

type Stdout struct {
}

func (*Stdout) Name() string {
	return "stdout"
}

func (*Stdout) Create(v interface{}) (repo.Interface, error) {
	return NewStdout(), nil
}

func NewStdout() *Stdout {
	return &Stdout{}
}

func (stdout *Stdout) Push(s *stats.Stats) error {
	fmt.Println(s)
	return nil
}

func (stdout *Stdout) Close() {

}
```

and then to include it in the main:

```go
cfg.AddRepository(&common.Stdout{}, nil)
```

That line says that we will register a repository of type `&common.Stdout`, and `nil`. The latter parameter is a bundle
of options for the flag parser (stdout does not use it, hence `nil`), as an example, the rest repository provides this
bundle:

```go
func CreateRestOpts() *RestOpts {
	o := &RestOpts{}

	flag.StringVar(&o.Address,
		"rest.address",
		":8080",
		"Address on which the Rest HTTP Server will publish data")

	flag.StringVar(&o.Path,
		"rest.path",
		"/stats",
		"Path on which data is served.")

	return o
}
```

RestOpts is a simple structure on which we will feed the parsed data into. Then, the same object created gets passed
to the `Create()` method of the repository, there you will have to assert type' it in order to use it (don't worry,
nobody will panic over this).

Finally, whenever your repository is needed, we will call the `Create()` method, in there you'll have to provide us
with an instance of your repository, and that's it!

Notice that your repository **has to be safe to be used concurrently**.
