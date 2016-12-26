package gops

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"gopkg.in/mgo.v2"
)

const (
	gamestopURLTemplate = "http://www.gamestop.com/browse/games/playstation-4?nav=2b%d,28-xu0,131dc-ffff2418-d3-162"
)

type GamestopGame struct {
	BasicGameInfo `bson:",inline"`
}

func NewGamestopGame(brand string, condition string, name string, price float64, t time.Time) *GamestopGame {
	return &GamestopGame{BasicGameInfo: BasicGameInfo{
		Brand:        brand,
		Condition:    condition,
		Name:         name,
		PriceHistory: []PricePoint{{Price: price, Timestamp: t}},
		Source:       ProductSourceGamestop,
	}}
}

type Gamestop struct{}

var _ GameFetcher = &Gamestop{}

func (gamestop *Gamestop) GetGamesInfo(c *mgo.Collection) ([]string, map[string]BasicGameInfoI, error) {
	t := time.Now()
	gamesByNameAndCondition := make(map[string]BasicGameInfoI)
	var gameNames []string
	startIndex := 0
	for {
		numProduct := 0
		url := fmt.Sprintf(gamestopURLTemplate, startIndex)
		document, err := goquery.NewDocument(url)
		if err != nil {
			return nil, nil, err
		}
		document.Find("div .product").Each(func(i int, s *goquery.Selection) {
			numProduct++

			var name, brand, condition string
			var price float64

			name = s.Find(".ats-product-title-lnk").Text()
			brand = strings.TrimPrefix(s.Find(".ats-product-publisher").Text(), "by ")
			conditionText := s.Find(".ats-product-condition > strong").Text()
			switch conditionText {
			case "NEW":
				condition = ProductConditionNew
			case "PRE-OWNED":
				condition = ProductConditionPreowned
			default:
				return
			}
			oldPrice := s.Find(".old_price").Text()
			mixedPrices := s.Find(".ats-product-price").Text()
			newPrice := strings.Replace(mixedPrices, oldPrice, "", -1)[1:]
			price, err = strconv.ParseFloat(newPrice, 64)
			if err != nil {
				fmt.Println(err)
				return
			}
			gameNames = append(gameNames, name)
			gamesByNameAndCondition[name+condition] = NewGamestopGame(brand, condition, name, price, t)
		})
		if numProduct == 0 {
			break
		}
		startIndex += numProduct
		time.Sleep(time.Second)
	}

	return gameNames, gamesByNameAndCondition, nil
}

func (gamestop *Gamestop) GetSource() string {
	return ProductSourceGamestop
}
