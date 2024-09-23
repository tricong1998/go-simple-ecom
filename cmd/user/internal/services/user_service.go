package services

import (
	"github.com/tricong1998/go-ecom/cmd/user/internal/repository"
	"github.com/tricong1998/go-ecom/cmd/user/pkg/models"
)

type UserService struct {
	UserRepo repository.IUserRepository
}

type IUserService interface {
	CreateUser(input *models.User) error
	ReadUser(id uint) (*models.User, error)
	ListUsers(
		perPage, page int32,
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

func (us *UserService) ReadUser(id uint) (*models.User, error) {
	user, err := us.UserRepo.ReadUser(id)
	return user, err
}

func (us *UserService) ListUsers(
	perPage, page int32,
	username *string,
) ([]models.User, int64, error) {
	return us.UserRepo.ListUsers(perPage, page, username)
}

func (us *UserService) UpdateUser(user *models.User) error {
	err := us.UserRepo.UpdateUser(user)
	return err
}

func (us *UserService) DeleteUser(id uint) error {
	return us.UserRepo.DeleteUser(id)
}
