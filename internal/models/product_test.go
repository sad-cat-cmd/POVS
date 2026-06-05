package models

import (
	"strings"
	"testing"
	"time"
)

type resultCase struct {
	productName string
	definition  string
	price       float64
	image       string
	wantErrMsg  string
}

var used_ID []string

func TestGenereteID(t *testing.T) {
	id1 := genereteID()
	if id1 == "" {
		t.Error("Generated ID returned empty string")
	}
	for _, ch := range id1 {
		if ch < '0' || ch > '9' {
			t.Errorf("ID contains non-digit character: %c", ch)
		}
	}
	id2 := genereteID()
	if id1 == id2 {
		t.Error("GenerateID returned the same value twice")
	}
}

func TestCheckName(t *testing.T) {
	test := []struct {
		name    string
		input   string
		wantErr bool
		check   func(t *testing.T, err error)
	}{
		// позитивный тест
		{
			name:    "normal test #1 (book)",
			input:   "book",
			wantErr: false,
			check:   nil,
		},
		// негативный тест
		{
			name:    "negative test #2 ()(empty name)",
			input:   "",
			wantErr: true,
			check: func(t *testing.T, err error) {
				if err != nil && err.Error() != "name is empty" {
					t.Errorf("Expected error: 'name is empty', got %q", err.Error())
					return
				}
			},
		},
		// негативный тест
		{
			name:    "negative test #3 ()(size name > 50)",
			input:   strings.Repeat("a", 51),
			wantErr: true,
			check: func(t *testing.T, err error) {
				if err != nil && err.Error() != "length name > 50" {
					t.Errorf("Excpected error: 'length name > 50', got '%q'", err.Error())
					return
				}
			},
		},
	}
	var errStr = ""
	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			err := checkName(tt.input)

			if (err != nil) != tt.wantErr {
				if tt.wantErr == true {
					errStr = "error (non-nil)"
				} else {
					errStr = "no error (nil)"
				}
				t.Errorf("checkName(%q)\n,"+
					"\tgot (error=%v)\n"+
					"\twant (error=%q)",
					tt.input,
					err,
					errStr)
				return
			}
			if tt.check != nil {
				tt.check(t, err)
			}
		})
	}
}

func TestCheckDefinition(t *testing.T) {
	test := []struct {
		name    string
		input   string
		wantStr string
		wantErr bool
		check   func(t *testing.T, str string, err error)
	}{
		// позитивный тест
		{
			name:    "normal test #1 (book about programming)",
			input:   "book about programming",
			wantStr: "book about programming",
			wantErr: false,
			check:   nil,
		},
		// позитивный тест
		{
			name:    "normal test #2 ()",
			input:   "",
			wantStr: "no data",
			wantErr: false,
			check:   nil,
		},
		// негативный тест
		{
			name:    "negative test #3 (aa...a)(difinition length > 500)",
			input:   strings.Repeat("a", 501),
			wantStr: "",
			wantErr: true,
			check: func(t *testing.T, str string, err error) {
				if err != nil && err.Error() != "length difinition > 500" {
					t.Errorf("Excpected error: 'length definition > 500', got %q", err.Error())
					return
				}
			},
		},
	}
	var outErr string = ""
	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			str, err := checkDefinition(tt.input)
			if (err != nil) != tt.wantErr || str != tt.wantStr {
				if tt.wantErr == true {
					outErr = "error (non-nil)"
				} else {
					outErr = "no error (nil)"
				}
				t.Errorf("checkDefinition(%q):\n"+
					"\tgot: (string=%q, error=%v)\n"+
					"\twant: (string=%q, error=%q)",
					tt.input,
					str, err,
					tt.wantStr, outErr)
				return
			}
			if tt.check != nil {
				tt.check(t, str, err)
			}
		})
	}
}

func TestCheckPrice(t *testing.T) {
	test := []struct {
		name    string
		input   float64
		wantErr bool
		check   func(t *testing.T, err error)
	}{
		// позитивный тест
		{
			name:    "normal test #1 (price = 700.00)",
			input:   700.00,
			wantErr: false,
			check:   nil,
		},
		// позитивный тест
		{
			name:    "normal test #2 (price == 500)",
			input:   500,
			wantErr: false,
			check:   nil,
		},
		// негативный тест
		{
			name:    "negitive test #3 (price = 0) (null price)",
			input:   0,
			wantErr: true,
			check: func(t *testing.T, err error) {
				if err.Error() != "value price <= 0" {
					t.Errorf("Excpected error: 'value price <= 0', got %q", err.Error())
					return
				}
			},
		},
		// негативный тест
		{
			name:    "negitive test #4 (price == -100) (price < 0)",
			input:   -100,
			wantErr: true,
			check: func(t *testing.T, err error) {
				if err.Error() != "value price <= 0" {
					t.Errorf("Excpected error: 'value price <= 0', got %q", err.Error())
					return
				}
			},
		},
		// негативный тест
		{
			name:    "negitive test #5 (price == 100000000000.01) (price > maxValue)",
			input:   100000000000.01,
			wantErr: true,
			check: func(t *testing.T, err error) {
				if err.Error() != "value price > 100000000000" {
					t.Errorf("Excpected error: 'value price > 100000000000', got %q", err.Error())
					return
				}
			},
		},
		// негативный тест
		{
			name:    "negitive test #6 (price == 100000000002) (price > maxValue)",
			input:   100000000002,
			wantErr: true,
			check: func(t *testing.T, err error) {
				if err.Error() != "value price > 100000000000" {
					t.Errorf("Excpected error: 'value price > 100000000000', got %q", err.Error())
					return
				}
			},
		},
	}
	var outErr string = ""
	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			err := checkPrice(tt.input)
			if (err != nil) != tt.wantErr {
				if tt.wantErr == true {
					outErr = "error (non-nill)"
				} else {
					outErr = "no error (nil)"
				}
				t.Errorf("checkPrice(%f)\n +"+
					"\tgot: (error=%v)\n"+
					"\twant: (error=%q)",
					tt.input,
					err,
					outErr)
				return
			}
			if tt.check != nil {
				tt.check(t, err)
			}
		})
	}
}

func TestCheckImage(t *testing.T) {
	test := []struct {
		name    string
		input   string
		wantStr string
		wantErr bool
		check   func(t *testing.T, str string, err error)
	}{
		// позитивный случай
		{
			name:    "normal test #1 (image.png)",
			input:   "image.png",
			wantStr: "image.png",
			wantErr: false,
			check:   nil,
		},
		// позитивный случай
		{
			name:    "normal test #2 (image.jpeg)",
			input:   "image.jpeg",
			wantStr: "image.jpeg",
			wantErr: false,
			check:   nil,
		},
		// позитивный случай
		{
			name:    "normal test #3 (image.jpg)",
			input:   "image.jpg",
			wantStr: "image.jpg",
			wantErr: false,
			check:   nil,
		},
		// позитивный случай
		{
			name:    "normal test #4 (12image12.png)",
			input:   "12image12.png",
			wantStr: "12image12.png",
			wantErr: false,
			check:   nil,
		},
		// позитивный случай
		{
			name:    "normal test #5 (изображение_12.png)",
			input:   "изображение_12.png",
			wantStr: "изображение_12.png",
			wantErr: false,
			check:   nil,
		},
		{
			name:    "normal test #6 (изображение-img.png)",
			input:   "изображение-img.png",
			wantStr: "изображение-img.png",
			wantErr: false,
			check:   nil,
		},
		// позитивный случай
		{
			name:    "normal test #7 ()(empty string)",
			input:   "",
			wantStr: "no data",
			wantErr: false,
			check:   nil,
		},
		// негативный случай (image.exe)
		{
			name:    "negative test #8 (image.exe)(uncorrect extensions)",
			input:   "image.exe",
			wantStr: "",
			wantErr: true,
			check: func(t *testing.T, str string, err error) {
				if err.Error() != "image contains unvalid extensions" {
					t.Errorf("Excpected 'image contains unvalid extensions', got %q", err.Error())
				}
			},
		},
		// негативный случай (image.exe.png)(multiple extensions)
		{
			name:    "negative test #9 (image.exe.png)(multiple extensions)",
			input:   "image.exe.png",
			wantStr: "",
			wantErr: true,
			check: func(t *testing.T, str string, err error) {
				if err.Error() != "image contains multiple extensions" {
					t.Errorf("Excpected 'image contains multiple extensions', got %q", err.Error())
				}
			},
		},
		// негативный случай (image length > 255)
		{
			name:    "negative test #10 (aa...a)(image length > 255)",
			input:   strings.Repeat("a", 256),
			wantStr: "",
			wantErr: true,
			check: func(t *testing.T, str string, err error) {
				if err.Error() != "length image > 255" {
					t.Errorf("Excpected 'length image > 255', got %q", err.Error())
				}
			},
		},
		// негативный случай (image contains unvalid chars)
		{
			name:    "negative test #11 (image{}.png)(image contains unvalid chars)",
			input:   "image{}.png",
			wantStr: "",
			wantErr: true,
			check: func(t *testing.T, str string, err error) {
				if err.Error() != "image contains unvalid chars" {
					t.Errorf("Excpected 'image contains unvalid chars', got %q", err.Error())
				}
			},
		},
	}

	var outErr string = ""
	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			str, err := checkImage(tt.input)
			if (err != nil) != tt.wantErr || str != tt.wantStr {
				if tt.wantErr == true {
					outErr = "error (non-nill)"
				} else {
					outErr = "no error (nil)"
				}
				t.Errorf("checkImage(%q):\n"+
					"\tgot: (string=%q, error=%v)\n"+
					"\t,want (string=%q, error=%q)",
					tt.input,
					str, err,
					tt.wantStr, outErr)
				return
			}
			if tt.check != nil {
				tt.check(t, str, err)
			}
		})
	}
}

func checkValidProduct(t *testing.T, p *Product, err error, res *resultCase) {
	if p.Name != res.productName {
		t.Errorf("Expected name : '%q', got name: '%q'",
			res.productName, p.Name)
		return
	}
	if p.Definition != res.definition {
		t.Errorf("Expected definition: '%q', got definition : '%q'",
			res.definition, p.Definition)
		return
	}
	if p.Price != res.price {
		t.Errorf("Expected price: %f, got price: %f",
			res.price, p.Price)
	}
	if p.Image != res.image {
		t.Errorf("Expected image: '%q', got image: '%q'",
			res.image, p.Image)
	}
	for _, id := range used_ID {
		if p.ID == id {
			t.Errorf("Created id already exists: '%q'", id)
			return
		}
	}
	used_ID = append(used_ID, p.ID)

	diff := p.UpdatedAt.Sub(p.CreateAt)
	if diff < 0 {
		diff = -diff
	}
	if diff > time.Second {
		t.Errorf("CreateAt and UpdatedAt differ by %v (more than 1 second)", diff)
	}
}
func CheckUnvalidProduct(t *testing.T, p *Product, err error, res *resultCase) {
	if err == nil {
		t.Errorf("Expected error, got nil")
		return
	}
	if err.Error() != res.wantErrMsg {
		t.Errorf("Expected error string: '%q' , got %q",
			res.wantErrMsg, err.Error())
		return
	}
}
func TestNewProduct(t *testing.T) {
	excepected_result := []*resultCase{
		// нормальный случай #1
		{
			"book", "book definition", 700, "image.png", "",
		},
		// нормальный случай (все доступные поля пустые) #2
		{
			"book", "no data", 700.0, "no data", "",
		},
		// нормальный случай (вариант image) #3
		{
			"book", "no data", 700, "image123-изображение352.png", "",
		},
		// нормальный случай (вариант image) #4
		{
			"book", "no data", 700, "image123-изображение.jpeg", "",
		},
		// нормальный случай (вариант image) #5
		{
			"book", "no data", 700, "image123-изображение.jpg", "",
		},
		// негативный случай #6
		// name is empty
		// value price <= 0
		// multiple extenstion
		{
			"", "", 0, "", "name is empty;value price <= 0;image contains multiple extensions;",
		},
		// негативный случай #7
		// name is empty
		// image with unvalid chars
		{
			"", "", 0, "", "name is empty;image contains unvalid chars;",
		},
		// негативный случай #8
		// all values > max valid value
		{
			"", "", 0, "", "length name > 50;length difinition > 500;value price > 100000000000;length image > 255;",
		},
		// негативный случай #9
		// name is empty
		// image with unvalid extensions
		{
			"", "", 0, "", "name is empty;image contains unvalid extensions;",
		},
	}
	in_test := []struct {
		name             string
		productName      string
		definition       string
		price            float64
		image            string
		wantValidProduct bool
		wantErr          bool
		checkFunc        func(t *testing.T, p *Product, err error, res *resultCase)
	}{
		{
			"normal test #1", "book", "book definition", 700, "image.png", true, false,
			checkValidProduct,
		},
		{
			"normal test #2", "book", "", 700.0, "", true, false, checkValidProduct,
		},
		{
			"normal test #3", "book", "", 700, "image123-изображение352.png", true, false, checkValidProduct,
		},
		{
			"normal test #4", "book", "", 700, "image123-изображение.jpeg", true, false, checkValidProduct,
		},
		{
			"normal test #5", "book", "", 700, "image123-изображение.jpg", true, false, checkValidProduct,
		},
		{
			"negative test #6", "", "definition", -100, "изображение.exe.rar", false, true, CheckUnvalidProduct,
		},
		{
			"negative test #7", "", "definition", 100, "изображение@1%.png", false, true, CheckUnvalidProduct,
		},
		{
			"negitive test #8", strings.Repeat("a", 51), strings.Repeat("b", 501), 100000000001,
			strings.Repeat("c", 256), false, true, CheckUnvalidProduct,
		},
		{
			"negative test #9", "", "difinition", 100, "image.exe", false, true, CheckUnvalidProduct,
		},
	}
	var index int = 0
	var stringOutProduct string = ""
	var stringOutErr string = ""

	for _, tt := range in_test {
		t.Run(tt.name, func(t *testing.T) {
			prod, err := NewProduct(tt.productName,
				tt.definition,
				tt.price,
				tt.image)
			if (err != nil) != tt.wantErr || (prod != nil) != tt.wantValidProduct {
				if tt.wantErr == true {
					stringOutErr = "error (non-nil)"
				} else {
					stringOutErr = "no error (nil)"
				}
				if tt.wantValidProduct == true {
					stringOutProduct = "non-nil Product"
				} else {
					stringOutProduct = "nil"
				}
				t.Errorf("NewProduct(%q, %q, %f, %q):\n"+
					"\tgot: (product=%v , error=%v)\n"+
					"\twant: (product=%q , error=%q)",
					tt.name, tt.definition, tt.price, tt.image,
					prod, err,
					stringOutProduct, stringOutErr)
				index++
				return
			}
			if tt.checkFunc != nil {
				tt.checkFunc(t, prod, err, excepected_result[index])
			}
			index++
		})
	}
}
