package server

import (
	"fmt"
	"runtime/debug"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"github.com/wangkekekexili/gops/model"
	"github.com/wangkekekexili/gops/reporter"
	"github.com/wangkekekexili/gops/util"
	"go.uber.org/zap"
)

type GOPS struct {
	DB       *DB
	Logger   *Logger
	Reporter *reporter.Module

	Gamestop *GamestopHandler
	Target   *TargetHandler
	Walmart  *WalmartHandler

	handlers []GameHandler
}

func (s *GOPS) Load() error {
	if s.DB == nil {
		return errors.New("cannot load server")
	}
	s.handlers = []GameHandler{
		s.Gamestop,
		s.Target,
		s.Walmart,
	}
	return nil
}

func (s *GOPS) Start() error {
	defer func() {
		err := s.DB.Close()
		if err != nil {
			s.Logger.Err(err.Error())
			s.Reporter.ErrSync(err)
		}
	}()

	var wg sync.WaitGroup
	for _, handler := range s.handlers {
		wg.Add(1)
		go func(handler GameHandler) {
			defer func() {
				if r := recover(); r != nil {
					err := fmt.Errorf("panic captured: %v\n%s", r, string(debug.Stack()))
					s.Logger.Info(err.Error())
					s.Reporter.Err(err)
				}
				wg.Done()
			}()
			if err := s.handleGames(handler); err != nil {
				s.Logger.Err(err.Error(), zap.String("source", handler.GetSource()))
				s.Reporter.Err(err)
			}
		}(handler)
	}
	wg.Wait()
	return nil
}

func (s *GOPS) handleGames(handler GameHandler) error {
	logger := s.Logger.With(zap.String("source", handler.GetSource()))
	logger.Info("start processing")
	tx, err := s.DB.Beginx()
	if err != nil {
		return errors.Wrap(err, "begin tx")
	}
	defer tx.Rollback()

	gamesWithPrice, err := handler.GetGames()
	if err != nil {
		return errors.WithMessage(err, "get games failed")
	}
	logger.Info("successfully get games", zap.Int("count", len(gamesWithPrice)))
	if len(gamesWithPrice) == 0 {
		return nil
	}
	var gameNames []interface{}
	for _, entry := range gamesWithPrice {
		gameNames = append(gameNames, entry.Name)
	}

	// Get existing games.
	var existingGames []*model.Game
	err = tx.Select(&existingGames, `SELECT * FROM game WHERE name IN `+util.QuestionMarks(len(gameNames)), gameNames...)
	if err != nil {
		return errors.Wrapf(err, "get existing games failed with game names: %v", gameNames)
	}

	existingGamesByKey := make(map[string]*model.Game)
	for i, game := range existingGames {
		existingGamesByKey[game.GetKey()] = existingGames[i]
	}

	var numNewGames, numPriceUpdate int
	for _, gameWithPrice := range gamesWithPrice {
		game, price := gameWithPrice.ToGameAndPrice()
		// Check if the game has an existing entry.
		if existingGame, hasExistingEntry := existingGamesByKey[game.GetKey()]; hasExistingEntry {
			// Get the last price for the game.
			var lastPrice float64
			err = tx.Get(&lastPrice, "SELECT value FROM price WHERE game_id = ? ORDER BY timestamp desc", existingGame.ID)
			if err != nil {
				return errors.Wrapf(err, "failed to get last price for game %v", game)
			}
			if price.Value == lastPrice {
				continue
			}

			// New price point.
			numPriceUpdate++
			price.GameID = existingGame.ID
			_, err = tx.Exec("INSERT INTO price (`game_id`, `value`) VALUES (?,?)", existingGame.ID, price.Value)
			if err != nil {
				return errors.Wrapf(err, "error inserting price with game_id %v value %v", existingGame.ID, price.Value)
			}
		} else {
			// It's new game.
			numNewGames++
			result, err := tx.Exec("INSERT INTO game (`name`, `condition`, `source`) VALUES (?,?,?)",
				game.Name, game.Condition, game.Source)
			if err != nil {
				return errors.Wrapf(err, "error inserting game %v", game)
			}
			gameID, err := result.LastInsertId()
			if err != nil {
				return errors.Wrap(err, "LastInsertId")
			}
			_, err = tx.Exec("INSERT INTO price (`game_id`, `value`) VALUES (?,?)", gameID, price.Value)
			if err != nil {
				return errors.Wrapf(err, "error inserting price with game_id %v value %v", gameID, price.Value)
			}
		}
	}

	logger.Info("price updated", zap.Int("count", numPriceUpdate))
	logger.Info("new games inserted", zap.Int("count", numNewGames))
	logger.Info("end")

	tx.Commit()
	return nil
}
