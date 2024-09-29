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
	"github.com/tricong1998/go-ecom/cmd/product/internal/mocks"
	"github.com/tricong1998/go-ecom/cmd/product/internal/services"
	"github.com/tricong1998/go-ecom/cmd/product/pkg/dto"
	"github.com/tricong1998/go-ecom/cmd/product/pkg/models"
)

func expectBodyProduct(t *testing.T, w *httptest.ResponseRecorder, mockResponse *models.Product) {
	var response dto.ProductResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, response.ID, mockResponse.ID)
	assert.Equal(t, response.Name, mockResponse.Name)
	assert.Equal(t, response.Price, mockResponse.Price)
	assert.WithinDuration(t, response.CreatedAt, mockResponse.CreatedAt, time.Second)
	assert.WithinDuration(t, response.UpdatedAt, mockResponse.UpdatedAt, time.Second)
}

func TestCreateProduct(t *testing.T) {
	testCases := []struct {
		name           string
		setupInputFunc func(input *dto.CreateProductDto, mockResponse *models.Product)
		mockFunc       func(userRepo *mocks.MockProductRepository, mockResponse *models.Product)
		expectFunc     func(w *httptest.ResponseRecorder, mockResponse *models.Product)
	}{
		{
			name: "OK",
			setupInputFunc: func(input *dto.CreateProductDto, mockResponse *models.Product) {
				input.Name = "Full name"
				input.Price = 1
				input.Quantity = 1
				mockResponse.ID = 1
				mockResponse.CreatedAt = time.Now()
				mockResponse.UpdatedAt = mockResponse.CreatedAt
				mockResponse.Name = input.Name
				mockResponse.Price = input.Price

			},
			mockFunc: func(userRepo *mocks.MockProductRepository, mockResponse *models.Product) {
				userRepo.On("CreateProduct", mock.AnythingOfType("*models.Product")).Return(nil).Run(func(args mock.Arguments) {
					arg := args.Get(0).(*models.Product)
					arg.ID = mockResponse.ID
					arg.CreatedAt = mockResponse.CreatedAt
					arg.UpdatedAt = mockResponse.UpdatedAt
				})
			},
			expectFunc: func(w *httptest.ResponseRecorder, mockResponse *models.Product) {
				assert.Equal(t, http.StatusCreated, w.Code)
				expectBodyProduct(t, w, mockResponse)
			},
		},
		{
			name: "BadInput",
			setupInputFunc: func(input *dto.CreateProductDto, mockResponse *models.Product) {
			},
			mockFunc: func(userRepo *mocks.MockProductRepository, mockResponse *models.Product) {
				userRepo.On("CreateProduct", mock.AnythingOfType("*models.Product")).Return(nil).Run(func(args mock.Arguments) {
					arg := args.Get(0).(*models.Product)
					arg.ID = mockResponse.ID
					arg.CreatedAt = mockResponse.CreatedAt
					arg.UpdatedAt = mockResponse.UpdatedAt
				})
			},
			expectFunc: func(w *httptest.ResponseRecorder, mockResponse *models.Product) {
				assert.Equal(t, http.StatusBadRequest, w.Code)
			},
		},
		{
			name: "CreateProductError",
			setupInputFunc: func(input *dto.CreateProductDto, mockResponse *models.Product) {
				input.Name = "Full name"
				input.Price = 1
				input.Quantity = 1
				mockResponse.ID = 1
				mockResponse.CreatedAt = time.Now()
				mockResponse.UpdatedAt = mockResponse.CreatedAt
			},
			mockFunc: func(userRepo *mocks.MockProductRepository, mockResponse *models.Product) {
				err := errors.New("Error")
				userRepo.On("CreateProduct", mock.AnythingOfType("*models.Product")).Return(err)
			},
			expectFunc: func(w *httptest.ResponseRecorder, mockResponse *models.Product) {
				assert.Equal(t, http.StatusInternalServerError, w.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			userRepo := new(mocks.MockProductRepository)
			userService := services.NewProductService(userRepo)
			userHandler := NewProductHandler(userService)
			var user dto.CreateProductDto
			var mockResponse models.Product
			tc.setupInputFunc(&user, &mockResponse)
			tc.mockFunc(userRepo, &mockResponse)
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			jsonProduct, _ := json.Marshal(user)
			c.Request, _ = http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(jsonProduct))
			c.Request.Header.Set("Content-Type", "application/json")
			userHandler.CreateProduct(c)

			tc.expectFunc(w, &mockResponse)
		})
	}
}

func TestReadProduct(t *testing.T) {
	testCases := []struct {
		name           string
		setupInputFunc func(input *dto.ReadProductRequest, mockResponse *models.Product)
		mockFunc       func(userRepo *mocks.MockProductRepository, mockResponse *models.Product)
		expectFunc     func(w *httptest.ResponseRecorder, mockResponse *models.Product)
	}{
		{
			name: "OK",
			setupInputFunc: func(input *dto.ReadProductRequest, mockResponse *models.Product) {
				input.ID = 1
				mockResponse.Name = "Full name"
				mockResponse.Price = 1
				mockResponse.ID = input.ID
				mockResponse.CreatedAt = time.Now()
				mockResponse.UpdatedAt = mockResponse.CreatedAt
			},
			mockFunc: func(userRepo *mocks.MockProductRepository, mockResponse *models.Product) {
				userRepo.On("ReadProduct", mock.AnythingOfType("uint")).Return(mockResponse, nil)
			},
			expectFunc: func(w *httptest.ResponseRecorder, mockResponse *models.Product) {
				assert.Equal(t, http.StatusOK, w.Code)
				expectBodyProduct(t, w, mockResponse)
			},
		},
		{
			name: "BadInput",
			setupInputFunc: func(input *dto.ReadProductRequest, mockResponse *models.Product) {
			},
			mockFunc: func(userRepo *mocks.MockProductRepository, mockResponse *models.Product) {
				userRepo.On("ReadProduct", mock.AnythingOfType("uint")).Return(mockResponse, nil)
			},
			expectFunc: func(w *httptest.ResponseRecorder, mockResponse *models.Product) {
				assert.Equal(t, http.StatusBadRequest, w.Code)
			},
		},
		{
			name: "ProductNotFound",
			setupInputFunc: func(input *dto.ReadProductRequest, mockResponse *models.Product) {
				input.ID = 1
			},
			mockFunc: func(userRepo *mocks.MockProductRepository, mockResponse *models.Product) {
				userRepo.On("ReadProduct", mock.AnythingOfType("uint")).Return(mockResponse, errors.New("Not found"))
			},
			expectFunc: func(w *httptest.ResponseRecorder, mockResponse *models.Product) {
				assert.Equal(t, http.StatusNotFound, w.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			userRepo := new(mocks.MockProductRepository)
			userService := services.NewProductService(userRepo)
			userHandler := NewProductHandler(userService)
			var input dto.ReadProductRequest
			var mockResponse models.Product
			tc.setupInputFunc(&input, &mockResponse)
			tc.mockFunc(userRepo, &mockResponse)
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Params = gin.Params{{Key: "id", Value: fmt.Sprint(input.ID)}}

			// Act
			userHandler.ReadProduct(c)

			// Assert
			tc.expectFunc(w, &mockResponse)
		})
	}
}

func TestListProduct(t *testing.T) {
	testCases := []struct {
		name           string
		setupInputFunc func(input *dto.ListProductQuery, total *int64) []models.Product
		mockFunc       func(userRepo *mocks.MockProductRepository, mockResponse []models.Product, input *dto.ListProductQuery, total *int64)
		expectFunc     func(
			w *httptest.ResponseRecorder,
			mockResponse []models.Product,
			input *dto.ListProductQuery,
			total *int64,
		)
	}{
		{
			name: "OK",
			setupInputFunc: func(input *dto.ListProductQuery, total *int64) []models.Product {
				var mockResponse []models.Product
				input.Page = int32(1)
				input.PerPage = int32(10)
				username := "username"
				input.Name = &username
				*total = 10
				now := time.Now()
				user1 := models.Product{
					Name:  "Full name 1",
					Price: 1,
				}
				user1.CreatedAt = now
				user1.UpdatedAt = now
				user1.ID = 1
				mockResponse = append(mockResponse, user1)

				user2 := models.Product{
					Name:  "Full name 2",
					Price: 1,
				}
				user2.CreatedAt = now
				user2.UpdatedAt = now
				user2.ID = 2
				mockResponse = append(mockResponse, user2)
				return mockResponse
			},
			mockFunc: func(userRepo *mocks.MockProductRepository, mockResponse []models.Product, input *dto.ListProductQuery, total *int64) {
				userRepo.On("ListProducts", input.PerPage, input.Page, input.Name).Return(mockResponse, *total, nil)
			},
			expectFunc: func(
				w *httptest.ResponseRecorder,
				mockResponse []models.Product,
				input *dto.ListProductQuery,
				total *int64,
			) {
				assert.Equal(t, http.StatusOK, w.Code)
				var response dto.ListProductResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, response.Metadata.Page, input.Page)
				assert.Equal(t, response.Metadata.PerPage, input.PerPage)
				assert.Equal(t, response.Metadata.Total, *total)
				assert.Len(t, response.Items, len(mockResponse))
				assert.Equal(t, response.Items[0].Price, mockResponse[0].Price)
				assert.Equal(t, response.Items[0].Name, mockResponse[0].Name)
			},
		},
		{
			name: "BadInput",
			setupInputFunc: func(input *dto.ListProductQuery, total *int64) []models.Product {
				var mockResponse []models.Product
				input.Page = 0
				input.PerPage = 0
				username := "username"
				input.Name = &username
				*total = 10
				return mockResponse
			},
			mockFunc: func(userRepo *mocks.MockProductRepository, mockResponse []models.Product, input *dto.ListProductQuery, total *int64) {
				userRepo.On("ListProducts").Return(mockResponse, total, nil)
			},
			expectFunc: func(
				w *httptest.ResponseRecorder,
				mockResponse []models.Product,
				input *dto.ListProductQuery,
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
			userRepo := new(mocks.MockProductRepository)
			userService := services.NewProductService(userRepo)
			userHandler := NewProductHandler(userService)
			var input dto.ListProductQuery
			var total int64
			mockResponse := tc.setupInputFunc(&input, &total)
			tc.mockFunc(userRepo, mockResponse, &input, &total)
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			url := fmt.Sprintf("/orders?page=%d&per_page=%d&username=%s", input.Page, input.PerPage, *input.Name)
			c.Request, _ = http.NewRequest(http.MethodGet, url, nil)

			// Act
			userHandler.ListProducts(c)

			// Assert
			tc.expectFunc(w, mockResponse, &input, &total)
		})
	}
}

func TestUpdateProduct(t *testing.T) {
	testCases := []struct {
		name           string
		setupInputFunc func(input *dto.CreateProductDto, mockResponse *models.Product)
		mockFunc       func(userRepo *mocks.MockProductRepository, mockResponse *models.Product)
		expectFunc     func(w *httptest.ResponseRecorder, mockResponse *models.Product)
	}{
		{
			name: "OK",
			setupInputFunc: func(input *dto.CreateProductDto, mockResponse *models.Product) {
				input.Name = "New full name"
				input.Price = 1
				mockResponse.ID = 1
				mockResponse.CreatedAt = time.Now()
				mockResponse.UpdatedAt = mockResponse.CreatedAt
				mockResponse.Name = input.Name
				mockResponse.Price = input.Price

			},
			mockFunc: func(userRepo *mocks.MockProductRepository, mockResponse *models.Product) {
				userRepo.On("UpdateProduct", mock.AnythingOfType("*models.Product")).Return(nil).Run(func(args mock.Arguments) {
					arg := args.Get(0).(*models.Product)
					arg.ID = mockResponse.ID
					arg.CreatedAt = mockResponse.CreatedAt
					arg.UpdatedAt = mockResponse.UpdatedAt
				})
			},
			expectFunc: func(w *httptest.ResponseRecorder, mockResponse *models.Product) {
				assert.Equal(t, http.StatusCreated, w.Code)
				expectBodyProduct(t, w, mockResponse)
			},
		},
		{
			name: "BadInput",
			setupInputFunc: func(input *dto.CreateProductDto, mockResponse *models.Product) {
			},
			mockFunc: func(userRepo *mocks.MockProductRepository, mockResponse *models.Product) {
				userRepo.On("UpdateProduct", mock.AnythingOfType("*models.Product")).Return(nil).Run(func(args mock.Arguments) {
					arg := args.Get(0).(*models.Product)
					arg.ID = mockResponse.ID
					arg.CreatedAt = mockResponse.CreatedAt
					arg.UpdatedAt = mockResponse.UpdatedAt
				})
			},
			expectFunc: func(w *httptest.ResponseRecorder, mockResponse *models.Product) {
				assert.Equal(t, http.StatusBadRequest, w.Code)
			},
		},
		{
			name: "UpdateProductError",
			setupInputFunc: func(input *dto.CreateProductDto, mockResponse *models.Product) {
				input.Name = "Full name"
				input.Price = 1
				mockResponse.ID = 1
				mockResponse.CreatedAt = time.Now()
				mockResponse.UpdatedAt = mockResponse.CreatedAt
			},
			mockFunc: func(userRepo *mocks.MockProductRepository, mockResponse *models.Product) {
				err := errors.New("Error")
				userRepo.On("UpdateProduct", mock.AnythingOfType("*models.Product")).Return(err)
			},
			expectFunc: func(w *httptest.ResponseRecorder, mockResponse *models.Product) {
				assert.Equal(t, http.StatusInternalServerError, w.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			userRepo := new(mocks.MockProductRepository)
			userService := services.NewProductService(userRepo)
			userHandler := NewProductHandler(userService)
			var user dto.CreateProductDto
			var mockResponse models.Product
			tc.setupInputFunc(&user, &mockResponse)
			tc.mockFunc(userRepo, &mockResponse)
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			jsonProduct, _ := json.Marshal(user)
			c.Request, _ = http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(jsonProduct))
			c.Request.Header.Set("Content-Type", "application/json")
			c.Params = gin.Params{{Key: "id", Value: fmt.Sprint(mockResponse.ID)}}

			userHandler.UpdateProduct(c)

			tc.expectFunc(w, &mockResponse)
		})
	}
}
