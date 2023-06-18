package game

import (
	"fmt"

	"github.com/version-1/bj-simulator/internal/card"
	"github.com/version-1/bj-simulator/internal/config"
	"github.com/version-1/bj-simulator/internal/player"
)

type Game struct {
	ctx *player.GameContext
}

func New() *Game {
	conf := config.New()
	pile := card.NewPile(conf.DeckCount)
	pile.Prepare()

	players := []player.Player{}
	for i := 0; i < conf.PlayCount; i++ {
		players = append(players, *player.New(conf.InitialAmount))
	}
	dealer := player.NewDealer()
	ctx := &player.GameContext{
		Config:  *conf,
		Pile:    *pile,
		Players: players,
		Dealer:  *dealer,
	}

	return &Game{
		ctx: ctx,
	}
}

func (g Game) Play() {
	fmt.Println("starting game")
	for g.ctx.Config.PlayCount > g.PlayCount() {
		fmt.Printf("round start, count: %d", g.PlayCount())
		g.playRound()
	}
}

func (g Game) GameContext() *player.GameContext {
	return g.ctx
}

func (g Game) PlayCount() int {
	return g.ctx.CurrentPlayCount
}

func (g Game) playRound() {
	ctx := g.ctx

	players := ctx.Players
	dealer := ctx.Dealer
	pile := ctx.Pile

	// betting
	for i := range players {
		players[i].Bet(*g.ctx)
	}

	// first hit
	for i := range players {
		c := pile.Pop()
		players[i].Hit(*c)

		c = pile.Pop()
		players[i].Hit(*c)
	}

	c := pile.Pop()
	dealer.Hit(*c)

	c = pile.Pop()
	dealer.Hit(*c)

	// hit or stand
	for i := range players {
		p := players[i]
		err := p.MakeAction(g.ctx)
		if err != nil {
			panic(fmt.Sprintf("got error for player %d: %s", i, err.Error()))
		}
	}

	dealer.MakeAction(g.ctx)
	for i := range players {
		p := players[i]
		r := p.CurrentRound()

		r.Result = dealer.Result(*r)
		p.Amount += r.Return()
	}

	g.ctx.IncrementPlayCount()

}
