package player

type Dealer struct {
	Player
}

func NewDealer() *Dealer {
	d := &Dealer{}
	d.HandStrategy(defaultDealerHandStrategy{})

	return d
}

func (d Dealer) Result(r Round) Result {
	if r.IsBust() {
		return Lose
	}

	if r.IsBlackjack() {
		return Win
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
