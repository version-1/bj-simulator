package player

type Reason string

const (
	ReasonIntial     Reason = "initial"
	ReasonDoubleDown        = "doubledown"
	ReasonSplit             = "split"
	ReasonSurrender         = "surrender"
	ReasonInsure            = "insure"
	ReasonReturn            = "return"
	ReasonHit               = "hit"
	ReasonStand             = "stand"
)

type Act struct {
	Reason Reason
	Value  int
}

func Bet(v int) Act {
	return Act{
		Reason: ReasonIntial,
		Value:  v,
	}
}

func Hit() Act {
	return Act{
		Reason: ReasonHit,
	}
}

func DoubleDown(v int) Act {
	return Act{
		Reason: ReasonDoubleDown,
		Value:  v,
	}
}

func Split(v int) Act {
	return Act{
		Reason: ReasonSplit,
		Value:  v,
	}
}

func Stand() Act {
	return Act{
		Reason: ReasonStand,
	}
}

func Return(num int) Act {
	return Act{
		Reason: ReasonReturn,
		Value:  num,
	}
}
