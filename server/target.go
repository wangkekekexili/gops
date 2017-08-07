package server

import (
	"html"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"

	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
	"github.com/wangkekekexili/gops/model"
	"go.uber.org/zap"
)

const (
	targetURL = "http://redsky.target.com/v1/plp/search?"
)

var (
	targetPS4Regex = regexp.MustCompile(`(.+)( \(PlayStation 4\)| - PlayStation 4)$`)
)

type TargetHandler struct {
	Logger *Logger
}

func (t *TargetHandler) Load() error { return nil }

var _ GameHandler = &TargetHandler{}

func (t *TargetHandler) GetGames() ([]*model.GamePrice, error) {
	logger := t.Logger.With(zap.String("source", model.ProductSourceTarget))

	var games []*model.GamePrice

	params := &url.Values{}
	params.Set("keyword", "playstation 4")
	params.Set("count", "24")
	params.Set("category", "55krz")
	offset := 0

	for {
		params.Set("offset", strconv.Itoa(offset))
		httpResponse, err := http.Get(targetURL + params.Encode())
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get from %v", targetURL+params.Encode())
		}
		responseBytes, err := ioutil.ReadAll(httpResponse.Body)
		if err != nil {
			return nil, errors.Wrap(err, "failed to read body")
		}
		httpResponse.Body.Close()

		// Get items from the json response.
		gameItemsResult := gjson.Get(string(responseBytes), "search_response.items.Item")
		if !gameItemsResult.Exists() {
			logger.Info("no games in the json", zap.String("json", string(responseBytes)))
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
				logger.Info("unrecognizable title", zap.String("json", gameInfo.Raw))
				continue
			}
			name, ok := t.extractName(title.String())
			if !ok {
				logger.Info("unrecognizable name", zap.String("name", title.String()))
				continue
			}
			price := gameInfo.Get("offer_price.price").Float()
			game := model.NewGamePriceBuilder().FromTarget().IsNew().SetName(name).SetPrice(price).Build()
			games = append(games, game)
		}

	}

	return games, nil
}

func (t *TargetHandler) GetSource() string {
	return model.ProductSourceTarget
}

func (t *TargetHandler) extractName(s string) (string, bool) {
	matches := targetPS4Regex.FindStringSubmatch(s)
	if len(matches) == 3 {
		return html.UnescapeString(matches[1]), true
	}
	return "", false
}
