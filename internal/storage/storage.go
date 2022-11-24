package storage

import (
	"team-task/internal/dto"
)

type UserGrade interface {
	Set(userGrade dto.UserGrade)
	Get(userID string) (dto.UserGrade, error)
	GetAll() map[string]dto.UserGrade
}

type Storage struct {
	UserGrade
}

func NewStorage() *Storage {
	return &Storage{UserGrade: NewUserGradeStorage()}
}
