package main

import (
	"fmt"
	"version-1/bj-simulator/internal/card"
	"version-1/bj-simulator/internal/config"
	"version-1/bj-simulator/internal/player"
)

func main() {
	playCount := 1
	conf := config.New()
	pile := card.Prepare(conf.DeckCount)
	fmt.Println("Starting game")
	players := []player.Player{}
	for i := 0; i < conf.PlayCount; i++ {
		players = append(players, *player.New(conf.InitialAmount))
	}
	dealer := player.Dealer{}

	for conf.PlayCount > playCount {
		fmt.Printf("round start, count: %d", playCount)
		// betting
		for i := range players {
			players[i].Bet(*conf, *pile, players, dealer)
		}

		// first hit
		for i := range players {
			c := pile.Pop()
			players[i].Hit(*c)
		}

		c := pile.Pop()
		dealer.Hit(*c)

		// second hit
		for i := range players {
			c := pile.Pop()
			players[i].Hit(*c)
		}

		c = pile.Pop()
		dealer.Hit(*c)

		// hit or stand
		for i := range players {
			p := players[i]
			err := p.MakeAction(*conf, *pile, p, players, dealer)
			if err != nil {
				panic(fmt.Sprintf("got error for player %d: %s", i, err.Error()))
			}
		}

		dealer.MakeAction(*conf, *pile, dealer.Player, players, dealer)
		for i := range players {
			p := players[i]
			r := p.CurrentRound()

			r.Result = dealer.Result(*r)
			p.Amount += r.Return()
		}

		playCount++
	}

}
