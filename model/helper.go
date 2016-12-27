package gops

import (
	"time"

	"encoding/json"
	"sync"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type GameFetcher interface {
	GetGamesInfo(*mgo.Collection) ([]string, map[string]BasicGameInfoI, error)
	GetSource() string
}

var allGameFetchers []GameFetcher
var allGameFetchersInit sync.Once

func GetAllGameFetchers() []GameFetcher {
	allGameFetchersInit.Do(func() {
		allGameFetchers = []GameFetcher{
			&Gamefly{},
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

var _ BasicGameInfoI = &BasicGameInfo{}

func (info *BasicGameInfo) String() string {
	bytes, _ := json.Marshal(&info)
	return string(bytes)
}

func (info *BasicGameInfo) GetID() bson.ObjectId {
	return info.ID
}

func (info *BasicGameInfo) GetCondition() string {
	return info.Condition
}

func (info *BasicGameInfo) GetName() string {
	return info.Name
}

func (info *BasicGameInfo) GetPriceHistory() []PricePoint {
	return info.PriceHistory
}

func (info *BasicGameInfo) GetRecentPrice() float64 {
	historyLen := len(info.PriceHistory)
	if historyLen == 0 {
		return 0
	}
	return info.PriceHistory[historyLen-1].Price
}

type BasicGameInfoI interface {
	GetID() bson.ObjectId
	GetCondition() string
	GetName() string
	GetPriceHistory() []PricePoint
	GetRecentPrice() float64
}
