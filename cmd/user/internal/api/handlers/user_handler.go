package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tricong1998/go-ecom/cmd/user/internal/services"
	"github.com/tricong1998/go-ecom/cmd/user/internal/util"
	"github.com/tricong1998/go-ecom/cmd/user/pkg/dto"
	"github.com/tricong1998/go-ecom/cmd/user/pkg/models"
	"github.com/tricong1998/go-ecom/pkg/gin/middleware"
	"github.com/tricong1998/go-ecom/pkg/token"
)

type UserHandler struct {
	UserService services.IUserService
	JwtService  services.IJwtService
}

func NewUserHandler(userService services.IUserService, jwtService services.IJwtService) *UserHandler {
	return &UserHandler{userService, jwtService}
}

func (userHandler *UserHandler) CreateUser(ctx *gin.Context) {
	var input dto.CreateUserDto
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPassword, err := util.HashPassword(input.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	user := models.User{
		Username: input.Username,
		FullName: input.FullName,
		Password: hashedPassword,
		Role:     "user",
	}
	if err := userHandler.UserService.CreateUser(&user); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, dto.ToUserResponse(&user))
}

func (userHandler *UserHandler) ReadUser(ctx *gin.Context) {
	var readUserRequest dto.ReadUserRequest
	if err := ctx.ShouldBindUri(&readUserRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := userHandler.UserService.ReadUser(uint(readUserRequest.ID))
	if err != nil {
		err := fmt.Errorf("user not found: %d", readUserRequest.ID)
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, dto.ToUserResponse(user))
}

func (userHandler *UserHandler) ReadMe(ctx *gin.Context) {
	userPayload, ok := ctx.Get(middleware.AuthorizationPayloadKey)
	if !ok {
		err := errors.New("user not found")
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}
	payload := userPayload.(*token.Payload)
	user, err := userHandler.UserService.ReadUser(payload.UserId)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, dto.ToUserResponse(user))
}

func (userHandler *UserHandler) UpdateMe(ctx *gin.Context) {
	userPayload, ok := ctx.Get(middleware.AuthorizationPayloadKey)
	if !ok {
		err := errors.New("user not found")
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}
	payload := userPayload.(*token.Payload)
	user, err := userHandler.UserService.ReadUser(payload.UserId)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	var input dto.CreateUserDto
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var readUserRequest dto.ReadUserRequest
	if err := ctx.ShouldBindUri(&readUserRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user.Username = input.Username
	user.FullName = input.FullName
	if err := userHandler.UserService.UpdateUser(user); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, dto.ToUserResponse(user))
}

func (userHandler *UserHandler) ListUsers(ctx *gin.Context) {
	var req dto.ListUserQuery
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	users, total, err := userHandler.UserService.ListUsers(req.PerPage, req.Page, req.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var usersResponse []dto.UserResponse
	for _, v := range users {
		usersResponse = append(usersResponse, *dto.ToUserResponse(&v))
	}

	ctx.JSON(http.StatusOK, dto.ListUserResponse{
		Items: usersResponse,
		Metadata: dto.MetadataDto{
			Total:   total,
			Page:    req.Page,
			PerPage: req.PerPage,
		},
	})
}

func (userHandler *UserHandler) DeleteUser(ctx *gin.Context) {
	var readUserRequest dto.ReadUserRequest
	if err := ctx.ShouldBindUri(&readUserRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := userHandler.UserService.DeleteUser(uint(readUserRequest.ID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

func (userHandler *UserHandler) Login(ctx *gin.Context) {
	var input dto.LoginUserDto
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := userHandler.UserService.GetUserByUsername(input.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	err = util.ComparePassword(user.Password, input.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	accessToken, _, err := userHandler.JwtService.CreateToken(user.Username, user.ID, user.Role)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	refreshToken, _, err := userHandler.JwtService.CreateToken(user.Username, user.ID, user.Role)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}
