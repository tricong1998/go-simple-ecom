package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tricong1998/go-ecom/cmd/payment/internal/mocks"
	"github.com/tricong1998/go-ecom/cmd/payment/internal/services"
	"github.com/tricong1998/go-ecom/cmd/payment/pkg/dto"
	"github.com/tricong1998/go-ecom/cmd/payment/pkg/models"
)

func expectBodyPayment(t *testing.T, w *httptest.ResponseRecorder, mockResponse *models.Payment) {
	var response dto.PaymentResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	fmt.Println("response", response, response.UserId)
	fmt.Println("mockResponse", mockResponse, mockResponse.UserID)
	assert.NoError(t, err)
	assert.Equal(t, response.ID, mockResponse.ID)
	assert.Equal(t, response.UserId, mockResponse.UserID)
	assert.Equal(t, response.OrderId, mockResponse.OrderID)
	assert.Equal(t, response.Error, mockResponse.Error)
	assert.Equal(t, response.Amount, mockResponse.Amount)
	assert.Equal(t, response.Method, mockResponse.Method)
	assert.WithinDuration(t, response.CreatedAt, mockResponse.CreatedAt, time.Second)
	assert.WithinDuration(t, response.UpdatedAt, mockResponse.UpdatedAt, time.Second)
}

func TestCreatePayment(t *testing.T) {
	testCases := []struct {
		name           string
		setupInputFunc func(input *dto.CreatePaymentDto, mockResponse *models.Payment)
		mockFunc       func(paymentRepo *mocks.MockPaymentRepository, mockResponse *models.Payment)
		expectFunc     func(w *httptest.ResponseRecorder, mockResponse *models.Payment)
	}{
		{
			name: "OK",
			setupInputFunc: func(input *dto.CreatePaymentDto, mockResponse *models.Payment) {
				input.OrderId = 1
				input.UserId = 1
				input.Amount = 100
				input.Method = "cash"
				mockResponse.ID = 1
				mockResponse.CreatedAt = time.Now()
				mockResponse.UpdatedAt = mockResponse.CreatedAt
				mockResponse.Status = "pending"
				mockResponse.Method = "cash"
				mockResponse.Amount = 100
				mockResponse.Error = ""
				mockResponse.UserID = 1
				mockResponse.OrderID = 1
			},
			mockFunc: func(paymentRepo *mocks.MockPaymentRepository, mockResponse *models.Payment) {
				paymentRepo.On("CreatePayment", mock.AnythingOfType("*models.Payment")).Return(nil).Run(func(args mock.Arguments) {
					arg := args.Get(0).(*models.Payment)
					arg.ID = mockResponse.ID
					arg.CreatedAt = mockResponse.CreatedAt
					arg.UpdatedAt = mockResponse.UpdatedAt
				})
				paymentRepo.On("UpdatePayment", mock.AnythingOfType("*models.Payment")).Return(nil)
			},
			expectFunc: func(w *httptest.ResponseRecorder, mockResponse *models.Payment) {
				assert.Equal(t, http.StatusCreated, w.Code)
				expectBodyPayment(t, w, mockResponse)
			},
		},
		{
			name: "BadInput",
			setupInputFunc: func(input *dto.CreatePaymentDto, mockResponse *models.Payment) {
			},
			mockFunc: func(paymentRepo *mocks.MockPaymentRepository, mockResponse *models.Payment) {
				paymentRepo.On("CreatePayment", mock.AnythingOfType("*models.Payment")).Return(nil).Run(func(args mock.Arguments) {
					arg := args.Get(0).(*models.Payment)
					arg.ID = mockResponse.ID
					arg.CreatedAt = mockResponse.CreatedAt
					arg.UpdatedAt = mockResponse.UpdatedAt
				})
			},
			expectFunc: func(w *httptest.ResponseRecorder, mockResponse *models.Payment) {
				assert.Equal(t, http.StatusBadRequest, w.Code)
			},
		},
		{
			name: "CreatePaymentError",
			setupInputFunc: func(input *dto.CreatePaymentDto, mockResponse *models.Payment) {
				input.OrderId = 1
				input.UserId = 1
				input.Amount = 100
				input.Method = "cash"
				mockResponse.ID = 1
				mockResponse.CreatedAt = time.Now()
				mockResponse.UpdatedAt = mockResponse.CreatedAt
				mockResponse.Status = "pending"
				mockResponse.Error = ""
			},
			mockFunc: func(paymentRepo *mocks.MockPaymentRepository, mockResponse *models.Payment) {
				err := errors.New("Error")
				paymentRepo.On("CreatePayment", mock.AnythingOfType("*models.Payment")).Return(err)
			},
			expectFunc: func(w *httptest.ResponseRecorder, mockResponse *models.Payment) {
				assert.Equal(t, http.StatusInternalServerError, w.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			paymentRepo := new(mocks.MockPaymentRepository)
			paymentService := services.NewPaymentService(paymentRepo)
			paymentHandler := NewPaymentHandler(paymentService)
			var payment dto.CreatePaymentDto
			var mockResponse models.Payment
			tc.setupInputFunc(&payment, &mockResponse)
			tc.mockFunc(paymentRepo, &mockResponse)
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			jsonPayment, _ := json.Marshal(payment)
			c.Request, _ = http.NewRequest(http.MethodPost, "/payments", bytes.NewBuffer(jsonPayment))
			c.Request.Header.Set("Content-Type", "application/json")
			paymentHandler.CreatePayment(c)
			tc.expectFunc(w, &mockResponse)
		})
	}
}

func TestReadPayment(t *testing.T) {
	testCases := []struct {
		name           string
		setupInputFunc func(input *dto.ReadPaymentRequest, mockResponse *models.Payment)
		mockFunc       func(paymentRepo *mocks.MockPaymentRepository, mockResponse *models.Payment)
		expectFunc     func(w *httptest.ResponseRecorder, mockResponse *models.Payment)
	}{
		{
			name: "OK",
			setupInputFunc: func(input *dto.ReadPaymentRequest, mockResponse *models.Payment) {
				input.ID = 1
				mockResponse.OrderID = 1
				mockResponse.UserID = 1
				mockResponse.Amount = 100
				mockResponse.Method = "cash"
				mockResponse.ID = 1
				mockResponse.CreatedAt = time.Now()
				mockResponse.UpdatedAt = mockResponse.CreatedAt
				mockResponse.Status = "pending"
				mockResponse.Error = ""
			},
			mockFunc: func(paymentRepo *mocks.MockPaymentRepository, mockResponse *models.Payment) {
				paymentRepo.On("ReadPayment", mock.AnythingOfType("uint")).Return(mockResponse, nil)
			},
			expectFunc: func(w *httptest.ResponseRecorder, mockResponse *models.Payment) {
				assert.Equal(t, http.StatusOK, w.Code)
				expectBodyPayment(t, w, mockResponse)
			},
		},
		{
			name: "BadInput",
			setupInputFunc: func(input *dto.ReadPaymentRequest, mockResponse *models.Payment) {
			},
			mockFunc: func(paymentRepo *mocks.MockPaymentRepository, mockResponse *models.Payment) {
				paymentRepo.On("ReadPayment", mock.AnythingOfType("uint")).Return(mockResponse, nil)
			},
			expectFunc: func(w *httptest.ResponseRecorder, mockResponse *models.Payment) {
				assert.Equal(t, http.StatusBadRequest, w.Code)
			},
		},
		{
			name: "PaymentNotFound",
			setupInputFunc: func(input *dto.ReadPaymentRequest, mockResponse *models.Payment) {
				input.ID = 1
			},
			mockFunc: func(paymentRepo *mocks.MockPaymentRepository, mockResponse *models.Payment) {
				paymentRepo.On("ReadPayment", mock.AnythingOfType("uint")).Return(mockResponse, errors.New("Not found"))
			},
			expectFunc: func(w *httptest.ResponseRecorder, mockResponse *models.Payment) {
				assert.Equal(t, http.StatusNotFound, w.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			paymentRepo := new(mocks.MockPaymentRepository)
			paymentService := services.NewPaymentService(paymentRepo)
			paymentHandler := NewPaymentHandler(paymentService)
			var input dto.ReadPaymentRequest
			var mockResponse models.Payment
			tc.setupInputFunc(&input, &mockResponse)
			tc.mockFunc(paymentRepo, &mockResponse)
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Params = gin.Params{{Key: "id", Value: fmt.Sprint(input.ID)}}

			// Act
			paymentHandler.ReadPayment(c)

			// Assert
			tc.expectFunc(w, &mockResponse)
		})
	}
}

func TestListPayment(t *testing.T) {
	testCases := []struct {
		name           string
		setupInputFunc func(input *dto.ListPaymentQuery, total *int64) []models.Payment
		mockFunc       func(paymentRepo *mocks.MockPaymentRepository, mockResponse []models.Payment, input *dto.ListPaymentQuery, total *int64)
		expectFunc     func(
			w *httptest.ResponseRecorder,
			mockResponse []models.Payment,
			input *dto.ListPaymentQuery,
			total *int64,
		)
	}{
		{
			name: "OK",
			setupInputFunc: func(input *dto.ListPaymentQuery, total *int64) []models.Payment {
				var mockResponse []models.Payment
				input.Page = int32(1)
				input.PerPage = int32(10)
				userId := uint(1)
				input.UserId = &userId
				*total = 10
				now := time.Now()
				payment1 := models.Payment{
					OrderID: 1,
					UserID:  userId,
					Amount:  100,
					Method:  "cash",
					Status:  "pending",
					Error:   "",
				}
				payment1.CreatedAt = now
				payment1.UpdatedAt = now
				payment1.ID = 1
				mockResponse = append(mockResponse, payment1)

				payment2 := models.Payment{
					OrderID: 2,
					UserID:  userId,
					Amount:  100,
					Method:  "cash",
					Status:  "pending",
					Error:   "",
				}
				payment2.CreatedAt = now
				payment2.UpdatedAt = now
				payment2.ID = 2
				mockResponse = append(mockResponse, payment2)
				return mockResponse
			},
			mockFunc: func(paymentRepo *mocks.MockPaymentRepository, mockResponse []models.Payment, input *dto.ListPaymentQuery, total *int64) {
				paymentRepo.On("ListPayments", input.PerPage, input.Page, input.UserId).Return(mockResponse, *total, nil)
			},
			expectFunc: func(
				w *httptest.ResponseRecorder,
				mockResponse []models.Payment,
				input *dto.ListPaymentQuery,
				total *int64,
			) {
				assert.Equal(t, http.StatusOK, w.Code)
				var response dto.ListPaymentResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, response.Metadata.Page, input.Page)
				assert.Equal(t, response.Metadata.PerPage, input.PerPage)
				assert.Equal(t, response.Metadata.Total, *total)
				assert.Len(t, response.Items, len(mockResponse))
				assert.Equal(t, response.Items[0].UserId, mockResponse[0].UserID)
				assert.Equal(t, response.Items[0].OrderId, mockResponse[0].OrderID)
			},
		},
		{
			name: "BadInput",
			setupInputFunc: func(input *dto.ListPaymentQuery, total *int64) []models.Payment {
				var mockResponse []models.Payment
				input.Page = 0
				input.PerPage = 0
				userId := uint(1)
				input.UserId = &userId
				*total = 10
				return mockResponse
			},
			mockFunc: func(paymentRepo *mocks.MockPaymentRepository, mockResponse []models.Payment, input *dto.ListPaymentQuery, total *int64) {
				paymentRepo.On("ListPayments").Return(mockResponse, total, nil)
			},
			expectFunc: func(
				w *httptest.ResponseRecorder,
				mockResponse []models.Payment,
				input *dto.ListPaymentQuery,
				total *int64,
			) {
				assert.Equal(t, http.StatusBadRequest, w.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			paymentRepo := new(mocks.MockPaymentRepository)
			paymentService := services.NewPaymentService(paymentRepo)
			paymentHandler := NewPaymentHandler(paymentService)
			var input dto.ListPaymentQuery
			var total int64
			mockResponse := tc.setupInputFunc(&input, &total)
			tc.mockFunc(paymentRepo, mockResponse, &input, &total)
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			url := fmt.Sprintf("/payments?page=%d&per_page=%d&user_id=%d", input.Page, input.PerPage, *input.UserId)
			c.Request, _ = http.NewRequest(http.MethodGet, url, nil)

			// Act
			paymentHandler.ListPayments(c)

			// Assert
			tc.expectFunc(w, mockResponse, &input, &total)
		})
	}
}

// func TestUpdatePayment(t *testing.T) {
// 	testCases := []struct {
// 		name           string
// 		setupInputFunc func(input *dto.CreatePaymentDto, mockResponse *models.Payment)
// 		mockFunc       func(paymentRepo *mocks.MockPaymentRepository, mockResponse *models.Payment)
// 		expectFunc     func(w *httptest.ResponseRecorder, mockResponse *models.Payment)
// 	}{
// 		{
// 			name: "OK",
// 			setupInputFunc: func(input *dto.CreatePaymentDto, mockResponse *models.Payment) {
// 				input.OrderId = 1
// 				input.UserId = 1
// 				mockResponse.ID = 1
// 				mockResponse.CreatedAt = time.Now()
// 				mockResponse.UpdatedAt = mockResponse.CreatedAt
// 				mockResponse.FullName = input.FullName
// 				mockResponse.Paymentname = input.Paymentname

// 			},
// 			mockFunc: func(paymentRepo *mocks.MockPaymentRepository, mockResponse *models.Payment) {
// 				paymentRepo.On("UpdatePayment", mock.AnythingOfType("*models.Payment")).Return(nil).Run(func(args mock.Arguments) {
// 					arg := args.Get(0).(*models.Payment)
// 					arg.ID = mockResponse.ID
// 					arg.CreatedAt = mockResponse.CreatedAt
// 					arg.UpdatedAt = mockResponse.UpdatedAt
// 				})
// 			},
// 			expectFunc: func(w *httptest.ResponseRecorder, mockResponse *models.Payment) {
// 				assert.Equal(t, http.StatusCreated, w.Code)
// 				expectBodyPayment(t, w, mockResponse)
// 			},
// 		},
// 		{
// 			name: "BadInput",
// 			setupInputFunc: func(input *dto.CreatePaymentDto, mockResponse *models.Payment) {
// 			},
// 			mockFunc: func(paymentRepo *mocks.MockPaymentRepository, mockResponse *models.Payment) {
// 				paymentRepo.On("UpdatePayment", mock.AnythingOfType("*models.Payment")).Return(nil).Run(func(args mock.Arguments) {
// 					arg := args.Get(0).(*models.Payment)
// 					arg.ID = mockResponse.ID
// 					arg.CreatedAt = mockResponse.CreatedAt
// 					arg.UpdatedAt = mockResponse.UpdatedAt
// 				})
// 			},
// 			expectFunc: func(w *httptest.ResponseRecorder, mockResponse *models.Payment) {
// 				assert.Equal(t, http.StatusBadRequest, w.Code)
// 			},
// 		},
// 		{
// 			name: "UpdatePaymentError",
// 			setupInputFunc: func(input *dto.CreatePaymentDto, mockResponse *models.Payment) {
// 				input.FullName = "Full name"
// 				input.Paymentname = "paymentname"
// 				mockResponse.ID = 1
// 				mockResponse.CreatedAt = time.Now()
// 				mockResponse.UpdatedAt = mockResponse.CreatedAt
// 			},
// 			mockFunc: func(paymentRepo *mocks.MockPaymentRepository, mockResponse *models.Payment) {
// 				err := errors.New("Error")
// 				paymentRepo.On("UpdatePayment", mock.AnythingOfType("*models.Payment")).Return(err)
// 			},
// 			expectFunc: func(w *httptest.ResponseRecorder, mockResponse *models.Payment) {
// 				assert.Equal(t, http.StatusInternalServerError, w.Code)
// 			},
// 		},
// 	}

// 	for i := range testCases {
// 		tc := testCases[i]
// 		t.Run(tc.name, func(t *testing.T) {
// 			// Arrange
// 			paymentRepo := new(mocks.MockPaymentRepository)
// 			paymentPointRepo := new(mocks.MockPaymentPointRepository)
// 			paymentPointService := services.NewPaymentPointService(paymentPointRepo)
// 			paymentService := services.NewPaymentService(paymentRepo, paymentPointService)
// 			paymentHandler := NewPaymentHandler(paymentService)
// 			var payment dto.CreatePaymentDto
// 			var mockResponse models.Payment
// 			tc.setupInputFunc(&payment, &mockResponse)
// 			tc.mockFunc(paymentRepo, &mockResponse)
// 			gin.SetMode(gin.TestMode)
// 			w := httptest.NewRecorder()
// 			c, _ := gin.CreateTestContext(w)

// 			jsonPayment, _ := json.Marshal(payment)
// 			c.Request, _ = http.NewRequest(http.MethodPost, "/payments", bytes.NewBuffer(jsonPayment))
// 			c.Request.Header.Set("Content-Type", "application/json")
// 			c.Params = gin.Params{{Key: "id", Value: fmt.Sprint(mockResponse.ID)}}

// 			paymentHandler.UpdatePayment(c)

// 			tc.expectFunc(w, &mockResponse)
// 		})
// 	}
// }
