package player

import (
	"fmt"

	"github.com/version-1/bj-simulator/internal/card"
	"github.com/version-1/bj-simulator/internal/config"
)

type Player struct {
	History []*Round
	Amount  int

	bettingStrategy BettingStrategy
	handStrategy    HandStrategy
}

func New(amount int) *Player {
	return &Player{
		Amount:          amount,
		History:         []*Round{},
		bettingStrategy: defaultBettingStrategy{},
		handStrategy:    defaultHandStrategy{},
	}
}

func (p *Player) BettingStrategy(strategy BettingStrategy) *Player {
	p.bettingStrategy = strategy
	return p
}

func (p *Player) HandStrategy(strategy HandStrategy) *Player {
	p.handStrategy = strategy
	return p
}

func (p *Player) MakeAction(ctx *GameContext) error {
	current := p.CurrentRound()

	for !current.Done() {
		reason := p.Act(*ctx)

		switch reason {
		case ReasonHit:
			c := ctx.Pile.Pop()
			current.Hit(*c)
		case ReasonDoubleDown:
			current.DoubleDown(p)
			c := ctx.Pile.Pop()
			current.Hit(*c)
		case ReasonSplit:
			current.Split(p, current.Hands)
			c := ctx.Pile.Pop()
			current.Hit(*c)
		case ReasonStand:
			current.Acts = append(current.Acts, Stand())
		default:
			return fmt.Errorf("unexpected act for player.")
		}
	}

	return nil
}

type GameContext struct {
	Config           config.Config
	Pile             card.Pile
	Players          []Player
	Dealer           Dealer
	CurrentPlayCount int
}

func (g *GameContext) IncrementPlayCount() {
	g.CurrentPlayCount += 1
}

func (p *Player) Act(c GameContext) Reason {
	re := p.handStrategy.Act(c.Config, c.Pile, *p, c.Players, c.Dealer)

	current := p.CurrentRound()
	return validateAct(current, re)
}

func validateAct(r *Round, re Reason) Reason {
	h := card.Hands(r.Hands)

	if re == ReasonSplit && !h.CanSplit() {
		return ReasonHit
	}

	return re
}

func (p *Player) Bet(c GameContext) (Act, error) {
	bettingAct := p.bettingStrategy.Bet(c.Config, c.Pile, *p, c.Players, c.Dealer)
	if bettingAct.Value > -c.Config.MinBet {
		return bettingAct, fmt.Errorf("betting amount must be greater equal than min bet. min bet: %d, bet: %d", c.Config.MinBet, -bettingAct.Value)
	}

	if bettingAct.Value < -c.Config.MaxBet {
		return bettingAct, fmt.Errorf("betting amount must be lesser equal than max bet. max bet: %d, bet: %d", c.Config.MaxBet, -bettingAct.Value)
	}

	betting := -bettingAct.Value
	if p.Amount < betting {
		return bettingAct, fmt.Errorf("bet exceeds player's amount. amount: %d, bet: %d", p.Amount, betting)
	}

	r := p.CurrentRound()
	r.Acts = append(r.Acts, bettingAct)

	p.Amount += bettingAct.Value

	return bettingAct, nil
}

func (p *Player) CurrentRound() *Round {
	r, ok := findCurrentRound(p.History)
	if ok {
		return r
	}

	rr := &Round{}
	p.History = append(p.History, rr)

	return p.History[len(p.History)-1]
}

func findCurrentRound(rounds []*Round) (*Round, bool) {
	if rounds == nil {
		return nil, false
	}

	for i := range rounds {
		if rounds[i].Result == "" {
			return rounds[i], true
		}

		if rounds[i].Result == Splitted {
			return findCurrentRound(rounds[i].Rounds)
		}
	}

	return nil, false
}

func (p *Player) Hit(c card.Card) {
	r := p.CurrentRound()

	r.Hit(c)
}
