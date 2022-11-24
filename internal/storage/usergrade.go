package storage

import (
	"fmt"
	"sync"
	"team-task/internal/dto"
)

type UserGradeStorage struct {
	data map[string]dto.UserGrade
	mu   sync.RWMutex
}

func (s *UserGradeStorage) GetAll() map[string]dto.UserGrade {
	return s.data
}

func NewUserGradeStorage() *UserGradeStorage {
	return &UserGradeStorage{
		data: make(map[string]dto.UserGrade),
		mu:   sync.RWMutex{},
	}
}

func (s *UserGradeStorage) Set(userGrade dto.UserGrade) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[userGrade.UserId] = userGrade
}

func (s *UserGradeStorage) Get(userID string) (dto.UserGrade, error) {
	var userGrade dto.UserGrade
	s.mu.RLock()
	defer s.mu.RUnlock()

	userGrade = s.data[userID]
	if len(userGrade.UserId) == 0 {
		return userGrade, fmt.Errorf("cannot find [UserGrade] with [user_id=%s]", userID)
	}

	return userGrade, nil
}
