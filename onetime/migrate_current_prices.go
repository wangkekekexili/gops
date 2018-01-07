package main

import (
	"log"

	"github.com/wangkekekexili/gops/db"
	"github.com/wangkekekexili/gops/model"
	"github.com/wangkekekexili/gops/model/tables"
	"github.com/wangkekekexili/module"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	d := &db.Module{}
	err := module.Load(d)
	if err != nil {
		return err
	}
	tx, err := d.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Get all games.
	var games []*model.Game
	err = tx.Select(&games, "SELECT * FROM "+tables.Games)
	if err != nil {
		return err
	}

	// For each game. Get all prices and put the latest one.
	for _, game := range games {
		var prices []*model.Price
		err = tx.Select(&prices, "SELECT * FROM "+tables.Prices+" WHERE game_id = ?", game.ID)
		if err != nil {
			return err
		}

		mostRecentPrice := prices[0]
		for i := 1; i < len(prices); i++ {
			price := prices[i]
			if price.Timestamp.After(mostRecentPrice.Timestamp) {
				mostRecentPrice = price
			}
		}

		_, err = tx.Exec("INSERT INTO "+tables.CurrentPrices+" (`game_id`, `price_id`, `value`)  VALUES (?,?,?)",
			game.ID, mostRecentPrice.ID, mostRecentPrice.Value)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
