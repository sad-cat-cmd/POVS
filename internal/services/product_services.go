package services

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
	"time"

	"github.com/sad-cat-cmd/WebApi/internal/config"
	"github.com/sad-cat-cmd/WebApi/internal/models"
)

type ProductService struct {
	cfg      *config.Configuration
	mutex    sync.Mutex
	products map[string]*models.Product
}

func (s *ProductService) writeProductsInFile() error {
	productsList := make([]*models.Product,
		0,
		len(s.products))
	for _, p := range s.products {
		productsList = append(productsList, p)
	}

	data, err := json.MarshalIndent(productsList, "", "  ")
	if err != nil {
		return err
	}
	err = os.WriteFile(s.cfg.DataBaseFilePath,
		data,
		0644)
	if err != nil {
		return err
	}
	return nil
}

func (s *ProductService) initProductsFromFile() error {
	filePath := s.cfg.DataBaseFilePath
	if filePath == "" {
		return errors.New("Error: file config path is empty string")
	}
	_, err := os.Stat(filePath)
	if err != nil {
		return err
	}
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	var pdoductList []*models.Product
	err = json.Unmarshal(data, &pdoductList)
	if err != nil {
		return err
	}
	s.mutex.Lock()
	for _, product := range pdoductList {
		s.products[product.ID] = product
	}

	s.mutex.Unlock()
	return nil
}

func NewProductService(c *config.Configuration) (*ProductService, error) {
	thisProductService := ProductService{
		products: make(map[string]*models.Product),
		cfg:      c,
	}
	err := thisProductService.initProductsFromFile()
	if err != nil {
		return nil, err
	}

	return &thisProductService, nil
}
func (s *ProductService) Add(product *models.Product) (*models.Product, error) {
	s.mutex.Lock()
	s.products[product.ID] = product
	errWriteData := s.writeProductsInFile()
	if errWriteData != nil {
		return nil, errors.New("Error write data")
	}
	s.mutex.Unlock()
	return product, nil
}
func (s *ProductService) Remove(id string) (*models.Product, error) {
	s.mutex.Lock()
	product, exist := s.products[id]

	if !exist {
		s.mutex.Unlock()
		return nil, errors.New("Object doesn't exist")
	}
	delete(s.products, id)
	errWriteData := s.writeProductsInFile()
	if errWriteData != nil {
		return nil, errors.New("Error write data")
	}
	s.mutex.Unlock()
	return product, nil
}
func (s *ProductService) Edit(new_product *models.Product) (*models.Product, error) {
	s.mutex.Lock()
	product, exist := s.products[new_product.ID]

	if !exist {
		s.mutex.Unlock()
		return nil, errors.New("object doesn't exist")
	}
	product.Name = new_product.Name
	product.ID = new_product.ID
	product.Definition = new_product.Definition
	product.Image = new_product.Image
	product.Price = new_product.Price
	product.UpdatedAt = time.Now()

	errWrite := s.writeProductsInFile()
	if errWrite != nil {
		return nil, errors.New("Error write data")
	}
	s.mutex.Unlock()
	return product, nil
}
func (s *ProductService) Search(id string) (*models.Product, error) {
	s.mutex.Lock()
	product, exist := s.products[id]
	if !exist {
		s.mutex.Unlock()
		return nil, errors.New("Product not found")
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

// func generateUUID() string {
// 	return fmt.Sprintf("%d", time.Now().UnixNano())
// }
