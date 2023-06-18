package strategy

import (
	"version-1/bj-simulator/internal/config"
	"version-1/bj-simulator/internal/player"
)

type Martingale struct{}

func (m Martingale) Bet(c config.Config, myself player.Player, players []player.Player, dealer player.Dealer) player.Act {
	if len(myself.History) == 0 {
		return player.Bet(-c.MinBet)
	}

	last := myself.History[len(myself.History)-1]
	if last.Result == player.Win {
		return player.Bet(c.MinBet)
	}

	bet := last.BetSummary() * 2
	if -c.MaxBet < bet {
		bet = c.MaxBet
	}

	return player.Bet(bet)
}
