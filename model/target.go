package gops

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	targetURL = "http://redsky.target.com/v1/plp/search?count=100&offset=0&sort_by=bestselling&category=55krz&faceted_value=5tdv1Z5tdv0"

	targetGameSuffix = " (PlayStation 4)"
)

type TargetGame struct {
	BasicGameInfo `bson:",inline"`
	DPCI          string `bson:"dpci"`
	TCIN          string `bson:"tcin"`
	UPC           string `bson:"upc"`
}

func NewTargetGame(brand string, condition string, name string, price float64, t time.Time,
	dpci, tcin, upc string) *TargetGame {
	return &TargetGame{
		BasicGameInfo: BasicGameInfo{
			Brand:        brand,
			Condition:    condition,
			Name:         name,
			PriceHistory: []PricePoint{{Price: price, Timestamp: t}},
			Source:       ProductSourceTarget,
		},
		DPCI: dpci,
		TCIN: tcin,
		UPC:  upc,
	}
}

type Target struct{}

var _ GameFetcher = &Target{}

type targetSearchResponse struct {
}

func (target *Target) GetGamesInfo(c *mgo.Collection) ([]interface{}, map[bson.ObjectId]interface{}, error) {
	t := time.Now()

	httpResponse, err := http.Get(targetURL)
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
	gameItems := searchResponse["search_response"].(map[string]interface{})["items"].(map[string]interface{})["Item"].([]interface{})

	gamesByName := make(map[string]*TargetGame)
	var gameNames []string
	for _, gameInfoInterface := range gameItems {
		gameInfo := gameInfoInterface.(map[string]interface{})
		title := gameInfo["title"].(string)
		name := strings.TrimSuffix(title, targetGameSuffix)
		gameNames = append(gameNames, name)
		brand := gameInfo["brand"].(string)
		price := gameInfo["offer_price"].(map[string]interface{})["price"].(float64)
		dpci := gameInfo["dpci"].(string)
		tcin := gameInfo["tcin"].(string)
		upc := gameInfo["upc"].(string)
		gamesByName[name] = NewTargetGame(brand, ProductConditionNew, name, price, t, dpci, tcin, upc)
	}
	// Get existing documents and see if we need to update them.
	namesSubQuery := make([]bson.M, len(gameNames))
	for i, gameName := range gameNames {
		namesSubQuery[i] = bson.M{"name": gameName}
	}
	cursor := c.Find(bson.M{"source": ProductSourceTarget, "$or": namesSubQuery}).Iter()
	var result TargetGame
	for cursor.Next(&result) {
		name := result.Name
		if gameNewPricePoint, ok := gamesByName[name]; ok {
			// Skip if the price is not changed.
			if result.PriceHistory[len(result.PriceHistory)-1].Price == gameNewPricePoint.PriceHistory[0].Price {
				delete(gamesByName, name)
				continue
			}
			// Make a copy of result.
			gameToUpdate := result
			gameToUpdate.PriceHistory = append(gameToUpdate.PriceHistory, gameNewPricePoint.PriceHistory[0])
			gamesByName[name] = &gameToUpdate
		}
	}
	cursor.Close()

	var gamesToInsert []interface{}
	gamesToUpdate := make(map[bson.ObjectId]interface{})

	for _, game := range gamesByName {
		if game.ID.Hex() == "" {
			gamesToInsert = append(gamesToInsert, game)
		} else {
			gamesToUpdate[game.ID] = game
		}
	}

	return gamesToInsert, gamesToUpdate, nil
}

func (target *Target) GetSource() string {
	return ProductSourceTarget
}
