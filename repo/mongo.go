package repo

import (
	"github.com/mijara/statspout/data"
	"gopkg.in/mgo.v2"
	"flag"
)

type Mongo struct {
	session *mgo.Session
	database string
	collection string
}

func NewMongo(address, database, collection string) (*Mongo, error) {
	session, err := mgo.Dial(address)
	if err != nil {
		return nil, err
	}

	return &Mongo{
		session: session,
		database: database,
		collection: collection,
	}, nil
}

func (mongo *Mongo) Push(stats *statspout.Stats) error {
	c := mongo.session.DB(mongo.database).C(mongo.collection)

	err := c.Insert(stats)
	if err != nil {
		return err
	}

	return nil
}

func (mongo *Mongo) Close() {
	mongo.session.Close()
}

func CreateMongoOpts() map[string]*string {
	return map[string]*string {
		"address": flag.String(
			"mongo.address",
			"localhost:27017",
			"Address of the MongoDB Endpoint"),

		"database": flag.String(
			"mongo.database",
			"statspout",
			"Database for the collection"),

		"collection": flag.String(
			"mongo.collection",
			"stats",
			"Collection for the stats"),
	}
}
