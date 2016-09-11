package main

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	// Get a list of price drop products of ps4 on gamestop.
	document, _ := goquery.NewDocument("http://www.gamestop.com/browse/playstation-4?nav=28-xu0,131dc-162")
	document.Find("div .product").Each(func(i int, s *goquery.Selection) {
		name := s.Find(".ats-product-title-lnk").Text()
		condition := s.Find(".ats-product-condition > strong").Text()
		oldPrice := s.Find(".old_price").Text()
		mixedPrices := s.Find(".ats-product-price").Text()
		newPrice := strings.Replace(mixedPrices, oldPrice, "", -1)
		fmt.Printf("%s\n%s\nOld Price: %s\nNow: %s\n\n", name, condition, oldPrice, newPrice)
	})
}
