package config

import (
	"os"
	"path/filepath"
	"testing"
)

type resultCase struct {
	dataConfig string
	errorMsg   string
}

func createFileConfig(filePath string, content string) error {
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return os.WriteFile(filePath, []byte(content), 0644)
}
func removeFileConfig(filePath string) error {
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return nil
	}
	return os.Remove(filePath)
}
func checkValidConfig(t *testing.T,
	config *Configuration,
	expected *resultCase,
	err error) {
	if config.DataBaseFilePath != expected.dataConfig {
		t.Errorf("Expected: config.pathFile=%q, got config.pathFile=%q;",
			expected.dataConfig,
			config.DataBaseFilePath)
		return
	}
}
func checkUnvalidConfig(t *testing.T,
	config *Configuration,
	expected *resultCase,
	err error) {
	if err.Error() != expected.errorMsg {
		t.Errorf("Expected: error=%q, got error=%q",
			err.Error(),
			expected.errorMsg)
	}
}
func TestNewConfiguration(t *testing.T) {
	expectedResult := []*resultCase{
		{"data/products.json", ""},
		{"", ""},
		{"", "unexpected end of JSON input"},
		{"", "open fileTestConfig.json: no such file or directory"},
		{"", "unexpected end of JSON input"},
		{"", "open : no such file or directory"},
	}
	testIn := []struct {
		name         string
		input        string
		content      string
		creatingFile bool
		wantConfig   bool
		wantErr      bool
		check        func(t *testing.T,
			config *Configuration,
			expected *resultCase,
			err error)
	}{
		{
			"normal test #1 (valid config)",
			"fileTestConfig.json",
			`{"DataBaseFilePath": "data/products.json"}`,
			true,
			true,
			false,
			checkValidConfig,
		},
		{
			"normal test #2 (empty object : missing field)",
			"fileTestConfig.json",
			`{"DataBaseFilePath": ""}`,
			true,
			true,
			false,
			checkValidConfig,
		},
		{
			"negative test #3 (empty file)",
			"fileTestConfig.json",
			"",
			true,
			false,
			true,
			checkUnvalidConfig,
		},
		{
			"negative test #4 (file not found)",
			"fileTestConfig.json",
			"",
			false,
			false,
			true,
			checkUnvalidConfig,
		},
		{
			"negative test #5 (invalid json)",
			"fileTestConfig.json",
			`{"DataBaseFilePath": "data/products.json"`,
			true,
			false,
			true,
			checkUnvalidConfig,
		},
		{
			"negative test #6 (empty path)",
			"",
			"",
			false,
			false,
			true,
			checkUnvalidConfig,
		},
	}
	var outConfig string = ""
	var outErr string = ""
	var indexTest int = 0
	for _, tt := range testIn {
		t.Run(tt.name, func(t *testing.T) {
			var errFileWork error
			if tt.creatingFile == true {
				errFileWork = createFileConfig(tt.input, tt.content)
			}
			if errFileWork != nil {
				t.Errorf("---- ERROR WITH FILE WORK")
				return
			}
			config, err := NewConfiguration(tt.input)
			if tt.creatingFile == true {
				errFileWork = removeFileConfig(tt.input)
			}
			if errFileWork != nil {
				t.Errorf("---- ERROR WITH FILE WORK")
				return
			}
			if (config != nil) != tt.wantConfig || (err != nil) != tt.wantErr {
				if tt.wantConfig == true {
					outConfig = "non-nil"
				} else {
					outConfig = "nil"
				}
				if tt.wantErr == true {
					outErr = "error(non-nil)"
				} else {
					outErr = "no error(nil)"
				}
				t.Errorf("NewConfiguration(%q)\n"+
					"\tgot : (Configuration=%q, error=%v)\n"+
					"\twant : (Configuration=%q, error=%v)",
					tt.input,
					config, err,
					outConfig, outErr)
				indexTest++
				return
			}
			if tt.check != nil {
				tt.check(t,
					config,
					expectedResult[indexTest],
					err)
			}
			indexTest++
		})
	}
}
