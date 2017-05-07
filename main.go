package main

import (
	"database/sql"
	"fmt"
	"os"
	"runtime/debug"
	"strings"
	"sync"

	"github.com/getsentry/raven-go"
	"github.com/go-gorp/gorp"
	_ "github.com/go-sql-driver/mysql"
	"github.com/wangkekekexili/gops/model"
	"github.com/wangkekekexili/gops/util"
	"go.uber.org/zap"
)

var (
	dbMap *gorp.DbMap
)

func init() {
	// Init error reporter.
	sentryDSN := strings.TrimSpace(os.Getenv("SENTRY_DSN"))
	if sentryDSN == "" {
		util.LogError("error reporter cannot be initialized: SENTRY_DSN is unset")
		os.Exit(1)
	}
	raven.SetDSN(sentryDSN)

	// Check db connection.
	var err error
	var db *sql.DB
	db, err = sql.Open("mysql", os.Getenv("MYSQL_DSN"))
	if err != nil {
		util.LogError(err.Error())
		raven.CaptureErrorAndWait(err, nil)
		os.Exit(1)
	}
	if err = db.Ping(); err != nil {
		util.LogError(err.Error())
		util.SendWarningSMS()
		raven.CaptureErrorAndWait(err, nil)
		os.Exit(1)
	}

	// Initialize dbMap.
	dbMap = &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{Engine: "innodb", Encoding: "ascii"}}
	dbMap.AddTableWithName(gops.Game{}, "game")
	dbMap.AddTableWithName(gops.Price{}, "price")
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			err := fmt.Errorf("panic captured: %v\n%s", r, string(debug.Stack()))
			util.LogError(err.Error())
			raven.CaptureError(err, nil)
		}
		if err := dbMap.Db.Close(); err != nil {
			util.LogError(err.Error())
			raven.CaptureError(err, nil, nil)
		}
		raven.Wait()
	}()

	var wg sync.WaitGroup
	for _, handler := range gops.AllGameHandlers {
		wg.Add(1)
		go func(handler gops.GameHandler) {
			defer func() {
				if r := recover(); r != nil {
					err := fmt.Errorf("panic captured: %v\n%s", r, string(debug.Stack()))
					util.LogError(err.Error())
					raven.CaptureError(err, nil)
				}
				wg.Done()
			}()
			if err := handleGames(handler); err != nil {
				util.LogError(err.Error(), zap.String("source", handler.GetSource()))
				raven.CaptureError(err, nil)
			}
		}(handler)
	}
	wg.Wait()
}

func handleGames(handler gops.GameHandler) error {
	util.LogInfo("start processing", zap.String("source", handler.GetSource()))

	gamesWithPrice, err := handler.GetGames()
	if err != nil {
		return err
	}
	util.LogInfo("successfully get games",
		zap.String("source", handler.GetSource()),
		zap.Int("game_count", len(gamesWithPrice)),
	)
	if len(gamesWithPrice) == 0 {
		return nil
	}
	var gameNames []interface{}
	for _, entry := range gamesWithPrice {
		gameNames = append(gameNames, entry.Game.Name)
	}

	// Get existing games.
	var existingGames []gops.Game
	if _, err := dbMap.Select(&existingGames, `SELECT * FROM game WHERE name IN `+util.QuestionMarks(len(gameNames)), gameNames...); err != nil {
		return err
	}
	existingGamesByKey := make(map[string]*gops.Game)
	for i, game := range existingGames {
		existingGamesByKey[game.GetKey()] = &existingGames[i]
	}

	var numNewGames, numPriceUpdate int
	for _, gameWithPrice := range gamesWithPrice {
		game := gameWithPrice.Game
		price := gameWithPrice.Price
		// Check if the game has an existing entry.
		if existingGame, hasExistingEntry := existingGamesByKey[game.GetKey()]; hasExistingEntry {
			// Get the last price for the game.
			lastPrice, err := dbMap.SelectFloat(`SELECT value FROM price WHERE game_id = ? ORDER BY timestamp desc`, existingGame.ID)
			if err != nil {
				return err
			}
			if price.Value == lastPrice {
				continue
			}

			// New price point.
			numPriceUpdate++
			price.GameID = existingGame.ID
			if err = dbMap.Insert(price); err != nil {
				return err
			}
		} else {
			// It's new game.
			numNewGames++
			if err = dbMap.Insert(game); err != nil {
				return err
			}
			price.GameID = game.ID
			if err = dbMap.Insert(price); err != nil {
				return err
			}
		}
	}
	util.LogInfo("price updated",
		zap.String("source", handler.GetSource()),
		zap.Int("number", numPriceUpdate),
	)
	util.LogInfo("new games inserted",
		zap.String("source", handler.GetSource()),
		zap.Int("number", numNewGames),
	)

	util.LogInfo("end", zap.String("source", handler.GetSource()))
	return nil
}
