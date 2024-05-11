package state

import (
	"errors"
	"sync"
)

type saver struct {
	userMutex           map[int64]*sync.Mutex
	userState           map[int64]interface{}
	tripMutex           map[string]*sync.Mutex
	globalTripState     map[string]interface{}
	globalTripObservers map[string][]int64
	tripGroup           map[string]int64
	activeTrip          map[int64]string
}

// RemoveUserState implements Saver.
func (s *saver) RemoveUserState(user int64) error {
	delete(s.userState, user)
	return nil
}

// DeleteActiveTrip implements Saver.
func (s *saver) DeleteActiveTrip(user int64) error {
	delete(s.activeTrip, user)
	return nil
}

// GetActiveTrip implements Saver.
func (s *saver) GetActiveTrip(user int64) (string, error) {
	v, ok := s.activeTrip[user]
	if !ok {
		return "", errors.New("no active trip")
	}
	return v, nil
}

// SetActiveTrip implements Saver.
func (s *saver) SetActiveTrip(user int64, tripCode string) error {
	s.activeTrip[user] = tripCode
	return nil
}

// GetTripGroup implements Saver.
func (s *saver) GetTripGroup(tripCode string) (int64, error) {
	v, ok := s.tripGroup[tripCode]
	if !ok {
		return -1, errors.New("no group found")
	}
	return v, nil
}

// SetTripGroup implements Saver.
func (s *saver) SetTripGroup(tripCode string, groupId int64) error {
	if len(tripCode) == 0 {
		return errors.New("trip code is empty")
	}
	if groupId == 0 {
		return errors.New("group id is empty")
	}
	_, ok := s.tripGroup[tripCode]
	if ok {
		return nil
	}
	s.tripGroup[tripCode] = groupId
	return nil
}

// GetTripObservers implements Saver.
func (s *saver) GetTripObservers(tripCode string) []int64 {
	v := s.getTripMutex(tripCode)
	v.Lock()
	obs, ok := s.globalTripObservers[tripCode]
	v.Unlock()
	if !ok {
		return []int64{}
	}
	return obs
}

// GetTripState implements Saver.
func (s *saver) GetTripState(tripCode string) (interface{}, error) {
	v, ok := s.globalTripState[tripCode]
	if !ok {
		return nil, errors.New("no trip state")
	}
	return v, nil
}

// AddTripObserver implements Saver.
func (s *saver) AddTripObserver(tripCode string, user int64) error {
	v := s.getTripMutex(tripCode)
	v.Lock()
	_, ok := s.globalTripObservers[tripCode]
	if !ok {
		s.globalTripObservers[tripCode] = []int64{}
	}
	mem := s.globalTripObservers[tripCode]
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
func (s *saver) ClearUserState(user int64) error {
	v := s.getUserMutex(user)
	v.Lock()
	delete(s.userState, user)
	v.Unlock()
	return nil
}

// GetUserState implements StateSaver.
func (s *saver) GetUserState(user int64) (interface{}, error) {
	v, ok := s.userState[user]
	if !ok {
		return nil, errors.New("no user data")
	}
	return v, nil
}

// SetUserState implements StateSaver.
func (s *saver) SetUserState(user int64, state interface{}) error {
	v := s.getUserMutex(user)
	v.Lock()
	s.userState[user] = state
	v.Unlock()
	return nil
}

func (s *saver) getUserMutex(user int64) *sync.Mutex {
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
	SetUserState(user int64, state interface{}) error
	RemoveUserState(user int64) error
	GetUserState(user int64) (interface{}, error)
	ClearUserState(user int64) error
	AddTripObserver(tripCode string, user int64) error
	GetTripObservers(tripCode string) []int64
	SetTripState(tripCode string, state interface{}) error
	GetTripState(tripCode string) (interface{}, error)
	SetTripGroup(tripCode string, groupId int64) error
	GetTripGroup(tripCode string) (int64, error)
	SetActiveTrip(user int64, tripCode string) error
	GetActiveTrip(user int64) (string, error)
	DeleteActiveTrip(user int64) error
}

func NewStateSaver() Saver {
	return &saver{
		userMutex:           map[int64]*sync.Mutex{},
		userState:           map[int64]interface{}{},
		tripMutex:           map[string]*sync.Mutex{},
		globalTripState:     map[string]interface{}{},
		globalTripObservers: map[string][]int64{},
		tripGroup:           map[string]int64{},
		activeTrip:          map[int64]string{},
	}
}
