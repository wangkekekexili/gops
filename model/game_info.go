package model

import (
	"fmt"
	"strings"
	"time"
)

const (
	ProductConditionNew      = "new"
	ProductConditionPreowned = "pre-owned"

	ProductSourceGamestop = "gamestop"
	ProductSourceTarget   = "target"
	ProductSourceWalmart  = "walmart"
)

// Game represents data in game table.
type Game struct {
	ID        int
	Name      string
	Condition string
	Source    string
}

// GetKey returns a key that uniquely identifies a game.
func (g *Game) GetKey() string {
	return fmt.Sprintf("%v\x00%v\x00%v", strings.ToLower(g.Name), g.Condition, g.Source)
}

func (g *Game) Stringer() string {
	return fmt.Sprintf("%v: %v %v from %v", g.ID, g.Condition, g.Name, g.Source)
}

// Price represents data in price table.
type Price struct {
	ID        int
	GameID    int
	Value     float64
	Timestamp time.Time
}

func NewPrice(gameID int, value float64) *Price {
	return &Price{
		GameID:    gameID,
		Value:     value,
		Timestamp: time.Now(),
	}
}

// GamePrice represents a combination of game and price data, without concerning unrelated database columns.
type GamePrice struct {
	Name      string
	Condition string
	Source    string
	Price     float64
}

func (g *GamePrice) ToGameAndPrice() (*Game, *Price) {
	return &Game{Name: g.Name, Condition: g.Condition, Source: g.Source}, &Price{Value: g.Price}
}

type GamePriceBuilder struct {
	*GamePrice
}

func NewGamePriceBuilder() *GamePriceBuilder {
	return &GamePriceBuilder{GamePrice: &GamePrice{}}
}

func (b *GamePriceBuilder) Build() *GamePrice {
	return b.GamePrice
}

func (b *GamePriceBuilder) SetName(name string) *GamePriceBuilder {
	b.GamePrice.Name = name
	return b
}

func (b *GamePriceBuilder) FromGamestop() *GamePriceBuilder {
	b.GamePrice.Source = ProductSourceGamestop
	return b
}

func (b *GamePriceBuilder) FromTarget() *GamePriceBuilder {
	b.GamePrice.Source = ProductSourceTarget
	return b
}

func (b *GamePriceBuilder) FromWalmart() *GamePriceBuilder {
	b.GamePrice.Source = ProductSourceWalmart
	return b
}

func (b *GamePriceBuilder) IsNew() *GamePriceBuilder {
	b.GamePrice.Condition = ProductConditionNew
	return b
}

func (b *GamePriceBuilder) IsPreOwned() *GamePriceBuilder {
	b.GamePrice.Condition = ProductConditionPreowned
	return b
}

func (b *GamePriceBuilder) SetCondition(condition string) *GamePriceBuilder {
	b.GamePrice.Condition = condition
	return b
}

func (b *GamePriceBuilder) SetPrice(p float64) *GamePriceBuilder {
	b.GamePrice.Price = p
	return b
}
