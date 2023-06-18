package player

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/version-1/bj-simulator/internal/card"
)

func TestHit(t *testing.T) {
	tests := []struct {
		name   string
		card   card.Card
		before Round
		after  Round
	}{
		{
			name:   "5 decks",
			before: Round{},
			card:   *card.NewDiamond(1),
			after: Round{
				Hands: []card.Card{*card.NewDiamond(1)},
				Acts:  []Act{Hit()},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.before.Hit(test.card)

			assert.Equal(t, test.after, test.before)
		})
	}
}

func TestReturn(t *testing.T) {
	tests := []struct {
		name   string
		input  Round
		expect int
	}{
		{
			name: "when lose",
			input: Round{
				Result: Lose,
				Acts: []Act{
					{Reason: ReasonIntial, Value: -10},
				},
			},
			expect: 0,
		},
		{
			name: "when draw",
			input: Round{
				Result: Draw,
				Acts: []Act{
					{Reason: ReasonIntial, Value: -10},
				},
			},
			expect: 10,
		},
		{
			name: "when win",
			input: Round{
				Result: Win,
				Acts: []Act{
					{Reason: ReasonIntial, Value: -10},
				},
			},
			expect: 20,
		},
		{
			name: "when win and double down",
			input: Round{
				Result: Win,
				Acts: []Act{
					{Reason: ReasonIntial, Value: -10},
					{Reason: ReasonDoubleDown, Value: -10},
				},
			},
			expect: 40,
		},
		{
			name: "when split",
			input: Round{
				Result: Splitted,
				Acts: []Act{
					{Reason: ReasonIntial, Value: -10},
					{Reason: ReasonSplit},
				},
				Hands: []card.Card{
					*card.NewDiamond(1),
					*card.NewSpade(1),
				},
				Rounds: []*Round{
					{
						Result: Win,
						Hands: []card.Card{
							*card.NewDiamond(1),
							*card.NewSpade(7),
							*card.NewSpade(3),
						},
						Acts: []Act{
							{Reason: ReasonIntial, Value: -10},
						},
					},
					{
						Result: Splitted,
						Hands: []card.Card{
							*card.NewSpade(1),
							*card.NewHeart(1),
						},
						Acts: []Act{
							{Reason: ReasonIntial, Value: -10},
							{Reason: ReasonSplit},
						},
						Rounds: []*Round{
							{
								Result: Lose,
								Hands: []card.Card{
									*card.NewSpade(1),
									*card.NewSpade(3),
									*card.NewSpade(3),
									*card.NewSpade(10),
									*card.NewSpade(10),
								},
								Acts: []Act{
									{Reason: ReasonIntial, Value: -10},
								},
							},
							{
								Result: Win,
								Hands: []card.Card{
									*card.NewHeart(1),
									*card.NewSpade(7),
									*card.NewSpade(3),
								},
								Acts: []Act{
									{Reason: ReasonIntial, Value: -10},
								},
							},
						},
					},
				},
			},
			expect: 40,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			re := test.input.Return()

			assert.Equal(t, test.expect, re)
		})
	}
}

func TestDone(t *testing.T) {
	tests := []struct {
		name   string
		input  Round
		expect bool
	}{
		{
			name: "not done",
			input: Round{
				Result: Win,
				Hands: []card.Card{
					*card.NewDiamond(2),
					*card.NewDiamond(10),
				},
				Acts: []Act{
					{Reason: ReasonIntial, Value: -10},
				},
			},
			expect: false,
		},
		{
			name: "blackjack",
			input: Round{
				Result: Win,
				Hands: []card.Card{
					*card.NewDiamond(1),
					*card.NewDiamond(10),
				},
				Acts: []Act{
					{Reason: ReasonIntial, Value: -10},
				},
			},
			expect: true,
		},
		{
			name: "bust",
			input: Round{
				Result: Lose,
				Hands: []card.Card{
					*card.NewDiamond(2),
					*card.NewDiamond(10),
					*card.NewClover(13),
				},
				Acts: []Act{
					{Reason: ReasonIntial, Value: -10},
					{Reason: ReasonHit, Value: 0},
					{Reason: ReasonHit, Value: 0},
					{Reason: ReasonHit, Value: 0},
				},
			},
			expect: true,
		},
		{
			name: "doubledown",
			input: Round{
				Result: Lose,
				Hands: []card.Card{
					*card.NewDiamond(7),
					*card.NewDiamond(4),
					*card.NewClover(13),
				},
				Acts: []Act{
					{Reason: ReasonIntial, Value: -10},
					{Reason: ReasonHit, Value: 0},
					{Reason: ReasonHit, Value: 0},
					{Reason: ReasonDoubleDown, Value: 0},
					{Reason: ReasonHit, Value: 0},
				},
			},
			expect: true,
		},
		{
			name: "stand",
			input: Round{
				Result: Lose,
				Hands: []card.Card{
					*card.NewDiamond(7),
					*card.NewDiamond(4),
					*card.NewClover(13),
				},
				Acts: []Act{
					{Reason: ReasonIntial, Value: -10},
					{Reason: ReasonHit, Value: 0},
					{Reason: ReasonHit, Value: 0},
					{Reason: ReasonHit, Value: 0},
					{Reason: ReasonStand, Value: 0},
				},
			},
			expect: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			done := test.input.Done()

			assert.Equal(t, test.expect, done)
		})
	}
}
