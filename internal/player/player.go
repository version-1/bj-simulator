package player

import (
	"fmt"
	"math"
	"version-1/bj-simulator/internal/card"
	"version-1/bj-simulator/internal/config"
)

type Player struct {
	History []Round
	Amount  int

	bettingStrategy BettingStrategy
	handStrategy    HandStrategy
}

func New(amount int) *Player {
	return &Player{
		Amount: amount,
	}
}

func (p *Player) MakeAction(c config.Config, pile card.Pile, myself Player, players []Player, dealer Dealer) error {
	last := p.CurrentRound()
	hands := card.Hands(last.Hands)
	for !hands.IsBust() || !last.Done() {
		act := myself.Act(c, pile, myself, players, dealer)
		last.Acts = append(last.Acts, act)

		switch act.Reason {
		case ReasonHit, ReasonDoubleDown:
			c := pile.Pop()
			p.Hit(*c)
		case ReasonSplit:
			last.Split(last.Hands)
		case ReasonStand:
		default:
			return fmt.Errorf("unexpected act for player.")
		}
	}

	return nil
}

func (p *Player) Act(c config.Config, pile card.Pile, players []Player, dealer Dealer) Act {
	return p.handStrategy.Act(c, pile, *p, players, dealer)
}

func (p *Player) Bet(c config.Config, pile card.Pile, players []Player, dealer Dealer) (Act, error) {
	bettingAct := p.bettingStrategy.Bet(c, pile, *p, players, dealer)
	if bettingAct.Value > -c.MinBet {
		return bettingAct, fmt.Errorf("betting amount must be greater than min bet. min bet: %d, bet: %d", c.MinBet, -bettingAct.Value)

	}

	r := p.CurrentRound()
	r.Acts = append(r.Acts, bettingAct)
	if p.Amount < bettingAct.Value {
		return bettingAct, fmt.Errorf("bet exceeds players amount. amount: %d, bet: %d", p.Amount, -bettingAct.Value)
	}

	p.Amount += bettingAct.Value

	return bettingAct, nil
}

func (p *Player) CurrentRound() *Round {
	r, ok := findCurrentRound(p.History)
	if ok {
		return r
	}

	rr := Round{}
	p.History = append(p.History, rr)

	return &rr
}

func findCurrentRound(rounds []Round) (*Round, bool) {
	for _, r := range rounds {
		if r.Result == "" {
			return &r, true
		}

		if r.Result == Splitted {
			return findCurrentRound(r.Rounds)
		}
	}

	return nil, true
}

func (p *Player) Hit(c card.Card) {
	r := p.CurrentRound()
	r.Hit(c)
}

type Dealer struct {
	Player
}

func (d Dealer) Result(r Round) Result {
	if r.IsBust() {
		return Lose
	}

	dealerRound := d.CurrentRound()

	if dealerRound.Sum() == r.Sum() {
		return Draw
	}

	if dealerRound.Sum() > r.Sum() {
		return Lose
	}

	return Win
}

type BettingStrategy interface {
	Bet(c config.Config, p card.Pile, myself Player, plyayers []Player, dealer Dealer) Act
}

type HandStrategy interface {
	Act(c config.Config, p card.Pile, myself Player, plyayers []Player, dealer Dealer) Act
}

type Result string

const (
	Win         Result = "win"
	Lose               = "lose"
	Draw               = "draw"
	Surrendered        = "surrendered"
	Insured            = "insured"
	Splitted           = "splitted"
)

type Round struct {
	Result    Result
	Hands     []card.Card
	Blackjack bool
	Acts      []Act
	Rounds    []Round
}

func (r *Round) IsBust() bool {
	hands := card.Hands(r.Hands)

	return hands.IsBust()
}

func (r *Round) IsBlackjack() bool {
	hands := card.Hands(r.Hands)

	return hands.IsBlackjack()
}

func (r *Round) Hit(c card.Card) {
	r.Hands = append(r.Hands, c)

	r.Acts = append(r.Acts, Hit())
}

func (r *Round) Return() int {
	if r.Result == Lose {
		return 0
	}

	bet := r.BetSummary()

	if r.Result == Draw {
		r.Acts = append(r.Acts, Return(bet))
		return bet
	}

	re := -bet * 2
	if r.IsBlackjack() {
		re = int(math.Floor(float64(-bet) * 2.5))
	}

	r.Acts = append(r.Acts, Return(re))
	return re
}

func (r Round) Sum() int {
	sum := 0
	for _, v := range r.Rounds {
		sum += v.Sum()
	}

	for _, v := range r.Acts {
		sum += v.Value
	}

	return sum
}

func (r Round) BetSummary() int {
	bet := 0
	for _, a := range r.Acts {
		bet += a.Value
	}

	return bet
}

func (r Round) Tail() *Act {
	if len(r.Acts) == 0 {
		return nil
	}

	return &r.Acts[len(r.Acts)-1]
}

func (r *Round) Split(cards []card.Card) error {
	if len(cards) != 2 {
		return fmt.Errorf("split must have one pair card. count: %d", len(cards))
	}

	r.Result = Splitted
	r.Rounds = []Round{
		{Hands: []card.Card{cards[0]}},
		{Hands: []card.Card{cards[1]}},
	}

	return nil
}

func (r Round) Done() bool {
	for i := len(r.Acts) - 1; i >= 0; i-- {
		if r.Acts[i].Reason == ReasonStand || r.Acts[i].Reason == ReasonDoubleDown {
			return true
		}
	}

	return false
}
