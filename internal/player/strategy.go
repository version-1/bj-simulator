package player

import (
	"github.com/version-1/bj-simulator/internal/card"
	"github.com/version-1/bj-simulator/internal/config"
)

type BettingStrategy interface {
	Bet(c config.Config, p card.Pile, myself Player, plyayers []Player, dealer Dealer) Act
}

type HandStrategy interface {
	Act(c config.Config, p card.Pile, myself Player, plyayers []Player, dealer Dealer) Reason
}

type defaultBettingStrategy struct{}

func (p defaultBettingStrategy) Bet(c config.Config, pile card.Pile, myself Player, players []Player, dealer Dealer) Act {
	return Bet(-c.MinBet)
}

type defaultHandStrategy struct{}

func (p defaultHandStrategy) Act(c config.Config, pile card.Pile, myself Player, players []Player, dealer Dealer) Reason {
	return ReasonStand
}

type defaultDealerHandStrategy struct{}

func (s defaultDealerHandStrategy) Act(c config.Config, pile card.Pile, myself Player, players []Player, dealer Dealer) Reason {
	r := myself.CurrentRound()
	mh := card.Hands(r.Hands)
	msum, _, _ := mh.Sum()

	if msum <= 16 {
		return ReasonHit
	}

	return ReasonHit
}
