package counter

type Counter struct {
	ID        uint
	CurrValue uint
}

type Storage interface {
	Increment() uint
}
