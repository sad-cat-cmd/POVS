package services

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	//"time"

	"github.com/sad-cat-cmd/WebApi/internal/config"
	"github.com/sad-cat-cmd/WebApi/internal/models"
	//"github.com/sad-cat-cmd/WebApi/internal/services"
	//"github.com/sad-cat-cmd/WebApi/internal/services"
	//"github.com/sad-cat-cmd/WebApi/internal/models"
)

type testingProduct struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	Definition string  `json:"definition"`
	Price      float64 `json:"price"`
	Image      string  `json:"image"`
	CreateAt   string  `json:"create_at"`
	UpdatedAt  string  `json:"updated_at"`
}
type resultCaseNewProductService struct {
	expectedProducts   map[string]*testingProduct
	expectedConfigPath string
	exptectedErrorMsg  string
}
type inputCaseNewProductService struct {
	name          string
	inConfig      *config.Configuration
	createdFileDb bool
	contentInDb   []*testingProduct
	wantService   bool
	wantErr       bool
	check         func(t *testing.T,
		expected *resultCaseNewProductService,
		service *ProductService,
		err error)
}

func createFileConfig(filePath string, content string) error {
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return os.WriteFile(filePath,
		[]byte(content),
		0644)
}
func removeFileConfig(filePath string) error {
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return nil
	}
	return os.Remove(filePath)
}

func CreateDb(filePath string,
	products []*testingProduct) error {
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	data, errMarsh := json.MarshalIndent(products, "", " ")
	if errMarsh != nil {
		return errMarsh
	}
	return os.WriteFile(filePath, data, 0644)
}
func RemoveDb(filePath string) error {
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return nil
	}
	return os.Remove(filePath)
}

func checkValidProductService(t *testing.T,
	expected *resultCaseNewProductService,
	service *ProductService,
	err error) {
	if service.cfg.DataBaseFilePath != expected.expectedConfigPath {
		t.Errorf("Expected service.configPath=%q, got service.configPath=%q",
			expected.expectedConfigPath, service.cfg.DataBaseFilePath)
		return
	}
	for id, expectedProduct := range expected.expectedProducts {
		actualProduct, exists := service.products[id]
		if exists == false {
			t.Errorf("Product with ID %q not found in service.products", id)
			continue
		}
		if actualProduct.Name != expectedProduct.Name {
			t.Errorf("Product %q: expected Name=%q, got %q",
				id, expectedProduct.Name, actualProduct.Name)
		}
		if actualProduct.Definition != expectedProduct.Definition {
			t.Errorf("Product %q: expected Definition=%q, got %q",
				id, expectedProduct.Definition, actualProduct.Definition)
		}
		if actualProduct.Price != expectedProduct.Price {
			t.Errorf("Product %q: expected Price=%f, got %f",
				id, expectedProduct.Price, actualProduct.Price)
		}
		if actualProduct.Image != expectedProduct.Image {
			t.Errorf("Product %q: expected Image=%q, got %q",
				id, expectedProduct.Image, actualProduct.Image)
		}
		actualCreateAtStr := actualProduct.CreateAt.Format(time.RFC3339Nano)
		if actualCreateAtStr != expectedProduct.CreateAt {
			t.Errorf("Product: %q expected CreateAt=%q, got %q",
				id, expectedProduct.CreateAt, actualProduct.CreateAt.String())
		}
		actualUpdateAtStr := actualProduct.UpdatedAt.Format(time.RFC3339Nano)
		if actualUpdateAtStr != expectedProduct.UpdatedAt {
			t.Errorf("Product: %q expected UpdatedAt=%q, got %q",
				id, expectedProduct.UpdatedAt, actualProduct.UpdatedAt.String())
		}
	}
}

func checkUnvalidProductService(t *testing.T,
	expected *resultCaseNewProductService,
	service *ProductService,
	err error) {
	if err.Error() != expected.exptectedErrorMsg {
		t.Errorf("Expected error:%q, got error %q",
			expected.exptectedErrorMsg, err)
		return
	}
}

func getAllConfig(filePathes []string,
	contents []string) ([]*config.Configuration, error) {
	configs := []*config.Configuration{}

	var index int = 0
	for _, path := range filePathes {
		errFileConfigWork := createFileConfig(path, contents[index])
		if errFileConfigWork != nil {
			return nil, errFileConfigWork
		}
		index++
	}
	for _, path := range filePathes {
		newConfig, errorNewConfig := config.NewConfiguration(path)
		if errorNewConfig != nil {
			return nil, errorNewConfig
		}
		configs = append(configs, newConfig)
	}
	for _, path := range filePathes {
		errFileConfigWork := removeFileConfig(path)
		if errFileConfigWork != nil {
			return nil, errFileConfigWork
		}
	}
	return configs, nil
}
func TestNewProductService(t *testing.T) {
	// init testing value
	filePathes := []string{
		"testNormalConfig.json",
		"testEmptyObectConfig.json",
	}
	fileContents := []string{
		`{"DataBaseFilePath": "data/test_products_1.txt"}`,
		`{"DataBaseFilePath": "data/test_products_2.txt"}`,
	}
	configurations, errorGetConfigs := getAllConfig(filePathes, fileContents)
	if errorGetConfigs != nil {
		t.Errorf("---- INTERNAL ERROR: FILE CONFIG :%q", errorGetConfigs.Error())
		return
	}

	var outErr string = ""
	var outService string = ""
	var index_test int = 0

	expectedResult := []*resultCaseNewProductService{
		{map[string]*testingProduct{
			"121241241": &testingProduct{
				"121241241",
				"name1",
				"definition1",
				100,
				"image1.png",
				"2026-05-19T09:02:45.660194838+03:00",
				"2026-05-19T09:02:45.660194838+03:00",
			},
			"1212412331241": &testingProduct{
				"1212412331241",
				"name2",
				"definition2",
				500,
				"image2.png",
				"2026-05-19T09:02:45.660194838+03:00",
				"2026-05-19T09:02:45.660194838+03:00",
			},
			"121241233124321": &testingProduct{
				"121241233124321",
				"name3",
				"no data",
				500,
				"no data",
				"2026-05-19T09:02:45.660194838+03:00",
				"2026-05-19T09:02:45.660194838+03:00",
			},
		}, configurations[0].DataBaseFilePath, "",
		},
		{map[string]*testingProduct{},
			configurations[1].DataBaseFilePath, "Error: file config path is empty string",
		},
		{map[string]*testingProduct{},
			configurations[0].DataBaseFilePath, "stat data/test_products_1.txt: no such file or directory",
		},
	}
	testIn := []*inputCaseNewProductService{
		{
			"normal test #1", configurations[0], true,
			[]*testingProduct{
				{
					"121241241",
					"name1",
					"definition1",
					100,
					"image1.png",
					"2026-05-19T09:02:45.660194838+03:00",
					"2026-05-19T09:02:45.660194838+03:00",
				},
				{
					"1212412331241",
					"name2",
					"definition2",
					500,
					"image2.png",
					"2026-05-19T09:02:45.660194838+03:00",
					"2026-05-19T09:02:45.660194838+03:00",
				},
				{
					"121241233124321",
					"name3",
					"no data",
					500,
					"no data",
					"2026-05-19T09:02:45.660194838+03:00",
					"2026-05-19T09:02:45.660194838+03:00",
				},
			},
			true, false, checkValidProductService,
		},
		{
			"normal test #2 (DB file is empty)", configurations[1], true,
			[]*testingProduct{},
			true, false, func(t *testing.T, expected *resultCaseNewProductService, service *ProductService, err error) {
				if len(service.products) != 0 {
					t.Errorf("Expected service.products length = 0, got service.products length %d",
						len(service.products))
					return
				}
			},
		},
		{
			"negative test #3 (DB file has been not existe)", configurations[0], false,
			[]*testingProduct{},
			false, true, checkUnvalidProductService,
		},
	}
	// end init testing value

	for _, tt := range testIn {
		t.Run(tt.name, func(t *testing.T) {
			var errFileDbWork error
			if tt.createdFileDb == true {
				errFileDbWork = CreateDb(tt.inConfig.DataBaseFilePath,
					tt.contentInDb)
			}
			if errFileDbWork != nil {
				t.Errorf("---- INTERNAL ERROR: FILE DB = %q", errFileDbWork.Error())
				return
			}
			service, err := NewProductService(tt.inConfig)
			if (err != nil) != tt.wantErr || (service != nil) != tt.wantService {
				if tt.wantErr == true {
					outErr = "error(non-nil)"
				} else {
					outErr = "no error(nil)"
				}
				if tt.wantService == true {
					outService = "non-nil"
				} else {
					outService = "nil"
				}
				t.Errorf("NewProductService(testConfig)\n"+
					"\tgot: (ProductService=%v, error=%v)\n"+
					"\twant: (ProductService=%q, error=%q)",
					service, err,
					outService, outErr)
				RemoveDb(tt.inConfig.DataBaseFilePath)
				index_test++
				return
			}
			if tt.check != nil {
				tt.check(t,
					expectedResult[index_test],
					service,
					err)
			}
			RemoveDb(tt.inConfig.DataBaseFilePath)
			index_test++
		})
	}
}
func setupTestService(t *testing.T) (*ProductService, string, error) {
	t.Helper()
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test_products.json")
	if err := os.WriteFile(dbPath, []byte("[]"), 0644); err != nil {
		return nil, "", err
	}
	cfg := &config.Configuration{DataBaseFilePath: dbPath}
	service, err := NewProductService(cfg)
	if err != nil {
		t.Fatalf("Failed to create service: %v", err)
		return nil, "", errors.New("Error:" + err.Error())
	}
	return service, dbPath, nil
}
func setupTestServiceWithData(t *testing.T, products []*models.Product) (*ProductService, string, error) {
	t.Helper()
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test_products.json")

	if len(products) > 0 {
		productsJSON := make([]models.Product, len(products))
		for i, p := range products {
			productsJSON[i] = *p
		}
		data, err := json.MarshalIndent(productsJSON, "", "  ")
		if err != nil {
			return nil, "", err
		}

		if err := os.WriteFile(dbPath, data, 0644); err != nil {
			return nil, "", err
		}
	}

	cfg := &config.Configuration{DataBaseFilePath: dbPath}

	service, err := NewProductService(cfg)
	if err != nil {
		return nil, "", err
	}

	return service, dbPath, nil
}
func areProductsEqual(t *testing.T,
	trueProduct *models.Product,
	p2 *models.Product,
	addInfo string) error {
	var flagError bool = false
	if trueProduct.ID != p2.ID {
		t.Errorf(addInfo+"Expected product.ID=%q, got=%q", trueProduct.ID, p2.ID)
		flagError = true
	}
	if trueProduct.Name != p2.Name {
		t.Errorf(addInfo+"Expected product.Name=%q, got=%q", trueProduct.Name, p2.Name)
		flagError = true
	}
	if trueProduct.Price != p2.Price {
		t.Errorf(addInfo+"Expected product.Price=%f, got=%f", trueProduct.Price, p2.Price)
		flagError = true
	}
	if trueProduct.Image != p2.Image {
		t.Errorf(addInfo+"Expected product.Image=%q, got=%q", trueProduct.Image, p2.Image)
		flagError = true
	}
	difCreateTime := trueProduct.CreateAt.Nanosecond() - p2.CreateAt.Nanosecond()
	if difCreateTime < 0 {
		difCreateTime = p2.CreateAt.Nanosecond() - trueProduct.CreateAt.Nanosecond()
	}
	if difCreateTime > int(time.Second) {
		t.Errorf(addInfo+"Expected product.CreateAt=%q, got=%q", trueProduct.CreateAt.String(), p2.CreateAt.String())
		flagError = true
	}
	difUpdateTime := trueProduct.CreateAt.Nanosecond() - p2.UpdatedAt.Nanosecond()
	if difUpdateTime < 0 {
		difUpdateTime = p2.UpdatedAt.Nanosecond() - trueProduct.UpdatedAt.Nanosecond()
	}
	if difUpdateTime > int(time.Second) {
		t.Errorf(addInfo+"Expected product.UpdatedAt=%q, got=%q", trueProduct.UpdatedAt.String(), p2.UpdatedAt.String())
		flagError = true
	}
	if flagError == true {
		return errors.New("Error")
	} else {
		return nil
	}
}
func verifyFileContainsProducts(t *testing.T,
	filePath string,
	expected map[string]*models.Product) {
	//t.Helper()
	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read DB file: %v", err)
	}

	var productsFromDisk []models.Product
	if err := json.Unmarshal(data, &productsFromDisk); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if len(productsFromDisk) != len(expected) {
		t.Errorf("File: expected %d products, got %d", len(expected), len(productsFromDisk))
		return
	}

	diskMap := make(map[string]models.Product)
	for _, p := range productsFromDisk {
		diskMap[p.ID] = p
	}

	for id, expectedProduct := range expected {
		actualProduct, exists := diskMap[id]
		if !exists {
			t.Errorf("File: product with ID %q not found", id)
			continue
		}
		areProductsEqual(t, expectedProduct, &actualProduct, "File:")
	}
}

func verifyMemoryContainsProducts(t *testing.T,
	service *ProductService,
	expected map[string]*models.Product) {
	//t.Helper()

	if len(service.products) != len(expected) {
		t.Errorf("Memory: expected %d products, got %d", len(expected), len(service.products))
		return
	}

	for id, expectedProduct := range expected {
		actualProduct, exists := service.products[id]
		if !exists {
			t.Errorf("Memory: product with ID %q not found", id)
			continue
		}
		areProductsEqual(t, expectedProduct, actualProduct, "Memory:")
	}
}
func TestAdd(t *testing.T) {
	service, dbPath, err := setupTestService(t)
	if err != nil {
		t.Fatalf("Failed to create service and DB: %v", err)
	}

	product1, err := models.NewProduct("name1", "definition1", 100, "image1.png")
	if err != nil {
		t.Fatalf("Failed to create product1: %v", err)
	}
	product2, err := models.NewProduct("name2", "", 200, "")
	if err != nil {
		t.Fatalf("Failed to create product2: %v", err)
	}

	expectedAfterFirst := map[string]*models.Product{
		product1.ID: product1,
	}
	expectedAfterSecond := map[string]*models.Product{
		product1.ID: product1,
		product2.ID: product2,
	}

	t.Run("normal test #1(add first product)", func(t *testing.T) {
		result, err := service.Add(product1)
		if err != nil || result == nil {
			t.Errorf("Add(product1)\n"+
				"\tgot (product=%v, err=%v)\n"+
				"\twant (product=non-nil, err=nil)",
				result, err)
			return
		}

		if err := areProductsEqual(t, product1, result, "Returned product:"); err != nil {
			return
		}

		verifyFileContainsProducts(t, dbPath, expectedAfterFirst)
		verifyMemoryContainsProducts(t, service, expectedAfterFirst)
	})

	t.Run("normal test #2(add second product)", func(t *testing.T) {
		result, err := service.Add(product2)
		if err != nil || result == nil {
			t.Errorf("Add(product2)\n"+
				"\tgot (product=%v, err=%v)\n"+
				"\twant (product=non-nil,err=nil)",
				result, err)
			return
		}

		if err := areProductsEqual(t, product2, result, "Returned product:"); err != nil {
			return
		}

		verifyFileContainsProducts(t, dbPath, expectedAfterSecond)
		verifyMemoryContainsProducts(t, service, expectedAfterSecond)
	})
}
func TestRemove(t *testing.T) {
	product1, err := models.NewProduct("name1", "definition1", 100, "image1.png")
	if err != nil {
		t.Fatalf("Failed to create product1: %v", err)
	}
	product2, err := models.NewProduct("name2", "", 200, "")
	if err != nil {
		t.Fatalf("Failed to create product2: %v", err)
	}

	products := []*models.Product{
		product1,
		product2,
	}
	expectedAfterFirstRemove := map[string]*models.Product{
		product2.ID: product2,
	}
	expectedAfterSecondRemove := map[string]*models.Product{}

	service, dbPath, errorSetup := setupTestServiceWithData(t, products)
	if errorSetup != nil {
		t.Fatalf("Failed to create service and DB: %v", errorSetup)
	}
	t.Run("normal test #1 (delition by first ID)", func(t *testing.T) {
		result, err := service.Remove(product1.ID)
		if err != nil || result == nil {
			t.Errorf("Remove(product1.ID)\n"+
				"\tgot (product=%v, err=%v)\n"+
				"\twant (product=non-nil, err=nil)",
				result, err)
			return
		}
		err1 := areProductsEqual(t, product1, result, "returned product")
		if err1 != nil {
			return
		}
		verifyFileContainsProducts(t, dbPath, expectedAfterFirstRemove)
		verifyMemoryContainsProducts(t, service, expectedAfterFirstRemove)
	})
	t.Run("normal test #2 (delition by second ID)", func(t *testing.T) {
		result, err := service.Remove(product2.ID)
		if err != nil || result == nil {
			t.Errorf("Remove(product2.id)\n"+
				"\tgot (product=%v, err=%v)\n"+
				"\twant (product=non-nil,err=nil)",
				result, err)
			return
		}
		err1 := areProductsEqual(t, product2, result, "returned product")
		if err1 != nil {
			return
		}
		verifyFileContainsProducts(t, dbPath, expectedAfterSecondRemove)
		verifyMemoryContainsProducts(t, service, expectedAfterSecondRemove)
	})
	t.Run("negative test #3 (delition by non-existent ID)", func(t *testing.T) {
		result, err := service.Remove("1241фыassaa")
		if err == nil || result != nil {
			t.Errorf("Remove(%q)\n"+
				"\tgot (product=%v, err=%v)\n"+
				"\twant (product=nil,err=error(non-nil))",
				"1241фыassaa",
				result, err)
			return
		}
		if err.Error() != "Object doesn't exist" {
			t.Errorf("Excpected error=Object doesn't exist, got=%q", err.Error())
		}
	})
}
func TestEdit(t *testing.T) {
	product1, err := models.NewProduct("name1", "definition1", 100, "image1.png")
	if err != nil {
		t.Fatalf("Failed to create product1: %v", err)
	}
	product2, err := models.NewProduct("name2", "", 200, "")
	if err != nil {
		t.Fatalf("Failed to create product2: %v", err)
	}
	product1Edit, err := models.NewProduct("newname1", "новое определение", 150, "image.png")
	product2Edit, err := models.NewProduct("newname2", "", 150, "image.png")
	product1Edit.ID = product1.ID
	product2Edit.ID = product2.ID

	products := []*models.Product{
		product1,
		product2,
	}
	expectedAfterFirstEdit := map[string]*models.Product{
		product1.ID: product1Edit,
		product2.ID: product2,
	}
	expectedAfterSecondEdit := map[string]*models.Product{
		product1.ID: product1Edit,
		product2.ID: product2Edit,
	}

	service, dbPath, errorSetup := setupTestServiceWithData(t, products)
	if errorSetup != nil {
		t.Fatalf("Failed to create service and DB: %v", errorSetup)
	}

	t.Run("normal test#1 (edits the first element)", func(t *testing.T) {
		edited, err := service.Edit(product1Edit)
		if edited == nil || err != nil {
			t.Errorf("Edit(product1)\n"+
				"\tgot (product=%v, err=%v)\n"+
				"\twant (product=nil,err=error(non-nil))",
				edited, err)
			return
		}
		err1 := areProductsEqual(t, product1Edit, edited, "returned product")
		if err1 != nil {
			return
		}
		verifyFileContainsProducts(t, dbPath, expectedAfterFirstEdit)
		verifyMemoryContainsProducts(t, service, expectedAfterFirstEdit)
	})

	t.Run("normal test#2 (edits the second element)", func(t *testing.T) {
		edited, err := service.Edit(product2Edit)
		if edited == nil || err != nil {
			t.Errorf("Edit(product2)\n"+
				"\tgot (product=%v, err=%v)\n"+
				"\twant (product=nil,err=error(non-nil))",
				edited, err)
			return
		}
		err1 := areProductsEqual(t, product2Edit, edited, "returned product")
		if err1 != nil {
			return
		}
		verifyFileContainsProducts(t, dbPath, expectedAfterSecondEdit)
		verifyMemoryContainsProducts(t, service, expectedAfterSecondEdit)
	})
	t.Run("negative test#3 (edits the product with non-existent id)", func(t *testing.T) {
		testNegativeProduct, errCreateNegativeTest := models.NewProduct("fakeName", "fakeDefinition", 100, "fakeImage.png")
		if errCreateNegativeTest != nil {
			t.Fatalf("Failed to create negative product")
			return
		}

		edited, err := service.Edit(testNegativeProduct)
		if edited != nil || err == nil {
			t.Errorf("Edit(product(non-existent id))\n"+
				"\tgot (product=%v, err=%v)\n"+
				"\twant (product=nil,err=error(non-nil))",
				edited, err)
			return
		}
		if err.Error() != "object doesn't exist" {
			t.Errorf("Excpected error=Object don't exist, got=%q", err.Error())
		}
	})
}
func TestSearch(t *testing.T) {
	product1, err := models.NewProduct("name1", "definition1", 100, "image1.png")
	if err != nil {
		t.Fatalf("Failed to create product1: %v", err)
	}
	product2, err := models.NewProduct("name2", "", 200, "")
	if err != nil {
		t.Fatalf("Failed to create product2: %v", err)
	}

	products := []*models.Product{
		product1,
		product2,
	}

	service, _, errorSetup := setupTestServiceWithData(t, products)
	if errorSetup != nil {
		t.Fatalf("Failed to create service and DB: %v", errorSetup)
	}
	t.Run("normal test#1 (searches at correcntess id)", func(t *testing.T) {
		searched, err := service.Search(product1.ID)
		if searched == nil || err != nil {
			t.Errorf("Edit(%q)\n"+
				"\tgot (product=%v, err=%v)\n"+
				"\twant (product=nil,err=error(non-nil))",
				product1.ID,
				searched, err)
			return
		}
		err1 := areProductsEqual(t, product1, searched, "returned product")
		if err1 != nil {
			return
		}
	})
	t.Run("negative test#2 (searches at non-existent id)", func(t *testing.T) {
		searched, err := service.Search("1111")
		if searched != nil || err == nil {
			t.Errorf("Edit(1111))\n"+
				"\tgot (product=%v, err=%v)\n"+
				"\twant (product=nil,err=error(non-nil))",
				searched, err)
			return
		}
		if err.Error() != "Product not found" {
			t.Errorf("Excpected error=Product not found, got=%q", err.Error())
		}
	})
}
func TestGetAll(t *testing.T) {
	product1, err := models.NewProduct("name1", "definition1", 100, "image1.png")
	if err != nil {
		t.Fatalf("Failed to create product1: %v", err)
	}
	product2, err := models.NewProduct("name2", "", 200, "")
	if err != nil {
		t.Fatalf("Failed to create product2: %v", err)
	}

	products := []*models.Product{
		product1,
		product2,
	}

	expectedSearchAll := map[string]*models.Product{
		product1.ID: product1,
		product2.ID: product2,
	}

	service, _, errorSetup := setupTestServiceWithData(t, products)
	if errorSetup != nil {
		t.Fatalf("Failed to create service and DB: %v", errorSetup)
	}
	t.Run("normal test#1", func(t *testing.T) {
		searchedall, err := service.GetAll()
		if len(searchedall) != 2 || err != nil {
			t.Errorf("GetAll(product1.ID))\n"+
				"\tgot (product.len=%d, err=%v)\n"+
				"\twant (product=%d,err=error(non-nil))",
				len(searchedall), err,
				2)
		}
		if len(searchedall) != 2 {
			t.Errorf("Expected 2 product, got %d", len(searchedall))
			return
		}
		actualMap := make(map[string]*models.Product)
		for _, p := range searchedall {
			actualMap[p.ID] = p
		}

		for id, product := range expectedSearchAll {
			actualProduct, exists := actualMap[id]
			if !exists {
				t.Errorf("Product with ID %q not found", id)
				continue
			}
			areProductsEqual(t, product, actualProduct, "GetAll result")
		}
	})
}
