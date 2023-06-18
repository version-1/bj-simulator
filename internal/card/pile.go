package card

import (
	"math/rand"
	"time"
)

const deckSize int = 52

type Pile struct {
	cards     []Card
	deckCount int
}

func NewPile(deckCount int) *Pile {
	return &Pile{
		deckCount: deckCount,
	}
}

func (p Pile) Length() int {
	return len(p.cards)
}

func (p Pile) ShouldShuffle() bool {
	return p.Length() <= (deckSize * p.deckCount / 3)
}

func (p *Pile) Add(c Card) {
	p.cards = append(p.cards, c)
}

func (p *Pile) Pop() *Card {
	if p.ShouldShuffle() {
		p.Prepare()
	}
	last := p.cards[p.Length()-1]

	p.cards = p.cards[:p.Length()-1]

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
	return []Kind{Spade, Heart, Diamond, Clover}
}

func PrepareDeck() *Pile {
	p := Pile{}
	for _, k := range kinds() {
		for i := 1; i <= 13; i++ {
			p.Add(Card{Kind: k, value: i})
		}
	}

	return &p
}

func (p *Pile) Prepare() *Pile {
	p.cards = []Card{}
	for i := 0; i < p.deckCount; i++ {
		adding := PrepareDeck()
		p.Append(*adding)
	}

	p.Shuffle()

	return p
}
