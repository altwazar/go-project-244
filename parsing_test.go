package code

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Тест функции CompareConfigs с разной комбинацией флагов.
func TestCompareConfigs(t *testing.T) {
	type tc struct {
		name       string
		file1      string
		file2      string
		format     string
		expOutPath string
		expErr     error
	}
	cases := []tc{
		{
			name:       "Two flat jsons 1 to 2 stylish",
			file1:      filepath.Join("testdata", "fixtures", "flatjson", "file1.json"),
			file2:      filepath.Join("testdata", "fixtures", "flatjson", "file2.json"),
			expOutPath: filepath.Join("testdata", "fixtures", "flatjson", "1to2expected.golden"),
			format:     "stylish",
			expErr:     nil,
		},
		{
			name:       "Two flat jsons 2 to 1 stylish",
			file1:      filepath.Join("testdata", "fixtures", "flatjson", "file2.json"),
			file2:      filepath.Join("testdata", "fixtures", "flatjson", "file1.json"),
			expOutPath: filepath.Join("testdata", "fixtures", "flatjson", "2to1expected.golden"),
			format:     "stylish",
			expErr:     nil,
		},
		{
			name:       "Two flat jsons 1 to 1 stylish",
			file1:      filepath.Join("testdata", "fixtures", "flatjson", "file1.json"),
			file2:      filepath.Join("testdata", "fixtures", "flatjson", "file1.json"),
			expOutPath: filepath.Join("testdata", "fixtures", "flatjson", "1to1expected.golden"),
			format:     "stylish",
			expErr:     nil,
		},
		{
			name:       "Two flat jsons 2 to 2 stylish",
			file1:      filepath.Join("testdata", "fixtures", "flatjson", "file2.json"),
			file2:      filepath.Join("testdata", "fixtures", "flatjson", "file2.json"),
			expOutPath: filepath.Join("testdata", "fixtures", "flatjson", "2to2expected.golden"),
			format:     "stylish",
			expErr:     nil,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			runCompareConfigs(
				t,
				c.file1,
				c.file2,
				c.expOutPath,
				c.format,
				c.expErr,
			)
		})
	}
}

func runCompareConfigs(
	t *testing.T,
	file1 string,
	file2 string,
	expOutPath string,
	format string,
	expErr error,
) {
	t.Helper()

	expOut, err := read(t, expOutPath)
	if err != nil {
		t.Fatalf("read %s: %v", expOutPath, err)
	}
	out, err := CompareConfigs(file1, file2, format)
	assert.Equal(t, expErr, err)
	assert.Equal(t, expOut, out)
}

// Для чтения ожидаемых результатов в строку
func read(t *testing.T, path string) (string, error) {
	t.Helper()
	b, err := os.ReadFile(path)
	return string(b), err
}
