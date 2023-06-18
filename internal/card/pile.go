package card

import (
	"math/rand"
	"time"
)

const deckSize int = 52

type Pile struct {
	cards []Card
}

func (p Pile) Length() int {
	return len(p.cards)
}

func (p Pile) ShouldShuffle(deckCount int) bool {
	return p.Length() <= (deckSize * deckCount / 3)
}

func (p *Pile) Add(c Card) {
	p.cards = append(p.cards, c)
}

func (p *Pile) Pop() *Card {
	last := p.cards[p.Length()-1]

	p.cards = p.cards[:p.Length()-2]

	return &last
}

func (p *Pile) Append(appending Pile) {
	p.cards = append(p.cards, appending.cards...)
}

func (p *Pile) Set(index int, value Card) Card {
	prev := p.cards[index]
	p.cards[index] = value

	return prev
}

func (p *Pile) Get(index int) Card {
	return p.cards[index]
}

func (p *Pile) Swap(i, j int) {
	tmp := p.Set(i, p.Get(j))
	p.Set(j, tmp)
}

func (p *Pile) Shuffle() {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(p.Length(), p.Swap)
}

func kinds() []Kind {
	return []Kind{Spade, Hart, Diamond, Clover}
}

func New() *Pile {
	p := Pile{}
	for _, k := range kinds() {
		for i := 1; i <= 13; i++ {
			p.Add(Card{Kind: k, value: i})
		}
	}

	return &p
}

func Prepare(deckCount int) *Pile {
	p := &Pile{}
	for i := 0; i < deckCount; i++ {
		adding := New()
		p.Append(*adding)
	}

	p.Shuffle()

	return p
}
