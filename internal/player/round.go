package player

import (
	"fmt"
	"math"

	"github.com/version-1/bj-simulator/internal/card"
)

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
	Rounds    []*Round
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
	sum := calcReturn(*r)
	if len(r.Rounds) > 0 {
		for _, rr := range r.Rounds {
			sum += rr.Return()
		}
	}

	return sum
}

func calcReturn(r Round) int {
	if r.Result == Lose {
		return 0
	}

	bet := r.BetSummary()

	if r.Result == Draw {
		r.Acts = append(r.Acts, Return(bet))
		return bet
	}

	re := bet * 2
	if r.IsBlackjack() {
		re = int(math.Floor(float64(bet) * 2.5))
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
	if r.Result == Splitted {
		return bet
	}

	for _, a := range r.Acts {
		if a.Value < 0 {
			bet += a.Value
		}
	}

	return -bet
}

func (r Round) Tail() *Act {
	if len(r.Acts) == 0 {
		return nil
	}

	return &r.Acts[len(r.Acts)-1]
}

func (r *Round) Split(p *Player, cards []card.Card) error {
	if len(cards) != 2 {
		return fmt.Errorf("split must have one pair card. count: %d", len(cards))
	}

	initialBet := r.FindBy(ReasonIntial)

	r.Acts = append(r.Acts, Split(initialBet.Value))
	p.Amount += initialBet.Value
	r.Result = Splitted

	r.Rounds = []*Round{
		{Hands: []card.Card{cards[0]}, Acts: []Act{*initialBet}},
		{Hands: []card.Card{cards[1]}, Acts: []Act{*initialBet}},
	}

	return nil
}

func (r *Round) InitialBet() int {
	initialBet := r.FindBy(ReasonIntial)

	return initialBet.Value
}

func (r *Round) DoubleDown(p *Player) error {
	p.Amount += r.InitialBet()

	r.Acts = append(r.Acts, DoubleDown(r.InitialBet()))

	return nil
}

func (r *Round) FindBy(reason Reason) *Act {
	for i := range r.Acts {
		if r.Acts[i].Reason == reason {
			return &r.Acts[i]
		}
	}

	return nil
}

func (r Round) Done() bool {
	hands := card.Hands(r.Hands)

	if hands.IsBlackjack() || hands.IsBust() {
		return true
	}

	for i := len(r.Acts) - 1; i >= 0; i-- {
		if r.Acts[i].Reason == ReasonStand || r.Acts[i].Reason == ReasonDoubleDown {
			return true
		}
	}

	return false
}
