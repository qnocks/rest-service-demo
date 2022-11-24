package storage

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"team-task/internal/dto"
	"testing"
)

func TestStorage_Set(t *testing.T) {
	testTable := []struct {
		name            string
		givenStorage    map[string]dto.UserGrade
		inputUserGrade  dto.UserGrade
		expectedStorage map[string]dto.UserGrade
	}{
		{
			name:         "OK",
			givenStorage: map[string]dto.UserGrade{},
			inputUserGrade: dto.UserGrade{
				UserId:        "1",
				PostpaidLimit: 1,
				Spp:           1,
				ShippingFee:   1,
				ReturnFee:     1,
			},
			expectedStorage: map[string]dto.UserGrade{
				"1": {
					UserId:        "1",
					PostpaidLimit: 1,
					Spp:           1,
					ShippingFee:   1,
					ReturnFee:     1,
				},
			},
		},
		{
			name: "Overwriting setting",
			givenStorage: map[string]dto.UserGrade{
				"1": {
					UserId:        "1",
					PostpaidLimit: 1,
					Spp:           1,
					ShippingFee:   1,
					ReturnFee:     1,
				},
			},
			inputUserGrade: dto.UserGrade{
				UserId:        "1",
				PostpaidLimit: 2,
			},
			expectedStorage: map[string]dto.UserGrade{
				"1": {
					UserId:        "1",
					PostpaidLimit: 2,
					Spp:           0,
					ShippingFee:   0,
					ReturnFee:     0,
				},
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			s := UserGradeStorage{
				data: testCase.givenStorage,
				mu:   sync.RWMutex{},
			}
			s.Set(testCase.inputUserGrade)

			assert.Equal(t, testCase.expectedStorage, s.data)
		})
	}

}

func TestStorage_Get(t *testing.T) {
	testTable := []struct {
		name              string
		givenStorage      map[string]dto.UserGrade
		inputUserID       string
		expectedUserGrade dto.UserGrade
		expectErr         bool
		errorMessage      string
	}{
		{
			name: "Getting existing UserGrade",
			givenStorage: map[string]dto.UserGrade{
				"2": {
					UserId:        "2",
					PostpaidLimit: 2,
					Spp:           2,
					ShippingFee:   0,
					ReturnFee:     0,
				},
			},
			inputUserID: "2",
			expectedUserGrade: dto.UserGrade{
				UserId:        "2",
				PostpaidLimit: 2,
				Spp:           2,
				ShippingFee:   0,
				ReturnFee:     0,
			},
			expectErr:    false,
			errorMessage: "",
		},
		{
			name: "Getting nonexistent UserGrade",
			givenStorage: map[string]dto.UserGrade{
				"2": {
					UserId:        "2",
					PostpaidLimit: 2,
					Spp:           2,
					ShippingFee:   2,
					ReturnFee:     2,
				},
			},
			inputUserID:       "3",
			expectedUserGrade: dto.UserGrade{},
			expectErr:         true,
			errorMessage:      "cannot find [UserGrade] with [user_id=3]",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			s := UserGradeStorage{testCase.givenStorage, sync.RWMutex{}}

			actual, err := s.Get(testCase.inputUserID)

			if testCase.expectErr {
				assert.Equal(t, err.Error(), testCase.errorMessage)
			} else {
				assert.Equal(t, err, nil)
				assert.Equal(t, testCase.expectedUserGrade, actual)
			}
		})
	}
}
