package services

import (
	"github.com/sad-cat-cmd/WebApi/internal/models"
)

type IProductService interface {
	Add(product *models.Product) (*models.Product, error)
	Remove(id string) (*models.Product, error)
	Edit(new_product *models.Product) (*models.Product, error)
	Search(id string) (*models.Product, error)
	GetAll() ([]*models.Product, error)
}
