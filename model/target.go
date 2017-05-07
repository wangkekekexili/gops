package gops

import (
	"html"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"

	"github.com/tidwall/gjson"
	"github.com/wangkekekexili/gops/util"
	"go.uber.org/zap"
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
		responseBytes, err := ioutil.ReadAll(httpResponse.Body)
		if err != nil {
			return nil, err
		}
		httpResponse.Body.Close()

		// Get items from the json response.
		gameItemsResult := gjson.Get(string(responseBytes), "search_response.items.Item")
		if !gameItemsResult.Exists() {
			util.LogInfo("no games in the json",
				zap.String("source", ProductSourceTarget),
				zap.String("json", string(responseBytes)),
			)
			break
		}
		gameItems := gameItemsResult.Array()
		if len(gameItems) == 0 {
			break
		}
		offset += len(gameItems)

		for _, gameInfo := range gameItems {
			title := gameInfo.Get("title")
			if !title.Exists() {
				util.LogInfo("uncognizable name",
					zap.String("source", ProductSourceTarget),
					zap.String("json", gameInfo.Raw),
				)
				continue
			}
			name, ok := t.extractName(title.String())
			if !ok {
				util.LogInfo("uncognizable name",
					zap.String("source", ProductSourceTarget),
					zap.String("name", title.String()),
				)
				continue
			}
			price := gameInfo.Get("offer_price.price").Float()
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
		return html.UnescapeString(matches[1]), true
	}
	return "", false
}
