package gops

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/uber-go/zap"
	"github.com/wangkekekexili/gops/util"
)

const (
	targetURL       = "http://redsky.target.com/v1/plp/search?"
	targetPS4Suffix = " (PlayStation 4)"
)

type TargetHandler struct{}

var _ GameHandler = &TargetHandler{}

func (t *TargetHandler) GetGames() ([]GamePrice, error) {
	var games []GamePrice

	params := &url.Values{}
	params.Set("keyword", "playstation 4")
	params.Set("count", "24")
	params.Set("category", "55krz")
	offset := 0

	for {
		params.Set("offset", strconv.Itoa(offset))
		httpResponse, err := http.Get(targetURL + params.Encode())
		if err != nil {
			return nil, err
		}
		bytes, err := ioutil.ReadAll(httpResponse.Body)
		if err != nil {
			return nil, err
		}
		httpResponse.Body.Close()

		searchResponse := make(map[string]interface{})
		if err := json.Unmarshal(bytes, &searchResponse); err != nil {
			return nil, err
		}
		gameItems := searchResponse["search_response"].(map[string]interface{})["items"].(map[string]interface{})["Item"].([]interface{})
		if len(gameItems) == 0 {
			break
		}
		offset += len(gameItems)

		for _, gameInfoInterface := range gameItems {
			gameInfo := gameInfoInterface.(map[string]interface{})
			title := gameInfo["title"].(string)
			if !strings.HasSuffix(title, targetPS4Suffix) {
				util.LogInfo("uncognizable name",
					zap.String("source", ProductSourceTarget),
					zap.String("name", title),
				)
				continue
			}
			name := strings.TrimSuffix(title, targetPS4Suffix)
			price := gameInfo["offer_price"].(map[string]interface{})["price"].(float64)
			game := NewGameBuilder().FromTarget().IsNew().SetName(name).Build()
			games = append(games, GamePrice{Game: game, Price: NewPrice(-1, price)})
		}

	}

	return games, nil
}

func (t *TargetHandler) GetSource() string {
	return ProductSourceTarget
}

//func (t *TargetHandler) extractNameAndCondition(s string) (string, string, bool) {
//	matches := walmartPreownedProductRegex.FindStringSubmatch(s)
//	if len(matches) == 3 {
//		return matches[1], ProductConditionPreowned, true
//	}
//	matches = walmartNewProductRegex.FindStringSubmatch(s)
//	if len(matches) == 2 {
//		return matches[1], ProductConditionNew, true
//	}
//	return "", "", false
//}

//
//type TargetGame struct {
//	BasicGameInfo `bson:",inline"`
//	DPCI          string `bson:"dpci"`
//	TCIN          string `bson:"tcin"`
//	UPC           string `bson:"upc"`
//}
//
//func NewTargetGame(brand string, name string, price float64, t time.Time,
//	dpci, tcin, upc string) *TargetGame {
//	return &TargetGame{
//		BasicGameInfo: BasicGameInfo{
//			Brand:        brand,
//			Condition:    ProductConditionNew,
//			Name:         name,
//			PriceHistory: []PricePoint{{Price: price, Timestamp: t}},
//			Source:       ProductSourceTarget,
//		},
//		DPCI: dpci,
//		TCIN: tcin,
//		UPC:  upc,
//	}
//}
//
//type Target struct{}
//
//var _ GameFetcher = &Target{}
//
//func (target *Target) GetGamesInfo(c *mgo.Collection) ([]string, map[string]BasicGameInfoI, error) {
//	t := time.Now()
//
//	httpResponse, err := http.Get(targetURL)
//	if err != nil {
//		return nil, nil, err
//	}
//	bytes, err := ioutil.ReadAll(httpResponse.Body)
//	if err != nil {
//		return nil, nil, err
//	}
//	httpResponse.Body.Close()
//
//	searchResponse := make(map[string]interface{})
//	if err := json.Unmarshal(bytes, &searchResponse); err != nil {
//		return nil, nil, err
//	}
//	gameItems := searchResponse["search_response"].(map[string]interface{})["items"].(map[string]interface{})["Item"].([]interface{})
//
//	gamesByNameAndCondition := make(map[string]BasicGameInfoI)
//	var gameNames []string
//	for _, gameInfoInterface := range gameItems {
//		gameInfo := gameInfoInterface.(map[string]interface{})
//		title := gameInfo["title"].(string)
//		name := strings.TrimSuffix(title, targetGameSuffix)
//		gameNames = append(gameNames, name)
//		brand := gameInfo["brand"].(string)
//		price := gameInfo["offer_price"].(map[string]interface{})["price"].(float64)
//		dpci := gameInfo["dpci"].(string)
//		tcin := gameInfo["tcin"].(string)
//		upc := gameInfo["upc"].(string)
//		gamesByNameAndCondition[name+ProductConditionNew] = NewTargetGame(brand, name, price, t, dpci, tcin, upc)
//	}
//
//	return gameNames, gamesByNameAndCondition, nil
//}
//
//func (target *Target) GetSource() string {
//	return ProductSourceTarget
//}
