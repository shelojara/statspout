package common

import (
	"flag"

	"github.com/mijara/statspout/data"
	"gopkg.in/mgo.v2"
)

type Mongo struct {
	session    *mgo.Session
	database   string
	collection string
}

type MongoOpts struct {
	Address    string
	Database   string
	Collection string
}

func NewMongo(opts MongoOpts) (*Mongo, error) {
	session, err := mgo.Dial(opts.Address)
	if err != nil {
		return nil, err
	}

	return &Mongo{
		session:    session,
		database:   opts.Database,
		collection: opts.Collection,
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

func CreateMongoOpts(opts *MongoOpts) {
	flag.StringVar(&opts.Address,
		"mongo.address",
		"localhost:27017",
		"Address of the MongoDB Endpoint")

	flag.StringVar(&opts.Database,
		"mongo.database",
		"statspout",
		"Database for the collection")

	flag.StringVar(& opts.Collection,
			"mongo.collection",
			"stats",
			"Collection for the stats")
}
