package service

import (
	"team-task/internal/dto"
	"team-task/internal/storage"
	"team-task/internal/stream"
)

type UserGrade interface {
	Set(userGrade dto.UserGrade) (dto.UserGrade, error)
	Get(userID string) (dto.UserGrade, error)
}

type Backupper interface {
	Backup() (string, error)
}

type Service struct {
	UserGrade
	Backupper
}

func NewService(storage *storage.Storage, stan *stream.STANClient) *Service {
	return &Service{
		UserGrade: NewUserGradeService(storage.UserGrade, *stan),
		Backupper: NewBackupperService(storage),
	}
}
