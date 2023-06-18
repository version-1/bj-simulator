package config

type Config struct {
	DeckCount  int
	PlayCount  int
	MaxBet     int
	MinBet     int
	MinBetUnit int
	Surrender  bool

	InitialAmount int

	PlayerCount int
}

func defaultConfig() *Config {
	return &Config{
		DeckCount:  5,
		PlayCount:  10,
		MaxBet:     5,
		MinBet:     50,
		MinBetUnit: 5,

		InitialAmount: 1000,

		PlayerCount: 5,
		Surrender:   true,
	}
}

func New() *Config {
	return defaultConfig()
}
