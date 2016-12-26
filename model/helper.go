package gops

import (
	"time"

	"encoding/json"
	"sync"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type GameFetcher interface {
	GetGamesInfo(*mgo.Collection) ([]interface{}, map[bson.ObjectId]interface{}, error)
	GetSource() string
}

var allGameFetchers []GameFetcher
var allGameFetchersInit sync.Once

func GetAllGameFetchers() []GameFetcher {
	allGameFetchersInit.Do(func() {
		allGameFetchers = []GameFetcher{
			&Gamestop{},
			&Target{},
		}
	})
	return allGameFetchers
}

type PricePoint struct {
	Price     float64   `bson:"price"`
	Timestamp time.Time `bson:"timestamp"`
}

type BasicGameInfo struct {
	ID           bson.ObjectId `bson:"_id,omitempty"`
	Brand        string        `bson:"brand"`
	Condition    string        `bson:"condition"`
	Name         string        `bson:"name"`
	PriceHistory []PricePoint  `bson:"prices"`
	Source       string        `bson:"source"`
}

func (info BasicGameInfo) String() string {
	bytes, _ := json.Marshal(&info)
	return string(bytes)
}
