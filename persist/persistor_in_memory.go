package persist

import "github.com/kiyutink/tipam/tipam"

type InMemory struct {
	state *tipam.State
}

func NewInMemory() *InMemory {
	state := tipam.NewState()
	return &InMemory{
		state: state,
	}
}

func (im *InMemory) Persist(s *tipam.State) error {
	im.state = s
	return nil
}

func (im *InMemory) Read() (*tipam.State, error) {
	if im.state == nil {
		return tipam.NewState(), nil
	}

	return im.state, nil
}

func (im *InMemory) Lock() error {
	return nil
}

func (im *InMemory) Unlock() error {
	return nil
}
