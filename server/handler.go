package server

import "github.com/wangkekekexili/gops/model"

// GameHandler defines an interface that can get a list of games and returns its source.
type GameHandler interface {
	GetGames() ([]*model.GamePrice, error)
	GetSource() string
}
