package gops

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"

	"github.com/wangkekekexili/gops/util"
	"go.uber.org/zap"
)

const (
	walmartSearchAPI     = "http://api.walmartlabs.com/v1/search?"
	walmartPS4CategoryID = "2636_1102672_1105671"
)

var (
	// Name like "Need for Speed: Rivals (PS4)" indicates it's a new game.
	walmartNewProductRegex = regexp.MustCompile(`(.+) \(PS4\)$`)
	// Name like "Metal Gear Solid V: The Phantom Pain (PS4) - Pre-Owned" indicates it's a pre-owned game.
	walmartPreownedProductRegex = regexp.MustCompile(`(.+) (- Pre-Owned \(PS4\)|\(PS4\) - Pre-Owned)$`)
)

type WalmartHandler struct{}

var _ GameHandler = &WalmartHandler{}

func (w *WalmartHandler) GetGames() ([]GamePrice, error) {
	var games []GamePrice

	params := &url.Values{}
	start := 1
	params.Set("apiKey", os.Getenv("WALMART_KEY"))
	params.Set("categoryId", walmartPS4CategoryID)
	params.Set("query", "ps4")
	params.Set("format", "json")
	params.Set("numItems", "25")
	for {
		// Make request and get JSON data.
		params.Set("start", strconv.Itoa(start))
		response, err := http.Get(walmartSearchAPI + params.Encode())
		if err != nil {
			return nil, err
		}
		bytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}
		response.Body.Close()

		// Extract from JSON.
		data := make(map[string]interface{})
		if err = json.Unmarshal(bytes, &data); err != nil {
			return nil, err
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
				util.LogInfo("unrecognizable name",
					zap.String("source", ProductSourceWalmart),
					zap.String("fullname", fullname),
				)
				continue
			}
			price := item["salePrice"].(float64)
			game := NewGameBuilder().FromWalmart().SetName(name).SetCondition(condition).Build()
			games = append(games, GamePrice{Game: game, Price: NewPrice(-1, price)})
		}
	}
	return games, nil
}

func (w *WalmartHandler) GetSource() string {
	return ProductSourceWalmart
}

func (w *WalmartHandler) extractNameAndCondition(s string) (string, string, bool) {
	matches := walmartPreownedProductRegex.FindStringSubmatch(s)
	if len(matches) == 3 {
		return matches[1], ProductConditionPreowned, true
	}
	matches = walmartNewProductRegex.FindStringSubmatch(s)
	if len(matches) == 2 {
		return matches[1], ProductConditionNew, true
	}
	return "", "", false
}
