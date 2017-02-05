package gops

import "fmt"

type Game struct {
	ID        int    `db:"id, primarykey, autoincrement"`
	Name      string `db:"name"`
	Condition string `db:"condition"`
	Source    string `db:"source"`
}

// GetKey returns a key that uniquely identifies a game.
func (g *Game) GetKey() string {
	return fmt.Sprintf("%v\x00%v\x00%v", g.Name, g.Condition, g.Source)
}

type GameBuilder struct {
	game *Game
}

func NewGameBuilder() *GameBuilder {
	return &GameBuilder{
		game: &Game{},
	}
}

func (builder *GameBuilder) Build() *Game {
	return builder.game
}

func (builder *GameBuilder) SetName(name string) *GameBuilder {
	builder.game.Name = name
	return builder
}

func (builder *GameBuilder) FromGamestop() *GameBuilder {
	builder.game.Source = ProductSourceGamestop
	return builder
}

func (builder *GameBuilder) FromTarget() *GameBuilder {
	builder.game.Source = ProductSourceTarget
	return builder
}

func (builder *GameBuilder) FromWalmart() *GameBuilder {
	builder.game.Source = ProductSourceWalmart
	return builder
}

func (builder *GameBuilder) IsNew() *GameBuilder {
	builder.game.Condition = ProductConditionNew
	return builder
}

func (builder *GameBuilder) IsPreOwned() *GameBuilder {
	builder.game.Condition = ProductConditionPreowned
	return builder
}

func (builder *GameBuilder) SetCondition(condition string) *GameBuilder {
	builder.game.Condition = condition
	return builder
}
