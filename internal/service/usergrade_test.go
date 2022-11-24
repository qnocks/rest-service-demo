package service

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"team-task/internal/dto"
	"team-task/internal/storage/mocks"
	"team-task/internal/stream"
	"testing"
)

func TestService_Set(t *testing.T) {
	type mockBehavior = func(s *mock_storage.MockUserGrade, userGrade dto.UserGrade)

	testTable := []struct {
		name              string
		userGrade         dto.UserGrade
		mockBehavior      mockBehavior
		expectedUserGrade dto.UserGrade
		expectErr         bool
		errorMessage      string
	}{
		{
			name: "OK",
			userGrade: dto.UserGrade{
				UserId:        "1",
				PostpaidLimit: 1,
				Spp:           1,
				ShippingFee:   1,
				ReturnFee:     1,
			},
			mockBehavior: func(s *mock_storage.MockUserGrade, userGrade dto.UserGrade) {
				s.EXPECT().Get(userGrade.UserId).Return(dto.UserGrade{}, errors.New("some error"))
				s.EXPECT().Set(userGrade)
			},
			expectedUserGrade: dto.UserGrade{
				UserId:        "1",
				PostpaidLimit: 1,
				Spp:           1,
				ShippingFee:   1,
				ReturnFee:     1,
			},
			expectErr: false,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)

			mockUserGrade := mock_storage.NewMockUserGrade(c)
			testCase.mockBehavior(mockUserGrade, testCase.userGrade)

			mockStanClient, _ := stream.NewSTANClient("test-cluster", "test", "0.0.0.0:4222", "wb")
			userGrade := NewUserGradeService(mockUserGrade, *mockStanClient)
			actual, err := userGrade.Set(testCase.userGrade)

			if testCase.expectErr {
				assert.Equal(t, err.Error(), testCase.errorMessage)
			} else {
				assert.Equal(t, err, nil)
				assert.Equal(t, actual, testCase.expectedUserGrade)
			}
		})
	}
}

func TestService_Get(t *testing.T) {
	type mockBehavior = func(s *mock_storage.MockUserGrade, userID string)

	testTable := []struct {
		name              string
		userID            string
		mockBehavior      mockBehavior
		expectedUserGrade dto.UserGrade
		expectErr         bool
		errorMessage      string
	}{
		{
			name:   "OK",
			userID: "1",
			mockBehavior: func(s *mock_storage.MockUserGrade, userID string) {
				s.EXPECT().Get(userID).Return(dto.UserGrade{
					UserId:        userID,
					PostpaidLimit: 1,
					Spp:           1,
				}, nil)
			},
			expectedUserGrade: dto.UserGrade{
				UserId:        "1",
				PostpaidLimit: 1,
				Spp:           1,
			},
			expectErr: false,
		},
		{
			name:   "Storage error",
			userID: "1",
			mockBehavior: func(s *mock_storage.MockUserGrade, userID string) {
				s.EXPECT().Get(userID).Return(dto.UserGrade{}, errors.New("some error"))
			},
			expectedUserGrade: dto.UserGrade{},
			expectErr:         true,
			errorMessage:      "some error",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)

			mockUserGrade := mock_storage.NewMockUserGrade(c)
			testCase.mockBehavior(mockUserGrade, testCase.userID)

			userGrade := NewUserGradeService(mockUserGrade, stream.STANClient{})
			actual, err := userGrade.Get(testCase.userID)

			if testCase.expectErr {
				assert.Equal(t, err.Error(), testCase.errorMessage)
			} else {
				assert.Equal(t, err, nil)
				assert.Equal(t, testCase.expectedUserGrade, actual)
			}
		})
	}
}
