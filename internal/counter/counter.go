package counter

type Counter struct {
	ID    uint
	Value uint
}

type Storage interface {
	Init() error
	Increment() (uint, error)
}
