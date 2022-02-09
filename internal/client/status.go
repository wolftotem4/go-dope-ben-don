package client

type Status struct {
	Logged bool
}

func NewStatus() *Status {
	return &Status{
		Logged: false,
	}
}
