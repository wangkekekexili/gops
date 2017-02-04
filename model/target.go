package gops

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"

	"github.com/uber-go/zap"
	"github.com/wangkekekexili/gops/util"
)

const (
	targetURL = "http://redsky.target.com/v1/plp/search?"
)

var (
	targetPS4Regex = regexp.MustCompile(`(.+)( \(PlayStation 4\)| - PlayStation 4)$`)
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
			name, ok := t.extractName(title)
			if !ok {
				util.LogInfo("uncognizable name",
					zap.String("source", ProductSourceTarget),
					zap.String("name", title),
				)
				continue
			}
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

func (t *TargetHandler) extractName(s string) (string, bool) {
	matches := targetPS4Regex.FindStringSubmatch(s)
	if len(matches) == 3 {
		return matches[1], true
	}
	return "", false
}
