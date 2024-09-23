package repository

import (
	"github.com/tricong1998/go-ecom/cmd/product/pkg/models"
	"gorm.io/gorm"
)

type ProductRepository struct {
	db *gorm.DB
}

type IProductRepository interface {
	CreateProduct(input *models.Product) error
	ReadProduct(id uint) (*models.Product, error)
	ListProducts(
		perPage, page int32,
		username *string,
	) ([]models.Product, int64, error)
	UpdateProduct(input *models.Product) error
	DeleteProduct(id uint) error
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db}
}

func (userRepo *ProductRepository) CreateProduct(input *models.Product) error {
	return userRepo.db.Create(input).Error
}

func (userRepo *ProductRepository) ReadProduct(id uint) (*models.Product, error) {
	var user *models.Product
	err := userRepo.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (userRepo *ProductRepository) ListProducts(
	perPage, page int32,
	username *string,
) ([]models.Product, int64, error) {
	var users []models.Product
	var total int64

	var query models.Product
	if username != nil {
		query.Name = *username
	}

	err := userRepo.db.Model(&models.Product{}).Where(query).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	userRepo.db.Where(query).Find(&users)

	return users, total, nil
}

func (userRepo *ProductRepository) UpdateProduct(input *models.Product) error {
	return userRepo.db.Save(input).Error
}

func (userRepo *ProductRepository) DeleteProduct(id uint) error {
	return userRepo.db.Delete(&models.Product{}, id).Error
}
