package repository

import (
	"github.com/tricong1998/go-ecom/cmd/user/pkg/models"
	"gorm.io/gorm"
)

type UserPointRepository struct {
	db *gorm.DB
}

type IUserPointRepository interface {
	CreateUserPoint(input *models.UserPoint) error
	ReadUserPoint(id uint) (*models.UserPoint, error)
	ListUserPoints(
		perPage, page int32,
		userId *uint,
	) ([]models.UserPoint, int64, error)
	UpdateUserPoint(input *models.UserPoint) error
	DeleteUserPoint(id uint) error
}

func NewUserPointRepository(db *gorm.DB) *UserPointRepository {
	return &UserPointRepository{db}
}

func (userRepo *UserPointRepository) CreateUserPoint(input *models.UserPoint) error {
	return userRepo.db.Create(input).Error
}

func (userRepo *UserPointRepository) ReadUserPoint(id uint) (*models.UserPoint, error) {
	var user *models.UserPoint
	err := userRepo.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (userRepo *UserPointRepository) ListUserPoints(
	perPage, page int32,
	userId *uint,
) ([]models.UserPoint, int64, error) {
	var users []models.UserPoint
	var total int64

	var query models.UserPoint
	if userId != nil {
		query.UserId = *userId
	}

	err := userRepo.db.Model(&models.UserPoint{}).Where(query).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	userRepo.db.Where(query).Find(&users)

	return users, total, nil
}

func (userRepo *UserPointRepository) UpdateUserPoint(input *models.UserPoint) error {
	return userRepo.db.Save(input).Error
}

func (userRepo *UserPointRepository) DeleteUserPoint(id uint) error {
	return userRepo.db.Delete(&models.UserPoint{}, id).Error
}
