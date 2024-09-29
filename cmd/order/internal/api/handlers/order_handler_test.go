package handlers

import (
	"bytes"
	"context"
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
	"github.com/tricong1998/go-ecom/cmd/order/internal/api/dto"
	"github.com/tricong1998/go-ecom/cmd/order/internal/mocks"
	"github.com/tricong1998/go-ecom/cmd/order/internal/models"
	"github.com/tricong1998/go-ecom/cmd/order/internal/services"
	paymentPb "github.com/tricong1998/go-ecom/cmd/payment/pkg/pb"
	productPb "github.com/tricong1998/go-ecom/cmd/product/pkg/pb"
	userPb "github.com/tricong1998/go-ecom/cmd/user/pkg/pb"
)

func expectBodyOrder(t *testing.T, w *httptest.ResponseRecorder, mockResponse *models.Order) {
	var response dto.OrderResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, response.ID, mockResponse.ID)
	assert.Equal(t, response.ProductId, mockResponse.ProductId)
	assert.Equal(t, response.UserId, mockResponse.UserId)
	assert.WithinDuration(t, response.CreatedAt, mockResponse.CreatedAt, time.Second)
	assert.WithinDuration(t, response.UpdatedAt, mockResponse.UpdatedAt, time.Second)
}

func TestCreateOrder(t *testing.T) {
	testCases := []struct {
		name           string
		setupInputFunc func(
			input *dto.CreateOrderDto,
			mockResponse *models.Order,
			userMock *userPb.User,
			mockProduct *productPb.ReadProductResponse,
			mockPayment *paymentPb.CreatePaymentResponse,
		)
		mockFunc func(
			userRepo *mocks.MockOrderRepository,
			mockResponse *models.Order,
			userGateway *mocks.MockUserGateway,
			productGateway *mocks.MockProductGateway,
			paymentGateway *mocks.MockPaymentGateway,
			mockUser *userPb.User,
			mockProduct *productPb.ReadProductResponse,
			mockPayment *paymentPb.CreatePaymentResponse,
			publisher *mocks.MockRabbitPublisher,
		)
		expectFunc func(
			w *httptest.ResponseRecorder, mockResponse *models.Order)
	}{
		{
			name: "OK",
			setupInputFunc: func(input *dto.CreateOrderDto,
				mockResponse *models.Order,
				userMock *userPb.User,
				mockProduct *productPb.ReadProductResponse,
				mockPayment *paymentPb.CreatePaymentResponse,
			) {
				input.UserId = 1
				input.ProductId = 1
				input.ProductCount = 1
				product := productPb.Product{
					Id:       uint64(input.ProductId),
					Name:     "product name",
					Price:    100,
					Quantity: 10,
				}
				mockProduct.Product = &product
				mockResponse.ID = 1
				mockResponse.CreatedAt = time.Now()
				mockResponse.UpdatedAt = mockResponse.CreatedAt
				mockResponse.UserId = input.UserId
				mockResponse.ProductId = input.ProductId
				mockResponse.ProductCount = input.ProductCount
				mockResponse.Amount = uint(mockProduct.Product.Price) * uint(input.ProductCount)
				payment := paymentPb.Payment{
					Id:     uint64(1),
					Amount: uint64(mockResponse.Amount),
				}
				mockPayment.Payment = &payment
				userMock.Id = uint64(input.UserId)
				userMock.Username = "test"
				userMock.FullName = "user full name"
			},
			mockFunc: func(
				userRepo *mocks.MockOrderRepository,
				mockResponse *models.Order,
				userGateway *mocks.MockUserGateway,
				productGateway *mocks.MockProductGateway,
				paymentGateway *mocks.MockPaymentGateway,
				mockUser *userPb.User,
				mockProduct *productPb.ReadProductResponse,
				mockPayment *paymentPb.CreatePaymentResponse,
				publisher *mocks.MockRabbitPublisher,
			) {
				userRepo.On("CreateOrder", mock.AnythingOfType("*models.Order")).Return(nil).Run(func(args mock.Arguments) {
					arg := args.Get(0).(*models.Order)
					arg.ID = mockResponse.ID
					arg.CreatedAt = mockResponse.CreatedAt
					arg.UpdatedAt = mockResponse.UpdatedAt
				})
				userRepo.On("UpdateOrderStatus", mock.AnythingOfType("uint"), mock.AnythingOfType("string")).Return(nil)
				userGateway.On("Get",
					context.Background(),
					mock.AnythingOfType("uint")).Return(mockUser, nil)
				productGateway.On("Get",
					context.Background(),
					mock.AnythingOfType("uint")).Return(mockProduct, nil)
				productGateway.On("UpdateProductQuantity",
					context.Background(),
					mock.AnythingOfType("uint"),
					mock.AnythingOfType("uint")).Return(true, nil)
				paymentGateway.On("Create",
					context.Background(),
					mock.AnythingOfType("*pb.CreatePaymentRequest")).
					Return(mockPayment, nil)
				publisher.On("PublishMessage", mock.AnythingOfType("dto.CreateUserPoint")).Return(nil)
			},
			expectFunc: func(w *httptest.ResponseRecorder, mockResponse *models.Order) {
				assert.Equal(t, http.StatusCreated, w.Code)
				expectBodyOrder(t, w, mockResponse)
			},
		},
		{
			name: "BadInput",
			setupInputFunc: func(input *dto.CreateOrderDto,
				mockResponse *models.Order,
				userMock *userPb.User,
				mockProduct *productPb.ReadProductResponse,
				mockPayment *paymentPb.CreatePaymentResponse,
			) {
			},
			mockFunc: func(
				userRepo *mocks.MockOrderRepository,
				mockResponse *models.Order,
				userGateway *mocks.MockUserGateway,
				productGateway *mocks.MockProductGateway,
				paymentGateway *mocks.MockPaymentGateway,
				mockUser *userPb.User,
				mockProduct *productPb.ReadProductResponse,
				mockPayment *paymentPb.CreatePaymentResponse,
				publisher *mocks.MockRabbitPublisher,
			) {
				userRepo.On("CreateOrder", mock.AnythingOfType("*models.Order")).Return(nil).Run(func(args mock.Arguments) {
					arg := args.Get(0).(*models.Order)
					arg.ID = mockResponse.ID
					arg.CreatedAt = mockResponse.CreatedAt
					arg.UpdatedAt = mockResponse.UpdatedAt
				})
				userRepo.On("UpdateOrderStatus", mock.AnythingOfType("uint"), mock.AnythingOfType("string")).Return(nil)
				userGateway.On("Get",
					context.Background(),
					mock.AnythingOfType("uint")).Return(mockUser, nil)
				productGateway.On("Get",
					context.Background(),
					mock.AnythingOfType("uint")).Return(mockProduct, nil)
				productGateway.On("UpdateProductQuantity",
					context.Background(),
					mock.AnythingOfType("uint"),
					mock.AnythingOfType("uint")).Return(true, nil)
				paymentGateway.On("Create",
					context.Background(),
					mock.AnythingOfType("*pb.CreatePaymentRequest")).
					Return(mockPayment, nil)
				publisher.On("PublishMessage", mock.AnythingOfType("dto.CreateUserPoint")).Return(nil)
			},
			expectFunc: func(w *httptest.ResponseRecorder, mockResponse *models.Order) {
				assert.Equal(t, http.StatusBadRequest, w.Code)
			},
		},
		{
			name: "CreateOrderError",
			setupInputFunc: func(
				input *dto.CreateOrderDto,
				mockResponse *models.Order,
				userMock *userPb.User,
				mockProduct *productPb.ReadProductResponse,
				mockPayment *paymentPb.CreatePaymentResponse,
			) {
				input.UserId = 1
				input.ProductId = 1
				input.ProductCount = 1
				product := productPb.Product{
					Id:       uint64(input.ProductId),
					Name:     "product name",
					Price:    100,
					Quantity: 10,
				}
				mockProduct.Product = &product
				mockResponse.ID = 1
				mockResponse.CreatedAt = time.Now()
				mockResponse.UpdatedAt = mockResponse.CreatedAt
				mockResponse.UserId = input.UserId
				mockResponse.ProductId = input.ProductId
				mockResponse.ProductCount = input.ProductCount
				mockResponse.Amount = uint(mockProduct.Product.Price) * uint(input.ProductCount)
				payment := paymentPb.Payment{
					Id:     uint64(1),
					Amount: uint64(mockResponse.Amount),
				}
				mockPayment.Payment = &payment
				userMock.Id = uint64(input.UserId)
				userMock.Username = "test"
				userMock.FullName = "user full name"
			},
			mockFunc: func(
				userRepo *mocks.MockOrderRepository,
				mockResponse *models.Order,
				userGateway *mocks.MockUserGateway,
				productGateway *mocks.MockProductGateway,
				paymentGateway *mocks.MockPaymentGateway,
				mockUser *userPb.User,
				mockProduct *productPb.ReadProductResponse,
				mockPayment *paymentPb.CreatePaymentResponse,
				publisher *mocks.MockRabbitPublisher,
			) {
				err := errors.New("Error")
				userRepo.On("CreateOrder", mock.AnythingOfType("*models.Order")).Return(err)
				userRepo.On("UpdateOrderStatus", mock.AnythingOfType("uint"), mock.AnythingOfType("string")).Return(nil)
				userGateway.On("Get",
					context.Background(),
					mock.AnythingOfType("uint")).Return(mockUser, nil)
				productGateway.On("Get",
					context.Background(),
					mock.AnythingOfType("uint")).Return(mockProduct, nil)
				productGateway.On("UpdateProductQuantity",
					context.Background(),
					mock.AnythingOfType("uint"),
					mock.AnythingOfType("uint")).Return(true, nil)
				paymentGateway.On("Create",
					context.Background(),
					mock.AnythingOfType("*pb.CreatePaymentRequest")).
					Return(mockPayment, nil)
				publisher.On("PublishMessage", mock.AnythingOfType("dto.CreateUserPoint")).Return(nil)
			},
			expectFunc: func(w *httptest.ResponseRecorder, mockResponse *models.Order) {
				assert.Equal(t, http.StatusInternalServerError, w.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			userRepo := new(mocks.MockOrderRepository)
			publisher := new(mocks.MockRabbitPublisher)
			userGateway := new(mocks.MockUserGateway)
			productGateway := new(mocks.MockProductGateway)
			paymentGateway := new(mocks.MockPaymentGateway)
			userService := services.NewOrderService(userRepo, userGateway, paymentGateway, publisher, productGateway)
			userHandler := NewOrderHandler(userService)
			var user dto.CreateOrderDto
			var mockResponse models.Order
			var userMock userPb.User
			var productMock productPb.ReadProductResponse
			var paymentMock paymentPb.CreatePaymentResponse
			tc.setupInputFunc(&user, &mockResponse, &userMock, &productMock, &paymentMock)
			tc.mockFunc(userRepo, &mockResponse, userGateway, productGateway, paymentGateway, &userMock, &productMock, &paymentMock, publisher)
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			jsonOrder, _ := json.Marshal(user)
			c.Request, _ = http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(jsonOrder))
			c.Request.Header.Set("Content-Type", "application/json")
			userHandler.CreateOrder(c)

			tc.expectFunc(w, &mockResponse)
		})
	}
}

func TestReadOrder(t *testing.T) {
	testCases := []struct {
		name           string
		setupInputFunc func(input *dto.ReadOrderRequest, mockResponse *models.Order)
		mockFunc       func(userRepo *mocks.MockOrderRepository, mockResponse *models.Order)
		expectFunc     func(w *httptest.ResponseRecorder, mockResponse *models.Order)
	}{
		{
			name: "OK",
			setupInputFunc: func(input *dto.ReadOrderRequest, mockResponse *models.Order) {
				input.ID = 1
				mockResponse.UserId = 1
				mockResponse.ProductId = 1
				mockResponse.ID = input.ID
				mockResponse.CreatedAt = time.Now()
				mockResponse.UpdatedAt = mockResponse.CreatedAt
			},
			mockFunc: func(userRepo *mocks.MockOrderRepository, mockResponse *models.Order) {
				userRepo.On("ReadOrder", mock.AnythingOfType("uint")).Return(mockResponse, nil)
			},
			expectFunc: func(w *httptest.ResponseRecorder, mockResponse *models.Order) {
				assert.Equal(t, http.StatusOK, w.Code)
				expectBodyOrder(t, w, mockResponse)
			},
		},
		{
			name: "BadInput",
			setupInputFunc: func(input *dto.ReadOrderRequest, mockResponse *models.Order) {
			},
			mockFunc: func(userRepo *mocks.MockOrderRepository, mockResponse *models.Order) {
				userRepo.On("ReadOrder", mock.AnythingOfType("uint")).Return(mockResponse, nil)
			},
			expectFunc: func(w *httptest.ResponseRecorder, mockResponse *models.Order) {
				assert.Equal(t, http.StatusBadRequest, w.Code)
			},
		},
		{
			name: "OrderNotFound",
			setupInputFunc: func(input *dto.ReadOrderRequest, mockResponse *models.Order) {
				input.ID = 1
			},
			mockFunc: func(userRepo *mocks.MockOrderRepository, mockResponse *models.Order) {
				userRepo.On("ReadOrder", mock.AnythingOfType("uint")).Return(mockResponse, errors.New("Not found"))
			},
			expectFunc: func(w *httptest.ResponseRecorder, mockResponse *models.Order) {
				assert.Equal(t, http.StatusNotFound, w.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			userRepo := new(mocks.MockOrderRepository)
			publisher := new(mocks.MockRabbitPublisher)
			userGateway := new(mocks.MockUserGateway)
			productGateway := new(mocks.MockProductGateway)
			paymentGateway := new(mocks.MockPaymentGateway)
			userService := services.NewOrderService(userRepo, userGateway, paymentGateway, publisher, productGateway)
			userHandler := NewOrderHandler(userService)
			var input dto.ReadOrderRequest
			var mockResponse models.Order
			tc.setupInputFunc(&input, &mockResponse)
			tc.mockFunc(userRepo, &mockResponse)
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Params = gin.Params{{Key: "id", Value: fmt.Sprint(input.ID)}}

			// Act
			userHandler.ReadOrder(c)

			// Assert
			tc.expectFunc(w, &mockResponse)
		})
	}
}

func TestListOrder(t *testing.T) {
	testCases := []struct {
		name           string
		setupInputFunc func(input *dto.ListOrderQuery, total *int64) []models.Order
		mockFunc       func(userRepo *mocks.MockOrderRepository, mockResponse []models.Order, input *dto.ListOrderQuery, total *int64)
		expectFunc     func(
			w *httptest.ResponseRecorder,
			mockResponse []models.Order,
			input *dto.ListOrderQuery,
			total *int64,
		)
	}{
		{
			name: "OK",
			setupInputFunc: func(input *dto.ListOrderQuery, total *int64) []models.Order {
				var mockResponse []models.Order
				input.Page = int32(1)
				input.PerPage = int32(10)
				input.UserId = 1
				*total = 10
				now := time.Now()
				user1 := models.Order{
					ProductId: 1,
					UserId:    1,
				}
				user1.CreatedAt = now
				user1.UpdatedAt = now
				user1.ID = 1
				mockResponse = append(mockResponse, user1)

				user2 := models.Order{
					ProductId: 2,
					UserId:    2,
				}
				user2.CreatedAt = now
				user2.UpdatedAt = now
				user2.ID = 2
				mockResponse = append(mockResponse, user2)
				return mockResponse
			},
			mockFunc: func(userRepo *mocks.MockOrderRepository, mockResponse []models.Order, input *dto.ListOrderQuery, total *int64) {
				userRepo.On("ListOrders", input.PerPage, input.Page, input.UserId).Return(mockResponse, *total, nil)
			},
			expectFunc: func(
				w *httptest.ResponseRecorder,
				mockResponse []models.Order,
				input *dto.ListOrderQuery,
				total *int64,
			) {
				assert.Equal(t, http.StatusOK, w.Code)
				var response dto.ListOrderResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, response.Metadata.Page, input.Page)
				assert.Equal(t, response.Metadata.PerPage, input.PerPage)
				assert.Equal(t, response.Metadata.Total, *total)
				assert.Len(t, response.Items, len(mockResponse))
				assert.Equal(t, response.Items[0].ProductId, mockResponse[0].ProductId)
				assert.Equal(t, response.Items[0].UserId, mockResponse[0].UserId)
			},
		},
		{
			name: "BadInput",
			setupInputFunc: func(input *dto.ListOrderQuery, total *int64) []models.Order {
				var mockResponse []models.Order
				input.Page = 0
				input.PerPage = 0
				input.UserId = 1
				*total = 10
				return mockResponse
			},
			mockFunc: func(userRepo *mocks.MockOrderRepository, mockResponse []models.Order, input *dto.ListOrderQuery, total *int64) {
				userRepo.On("ListOrders").Return(mockResponse, total, nil)
			},
			expectFunc: func(
				w *httptest.ResponseRecorder,
				mockResponse []models.Order,
				input *dto.ListOrderQuery,
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
			userRepo := new(mocks.MockOrderRepository)
			publisher := new(mocks.MockRabbitPublisher)
			userGateway := new(mocks.MockUserGateway)
			productGateway := new(mocks.MockProductGateway)
			paymentGateway := new(mocks.MockPaymentGateway)
			userService := services.NewOrderService(userRepo, userGateway, paymentGateway, publisher, productGateway)
			userHandler := NewOrderHandler(userService)
			var input dto.ListOrderQuery
			var total int64
			mockResponse := tc.setupInputFunc(&input, &total)
			tc.mockFunc(userRepo, mockResponse, &input, &total)
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			url := fmt.Sprintf("/orders?page=%d&per_page=%d&username=%d", input.Page, input.PerPage, input.UserId)
			c.Request, _ = http.NewRequest(http.MethodGet, url, nil)

			// Act
			userHandler.ListOrders(c)

			// Assert
			tc.expectFunc(w, mockResponse, &input, &total)
		})
	}
}

func TestUpdateOrder(t *testing.T) {
	testCases := []struct {
		name           string
		setupInputFunc func(input *dto.CreateOrderDto, mockResponse *models.Order)
		mockFunc       func(userRepo *mocks.MockOrderRepository, mockResponse *models.Order)
		expectFunc     func(w *httptest.ResponseRecorder, mockResponse *models.Order)
	}{
		{
			name: "OK",
			setupInputFunc: func(input *dto.CreateOrderDto, mockResponse *models.Order) {
				input.ProductId = 1
				input.UserId = 1
				mockResponse.ID = 1
				mockResponse.CreatedAt = time.Now()
				mockResponse.UpdatedAt = mockResponse.CreatedAt
				mockResponse.ProductId = input.ProductId
				mockResponse.UserId = input.UserId

			},
			mockFunc: func(userRepo *mocks.MockOrderRepository, mockResponse *models.Order) {
				userRepo.On("UpdateOrder", mock.AnythingOfType("*models.Order")).Return(nil).Run(func(args mock.Arguments) {
					arg := args.Get(0).(*models.Order)
					arg.ID = mockResponse.ID
					arg.CreatedAt = mockResponse.CreatedAt
					arg.UpdatedAt = mockResponse.UpdatedAt
				})
			},
			expectFunc: func(w *httptest.ResponseRecorder, mockResponse *models.Order) {
				assert.Equal(t, http.StatusCreated, w.Code)
				expectBodyOrder(t, w, mockResponse)
			},
		},
		{
			name: "BadInput",
			setupInputFunc: func(input *dto.CreateOrderDto, mockResponse *models.Order) {
			},
			mockFunc: func(userRepo *mocks.MockOrderRepository, mockResponse *models.Order) {
				userRepo.On("UpdateOrder", mock.AnythingOfType("*models.Order")).Return(nil).Run(func(args mock.Arguments) {
					arg := args.Get(0).(*models.Order)
					arg.ID = mockResponse.ID
					arg.CreatedAt = mockResponse.CreatedAt
					arg.UpdatedAt = mockResponse.UpdatedAt
				})
			},
			expectFunc: func(w *httptest.ResponseRecorder, mockResponse *models.Order) {
				assert.Equal(t, http.StatusBadRequest, w.Code)
			},
		},
		{
			name: "UpdateOrderError",
			setupInputFunc: func(input *dto.CreateOrderDto, mockResponse *models.Order) {
				input.UserId = 1
				input.ProductId = 1
				mockResponse.ID = 1
				mockResponse.CreatedAt = time.Now()
				mockResponse.UpdatedAt = mockResponse.CreatedAt
			},
			mockFunc: func(userRepo *mocks.MockOrderRepository, mockResponse *models.Order) {
				err := errors.New("Error")
				userRepo.On("UpdateOrder", mock.AnythingOfType("*models.Order")).Return(err)
			},
			expectFunc: func(w *httptest.ResponseRecorder, mockResponse *models.Order) {
				assert.Equal(t, http.StatusInternalServerError, w.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			userRepo := new(mocks.MockOrderRepository)
			publisher := new(mocks.MockRabbitPublisher)
			userGateway := new(mocks.MockUserGateway)
			productGateway := new(mocks.MockProductGateway)
			paymentGateway := new(mocks.MockPaymentGateway)
			userService := services.NewOrderService(userRepo, userGateway, paymentGateway, publisher, productGateway)
			userHandler := NewOrderHandler(userService)
			var user dto.CreateOrderDto
			var mockResponse models.Order
			tc.setupInputFunc(&user, &mockResponse)
			tc.mockFunc(userRepo, &mockResponse)
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			jsonOrder, _ := json.Marshal(user)
			c.Request, _ = http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(jsonOrder))
			c.Request.Header.Set("Content-Type", "application/json")
			c.Params = gin.Params{{Key: "id", Value: fmt.Sprint(mockResponse.ID)}}

			userHandler.UpdateOrder(c)

			tc.expectFunc(w, &mockResponse)
		})
	}
}
