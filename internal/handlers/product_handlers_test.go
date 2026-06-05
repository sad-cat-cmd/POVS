package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"

	//"os"
	"path/filepath"
	"testing"

	"github.com/sad-cat-cmd/WebApi/internal/config"

	//"github.com/sad-cat-cmd/WebApi/internal/handlers"
	"github.com/sad-cat-cmd/WebApi/internal/models"
	"github.com/sad-cat-cmd/WebApi/internal/services/data"
)

type stubProductService struct {
	addResult    *models.Product
	addError     error
	getAllResult []*models.Product
	getAllError  error
	searchResult *models.Product
	searchError  error
	editResult   *models.Product
	editError    error
	removeResult *models.Product
	removeError  error
}

func (s *stubProductService) Add(p *models.Product) (*models.Product, error) {
	return s.addResult, s.addError
}
func (s *stubProductService) GetAll() ([]*models.Product, error) {
	return s.getAllResult, s.getAllError
}
func (s *stubProductService) Search(id string) (*models.Product, error) {
	return s.searchResult, s.searchError
}
func (s *stubProductService) Edit(p *models.Product) (*models.Product, error) {
	return s.editResult, s.editError
}
func (s *stubProductService) Remove(id string) (*models.Product, error) {
	return s.removeResult, s.removeError
}

type mockProductService struct {
	calledAddWith    *models.Product
	calledGetAll     bool
	calledSearchWith string
	calledEditWith   *models.Product
	calledRemoveWith string

	addCount    int
	getAllCount int
	searchCount int
	editCount   int
	removeCount int
}

func (m *mockProductService) Add(p *models.Product) (*models.Product, error) {
	m.addCount++
	m.calledAddWith = p
	return p, nil
}

func (m *mockProductService) GetAll() ([]*models.Product, error) {
	m.getAllCount++
	m.calledGetAll = true
	return nil, nil
}

func (m *mockProductService) Search(id string) (*models.Product, error) {
	m.searchCount++
	m.calledSearchWith = id
	return nil, nil
}

func (m *mockProductService) Edit(p *models.Product) (*models.Product, error) {
	m.editCount++
	m.calledEditWith = p
	return p, nil
}

func (m *mockProductService) Remove(id string) (*models.Product, error) {
	m.removeCount++
	m.calledRemoveWith = id
	return nil, nil
}
func TestProductHandler_Integration(t *testing.T) {
	tempDir := t.TempDir()

	dbPath := filepath.Join(tempDir, "test_products.db")

	cfg := &config.Configuration{
		DataBaseFilePath: dbPath,
	}

	realService, err := data.NewProductServiceSQLlite(cfg)
	if err != nil {
		t.Fatalf("Failed to create SQLite service: %v", err)
	}
	defer realService.CloseDB()

	handler := NewProductHandler(realService)
	var created_Product1 models.Product
	t.Run("#1-test-AddHandler-POST-request(correctness_product)", func(t *testing.T) {
		test_product, errNewProduct := models.NewProduct("name1",
			"definition1",
			500,
			"image1.jpeg")
		if errNewProduct != nil {
			t.Fatalf("Failed to create Product")
			return
		}
		bodyRequest, err := json.Marshal(test_product)
		if err != nil {
			t.Fatalf("Failed to Marhaling Body request")
			return
		}
		req := httptest.NewRequest("POST",
			"/products",
			bytes.NewReader(bodyRequest))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		handler.AddHandler(w, req)
		if w.Code != http.StatusCreated {
			t.Errorf("Add: \n\tExcepted status 201\n\tgot=%d\n\tbody=%q", w.Code, w.Body.String())
		}
		json.NewDecoder(w.Body).Decode(&created_Product1)
		if created_Product1.Name != "name1" {
			t.Errorf("Add: \n\tExpected name='name1'\n\tgot=%q", created_Product1.Name)
		}
	})
	t.Run("#2-test-AddHandler-POST-request(invalid_body_request)", func(t *testing.T) {
		invalidJson := `{"name":"badName", "price":"100"`
		bodyReader := bytes.NewReader([]byte(invalidJson))
		invalidReq := httptest.NewRequest("POST", "/products", bodyReader)
		invalidReq.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		handler.AddHandler(w, invalidReq)

		if (w.Code != 400) || (w.Body.String() != "Invalid request body\n") {
			t.Errorf("Add:\n\texpected status=400\n\tgot=%d\n\texpected body=Invalid request body\n\tgot=%q", w.Code, w.Body.String())
		}
	})
	t.Run("#3-test-AddHandler-Post-request(another_struct_request_and_unvalid_data)", func(t *testing.T) {
		invalidDataStruct := struct {
			name       string  `json:"name"`
			definition string  `json:"describe"`
			price      float32 `json:"price"`
			image      string  `json:"picture"`
		}{
			name:       "badName",
			definition: "badDefinition",
			price:      -100,
			image:      "bedImage.png",
		}
		invalidBodyRequest, err := json.Marshal(invalidDataStruct)
		if err != nil {
			t.Fatalf("Failed to Marhaling Body request")
			return
		}
		badDataRequest := httptest.NewRequest("POST", "/products", bytes.NewReader(invalidBodyRequest))
		badDataRequest.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		handler.AddHandler(w, badDataRequest)
		if w.Code != 400 || w.Body.String() != "Invalid request data\n" {
			t.Errorf("Add: \n\t expected status: 400\n\tgot=%d\n\t expected body=Invalid request data\n\tgot=%q", w.Code, w.Body.String())
		}
	})
	t.Run("#4-test-SearchHandler-GET-OnlyProduct-request(correctness-id--first-addded-product)", func(t *testing.T) {
		reqSearch := httptest.NewRequest("GET",
			"/products/"+created_Product1.ID,
			nil)
		wSearch := httptest.NewRecorder()

		handler.SearchHandler(wSearch, reqSearch)
		if wSearch.Code != http.StatusOK {
			t.Errorf("Search:\n\t expected status=200\n\tgot=%d\n\tbody=%q", wSearch.Code, wSearch.Body.String())
		}
		var found models.Product
		json.NewDecoder(wSearch.Body).Decode(&found)
		if found.Name != "name1" {
			t.Errorf("Search:\n\texpected name='name1'\n\tgot=%q", found.Name)
		}
	})
	t.Run("#5-test-SearchHandler-GET-OnlyProduct-request(id_is_empty)", func(t *testing.T) {
		reqSearch := httptest.NewRequest("GET",
			"/products/",
			nil)
		wSearch := httptest.NewRecorder()
		handler.SearchHandler(wSearch, reqSearch)
		if wSearch.Code != 400 || wSearch.Body.String() != "ID required\n" {
			t.Errorf("Search: \n\t expected status: 400\n\tgot=%d\n\texpected body=ID required\n\tgot=%q", wSearch.Code, wSearch.Body.String())
		}
	})
	t.Run("#6-test-SearchHandler-GET-OnlyProduct-request(id_is_invalid)", func(t *testing.T) {
		reqSearch := httptest.NewRequest("GET",
			"/products/"+"2131fwef",
			nil)
		wSearch := httptest.NewRecorder()
		handler.SearchHandler(wSearch, reqSearch)
		if wSearch.Code != 404 || wSearch.Body.String() != "Resource not found\n" {
			t.Errorf("Search: \n\t expected status: 404\n\tgot=%d\n\texpected body=Resource not found\n\tgot=%q", wSearch.Code, wSearch.Body.String())
		}
	})
	t.Run("#7-test-EditHandler-PUT-request(correctnessEdit)", func(t *testing.T) {
		updatedProduct, errEditNewProduct := models.NewProduct("Updated product",
			"",
			300,
			"image2.png")
		if errEditNewProduct != nil {
			t.Fatalf("Failed to create Product")
			return
		}
		bodyEdit, errEditMarsh := json.Marshal(updatedProduct)
		if errEditMarsh != nil {
			t.Fatalf("Failed to Marshaling")
			return
		}
		reqEdit := httptest.NewRequest("PUT",
			"/products/"+created_Product1.ID,
			bytes.NewReader(bodyEdit))
		reqEdit.Header.Set("Content-Type", "application/json")
		wEdit := httptest.NewRecorder()
		handler.EditHandler(wEdit, reqEdit)

		if wEdit.Code != http.StatusOK {
			t.Errorf("Edit:\n\texpected status=200\n\tgot=%d\n\tbody=%q", wEdit.Code, wEdit.Body.String())
		}
		var edited models.Product
		json.NewDecoder(wEdit.Body).Decode(&edited)
		if edited.Name != "Updated product" {
			t.Errorf("Edit:\n\texpected name=Updated product\n\tgot=%q",
				edited.Name)
		}
		if edited.Definition != "no data" {
			t.Errorf("Edit:\n\texpected definition=no data\n\tgot=%q",
				edited.Definition)
		}
		if edited.Price != 300 {
			t.Errorf("Edit:\n\texpected price=300\n\tgot=%f",
				edited.Price)
		}
		created_Product1.ID = edited.ID
	})
	t.Run("#8-test-EditHandler-PUT-request(id_is_empty)", func(t *testing.T) {
		updatedProduct, errEditNewProduct := models.NewProduct("Updated product",
			"",
			300,
			"image2.png")
		if errEditNewProduct != nil {
			t.Fatalf("Failed to create Product")
			return
		}
		bodyEdit, errEditMarsh := json.Marshal(updatedProduct)
		if errEditMarsh != nil {
			t.Fatalf("Failed to Marshaling")
			return
		}
		reqEdit := httptest.NewRequest("PUT",
			"/products/",
			bytes.NewReader(bodyEdit))
		reqEdit.Header.Set("Content-Type", "application/json")
		wEdit := httptest.NewRecorder()
		handler.EditHandler(wEdit, reqEdit)
		if wEdit.Code != 400 || wEdit.Body.String() != "ID required\n" {
			t.Errorf("Edit:\n\texpected status=400\n\tgot=%d\n\texpected body=ID required\n\tgot=%q", wEdit.Code, wEdit.Body.String())
		}
	})
	t.Run("#9-test-EditHandler-PUT-request(id_is_non-existent)", func(t *testing.T) {
		updatedProduct, errEditNewProduct := models.NewProduct("Updated product",
			"",
			300,
			"image2.png")
		if errEditNewProduct != nil {
			t.Fatalf("Failed to create Product")
			return
		}
		bodyEdit, errEditMarsh := json.Marshal(updatedProduct)
		if errEditMarsh != nil {
			t.Fatalf("Failed to Marshaling")
			return
		}
		reqEdit := httptest.NewRequest("PUT",
			"/products/1231241",
			bytes.NewReader(bodyEdit))
		reqEdit.Header.Set("Content-Type", "application/json")
		wEdit := httptest.NewRecorder()
		handler.EditHandler(wEdit, reqEdit)
		if wEdit.Code != 404 || wEdit.Body.String() != "Resource not found\n" {
			t.Errorf("Edit:\n\texpected status=400\n\tgot=%d\n\texpected body=Resource not found\n\tgot=%q", wEdit.Code, wEdit.Body.String())
		}
	})
	t.Run("#10-test-EditHandler-PUT-request(body_contains_invalid_data)", func(t *testing.T) {
		updatedProduct, errEditNewProduct := models.NewProduct("Updated product",
			"",
			300,
			"image2.png")
		if errEditNewProduct != nil {
			t.Fatalf("Failed to create Product")
			return
		}
		updatedProduct.Name = ""
		updatedProduct.Image = "221.exe"
		updatedProduct.Price = -100
		bodyEdit, errEditMarsh := json.Marshal(updatedProduct)
		if errEditMarsh != nil {
			t.Fatalf("Failed to Marshaling")
			return
		}
		reqEdit := httptest.NewRequest("PUT",
			"/products/"+created_Product1.ID,
			bytes.NewReader(bodyEdit))
		reqEdit.Header.Set("Content-Type", "application/json")
		wEdit := httptest.NewRecorder()
		handler.EditHandler(wEdit, reqEdit)
		if wEdit.Code != 400 || wEdit.Body.String() != "Invalid request data\n" {
			t.Errorf("Edit:\n\texpected status=400\n\tgot=%d\n\texpected body=Invalid request data\n\tgot=%q", wEdit.Code, wEdit.Body.String())
		}
	})
	t.Run("#11-test-EditHandler-PUT-request(invalid_body_request)", func(t *testing.T) {
		invalidJson := `{"name":"badName", "price":"100"`
		bodyReader := bytes.NewReader([]byte(invalidJson))
		invalidReq := httptest.NewRequest("PUT", "/products/"+created_Product1.ID, bodyReader)
		invalidReq.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		handler.AddHandler(w, invalidReq)

		if (w.Code != 400) || (w.Body.String() != "Invalid request body\n") {
			t.Errorf("Add:\n\t, expected status=400\n\t, got=%d\n\texpected body=Invalid request body\n\tgot=%q", w.Code, w.Body.String())
		}

	})
	t.Run("#12-test-GetAllHandler-GET-allProduct(check_correctness_edit)", func(t *testing.T) {
		reqGet := httptest.NewRequest("GET",
			"/products",
			nil)
		wGet := httptest.NewRecorder()
		handler.GetAllHandler(wGet, reqGet)

		var products []models.Product
		json.NewDecoder(wGet.Body).Decode(&products)
		if len(products) != 1 {
			t.Errorf("GetAll:\n\tExpected count products=1\n\tgot=%d", len(products))
		}
		if products[0].Name != "Updated product" {
			t.Errorf("GetAll:\n\tEpected name=Updated product\n\tgot=%q", products[0].Name)
		}
		if products[0].Definition != "no data" {
			t.Errorf("GetAll:\n\tEpected definition=no data\n\tgot=%q", products[0].Definition)
		}
		if products[0].Price != 300 {
			t.Errorf("GetAll:\n\tExpected price=300\n\tgot=%f", products[0].Price)
		}
		if products[0].Image != "image2.png" {
			t.Errorf("GetAll:\n\tEpected image=image2.png\n\tgot=%q", products[0].Image)
		}
	})
	// Remove
	t.Run("#13-test-RemoveHandler-DELETE-(normal_delete_first_added_product)", func(t *testing.T) {
		reqDelete := httptest.NewRequest("DELETE",
			"/products/"+created_Product1.ID,
			nil)
		wDelete := httptest.NewRecorder()
		handler.RemoveHandler(wDelete, reqDelete)

		if wDelete.Code != http.StatusOK {
			t.Errorf("Remove:\n\texpected status=200\n\tgot=%d\n\tbody=%q", wDelete.Code, wDelete.Body.String())
		}
	})
	t.Run("#14-test-RemoveHandler-DELETE-(emptyId)", func(t *testing.T) {
		reqDelete := httptest.NewRequest("DELETE",
			"/products/",
			nil)
		wDelete := httptest.NewRecorder()
		handler.RemoveHandler(wDelete, reqDelete)
		if wDelete.Code != 400 || wDelete.Body.String() != "ID required\n" {
			t.Errorf("Remove:\n\tExpected status=400\n\tgot=%d\n\tExpected body=ID required\n\tgot=%q", wDelete.Code, wDelete.Body.String())
		}
	})
	t.Run("#15-test-RemoveHandler-DELETE-(non-existens_id)", func(t *testing.T) {
		reqDelete := httptest.NewRequest("DELETE",
			"/products/1232131",
			nil)
		wDelete := httptest.NewRecorder()
		handler.RemoveHandler(wDelete, reqDelete)
		if wDelete.Code != 404 || wDelete.Body.String() != "Resource not found\n" {
			t.Errorf("Remove:\n\tExpected status=404\n\tgot=%d\n\tExpected body=Resource not found\n\tgot=%q", wDelete.Code, wDelete.Body.String())
		}
	})
	t.Run("#16-test-GetAllHandler-GET-allProduct(check_correctness_delete_and_products's_DB_is_empty)", func(t *testing.T) {
		reqGetAllFinal := httptest.NewRequest("GET",
			"/products",
			nil)
		wGetALlFinal := httptest.NewRecorder()
		handler.GetAllHandler(wGetALlFinal, reqGetAllFinal)

		var finalProducts []models.Product
		json.NewDecoder(wGetALlFinal.Body).Decode(&finalProducts)
		if len(finalProducts) != 0 {
			t.Errorf("GetAll after delete:\n\texpected count product=0\n\tgot=%d", len(finalProducts))
		}
	})
}

func TestProductHandler_Stub_Request(t *testing.T) {
	t.Run("#1-test-GetAllStub-Without-Error", func(t *testing.T) {
		stub := &stubProductService{
			getAllResult: []*models.Product{
				{ID: "1", Name: "Stub Book 1", Definition: "definition 1", Price: 100},
				{ID: "2", Name: "Stub Book 2", Definition: "definition 2", Price: 200},
			},
			getAllError: nil,
		}
		handler := NewProductHandler(stub)

		req := httptest.NewRequest("GET", "/products", nil)
		w := httptest.NewRecorder()
		handler.GetAllHandler(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("GetAll:\n\t Expected status=200\n\tgot=%d", w.Code)
		}
		var products []models.Product
		json.NewDecoder(w.Body).Decode(&products)
		if len(products) != 2 {
			t.Errorf("GetAll:\n\tExpected 2 products\n\tgot %d", len(products))
		}
		if products[0].Name != "Stub Book 1" {
			t.Errorf("GetAll:\n\tExpected products[0].Name=Stub Book 1\n\tgot=%q", products[0].Name)
		}
	})
	t.Run("#2-test-GetAllStub-With-Error", func(t *testing.T) {
		stub := &stubProductService{
			getAllResult: nil,
			getAllError:  errors.New("Specific-Error"),
		}
		handler := NewProductHandler(stub)
		req := httptest.NewRequest("GET", "/products", nil)
		w := httptest.NewRecorder()
		handler.GetAllHandler(w, req)
		if w.Code != http.StatusInternalServerError || w.Body.String() != "Internal server error\n" {
			t.Errorf("GetAll:\n\tExpected status=500\n\tgot=%d\n\tExpected body=Internal server error\n\tgot=%q", w.Code, w.Body.String())
		}
	})
	t.Run("#3-test-AddHandlerStub-Without-Error", func(t *testing.T) {
		stub := &stubProductService{
			addResult: &models.Product{ID: "12312", Definition: "definition", Name: "New book", Price: 500},
			addError:  nil,
		}
		handler := NewProductHandler(stub)
		reqBody := `{"name":"New book","definition":"definition","Price":500 }`
		req := httptest.NewRequest("POST", "/products", bytes.NewBuffer([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		handler.AddHandler(w, req)
		if w.Code != http.StatusCreated {
			t.Errorf("ADD:\n\tExpected status=201\n\tgot=%d", w.Code)
		}
		var created models.Product
		json.NewDecoder(w.Body).Decode(&created)
		if created.ID != "12312" {
			t.Errorf("ADD:\n\tExpected ID=12312\n\tgot=%q", created.ID)
		}
	})
	t.Run("#4-test-AddHandlerStub-With-Internal-Error", func(t *testing.T) {
		stub := &stubProductService{
			addResult: nil,
			addError:  errors.New("Error write data"),
		}
		handler := NewProductHandler(stub)
		reqBody := `{"name":"New book","definition":"definition","Price":500 }`
		req := httptest.NewRequest("POST", "/products", bytes.NewBuffer([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.AddHandler(w, req)
		if w.Code != 500 || w.Body.String() != "Internal Server Error\n" {
			t.Errorf("Add:\n\tExpected status=500\n\tgot=%d\n\tExpected body=Internal server error\n\tgot=%q", w.Code, w.Body.String())
		}
	})
	t.Run("#5-test-SearchHandlerStub-Without-Error", func(t *testing.T) {
		stub := &stubProductService{
			searchResult: &models.Product{ID: "123", Name: "Found Book", Definition: "definition", Price: 100},
			searchError:  nil,
		}
		handler := NewProductHandler(stub)

		req := httptest.NewRequest("GET", "/products/123", nil)
		w := httptest.NewRecorder()

		handler.SearchHandler(w, req)

		if w.Code != 200 {
			t.Errorf("Search:\n\tExpected status=200\n\tgot=%d", w.Code)
		}
		var searched models.Product
		json.NewDecoder(w.Body).Decode(&searched)
		if searched.Name != "Found Book" {
			t.Errorf("Search:\n\tExpected Name=200\n\tgot=%q", searched.Name)
		}
	})
	t.Run("#6-test-SearchHandlerStub-With-Found-Error", func(t *testing.T) {
		stub := &stubProductService{
			searchResult: nil,
			searchError:  errors.New("object doesn't exist"),
		}
		handler := NewProductHandler(stub)

		req := httptest.NewRequest("GET", "/products/999", nil)
		w := httptest.NewRecorder()
		handler.SearchHandler(w, req)

		if w.Code != 404 || w.Body.String() != "Resource not found\n" {
			t.Errorf("Search:\n\tExpected status=404\n\tgot=%d\n\tExected body=Resource not found\n\tgot=%q", w.Code, w.Body.String())
		}
	})
	t.Run("#7-test-RemoveHandlerStub-Without-Error", func(t *testing.T) {
		stub := &stubProductService{
			removeResult: &models.Product{ID: "123", Name: "Found Book", Definition: "definition", Price: 100},
			removeError:  nil,
		}
		handler := NewProductHandler(stub)
		req := httptest.NewRequest("DELETE", "/products/123", nil)
		w := httptest.NewRecorder()

		handler.RemoveHandler(w, req)
		if w.Code != 200 {
			t.Errorf("DELETE:\n\tExpected status=200\n\tgot=%d", w.Code)
		}
		var removed models.Product
		json.NewDecoder(w.Body).Decode(&removed)
		if removed.Name != "Found Book" {
			t.Errorf("DELETE:\n\tExpected Name=Found Book\n\tgot=%q", removed.Name)
		}
	})
	t.Run("#8-test-RemoveHandlerStub-With-Found-Error", func(t *testing.T) {
		stub := &stubProductService{
			removeResult: nil,
			removeError:  errors.New("Object doesn't exist"),
		}
		handler := NewProductHandler(stub)
		req := httptest.NewRequest("DELETE", "/products/123", nil)
		w := httptest.NewRecorder()
		handler.RemoveHandler(w, req)

		if w.Code != 404 || w.Body.String() != "Resource not found\n" {
			t.Errorf("DELETE:\n\tExpected status=404\n\tgot=%d\n\tExpected body=Resource not found\n\tgot=%q", w.Code, w.Body.String())
		}
	})
	t.Run("#9-test-EditHandlerStub-Without-Error", func(t *testing.T) {
		stub := &stubProductService{
			editResult: &models.Product{ID: "123", Name: "Found Book", Definition: "definition", Price: 100},
			editError:  nil,
		}
		handler := NewProductHandler(stub)
		reqBody := `{"name":"Found Book","definition":"definition","Price":100}`
		req := httptest.NewRequest("PUT", "/products/123", bytes.NewBuffer([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		handler.EditHandler(w, req)
		if w.Code != 200 {
			t.Errorf("EDIT:\n\tExpected status=200\n\tgot=%d", w.Code)
		}
	})
	t.Run("#10-test-EditHandlerStub-With-Found-Error", func(t *testing.T) {
		stub := &stubProductService{
			editResult: nil,
			editError:  errors.New("object doesn't exist"),
		}
		handler := NewProductHandler(stub)
		reqBody := `{"name":"Found Book","definition":"definition","Price":100}`
		req := httptest.NewRequest("PUT", "/products/123", bytes.NewBuffer([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		handler.EditHandler(w, req)

		if w.Code != 404 || w.Body.String() != "Resource not found\n" {
			t.Errorf("EDIT:\n\tExpected status=404\n\tgot=%d\n\tExpected body=Resource not found\n\tgot=%q", w.Code, w.Body.String())
		}
	})
}
func TestProductHandler_Mock_Manage_Service(t *testing.T) {
	t.Run("#1-test-AddHandler-Call-And-With-Correct-Product", func(t *testing.T) {
		mock := &mockProductService{}
		handler := NewProductHandler(mock)

		reqBody := `{"name":"Mock Book","definition":"test definition","price":500,"image":"test.jpg"}`
		req := httptest.NewRequest("POST", "/products", bytes.NewBuffer([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.AddHandler(w, req)

		if mock.addCount != 1 {
			t.Errorf("Add:\n\tExpected Add to be called once\n\tgot=%d", mock.addCount)
		}
		if mock.calledAddWith == nil {
			t.Errorf("Add:\n\tAdd was called with nil")
		}
		if mock.calledAddWith.Name != "Mock Book" {
			t.Errorf("Add:\n\tExpected Name=Mock Book\n\tgot=%q", mock.calledAddWith.Name)
		}
		if mock.calledAddWith.Definition != "test definition" {
			t.Errorf("Add:\n\tExpected Definition=test definition\n\tgot %q", mock.calledAddWith.Definition)
		}
		if mock.calledAddWith.Price != 500 {
			t.Errorf("Add:\n\tExpected Price=500\n\tgot %f", mock.calledAddWith.Price)
		}
		if mock.calledAddWith.Image != "test.jpg" {
			t.Errorf("Add:\n\tExpected Image=test.jpg\n\tgot %q", mock.calledAddWith.Image)
		}
	})
	t.Run("#2-test-GetAllHandler-Mock-Calls-GetAll", func(t *testing.T) {
		mock := &mockProductService{}
		handler := NewProductHandler(mock)

		req := httptest.NewRequest("GET", "/products", nil)
		w := httptest.NewRecorder()
		handler.GetAllHandler(w, req)

		if mock.getAllCount != 1 {
			t.Errorf("GetAll:\n\tExpected GetAll to be called once\n\tgot %d", mock.getAllCount)
		}
	})
	t.Run("#3-test-SearchHandler-Mock-Calls-Search-With-Correct-ID", func(t *testing.T) {
		mock := &mockProductService{}
		handler := NewProductHandler(mock)

		req := httptest.NewRequest("GET", "/products/test-id-123", nil)
		w := httptest.NewRecorder()
		handler.SearchHandler(w, req)

		if mock.searchCount != 1 {
			t.Errorf("Search:\n\tExpected Search to be called once\n\tgot %d", mock.searchCount)
		}
		if mock.calledSearchWith != "test-id-123" {
			t.Errorf("Search:\n\tExpected Search called with 'test-id-123'\n\tgot %q", mock.calledSearchWith)
		}
	})
	t.Run("#4-test-EditHandler-Mock-Calls-Edit-With-Correct-Product", func(t *testing.T) {
		mock := &mockProductService{}
		handler := NewProductHandler(mock)

		reqBody := `{"name":"Updated Book","definition":"new desc","price":250,"image":"new.jpg"}`
		req := httptest.NewRequest("PUT", "/products/edit-id-456", bytes.NewBuffer([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		handler.EditHandler(w, req)

		if mock.editCount != 1 {
			t.Errorf("Edit:\n\tExpected Edit to be called once\n\tgot %d", mock.editCount)
		}
		if mock.calledEditWith == nil {
			t.Fatal("Edit:\n\tEdit was called with nil")
		}
		if mock.calledEditWith.ID != "edit-id-456" {
			t.Errorf("Edit:\n\tExpected ID=edit-id-456\n\tgot %q", mock.calledEditWith.ID)
		}
		if mock.calledEditWith.Name != "Updated Book" {
			t.Errorf("Edit:\n\tExpected Name=Updated Book\n\tgot %q", mock.calledEditWith.Name)
		}
		if mock.calledEditWith.Price != 250 {
			t.Errorf("Edit:\n\tExpected Price=250\n\tgot %f", mock.calledEditWith.Price)
		}
	})
	t.Run("#5-test-RemoveHandler-Mock-Calls-Remove-With-Correct-ID", func(t *testing.T) {
		mock := &mockProductService{}
		handler := NewProductHandler(mock)

		req := httptest.NewRequest("DELETE", "/products/remove-id-789", nil)
		w := httptest.NewRecorder()
		handler.RemoveHandler(w, req)

		if mock.removeCount != 1 {
			t.Errorf("Edit\n\tExpected Remove to be called once\n\tgot %d", mock.removeCount)
		}
		if mock.calledRemoveWith != "remove-id-789" {
			t.Errorf("Edit\n\tExpected Remove called with remove-id-789\n\tgot %q", mock.calledRemoveWith)
		}
	})
	t.Run("#6-test-AddHandler-Mock-Not-Called-On-Invalid-JSON", func(t *testing.T) {
		mock := &mockProductService{}
		handler := NewProductHandler(mock)

		req := httptest.NewRequest("POST", "/products", bytes.NewBuffer([]byte(`{"name":`)))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		handler.AddHandler(w, req)

		if mock.addCount != 0 {
			t.Errorf("Add:\n\tAdd should not be called on invalid JSON, but was called %d times", mock.addCount)
		}
	})
	t.Run("#7-test-AddHandler-Mock-Not-Called-On-Validation-Error", func(t *testing.T) {
		mock := &mockProductService{}
		handler := NewProductHandler(mock)

		reqBody := `{"name":"","price":500}`
		req := httptest.NewRequest("POST", "/products", bytes.NewBuffer([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		handler.AddHandler(w, req)

		if mock.addCount != 0 {
			t.Errorf("Add:\n\tAdd should not be called on validation error, but was called %d times", mock.addCount)
		}
	})
	t.Run("#8-test-EditHandler-Mock-Not-Called-On-Empty-ID", func(t *testing.T) {
		mock := &mockProductService{}
		handler := NewProductHandler(mock)

		reqBody := `{"name":"Book","price":100}`
		req := httptest.NewRequest("PUT", "/products/", bytes.NewBuffer([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		handler.EditHandler(w, req)

		if mock.editCount != 0 {
			t.Errorf("Edit:\n\tEdit should not be called on empty ID, but was called %d times", mock.editCount)
		}
	})
	t.Run("#9-test-RemoveHandler-Mock-Not-Called-On-Empty-ID", func(t *testing.T) {
		mock := &mockProductService{}
		handler := NewProductHandler(mock)

		req := httptest.NewRequest("DELETE", "/products/", nil)
		w := httptest.NewRecorder()
		handler.RemoveHandler(w, req)

		if mock.removeCount != 0 {
			t.Errorf("Remove\n\tRemove should not be called on empty ID, but was called %d times", mock.removeCount)
		}
	})
	t.Run("#10-test-SearchHandler-Mock-Not-Called-On-Empty-ID", func(t *testing.T) {
		mock := &mockProductService{}
		handler := NewProductHandler(mock)

		req := httptest.NewRequest("GET", "/products/", nil)
		w := httptest.NewRecorder()
		handler.SearchHandler(w, req)

		if mock.searchCount != 0 {
			t.Errorf("Search:\n\tSearch should not be called on empty ID, but was called %d times", mock.searchCount)
		}
	})

}
