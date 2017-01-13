package opts

import (
	"github.com/mijara/statspout/repo"
	"github.com/prometheus/common/log"
)

type Pair struct {
	Repository repo.Interface
	Options    interface{}
}

type Config struct {
	Repositories map[string]*Pair
}

func NewConfig() *Config {
	return &Config{
		Repositories: make(map[string]*Pair),
	}
}

// Adds a repository to the registry, it should not collide with other repository names.
func (cfg *Config) AddRepository(repo repo.Interface, options interface{}) {
	// check that the name does not exist, this is considered an error because could lead
	// to mistakes.
	if _, ok := cfg.Repositories[repo.Name()]; ok {
		log.Fatal("Repository name taken: " + repo.Name())
	}

	if repo.Name() == "" {
		log.Fatal("Got empty repository name.")
	}

	cfg.Repositories[repo.Name()] = &Pair{
		Repository: repo,
		Options: options,
	}
}
