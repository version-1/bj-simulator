package card

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrepare(t *testing.T) {
	tests := []struct {
		name      string
		deckCount int
		popCount  int
		expect    int
	}{
		{
			name:      "5 decks",
			deckCount: 5,
			popCount:  0,
			expect:    52 * 5,
		},
		{
			name:      "1 deck",
			deckCount: 1,
			popCount:  0,
			expect:    52 * 1,
		},
		{
			name:      "1 deck, 1 pop",
			deckCount: 1,
			popCount:  1,
			expect:    52 - 1,
		},
		{
			name:      "1 deck, 10 pop",
			deckCount: 1,
			popCount:  10,
			expect:    42,
		},
		{
			name:      "1 deck, 36 pop",
			deckCount: 1,
			popCount:  36,
			expect:    51,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := NewPile(test.deckCount)
			p.Prepare()

			for i := 0; i < test.popCount; i++ {
				p.Pop()
			}

			assert.Equal(t, test.expect, p.Length())
		})
	}

}
