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
	"github.com/tricong1998/go-ecom/cmd/user/internal/mocks"
	"github.com/tricong1998/go-ecom/cmd/user/internal/services"
	"github.com/tricong1998/go-ecom/cmd/user/pkg/dto"
	"github.com/tricong1998/go-ecom/cmd/user/pkg/models"
)

func expectBodyUser(t *testing.T, w *httptest.ResponseRecorder, mockResponse *models.User) {
	var response dto.UserResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, response.ID, mockResponse.ID)
	assert.Equal(t, response.Username, mockResponse.Username)
	assert.Equal(t, response.FullName, mockResponse.FullName)
	assert.Equal(t, response.FullName, mockResponse.FullName)
	assert.WithinDuration(t, response.CreatedAt, mockResponse.CreatedAt, time.Second)
	assert.WithinDuration(t, response.UpdatedAt, mockResponse.UpdatedAt, time.Second)
}

func TestCreateUser(t *testing.T) {
	testCases := []struct {
		name           string
		setupInputFunc func(input *dto.CreateUserDto, mockResponse *models.User)
		mockFunc       func(userRepo *mocks.MockUserRepository, mockResponse *models.User)
		expectFunc     func(w *httptest.ResponseRecorder, mockResponse *models.User)
	}{
		{
			name: "OK",
			setupInputFunc: func(input *dto.CreateUserDto, mockResponse *models.User) {
				input.FullName = "Full name"
				input.Username = "username"
				mockResponse.ID = 1
				mockResponse.CreatedAt = time.Now()
				mockResponse.UpdatedAt = mockResponse.CreatedAt
				mockResponse.FullName = input.FullName
				mockResponse.Username = input.Username

			},
			mockFunc: func(userRepo *mocks.MockUserRepository, mockResponse *models.User) {
				userRepo.On("CreateUser", mock.AnythingOfType("*models.User")).Return(nil).Run(func(args mock.Arguments) {
					arg := args.Get(0).(*models.User)
					arg.ID = mockResponse.ID
					arg.CreatedAt = mockResponse.CreatedAt
					arg.UpdatedAt = mockResponse.UpdatedAt
				})
			},
			expectFunc: func(w *httptest.ResponseRecorder, mockResponse *models.User) {
				assert.Equal(t, http.StatusCreated, w.Code)
				expectBodyUser(t, w, mockResponse)
			},
		},
		{
			name: "BadInput",
			setupInputFunc: func(input *dto.CreateUserDto, mockResponse *models.User) {
			},
			mockFunc: func(userRepo *mocks.MockUserRepository, mockResponse *models.User) {
				userRepo.On("CreateUser", mock.AnythingOfType("*models.User")).Return(nil).Run(func(args mock.Arguments) {
					arg := args.Get(0).(*models.User)
					arg.ID = mockResponse.ID
					arg.CreatedAt = mockResponse.CreatedAt
					arg.UpdatedAt = mockResponse.UpdatedAt
				})
			},
			expectFunc: func(w *httptest.ResponseRecorder, mockResponse *models.User) {
				assert.Equal(t, http.StatusBadRequest, w.Code)
			},
		},
		{
			name: "CreateUserError",
			setupInputFunc: func(input *dto.CreateUserDto, mockResponse *models.User) {
				input.FullName = "Full name"
				input.Username = "username"
				mockResponse.ID = 1
				mockResponse.CreatedAt = time.Now()
				mockResponse.UpdatedAt = mockResponse.CreatedAt
			},
			mockFunc: func(userRepo *mocks.MockUserRepository, mockResponse *models.User) {
				err := errors.New("Error")
				userRepo.On("CreateUser", mock.AnythingOfType("*models.User")).Return(err)
			},
			expectFunc: func(w *httptest.ResponseRecorder, mockResponse *models.User) {
				assert.Equal(t, http.StatusInternalServerError, w.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			userRepo := new(mocks.MockUserRepository)
			userPointRepo := new(mocks.MockUserPointRepository)
			userPointService := services.NewUserPointService(userPointRepo)
			userService := services.NewUserService(userRepo, userPointService)
			userHandler := NewUserHandler(userService)
			var user dto.CreateUserDto
			var mockResponse models.User
			tc.setupInputFunc(&user, &mockResponse)
			tc.mockFunc(userRepo, &mockResponse)
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			jsonUser, _ := json.Marshal(user)
			c.Request, _ = http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(jsonUser))
			c.Request.Header.Set("Content-Type", "application/json")
			userHandler.CreateUser(c)

			tc.expectFunc(w, &mockResponse)
		})
	}
}

func TestReadUser(t *testing.T) {
	testCases := []struct {
		name           string
		setupInputFunc func(input *dto.ReadUserRequest, mockResponse *models.User)
		mockFunc       func(userRepo *mocks.MockUserRepository, mockResponse *models.User)
		expectFunc     func(w *httptest.ResponseRecorder, mockResponse *models.User)
	}{
		{
			name: "OK",
			setupInputFunc: func(input *dto.ReadUserRequest, mockResponse *models.User) {
				input.ID = 1
				mockResponse.FullName = "Full name"
				mockResponse.Username = "username"
				mockResponse.ID = input.ID
				mockResponse.CreatedAt = time.Now()
				mockResponse.UpdatedAt = mockResponse.CreatedAt
			},
			mockFunc: func(userRepo *mocks.MockUserRepository, mockResponse *models.User) {
				userRepo.On("ReadUser", mock.AnythingOfType("uint")).Return(mockResponse, nil)
			},
			expectFunc: func(w *httptest.ResponseRecorder, mockResponse *models.User) {
				assert.Equal(t, http.StatusOK, w.Code)
				expectBodyUser(t, w, mockResponse)
			},
		},
		{
			name: "BadInput",
			setupInputFunc: func(input *dto.ReadUserRequest, mockResponse *models.User) {
			},
			mockFunc: func(userRepo *mocks.MockUserRepository, mockResponse *models.User) {
				userRepo.On("ReadUser", mock.AnythingOfType("uint")).Return(mockResponse, nil)
			},
			expectFunc: func(w *httptest.ResponseRecorder, mockResponse *models.User) {
				assert.Equal(t, http.StatusBadRequest, w.Code)
			},
		},
		{
			name: "UserNotFound",
			setupInputFunc: func(input *dto.ReadUserRequest, mockResponse *models.User) {
				input.ID = 1
			},
			mockFunc: func(userRepo *mocks.MockUserRepository, mockResponse *models.User) {
				userRepo.On("ReadUser", mock.AnythingOfType("uint")).Return(mockResponse, errors.New("Not found"))
			},
			expectFunc: func(w *httptest.ResponseRecorder, mockResponse *models.User) {
				assert.Equal(t, http.StatusNotFound, w.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			userRepo := new(mocks.MockUserRepository)
			userPointRepo := new(mocks.MockUserPointRepository)
			userPointService := services.NewUserPointService(userPointRepo)
			userService := services.NewUserService(userRepo, userPointService)
			userHandler := NewUserHandler(userService)
			var input dto.ReadUserRequest
			var mockResponse models.User
			tc.setupInputFunc(&input, &mockResponse)
			tc.mockFunc(userRepo, &mockResponse)
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Params = gin.Params{{Key: "id", Value: fmt.Sprint(input.ID)}}

			// Act
			userHandler.ReadUser(c)

			// Assert
			tc.expectFunc(w, &mockResponse)
		})
	}
}

func TestListUser(t *testing.T) {
	testCases := []struct {
		name           string
		setupInputFunc func(input *dto.ListUserQuery, total *int64) []models.User
		mockFunc       func(userRepo *mocks.MockUserRepository, mockResponse []models.User, input *dto.ListUserQuery, total *int64)
		expectFunc     func(
			w *httptest.ResponseRecorder,
			mockResponse []models.User,
			input *dto.ListUserQuery,
			total *int64,
		)
	}{
		{
			name: "OK",
			setupInputFunc: func(input *dto.ListUserQuery, total *int64) []models.User {
				var mockResponse []models.User
				input.Page = int32(1)
				input.PerPage = int32(10)
				username := "username"
				input.Username = &username
				*total = 10
				now := time.Now()
				user1 := models.User{
					FullName: "Full name 1",
					Username: "username1",
				}
				user1.CreatedAt = now
				user1.UpdatedAt = now
				user1.ID = 1
				mockResponse = append(mockResponse, user1)

				user2 := models.User{
					FullName: "Full name 2",
					Username: "username2",
				}
				user2.CreatedAt = now
				user2.UpdatedAt = now
				user2.ID = 2
				mockResponse = append(mockResponse, user2)
				return mockResponse
			},
			mockFunc: func(userRepo *mocks.MockUserRepository, mockResponse []models.User, input *dto.ListUserQuery, total *int64) {
				userRepo.On("ListUsers", input.PerPage, input.Page, input.Username).Return(mockResponse, *total, nil)
			},
			expectFunc: func(
				w *httptest.ResponseRecorder,
				mockResponse []models.User,
				input *dto.ListUserQuery,
				total *int64,
			) {
				assert.Equal(t, http.StatusOK, w.Code)
				var response dto.ListUserResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, response.Metadata.Page, input.Page)
				assert.Equal(t, response.Metadata.PerPage, input.PerPage)
				assert.Equal(t, response.Metadata.Total, *total)
				assert.Len(t, response.Items, len(mockResponse))
				assert.Equal(t, response.Items[0].Username, mockResponse[0].Username)
				assert.Equal(t, response.Items[0].FullName, mockResponse[0].FullName)
			},
		},
		{
			name: "BadInput",
			setupInputFunc: func(input *dto.ListUserQuery, total *int64) []models.User {
				var mockResponse []models.User
				input.Page = 0
				input.PerPage = 0
				username := "username"
				input.Username = &username
				*total = 10
				return mockResponse
			},
			mockFunc: func(userRepo *mocks.MockUserRepository, mockResponse []models.User, input *dto.ListUserQuery, total *int64) {
				userRepo.On("ListUsers").Return(mockResponse, total, nil)
			},
			expectFunc: func(
				w *httptest.ResponseRecorder,
				mockResponse []models.User,
				input *dto.ListUserQuery,
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
			userRepo := new(mocks.MockUserRepository)
			userPointRepo := new(mocks.MockUserPointRepository)
			userPointService := services.NewUserPointService(userPointRepo)
			userService := services.NewUserService(userRepo, userPointService)
			userHandler := NewUserHandler(userService)
			var input dto.ListUserQuery
			var total int64
			mockResponse := tc.setupInputFunc(&input, &total)
			tc.mockFunc(userRepo, mockResponse, &input, &total)
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			url := fmt.Sprintf("/orders?page=%d&per_page=%d&username=%s", input.Page, input.PerPage, *input.Username)
			c.Request, _ = http.NewRequest(http.MethodGet, url, nil)

			// Act
			userHandler.ListUsers(c)

			// Assert
			tc.expectFunc(w, mockResponse, &input, &total)
		})
	}
}

func TestUpdateUser(t *testing.T) {
	testCases := []struct {
		name           string
		setupInputFunc func(input *dto.CreateUserDto, mockResponse *models.User)
		mockFunc       func(userRepo *mocks.MockUserRepository, mockResponse *models.User)
		expectFunc     func(w *httptest.ResponseRecorder, mockResponse *models.User)
	}{
		{
			name: "OK",
			setupInputFunc: func(input *dto.CreateUserDto, mockResponse *models.User) {
				input.FullName = "New full name"
				input.Username = "newUsername"
				mockResponse.ID = 1
				mockResponse.CreatedAt = time.Now()
				mockResponse.UpdatedAt = mockResponse.CreatedAt
				mockResponse.FullName = input.FullName
				mockResponse.Username = input.Username

			},
			mockFunc: func(userRepo *mocks.MockUserRepository, mockResponse *models.User) {
				userRepo.On("UpdateUser", mock.AnythingOfType("*models.User")).Return(nil).Run(func(args mock.Arguments) {
					arg := args.Get(0).(*models.User)
					arg.ID = mockResponse.ID
					arg.CreatedAt = mockResponse.CreatedAt
					arg.UpdatedAt = mockResponse.UpdatedAt
				})
			},
			expectFunc: func(w *httptest.ResponseRecorder, mockResponse *models.User) {
				assert.Equal(t, http.StatusCreated, w.Code)
				expectBodyUser(t, w, mockResponse)
			},
		},
		{
			name: "BadInput",
			setupInputFunc: func(input *dto.CreateUserDto, mockResponse *models.User) {
			},
			mockFunc: func(userRepo *mocks.MockUserRepository, mockResponse *models.User) {
				userRepo.On("UpdateUser", mock.AnythingOfType("*models.User")).Return(nil).Run(func(args mock.Arguments) {
					arg := args.Get(0).(*models.User)
					arg.ID = mockResponse.ID
					arg.CreatedAt = mockResponse.CreatedAt
					arg.UpdatedAt = mockResponse.UpdatedAt
				})
			},
			expectFunc: func(w *httptest.ResponseRecorder, mockResponse *models.User) {
				assert.Equal(t, http.StatusBadRequest, w.Code)
			},
		},
		{
			name: "UpdateUserError",
			setupInputFunc: func(input *dto.CreateUserDto, mockResponse *models.User) {
				input.FullName = "Full name"
				input.Username = "username"
				mockResponse.ID = 1
				mockResponse.CreatedAt = time.Now()
				mockResponse.UpdatedAt = mockResponse.CreatedAt
			},
			mockFunc: func(userRepo *mocks.MockUserRepository, mockResponse *models.User) {
				err := errors.New("Error")
				userRepo.On("UpdateUser", mock.AnythingOfType("*models.User")).Return(err)
			},
			expectFunc: func(w *httptest.ResponseRecorder, mockResponse *models.User) {
				assert.Equal(t, http.StatusInternalServerError, w.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			userRepo := new(mocks.MockUserRepository)
			userPointRepo := new(mocks.MockUserPointRepository)
			userPointService := services.NewUserPointService(userPointRepo)
			userService := services.NewUserService(userRepo, userPointService)
			userHandler := NewUserHandler(userService)
			var user dto.CreateUserDto
			var mockResponse models.User
			tc.setupInputFunc(&user, &mockResponse)
			tc.mockFunc(userRepo, &mockResponse)
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			jsonUser, _ := json.Marshal(user)
			c.Request, _ = http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(jsonUser))
			c.Request.Header.Set("Content-Type", "application/json")
			c.Params = gin.Params{{Key: "id", Value: fmt.Sprint(mockResponse.ID)}}

			userHandler.UpdateUser(c)

			tc.expectFunc(w, &mockResponse)
		})
	}
}
