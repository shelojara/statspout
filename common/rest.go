package common

import (
	"net/http"
	"encoding/json"

	"github.com/mijara/statspout/data"
	"github.com/prometheus/common/log"
	"flag"
)

type Rest struct {
	stats map[string]statspout.Stats
}

type RestOpts struct {
	Address string
	Path    string
}

// instance of this repository, due to the handler callback limitations.
var rest Rest

func NewRest(opts RestOpts) (*Rest, error) {
	http.HandleFunc(checkAndFixPrefixSlash(opts.Path), handler)

	rest.stats = map[string]statspout.Stats{}

	go serveRest(opts.Address)

	return &rest, nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(rest.asListOfValues())
}

func (rest *Rest) asListOfValues() []statspout.Stats {
	var list []statspout.Stats

	for _, value := range rest.stats {
		list = append(list, value)
	}

	return list
}

func (rest *Rest) Push(stats *statspout.Stats) error {
	rest.stats[stats.Name] = *stats
	return nil
}

func (rest *Rest) Close() {
}

func CreateRestOpts(opts *RestOpts) {
	flag.StringVar(&opts.Address,
		"rest.address",
		":8080",
		"Address on which the Rest HTTP Server will publish data")

	flag.StringVar(&opts.Path,
		"rest.path",
		"/stats",
		"Path on which data is served.")
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
