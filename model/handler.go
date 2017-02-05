package gops

type GamePrice struct {
	Game  *Game
	Price *Price
}

type GameHandler interface {
	GetGames() ([]GamePrice, error)
	GetSource() string
}

var AllGameHandlers = []GameHandler{
	&GamestopHandler{},
	&TargetHandler{},
	&WalmartHandler{},
}
