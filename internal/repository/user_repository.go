package repository

import (
	"github.com/tricong1998/go-ecom/internal/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

type IUserRepository interface {
	CreateUser(input *models.User) error
	ReadUser(id uint) *models.User
	GetUsers(
		perPage, page int,
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

func (userRepo *UserRepository) ReadUser(id uint) *models.User {
	var user *models.User
	userRepo.db.First(&user, id)

	return user
}

func (userRepo *UserRepository) GetUsers(
	perPage, page int,
	fullName *string,
) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	var query models.User
	if fullName != nil {
		query.FullName = *fullName
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
