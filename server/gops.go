package server

import (
	"fmt"
	"runtime/debug"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"github.com/wangkekekexili/gops/model"
	"github.com/wangkekekexili/gops/util"
	"go.uber.org/zap"
)

type GOPS struct {
	DB       *DB
	Logger   *Logger
	Reporter *Reporter

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

	gamesWithPrice, err := handler.GetGames()
	if err != nil {
		return err
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
	rows, err := s.DB.Query(`SELECT * FROM game WHERE name IN `+util.QuestionMarks(len(gameNames)), gameNames...)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var g model.Game
		err = rows.Scan(&g.ID, &g.Name, &g.Condition, &g.Source)
		if err != nil {
			return err
		}
		existingGames = append(existingGames, &g)
	}
	rows.Close()
	if rows.Err() != nil {
		return rows.Err()
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
			err = s.DB.QueryRow(`SELECT value FROM price WHERE game_id = ? ORDER BY timestamp desc`, existingGame.ID).
				Scan(&lastPrice)
			if err != nil {
				return err
			}
			if price.Value == lastPrice {
				continue
			}

			// New price point.
			numPriceUpdate++
			price.GameID = existingGame.ID
			_, err = s.DB.Exec(`INSERT INTO price (game_id, value) VALUES `+util.QuestionMarks(2), existingGame.ID, price.Value)
			if err != nil {
				return err
			}
		} else {
			// It's new game.
			numNewGames++
			result, err := s.DB.Exec("INSERT INTO game (`name`, `condition`, `source`) VALUES "+util.QuestionMarks(3),
				game.Name, game.Condition, game.Source)
			if err != nil {
				return err
			}
			gameID, err := result.LastInsertId()
			if err != nil {
				return err
			}
			_, err = s.DB.Exec(`INSERT INTO price (game_id, values) VALUES `+util.QuestionMarks(2), gameID, price.Value)
			if err != nil {
				return err
			}
		}
	}

	logger.Info("price updated", zap.Int("count", numPriceUpdate))
	logger.Info("new games inserted", zap.Int("count", numNewGames))
	logger.Info("end")
	return nil
}
