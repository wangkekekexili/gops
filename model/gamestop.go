package gops

import (
	"fmt"
	"strconv"
	"strings"

	"net"

	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/wangkekekexili/gops/util"
	"go.uber.org/zap"
	"github.com/pkg/errors"
)

const (
	gamestopURL = "http://www.gamestop.com/browse/games/playstation-4?"

	gamestopPreownedClass = "preowned_product"
	gamestopNewClass      = "new_product"
)

type GamestopHandler struct{}

var _ GameHandler = &GamestopHandler{}

func (g *GamestopHandler) GetGames() ([]GamePrice, error) {
	var games []GamePrice
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
				condition = ProductConditionPreowned
			} else if s.HasClass(gamestopNewClass) {
				condition = ProductConditionNew
			} else if s.HasClass("digital_product") {
				// We currently don't handle digital products.
				return
			} else {
				h, _ := s.Html()
				util.LogInfo("the product doesn't have condition class",
					zap.String("source", ProductSourceGamestop),
					zap.String("content", h),
				)
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
				util.LogInfo("cannot parse price",
					zap.String("source", ProductSourceGamestop),
					zap.String("content", priceNode.Text()),
				)
				return
			}
			game := NewGameBuilder().SetName(name).SetCondition(condition).FromGamestop().Build()
			games = append(games, GamePrice{Game: game, Price: NewPrice(-1, price)})
		})
	}

	return games, nil
}

func (g *GamestopHandler) GetSource() string {
	return ProductSourceGamestop
}
