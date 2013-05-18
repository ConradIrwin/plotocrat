package db

import (
    "github.com/ConradIrwin/plotocrat/models"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"sync"
)

var connection = (*mgo.Session)(nil)
var indexOnce sync.Once

func database() *mgo.Database {
	return connection.DB("")
}

func plotCollection() (collection *mgo.Collection) {
	collection = database().C("plots");

	indexOnce.Do(func () {
		err := collection.EnsureIndexKey("Uid")
		if err != nil {
			panic(err);
		}
	});

	return;
}

func Setup(url string) {
	session, err := mgo.Dial(url)

	if err != nil {
		panic(err)
	}
	connection = session
}

func SavePlot(plot *models.Plot) error {
	_, err := plotCollection().Upsert(bson.M{"uid": plot.Uid}, plot);
	return err;
}

func LoadPlot(uid string) (*models.Plot, error) {
	plot := new(models.Plot)
	err := plotCollection().Find(bson.M{"uid": uid}).One(plot);

	if err != nil {
		return nil, err
	}

	return plot, nil
}
