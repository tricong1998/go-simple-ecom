package services

import (
	"github.com/tricong1998/go-ecom/internal/models"
	"github.com/tricong1998/go-ecom/internal/repository"
)

type UserService struct {
	UserRepo repository.IUserRepository
}

type IUserService interface {
	CreateUser(input *models.User) error
	ReadUser(id uint) *models.User
	GetUsers(
		perPage, page int,
		username *string,
	) ([]models.User, int64, error)
	UpdateUser(user *models.User) error
	DeleteUser(id uint) error
}

func NewUserService(userRepo repository.IUserRepository) *UserService {
	return &UserService{userRepo}
}

func (us *UserService) CreateUser(user *models.User) error {
	err := us.UserRepo.CreateUser(user)
	return err
}

func (us *UserService) ReadUser(id uint) *models.User {
	user := us.UserRepo.ReadUser(id)
	return user
}

func (us *UserService) GetUsers(
	perPage, page int,
	fullName *string,
) ([]models.User, int64, error) {
	return us.UserRepo.GetUsers(perPage, page, fullName)
}

func (us *UserService) UpdateUser(user *models.User) error {
	err := us.UserRepo.UpdateUser(user)
	return err
}

func (us *UserService) DeleteUser(id uint) error {
	return us.UserRepo.DeleteUser(id)
}
