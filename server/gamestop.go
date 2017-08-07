package server

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
	"github.com/wangkekekexili/gops/model"
	"go.uber.org/zap"
)

const (
	gamestopURL = "http://www.gamestop.com/browse/games/playstation-4?"

	gamestopPreownedClass = "preowned_product"
	gamestopNewClass      = "new_product"
)

type GamestopHandler struct {
	Logger *Logger
}

func (g *GamestopHandler) Load() error { return nil }

var _ GameHandler = &GamestopHandler{}

func (g *GamestopHandler) GetGames() ([]*model.GamePrice, error) {
	logger := g.Logger.With(zap.String("source", model.ProductSourceGamestop))
	var games []*model.GamePrice
	offset := 0
	for {
		var document *goquery.Document
		var err error
		for retry := 0; retry != 2; retry++ {
			u := fmt.Sprintf("%vnav=2b%d,28-xu0,131dc-ffff2418", gamestopURL, offset)
			document, err = goquery.NewDocument(u)
			if err == nil {
				break
			}
			if _, ok := err.(*net.OpError); ok {
				time.Sleep(time.Second)
				continue
			}
			return nil, errors.Wrapf(err, "failed to get from %v", u)
		}
		products := document.Find(".product")
		if products.Size() == 0 {
			break
		}
		offset += products.Size()

		products.Each(func(index int, s *goquery.Selection) {
			var name, condition string
			var price float64
			if s.HasClass(gamestopPreownedClass) {
				condition = model.ProductConditionPreowned
			} else if s.HasClass(gamestopNewClass) {
				condition = model.ProductConditionNew
			} else if s.HasClass("digital_product") {
				// We currently don't handle digital products.
				return
			} else {
				h, _ := s.Html()
				logger.Info("the product doesn't have condition class", zap.String("content", h))
				return
			}

			name = s.Find(".ats-product-title-lnk").First().Text()
			if strings.HasPrefix(name, "PlayStation 4") {
				// It's probably a bundle product.
				return
			}
			priceNode := s.Find(".pricing").First()
			priceNode.Children().First().Remove()
			price, err = strconv.ParseFloat(priceNode.Text()[1:], 64)
			if err != nil {
				logger.Info("cannot parse price", zap.String("content", priceNode.Text()))
				return
			}
			game := model.NewGamePriceBuilder().SetName(name).SetCondition(condition).FromGamestop().SetPrice(price).Build()
			games = append(games, game)
		})
	}

	return games, nil
}

func (g *GamestopHandler) GetSource() string {
	return model.ProductSourceGamestop
}
