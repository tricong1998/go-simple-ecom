package repository

import (
	"errors"

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
	UpdateProductQuantity(productId, quantity uint) (bool, error)
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

func (userRepo *ProductRepository) UpdateProductQuantity(productId, quantity uint) (bool, error) {
	tx := userRepo.db.Begin()
	if tx.Error != nil {
		return false, tx.Error
	}

	var product models.Product
	err := tx.First(&product, productId).Error
	if err != nil {
		return false, err
	}
	if product.Quantity < quantity {
		return false, errors.New("product quantity is not enough")
	}

	oldQuantity := product.Quantity
	product.Quantity = product.Quantity - quantity
	err = tx.Model(&product).
		Where("id = ?", productId).
		Where("quantity = ?", oldQuantity).
		Update("quantity", product.Quantity).Error
	if err != nil {
		tx.Rollback()
		return false, err
	}

	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return false, err
	}

	return true, nil
}
