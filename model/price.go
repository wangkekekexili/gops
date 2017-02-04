package gops

import "time"

type Price struct {
	ID        int       `db:"id, primarykey, autoincrement"`
	GameID    int       `db:"game_id"`
	Value     float64   `db:"value"`
	Timestamp time.Time `db:"timestamp"`
}

func NewPrice(gameID int, value float64) *Price {
	return &Price{
		GameID:    gameID,
		Value:     value,
		Timestamp: time.Now(),
	}
}
