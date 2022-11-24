package service

import (
	"encoding/json"
	"team-task/internal/dto"
	"team-task/internal/storage"
	"team-task/internal/stream"
)

type UserGradeService struct {
	storage storage.UserGrade
	stan    stream.STANClient
}

func NewUserGradeService(storage storage.UserGrade, stan stream.STANClient) *UserGradeService {
	return &UserGradeService{storage: storage, stan: stan}
}

func (s *UserGradeService) Set(userGrade dto.UserGrade) (dto.UserGrade, error) {
	existingUserGrade, err := s.storage.Get(userGrade.UserId)
	if err == nil {
		userGrade = copyFields(userGrade, existingUserGrade)
	}

	s.storage.Set(userGrade)
	if err := replicate(s.stan, userGrade); err != nil {
		return userGrade, err
	}

	return userGrade, nil
}

func (s *UserGradeService) Get(userID string) (dto.UserGrade, error) {
	return s.storage.Get(userID)
}

func copyFields(src, dest dto.UserGrade) dto.UserGrade {
	res := dest
	if src.PostpaidLimit != 0 && dest.PostpaidLimit == 0 {
		res.PostpaidLimit = src.PostpaidLimit
	}
	if src.Spp != 0 && dest.Spp == 0 {
		res.Spp = src.Spp
	}
	if src.ShippingFee != 0 && dest.ShippingFee == 0 {
		res.ShippingFee = src.ShippingFee
	}
	if src.ReturnFee != 0 && dest.ReturnFee == 0 {
		res.ReturnFee = src.ReturnFee
	}

	return res
}

func replicate(stanClient stream.STANClient, userGrade dto.UserGrade) error {
	bytes, err := json.Marshal(userGrade)
	if err != nil {
		return err
	}

	if err = stanClient.Publish(stanClient.Subject, bytes); err != nil {
		return err
	}

	return nil
}
