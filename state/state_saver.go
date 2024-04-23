package state

import (
	"errors"
	"sync"
)

type saver struct {
	userMutex           map[string]*sync.Mutex
	userState           map[string]interface{}
	tripMutex           map[string]*sync.Mutex
	globalTripState     map[string]interface{}
	globalTripObservers map[string][]string
}

// GetTripState implements Saver.
func (s *saver) GetTripState(tripCode string) (interface{}, error) {
	v,ok :=  s.globalTripState[tripCode]
	if !ok {
		return nil,errors.New("no trip state")
	}
	return v,nil
}

// AddTripObserver implements Saver.
func (s *saver) AddTripObserver(tripCode string, user string) error {
	v := s.getTripMutex(tripCode)
	v.Lock()
	mem, ok := s.globalTripObservers[tripCode]
	if !ok {
		s.globalTripObservers[tripCode] = []string{}
	}
	mem = s.globalTripObservers[tripCode]
	mem = append(mem, user)

	s.globalTripObservers[tripCode] = mem
	v.Unlock()
	return nil
}

// SetTripState implements Saver.
func (s *saver) SetTripState(tripCode string, state interface{}) error {
	v := s.getTripMutex(tripCode)
	v.Lock()
	s.globalTripState[tripCode] = state
	v.Unlock()
	return nil
}

// ClearUserState implements StateSaver.
func (s *saver) ClearUserState(user string) error {
	v := s.getUserMutex(user)
	v.Lock()
	delete(s.userState, user)
	v.Unlock()
	return nil
}

// GetUserState implements StateSaver.
func (s *saver) GetUserState(user string) (interface{}, error) {
	v, ok := s.userState[user]
	if !ok {
		return nil, errors.New("no user data")
	}
	return v, nil
}

// SetUserState implements StateSaver.
func (s *saver) SetUserState(user string, state interface{}) error {
	v := s.getUserMutex(user)
	v.Lock()
	s.userState[user] = state
	v.Unlock()
	return nil
}

func (s *saver) getUserMutex(user string) *sync.Mutex {
	_, ok := s.userMutex[user]
	if !ok {
		s.userMutex[user] = &sync.Mutex{}
	}
	return s.userMutex[user]
}
func (s *saver) getTripMutex(tripCode string) *sync.Mutex {
	_, ok := s.tripMutex[tripCode]
	if !ok {
		s.tripMutex[tripCode] = &sync.Mutex{}
	}
	return s.tripMutex[tripCode]
}

type Saver interface {
	SetUserState(user string, state interface{}) error
	GetUserState(user string) (interface{}, error)
	ClearUserState(user string) error
	AddTripObserver(tripCode string, user string) error
	SetTripState(tripCode string, state interface{}) error
	GetTripState(tripCode string) (interface{}, error)
}

func NewStateSaver() Saver {
	return &saver{
		userMutex:           map[string]*sync.Mutex{},
		userState:           map[string]interface{}{},
		tripMutex:           map[string]*sync.Mutex{},
		globalTripState:     map[string]interface{}{},
		globalTripObservers: map[string][]string{},
	}
}
