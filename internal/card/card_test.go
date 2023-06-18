package card

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValue(t *testing.T) {

	tests := []struct {
		name   string
		card   Card
		expect int
	}{
		{
			name: "value is less than 10",
			card: Card{
				Kind:  Diamond,
				value: 9,
			},
			expect: 9,
		},
		{
			name: "value is 10",
			card: Card{
				Kind:  Diamond,
				value: 10,
			},
			expect: 10,
		},
		{
			name: "value is greater than 10",
			card: Card{
				Kind:  Diamond,
				value: 13,
			},
			expect: 10,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expect, test.card.Value())
		})
	}
}

func TestIsBlackjack(t *testing.T) {
	tests := []struct {
		name   string
		hands  Hands
		expect bool
	}{
		{
			name: "blackjack 1",
			hands: Hands([]Card{
				{
					Kind:  Diamond,
					value: 13,
				},
				{
					Kind:  Diamond,
					value: 1,
				},
			}),
			expect: true,
		},
		{
			name: "blackjack 2",
			hands: Hands([]Card{
				{
					Kind:  Diamond,
					value: 1,
				},
				{
					Kind:  Diamond,
					value: 10,
				},
			}),
			expect: true,
		},
		{
			name: "not blackjack 1",
			hands: Hands([]Card{
				{
					Kind:  Diamond,
					value: 10,
				},
				{
					Kind:  Diamond,
					value: 10,
				},
			}),
			expect: false,
		},
		{
			name: "not blackjack 2",
			hands: Hands([]Card{
				{
					Kind:  Diamond,
					value: 11,
				},
				{
					Kind:  Diamond,
					value: 5,
				},
				{
					Kind:  Diamond,
					value: 6,
				},
			}),
			expect: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expect, test.hands.IsBlackjack())
		})
	}
}

func TestIsBust(t *testing.T) {
	tests := []struct {
		name   string
		hands  Hands
		expect bool
	}{
		{
			name: "bust 1",
			hands: Hands([]Card{
				{
					Kind:  Diamond,
					value: 13,
				},
				{
					Kind:  Diamond,
					value: 2,
				},
				{
					Kind:  Spade,
					value: 10,
				},
			}),
			expect: true,
		},
		{
			name: "not bust",
			hands: Hands([]Card{
				{
					Kind:  Clover,
					value: 10,
				},
				{
					Kind:  Diamond,
					value: 6,
				},
			}),
			expect: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expect, test.hands.IsBust())
		})
	}
}

// TODO
// func TestSum(t *testing.T) {
// 	tests := []struct {
// 		name      string
// 		hands     Hands
// 		expectSum int
// 		bust      bool
// 		blackjack bool
// 	}{
// 		{
// 			name: "bust 1",
// 			hands: Hands([]Card{
// 				{
// 					Kind:  Diamond,
// 					value: 10,
// 				},
// 				{
// 					Kind:  Diamond,
// 					value: 1,
// 				},
// 			}),
// 			expect: true,
// 		},
// 		{
// 			name: "not bust",
// 			hands: Hands([]Card{
// 				{
// 					Kind:  Clover,
// 					value: 10,
// 				},
// 				{
// 					Kind:  Diamond,
// 					value: 6,
// 				},
// 			}),
// 			expect: false,
// 		},
// 	}
//
// 	for _, test := range tests {
// 		t.Run(test.name, func(t *testing.T) {
// 			assert.Equal(t, test.expect, test.hands.IsBust())
// 		})
// 	}
// }
