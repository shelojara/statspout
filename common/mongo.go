package common

import (
	"flag"

	"gopkg.in/mgo.v2"

	"github.com/mijara/statspout/repo"
	"github.com/mijara/statspout/stats"
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

func NewMongo(opts *MongoOpts) (*Mongo, error) {
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

func (*Mongo) Create(v interface{}) (repo.Interface, error) {
	return NewMongo(v.(*MongoOpts))
}

func (mongo *Mongo) Push(s *stats.Stats) error {
	c := mongo.session.DB(mongo.database).C(mongo.collection)

	err := c.Insert(s)
	if err != nil {
		return err
	}

	return nil
}

func (*Mongo) Name() string {
	return "mongodb"
}

func (mongo *Mongo) Close() {
	mongo.session.Close()
}

func (mongo *Mongo) Clear(name string) {
	// not used.
}

func CreateMongoOpts() *MongoOpts {
	o := &MongoOpts{}

	flag.StringVar(&o.Address,
		"mongo.address",
		"localhost:27017",
		"Address of the MongoDB Endpoint")

	flag.StringVar(&o.Database,
		"mongo.database",
		"statspout",
		"Database for the collection")

	flag.StringVar(&o.Collection,
		"mongo.collection",
		"stats",
		"Collection for the stats")

	return o
}
