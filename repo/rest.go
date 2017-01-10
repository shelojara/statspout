package repo

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

// instance of this repository, due to the handler callback limitations.
var rest Rest

func NewRest(address, path string) (*Rest, error) {
	http.HandleFunc(checkAndFixPrefixSlash(path), handler)

	rest.stats = map[string]statspout.Stats{}

	go serveRest(address)

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

func CreateRestOpts() map[string]*string {
	return map[string]*string {
		"address": flag.String(
			"rest.address",
			":8080",
			"Address on which the Rest HTTP Server will publish data"),

		"path": flag.String(
			"rest.path",
			"/stats",
			"Path on which data is served."),
	}
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
