package http

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"team-task/internal/dto"
	"team-task/internal/service"
	"team-task/internal/service/mocks"
	"testing"
)

func TestHandler_set(t *testing.T) {
	type mockBehavior = func(s *mock_service.MockUserGrade, userGrade dto.UserGrade)

	testTable := []struct {
		name                 string
		inputBody            string
		inputUserGrade       dto.UserGrade
		mockBehavior         mockBehavior
		expectErr            bool
		expectedStatus       int
		expectedResponseBody string
	}{
		{
			name:      "OK",
			inputBody: `{"user_id": "1","shipping_fee": 1,"return_fee": 1}`,
			inputUserGrade: dto.UserGrade{
				UserId:      "1",
				ShippingFee: 1,
				ReturnFee:   1,
			},
			mockBehavior: func(s *mock_service.MockUserGrade, userGrade dto.UserGrade) {
				s.EXPECT().Set(userGrade).Return(userGrade, nil)
			},
			expectErr:            false,
			expectedStatus:       200,
			expectedResponseBody: `{"user_id":"1","postpaid_limit":0,"spp":0,"shipping_fee":1,"return_fee":1}`,
		},
		{
			name:           "Request parsing error",
			inputBody:      `{"user_id": "1","shipping_fee": 1,}`,
			mockBehavior:   func(s *mock_service.MockUserGrade, userGrade dto.UserGrade) {},
			expectErr:      true,
			expectedStatus: 400,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			mockUserGrade := mock_service.NewMockUserGrade(c)
			testCase.mockBehavior(mockUserGrade, testCase.inputUserGrade)

			s := &service.Service{UserGrade: mockUserGrade}
			handler := NewHandler(s)

			r := gin.New()
			r.POST("/", handler.set)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(testCase.inputBody))
			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatus, w.Code)
			if !testCase.expectErr {
				assert.Equal(t, testCase.expectedResponseBody, w.Body.String())
			}
		})
	}
}

func TestHandler_get(t *testing.T) {
	type mockBehavior = func(s *mock_service.MockUserGrade, userID string)

	testTable := []struct {
		name                 string
		param                string
		mockBehavior         mockBehavior
		expectErr            bool
		expectedStatus       int
		expectedResponseBody string
	}{
		{
			name:  "OK",
			param: "1",
			mockBehavior: func(s *mock_service.MockUserGrade, userID string) {
				s.EXPECT().Get(userID).Return(dto.UserGrade{
					UserId:        "1",
					PostpaidLimit: 1,
					Spp:           1,
					ShippingFee:   1,
					ReturnFee:     1,
				}, nil)
			},
			expectErr:            false,
			expectedStatus:       200,
			expectedResponseBody: `{"user_id":"1","postpaid_limit":1,"spp":1,"shipping_fee":1,"return_fee":1}`,
		},
		{
			name:  "Missing query param",
			param: "",
			mockBehavior: func(s *mock_service.MockUserGrade, userID string) {
			},
			expectErr:      true,
			expectedStatus: 400,
		},
		{
			name:  "Service error",
			param: "1",
			mockBehavior: func(s *mock_service.MockUserGrade, userID string) {
				s.EXPECT().Get(userID).Return(dto.UserGrade{}, errors.New("some error"))
			},
			expectErr:      true,
			expectedStatus: 404,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			mockUserGrade := mock_service.NewMockUserGrade(c)
			testCase.mockBehavior(mockUserGrade, testCase.param)

			s := &service.Service{UserGrade: mockUserGrade}
			handler := NewHandler(s)

			r := gin.New()
			r.GET("/", handler.get)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/?user_id=%s", testCase.param),
				bytes.NewBufferString(testCase.param))
			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatus, w.Code)
			if !testCase.expectErr {
				assert.Equal(t, testCase.expectedResponseBody, w.Body.String())
			}
		})
	}
}
