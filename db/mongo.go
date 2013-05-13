package db

import (
	"labix.org/v2/mgo"
)

var connection = (*mgo.Session)(nil)

func database() *mgo.Database {
    return connection.DB("");
}

func Setup(url string) {
    session, err := mgo.Dial(url)

    if err != nil {
        panic(err)
    }
    connection = session
}
