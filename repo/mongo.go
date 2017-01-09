package repo

import (
	"github.com/mijara/statspout/data"
	"gopkg.in/mgo.v2"
)

type Mongo struct {
	session *mgo.Session
}

func NewMongo() (*Mongo, error) {
	session, err := mgo.Dial("localhost:27017")
	if err != nil {
		return nil, err
	}

	return &Mongo{
		session: session,
	}, nil
}

func (mongo *Mongo) Push(stats *statspout.Stats) error {
	c := mongo.session.DB("test").C("stats")

	err := c.Insert(stats)
	if err != nil {
		return err
	}

	return nil
}

func (mongo *Mongo) Close() {
	mongo.session.Close()
}
