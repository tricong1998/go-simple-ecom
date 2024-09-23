package repository

import (
	"github.com/tricong1998/go-ecom/cmd/user/pkg/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

type IUserRepository interface {
	CreateUser(input *models.User) error
	ReadUser(id uint) (*models.User, error)
	ListUsers(
		perPage, page int32,
		username *string,
	) ([]models.User, int64, error)
	UpdateUser(input *models.User) error
	DeleteUser(id uint) error
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db}
}

func (userRepo *UserRepository) CreateUser(input *models.User) error {
	return userRepo.db.Create(input).Error
}

func (userRepo *UserRepository) ReadUser(id uint) (*models.User, error) {
	var user *models.User
	err := userRepo.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (userRepo *UserRepository) ListUsers(
	perPage, page int32,
	username *string,
) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	var query models.User
	if username != nil {
		query.Username = *username
	}

	err := userRepo.db.Model(&models.User{}).Where(query).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	userRepo.db.Where(query).Find(&users)

	return users, total, nil
}

func (userRepo *UserRepository) UpdateUser(input *models.User) error {
	return userRepo.db.Save(input).Error
}

func (userRepo *UserRepository) DeleteUser(id uint) error {
	return userRepo.db.Delete(&models.User{}, id).Error
}
