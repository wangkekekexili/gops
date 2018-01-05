package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"

	"github.com/pkg/errors"
	"github.com/wangkekekexili/gops/logger"
	"github.com/wangkekekexili/gops/model"
	"go.uber.org/zap"
)

const (
	walmartSearchAPI     = "http://api.walmartlabs.com/v1/search?"
	walmartPS4CategoryID = "2636_1102672_1105671"
)

var (
	// Name like "Need for Speed: Rivals (PS4)" indicates it's a new game.
	walmartNewProductRegex = regexp.MustCompile(`(.+)( \(PS4\)| \(PlayStation 4\))$`)
	// Name like "Metal Gear Solid V: The Phantom Pain (PS4) - Pre-Owned" indicates it's a pre-owned game.
	walmartPreownedProductRegex = regexp.MustCompile(`(.+) (- Pre-Owned \(PS4\)|\(PS4\) - Pre-Owned)$`)
)

type WalmartHandler struct {
	Logger *logger.Module

	ok     bool
	params url.Values
}

func (w *WalmartHandler) Load() error {
	k := os.Getenv("WALMART_KEY")
	if k == "" {
		return nil
	}

	w.ok = true
	w.params = url.Values{}
	w.params.Set("apiKey", k)
	w.params.Set("categoryId", walmartPS4CategoryID)
	w.params.Set("query", "ps4")
	w.params.Set("format", "json")
	w.params.Set("numItems", "25")
	return nil
}

var _ GameHandler = &WalmartHandler{}

func (w *WalmartHandler) GetGames() ([]*model.GamePrice, error) {
	if !w.ok {
		return nil, nil
	}

	var games []*model.GamePrice
	start := 1
	for {
		// Make request and get JSON data.
		params := w.params
		params.Set("start", strconv.Itoa(start))
		response, err := http.Get(walmartSearchAPI + params.Encode())
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get from %v", walmartSearchAPI+params.Encode())
		}
		bytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, errors.Wrap(err, "failed to read body")
		}
		response.Body.Close()

		// Extract from JSON.
		data := make(map[string]interface{})
		if err = json.Unmarshal(bytes, &data); err != nil {
			return nil, errors.Wrapf(err, "failed to unmarshal from %v", string(bytes))
		}
		numItems := int(data["numItems"].(float64))
		if numItems == 0 {
			break
		}
		start += numItems

		items := data["items"].([]interface{})
		for _, itemInterface := range items {
			item := itemInterface.(map[string]interface{})
			// Don't record martketplace products.
			isMarketplace, ok := item["marketplace"].(bool)
			if isMarketplace || !ok {
				continue
			}

			fullname := item["name"].(string)
			name, condition, ok := w.extractNameAndCondition(fullname)
			if !ok {
				w.Logger.Info("unrecognizable name", zap.String("source", model.ProductSourceWalmart), zap.String("fullname", fullname))
				continue
			}
			price := item["salePrice"].(float64)
			game := model.NewGamePriceBuilder().FromWalmart().SetName(name).SetCondition(condition).SetPrice(price).Build()
			games = append(games, game)
		}
	}
	return games, nil
}

func (w *WalmartHandler) GetSource() string {
	return model.ProductSourceWalmart
}

func (w *WalmartHandler) extractNameAndCondition(s string) (string, string, bool) {
	matches := walmartPreownedProductRegex.FindStringSubmatch(s)
	if len(matches) == 3 {
		return matches[1], model.ProductConditionPreowned, true
	}
	matches = walmartNewProductRegex.FindStringSubmatch(s)
	if len(matches) == 3 {
		return matches[1], model.ProductConditionNew, true
	}
	return "", "", false
}
