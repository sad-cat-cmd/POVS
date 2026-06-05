package data

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/sad-cat-cmd/WebApi/internal/config"
	"github.com/sad-cat-cmd/WebApi/internal/models"
)

func setupTestService(t *testing.T) (*ProductServiceSQLlite, string) {
	t.Helper()
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	cfg := &config.Configuration{
		DataBaseFilePath: dbPath,
	}

	service, err := NewProductServiceSQLlite(cfg)
	if err != nil {
		t.Fatalf("Failed to create service: %v", err)
	}

	return service, dbPath
}

func createTestProduct(t *testing.T, id, name, definition string, price float64, image string) *models.Product {
	t.Helper()
	now := time.Now()
	return &models.Product{
		ID:         id,
		Name:       name,
		Definition: definition,
		Price:      price,
		Image:      image,
		CreateAt:   now,
		UpdatedAt:  now,
	}
}

func TestAdd_Success(t *testing.T) {
	service, _ := setupTestService(t)
	defer service.CloseDB()

	product := createTestProduct(t, "test-1", "Test Book", "Interesting story", 500, "book.jpg")

	added, err := service.Add(product)
	if err != nil {
		t.Fatalf("Add failed: %v", err)
	}
	if added == nil {
		t.Fatal("Add returned nil")
	}
	if added.ID != product.ID {
		t.Errorf("Expected ID=%q, got %q", product.ID, added.ID)
	}
	if added.Name != "Test Book" {
		t.Errorf("Expected Name='Test Book', got %q", added.Name)
	}
}

func TestAdd_DuplicateID(t *testing.T) {
	service, _ := setupTestService(t)
	defer service.CloseDB()

	product := createTestProduct(t, "same-id", "Book", "Desc", 100, "img.jpg")

	_, err := service.Add(product)
	if err != nil {
		t.Fatalf("First add failed: %v", err)
	}

	duplicate := createTestProduct(t, "same-id", "Duplicate", "Duplicate", 200, "dup.jpg")

	_, err = service.Add(duplicate)
	if err == nil {
		t.Error("Expected error for duplicate ID, got nil")
	}
}

func TestSearch_ExistingProduct(t *testing.T) {
	service, _ := setupTestService(t)
	defer service.CloseDB()

	product := createTestProduct(t, "search-1", "Searchable Book", "Description", 300, "search.jpg")
	_, err := service.Add(product)
	if err != nil {
		t.Fatalf("Add failed: %v", err)
	}

	found, err := service.Search(product.ID)
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}
	if found == nil {
		t.Fatal("Search returned nil for existing product")
	}
	if found.Name != "Searchable Book" {
		t.Errorf("Expected Name='Searchable Book', got %q", found.Name)
	}
}

func TestSearch_NonExistentProduct(t *testing.T) {
	service, _ := setupTestService(t)
	defer service.CloseDB()

	found, err := service.Search("non-existent-id")

	if found != nil {
		t.Errorf("Expected nil product, got %v", found)
	}

	if err == nil {
		t.Error("Expected error for non-existent product, got nil")
	}
	if err.Error() != "object doesn't exist" {
		t.Errorf("Expected error 'object doesn't exist', got %q", err.Error())
	}
}

func TestGetAll_Empty(t *testing.T) {
	service, _ := setupTestService(t)
	defer service.CloseDB()

	products, err := service.GetAll()
	if err != nil {
		t.Fatalf("GetAll failed: %v", err)
	}
	if len(products) != 0 {
		t.Errorf("Expected 0 products, got %d", len(products))
	}
}

func TestGetAll_MultipleProducts(t *testing.T) {
	service, _ := setupTestService(t)
	defer service.CloseDB()

	p1 := createTestProduct(t, "id-1", "Book1", "Desc1", 100, "img1.jpg")
	p2 := createTestProduct(t, "id-2", "Book2", "Desc2", 200, "img2.jpg")
	p3 := createTestProduct(t, "id-3", "Book3", "Desc3", 300, "img3.jpg")

	service.Add(p1)
	service.Add(p2)
	service.Add(p3)

	products, err := service.GetAll()
	if err != nil {
		t.Fatalf("GetAll failed: %v", err)
	}
	if len(products) != 3 {
		t.Errorf("Expected 3 products, got %d", len(products))
	}

	ids := make(map[string]bool)
	for _, p := range products {
		ids[p.ID] = true
	}
	if !ids[p1.ID] || !ids[p2.ID] || !ids[p3.ID] {
		t.Error("Not all products found in result")
	}
}

func TestEdit_Success(t *testing.T) {
	service, _ := setupTestService(t)
	defer service.CloseDB()

	product := createTestProduct(t, "edit-1", "Original", "Original desc", 100, "original.jpg")
	added, _ := service.Add(product)

	// Обновляем поля
	added.Name = "Updated Name"
	added.Definition = "Updated desc"
	added.Price = 200
	added.Image = "updated.jpg"
	added.UpdatedAt = time.Now()

	edited, err := service.Edit(added)
	if err != nil {
		t.Fatalf("Edit failed: %v", err)
	}
	if edited.Name != "Updated Name" {
		t.Errorf("Expected Name='Updated Name', got %q", edited.Name)
	}
	if edited.Price != 200 {
		t.Errorf("Expected Price=200, got %f", edited.Price)
	}

	found, _ := service.Search(added.ID)
	if found.Name != "Updated Name" {
		t.Errorf("After edit: expected 'Updated Name', got %q", found.Name)
	}
}

func TestEdit_NonExistentProduct(t *testing.T) {
	service, _ := setupTestService(t)
	defer service.CloseDB()

	product := createTestProduct(t, "non-existent-id", "Ghost", "Ghost desc", 100, "ghost.jpg")

	_, err := service.Edit(product)
	if err == nil {
		t.Error("Expected error for non-existent product, got nil")
	}
	if err.Error() != "object doesn't exist" {
		t.Errorf("Expected 'object doesn't exist', got %q", err.Error())
	}
}

func TestRemove_Success(t *testing.T) {
	service, _ := setupTestService(t)
	defer service.CloseDB()

	product := createTestProduct(t, "remove-1", "To Delete", "Will be deleted", 100, "delete.jpg")
	added, _ := service.Add(product)

	removed, err := service.Remove(added.ID)
	if err != nil {
		t.Fatalf("Remove failed: %v", err)
	}
	if removed == nil {
		t.Fatal("Remove returned nil")
	}
	if removed.ID != added.ID {
		t.Errorf("Expected ID=%q, got %q", added.ID, removed.ID)
	}

	found, _ := service.Search(added.ID)
	if found != nil {
		t.Error("Product still exists after removal")
	}
}

func TestRemove_NonExistentProduct(t *testing.T) {
	service, _ := setupTestService(t)
	defer service.CloseDB()

	_, err := service.Remove("non-existent-id")
	if err == nil {
		t.Error("Expected error for non-existent product, got nil")
	}
	if err.Error() != "object doesn't exist" {
		t.Errorf("Expected 'object doesn't exist', got %q", err.Error())
	}
}
func TestAdd_TableDriven(t *testing.T) {
	tests := []struct {
		name        string
		id          string
		productName string
		definition  string
		price       float64
		image       string
		wantErr     bool
	}{
		{
			name:        "valid product",
			id:          "valid-1",
			productName: "Valid Book",
			definition:  "Good book",
			price:       500,
			image:       "book.jpg",
			wantErr:     false,
		},
		{
			name:        "product with empty name (БД не валидирует)",
			id:          "empty-name-1",
			productName: "",
			definition:  "No name",
			price:       100,
			image:       "noname.jpg",
			wantErr:     false,
		},
		{
			name:        "product with negative price (БД не валидирует)",
			id:          "negative-price-1",
			productName: "Negative Price",
			definition:  "Desc",
			price:       -100,
			image:       "neg.jpg",
			wantErr:     false,
		},
		{
			name:        "product with zero price",
			id:          "zero-price-1",
			productName: "Free Book",
			definition:  "Free",
			price:       0,
			image:       "free.jpg",
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, _ := setupTestService(t)
			defer service.CloseDB()

			product := createTestProduct(t, tt.id, tt.productName, tt.definition, tt.price, tt.image)

			added, err := service.Add(product)
			if tt.wantErr {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("Add failed: %v", err)
			}
			if added == nil {
				t.Fatal("Add returned nil")
			}

			found, _ := service.Search(tt.id)
			if found == nil {
				t.Error("Product not found in database after add")
			}
		})
	}
}

func TestConcurrentAdd(t *testing.T) {
	service, _ := setupTestService(t)
	defer service.CloseDB()

	const goroutines = 10
	done := make(chan bool, goroutines)

	for i := 0; i < goroutines; i++ {
		go func(idx int) {
			id := generateTestID(idx)
			product := createTestProduct(t, id, "Concurrent", "Test", float64(idx), "img.jpg")
			_, err := service.Add(product)
			if err != nil {
				t.Errorf("Concurrent add failed for id=%s: %v", id, err)
			}
			done <- true
		}(i)
	}

	for i := 0; i < goroutines; i++ {
		<-done
	}

	products, _ := service.GetAll()
	if len(products) != goroutines {
		t.Errorf("Expected %d products, got %d", goroutines, len(products))
	}
}

func generateTestID(idx int) string {
	return time.Now().Format("20060102150405") + string(rune('0'+idx%10))
}
