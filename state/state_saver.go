package state

import (
	"errors"
	"sync"
)

type saver struct {
	userMutex map[string]*sync.Mutex
	userState map[string]interface{}
}

// ClearState implements StateSaver.
func (s *saver) ClearState(user string) error {
	v := s.getUserMutex(user)
	v.Lock()
	delete(s.userState,user)
	v.Unlock()
	return nil
}

// GetState implements StateSaver.
func (s *saver) GetState(user string) (interface{}, error) {
	v,ok := s.userState[user]
	if !ok{
		return nil,errors.New("no user data")
	}
	return v,nil
}

// SetState implements StateSaver.
func (s *saver) SetState(user string, state interface{}) error {
	v := s.getUserMutex(user)
	v.Lock()
	s.userState[user] = state
	v.Unlock()
	return nil
}

func (s *saver) getUserMutex(user string)(*sync.Mutex){
	_ ,ok := s.userMutex[user]
	if !ok {
		s.userMutex[user] = &sync.Mutex{}
	}
	return s.userMutex[user]	
}

type Saver interface {
	SetState(user string, state interface{}) error
	GetState(user string) (interface{}, error)
	ClearState(user string) error
}

func NewStateSaver() Saver {
	return &saver{
		userMutex:     map[string]*sync.Mutex{},
		userState: map[string]interface{}{},
	}
}
