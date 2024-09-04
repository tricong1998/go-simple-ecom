package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tricong1998/go-ecom/internal/api/dto"
	"github.com/tricong1998/go-ecom/internal/models"
	"github.com/tricong1998/go-ecom/internal/services"
)

type UserHandler struct {
	UserService services.IUserService
}

func NewUserHandler(userService services.IUserService) *UserHandler {
	return &UserHandler{userService}
}

func (userHandler *UserHandler) CreateUser(ctx *gin.Context) {
	var input dto.CreateUserDto
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user := models.User{
		Username: input.Username,
		FullName: input.FullName,
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

func (userHandler *UserHandler) UpdateUser(ctx *gin.Context) {
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

	user := models.User{
		Username: input.Username,
		FullName: input.FullName,
	}
	user.ID = readUserRequest.ID
	if err := userHandler.UserService.UpdateUser(&user); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, dto.ToUserResponse(&user))
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
