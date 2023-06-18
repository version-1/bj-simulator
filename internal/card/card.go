package card

type Kind string

const (
	Spade   Kind = "spade"
	Hart         = "heart"
	Diamond      = "diamond"
	Clover       = "clover"
)

type Card struct {
	Kind  Kind
	value int
}

func (c Card) Value() int {
	if c.value >= 10 {
		return 10
	}

	return c.value
}

func (c Card) Equal(v Card) bool {
	return c.value == v.value && c.Kind == v.Kind
}

type Hands []Card

func (h Hands) Find(i int) int {
	for j, v := range h {
		if v.Value() == i {
			return j
		}
	}

	return -1
}

func (h Hands) Sum() (sum int, bust, blackjack bool) {
	if h.IsBlackjack() {
		return 21, false, true
	}

	for _, v := range h {
		sum += v.Value()
		if sum > 21 {
			return sum, true, false
		}
	}

	return sum, false, false
}

func (h Hands) IsBlackjack() bool {
	if len(h) == 2 {
		ace := h.Find(1)
		ten := h.Find(10)
		return ace >= 0 && ten >= 0 && ace != ten
	}

	return false
}

func (h Hands) IsBust() bool {
	_, bust, _ := h.Sum()
	return bust
}

func (h Hands) CanSplit() bool {
	if len(h) == 2 {
		return h[0] == h[1]
	}

	return false
}
