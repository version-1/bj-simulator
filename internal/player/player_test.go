package player

import (
	"fmt"
	"testing"

	"github.com/version-1/bj-simulator/internal/card"
	"github.com/version-1/bj-simulator/internal/config"

	"github.com/stretchr/testify/assert"
)

type underMinBetStrategy struct{}

func (s underMinBetStrategy) Bet(c config.Config, pile card.Pile, myself Player, players []Player, dealer Dealer) Act {
	return Bet(-c.MinBet + 1)
}

type exceedsMaxBetStrategy struct{}

func (s exceedsMaxBetStrategy) Bet(c config.Config, pile card.Pile, myself Player, players []Player, dealer Dealer) Act {
	return Bet(-c.MaxBet - 1)
}

func TestBet(t *testing.T) {
	conf := config.New()
	pile := &card.Pile{}
	dealer := NewDealer()

	gameContext := GameContext{
		Config:  *conf,
		Pile:    *pile,
		Players: []Player{},
		Dealer:  *dealer,
	}

	p1 := New(conf.InitialAmount)
	p2 := New(0)

	p3 := New(conf.InitialAmount)
	p3.BettingStrategy(underMinBetStrategy{})
	p4 := New(conf.InitialAmount)
	p4.BettingStrategy(exceedsMaxBetStrategy{})

	tests := []struct {
		name           string
		player         *Player
		ctx            func(ctx GameContext) GameContext
		expectedResult *Player
		expectedReturn Act
		expectedError  error
	}{
		{
			name:   "bet with default betting strategy",
			player: p1,
			ctx: func(ctx GameContext) GameContext {
				return ctx
			},
			expectedResult: &Player{
				Amount: 995,
				History: []*Round{
					{
						Acts: []Act{
							{
								Reason: ReasonIntial,
								Value:  -5,
							},
						},
					},
				},
				bettingStrategy: defaultBettingStrategy{},
				handStrategy:    defaultHandStrategy{},
			},
			expectedReturn: Act{
				Reason: ReasonIntial,
				Value:  -5,
			},
		},
		{
			name:   "bet with default betting strategy, the player doesn't have enough chip",
			player: p2,
			ctx: func(ctx GameContext) GameContext {
				return ctx
			},
			expectedResult: &Player{
				Amount:          0,
				History:         []*Round{},
				bettingStrategy: defaultBettingStrategy{},
				handStrategy:    defaultHandStrategy{},
			},
			expectedReturn: Act{
				Reason: ReasonIntial,
				Value:  -5,
			},
			expectedError: fmt.Errorf("bet exceeds player's amount. amount: 0, bet: 5"),
		},
		{
			name:   "bet with custom strategy, the betting is less than min bet.",
			player: p3,
			ctx: func(ctx GameContext) GameContext {
				return ctx
			},
			expectedResult: &Player{
				Amount:          1000,
				History:         []*Round{},
				bettingStrategy: underMinBetStrategy{},
				handStrategy:    defaultHandStrategy{},
			},
			expectedReturn: Act{
				Reason: ReasonIntial,
				Value:  -4,
			},
			expectedError: fmt.Errorf("betting amount must be greater equal than min bet. min bet: 5, bet: 4"),
		},
		{
			name:   "bet with custom strategy, the betting exceeds max bet.",
			player: p4,
			ctx: func(ctx GameContext) GameContext {
				return ctx
			},
			expectedResult: &Player{
				Amount:          1000,
				History:         []*Round{},
				bettingStrategy: exceedsMaxBetStrategy{},
				handStrategy:    defaultHandStrategy{},
			},
			expectedReturn: Act{
				Reason: ReasonIntial,
				Value:  -51,
			},
			expectedError: fmt.Errorf("betting amount must be lesser equal than max bet. max bet: 50, bet: 51"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := test.ctx(gameContext)
			act, err := test.player.Bet(ctx)

			assert.Equal(t, test.expectedError, err)
			assert.Equal(t, test.expectedReturn, act)
			assert.Equal(t, test.expectedResult, test.player)
		})
	}
}

type dummyHitStrategy struct{}

func (p dummyHitStrategy) Act(c config.Config, pile card.Pile, myself Player, players []Player, dealer Dealer) Reason {
	r := myself.CurrentRound()
	if len(r.Hands) == 2 {
		return ReasonHit
	}

	return ReasonStand
}

type dummyBustStrategy struct{}

func (p dummyBustStrategy) Act(c config.Config, pile card.Pile, myself Player, players []Player, dealer Dealer) Reason {
	return ReasonHit
}

type dummyDDStrategy struct{}

func (p dummyDDStrategy) Act(c config.Config, pile card.Pile, myself Player, players []Player, dealer Dealer) Reason {
	return ReasonDoubleDown
}

func TestMakeAction(t *testing.T) {
	tests := []struct {
		name          string
		before        func() (*Player, *GameContext)
		expect        *Player
		expectedError error
	}{
		{
			name: "when hit",
			before: func() (*Player, *GameContext) {
				p := card.NewPile(1)
				p.Prepare()
				p.Add(*card.NewDiamond(13))

				ctx := &GameContext{
					Pile: *p,
				}
				player := &Player{
					handStrategy: dummyHitStrategy{},
					Amount:       1000,
					History: []*Round{
						{
							Hands: []card.Card{
								*card.NewDiamond(5),
								*card.NewDiamond(6),
							},
							Acts: []Act{
								{
									Reason: ReasonIntial,
									Value:  -10,
								},
								Hit(),
								Hit(),
							},
						},
					},
				}

				return player, ctx
			},
			expect: &Player{
				handStrategy: dummyHitStrategy{},
				Amount:       1000,
				History: []*Round{
					{
						Hands: []card.Card{
							*card.NewDiamond(5),
							*card.NewDiamond(6),
							*card.NewDiamond(13),
						},
						Acts: []Act{
							{
								Reason: ReasonIntial,
								Value:  -10,
							},
							Hit(),
							Hit(),
							Hit(),
							Stand(),
						},
					},
				},
			},
		},
		{
			name: "when bust",
			before: func() (*Player, *GameContext) {
				p := card.NewPile(1)
				p.Prepare()
				p.Add(*card.NewDiamond(13))
				p.Add(*card.NewDiamond(5))

				ctx := &GameContext{
					Pile: *p,
				}
				player := &Player{
					handStrategy: dummyBustStrategy{},
					History: []*Round{
						{
							Hands: []card.Card{
								*card.NewDiamond(5),
								*card.NewDiamond(6),
							},
							Acts: []Act{
								{
									Reason: ReasonIntial,
									Value:  -10,
								},
								Hit(),
								Hit(),
							},
						},
					},
				}

				return player, ctx
			},
			expect: &Player{
				handStrategy: dummyBustStrategy{},
				History: []*Round{
					{
						Hands: []card.Card{
							*card.NewDiamond(5),
							*card.NewDiamond(6),
							*card.NewDiamond(5),
							*card.NewDiamond(13),
						},
						Acts: []Act{
							{
								Reason: ReasonIntial,
								Value:  -10,
							},
							Hit(),
							Hit(),
							Hit(),
							Hit(),
						},
					},
				},
			},
		},
		{
			name: "when double down",
			before: func() (*Player, *GameContext) {
				p := card.NewPile(1)
				p.Prepare()
				p.Add(*card.NewDiamond(13))

				ctx := &GameContext{
					Pile: *p,
				}
				player := &Player{
					handStrategy: dummyDDStrategy{},
					Amount:       1000,
					History: []*Round{
						{
							Hands: []card.Card{
								*card.NewDiamond(5),
								*card.NewDiamond(6),
							},
							Acts: []Act{
								{
									Reason: ReasonIntial,
									Value:  -10,
								},
								Hit(),
								Hit(),
							},
						},
					},
				}

				return player, ctx
			},
			expect: &Player{
				handStrategy: dummyDDStrategy{},
				Amount:       990,
				History: []*Round{
					{
						Hands: []card.Card{
							*card.NewDiamond(5),
							*card.NewDiamond(6),
							*card.NewDiamond(13),
						},
						Acts: []Act{
							{
								Reason: ReasonIntial,
								Value:  -10,
							},
							Hit(),
							Hit(),
							DoubleDown(-10),
							Hit(),
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p, ctx := test.before()
			err := p.MakeAction(ctx)

			assert.Equal(t, test.expectedError, err)
			assert.Equal(t, test.expect, p)
		})
	}
}

// func TestAct(t *testing.T) {
// 	tests := []struct {
// 	}
// }
