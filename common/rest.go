package common

import (
	"net/http"
	"encoding/json"
	"flag"

	"github.com/prometheus/common/log"
	"github.com/mijara/statspout/repo"
	"github.com/mijara/statspout/stats"
)

type Rest struct {
	registry map[string]stats.Stats
}

type RestOpts struct {
	Address string
	Path    string
}

// instance of this repository, due to the handler callback limitations.
var rest Rest

func (*Rest) Name() string {
	return "rest"
}

func (*Rest) Create(v interface{}) (repo.Interface, error) {
	return NewRest(v.(*RestOpts))
}

func NewRest(opts *RestOpts) (*Rest, error) {
	http.HandleFunc(checkAndFixPrefixSlash(opts.Path), handler)

	rest.registry = map[string]stats.Stats{}

	go serveRest(opts.Address)

	return &rest, nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(rest.asListOfValues())
}

func (rest *Rest) asListOfValues() []stats.Stats {
	var list []stats.Stats

	for _, value := range rest.registry {
		list = append(list, value)
	}

	return list
}

func (rest *Rest) Push(s *stats.Stats) error {
	rest.registry[s.Name] = *s
	return nil
}

func (rest *Rest) Close() {
}

func (rest *Rest) Clear(name string) {
	delete(rest.registry, name)
}

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

func serveRest(address string) {
	log.Fatal(http.ListenAndServe(address, nil))
}

func checkAndFixPrefixSlash(path string) string {
	if len(path) == 0 {
		return "/"
	}

	if path[0] == '/' {
		return path
	}

	return "/" + path
}
