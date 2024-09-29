package services

import (
	"github.com/tricong1998/go-ecom/cmd/product/internal/repository"
	"github.com/tricong1998/go-ecom/cmd/product/pkg/models"
)

type ProductService struct {
	ProductRepo repository.IProductRepository
}

type IProductService interface {
	CreateProduct(input *models.Product) error
	ReadProduct(id uint) (*models.Product, error)
	UpdateProductQuantity(productId, quantity uint) (bool, error)
	ListProducts(
		perPage, page int32,
		username *string,
	) ([]models.Product, int64, error)
	UpdateProduct(user *models.Product) error
	DeleteProduct(id uint) error
}

func NewProductService(userRepo repository.IProductRepository) *ProductService {
	return &ProductService{userRepo}
}

func (us *ProductService) CreateProduct(user *models.Product) error {
	err := us.ProductRepo.CreateProduct(user)
	return err
}

func (us *ProductService) ReadProduct(id uint) (*models.Product, error) {
	user, err := us.ProductRepo.ReadProduct(id)
	return user, err
}

func (us *ProductService) ListProducts(
	perPage, page int32,
	username *string,
) ([]models.Product, int64, error) {
	return us.ProductRepo.ListProducts(perPage, page, username)
}

func (us *ProductService) UpdateProduct(user *models.Product) error {
	err := us.ProductRepo.UpdateProduct(user)
	return err
}

func (us *ProductService) DeleteProduct(id uint) error {
	return us.ProductRepo.DeleteProduct(id)
}

func (us *ProductService) UpdateProductQuantity(productId, quantity uint) (bool, error) {
	return us.ProductRepo.UpdateProductQuantity(productId, quantity)
}
