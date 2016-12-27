package gops

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"gopkg.in/mgo.v2"
)

const (
	gameflyURL = "https://api.gamefly.com/api/productquery/findpage?bargainbin=true&pageindex=1&pagesize=24&platform=1225&productType=ConsoleGame&rentOnly=false&sort.direction=asc&sort.field=bestprice"
)

type GameflyGame struct {
	BasicGameInfo `bson:",inline"`
}

func NewGameflyGame(brand string, name string, price float64, t time.Time) *GameflyGame {
	return &GameflyGame{
		BasicGameInfo: BasicGameInfo{
			Brand:        brand,
			Condition:    ProductConditionPreowned,
			Name:         name,
			PriceHistory: []PricePoint{{Price: price, Timestamp: t}},
			Source:       ProductSourceGamefly,
		},
	}
}

type Gamefly struct{}

var _ GameFetcher = &Gamefly{}

func (gamefly *Gamefly) GetGamesInfo(c *mgo.Collection) ([]string, map[string]BasicGameInfoI, error) {
	t := time.Now()

	httpResponse, err := http.Get(gameflyURL)
	if err != nil {
		return nil, nil, err
	}
	bytes, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, nil, err
	}
	httpResponse.Body.Close()

	searchResponse := make(map[string]interface{})
	if err := json.Unmarshal(bytes, &searchResponse); err != nil {
		return nil, nil, err
	}
	gameItems := searchResponse["products"].(map[string]interface{})["items"].([]interface{})

	gamesByNameAndCondition := make(map[string]BasicGameInfoI)
	var gameNames []string
	for _, gameInfoInterface := range gameItems {
		gameInfo := gameInfoInterface.(map[string]interface{})
		name := gameInfo["title"].(string)
		gameNames = append(gameNames, name)
		brand := gameInfo["publisher"].(string)
		offerActions := gameInfo["offerActions"].([]interface{})
		price := offerActions[len(offerActions)-1].(map[string]interface{})["purchasePrice"].(map[string]interface{})["amount"].(float64)
		gamesByNameAndCondition[name+ProductConditionPreowned] = NewGameflyGame(brand, name, price, t)
	}

	return gameNames, gamesByNameAndCondition, nil
}

func (gamefly *Gamefly) GetSource() string {
	return ProductSourceGamefly
}
