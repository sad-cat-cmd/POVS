package services

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/sad-cat-cmd/WebApi/internal/models"
)

type ProductService struct {
	mutex    sync.Mutex
	products map[string]*models.Product
}

func NewProductServices() *ProductService {
	return &ProductService{
		products: make(map[string]*models.Product),
	}
}
func (s *ProductService) Add(product *models.Product) (*models.Product, error) {
	s.mutex.Lock()
	if product.ID == "" {
		product.ID = generateUUID()
	}
	s.products[product.ID] = product
	s.mutex.Unlock()
	return product, nil
}
func (s *ProductService) Remove(id string) (*models.Product, error) {
	s.mutex.Lock()
	product, exist := s.products[id]

	if !exist {
		s.mutex.Unlock()
		return nil, errors.New("Error with remove : object don't exist")
	}
	delete(s.products, id)
	s.mutex.Unlock()
	return product, nil
}
func (s *ProductService) Edit(new_product *models.Product) (*models.Product, error) {
	s.mutex.Lock()
	product, exist := s.products[new_product.ID]

	if !exist {
		s.mutex.Unlock()
		return nil, errors.New("Error with edit : object don't exist")
	}
	product.Name = new_product.Name
	product.ID = new_product.ID
	product.Definition = new_product.Definition
	product.Image = new_product.Image
	product.Price = new_product.Price
	product.UpdatedAt = time.Now()

	s.mutex.Unlock()
	return product, nil
}
func (s *ProductService) Search(id string) (*models.Product, error) {
	s.mutex.Lock()
	product, exist := s.products[id]

	if !exist {
		s.mutex.Unlock()
		return nil, nil
	}

	s.mutex.Unlock()
	return product, nil
}
func (s *ProductService) GetAll() ([]*models.Product, error) {
	s.mutex.Lock()

	result := make([]*models.Product, 0, len(s.products))

	for _, p := range s.products {
		result = append(result, p)
	}
	s.mutex.Unlock()
	return result, nil
}
func generateUUID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
