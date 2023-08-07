package aisleriot

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConstructorWithBogusFilename(t *testing.T) {
	_, err := NewDataProvider("Bogus")
	assert.NotNil(t, err)
}

func TestConstructorWithDefault(t *testing.T) {
	pdp, err := NewDataProvider()
	_ = pdp
	_ = err
}

func TestDefaultFileName(t *testing.T) {
	userId, err := user.Current()
	username := userId.Username
	assert.Nil(t, err)
	filename := DefaultFileName()
	assert.NotNil(t, filename)
	assert.NotEmpty(t, filename)
	assert.Contains(t, filename, username)
}

func TestParseData(t *testing.T) {
	stoogeFile := filepath.Join("testdata", "stooges.ini")
	stooges, err := os.ReadFile(stoogeFile)
	assert.Nil(t, err)
	tests := []struct {
		name string
		data []byte
		want map[string]map[string]string
	}{
		{"stooges", stooges, map[string]map[string]string{
			"Moe": {
				"rank":   "1",
				"saying": "Why, I oughta...",
			},
			"Larry": {
				"rank":   "2",
				"saying": "Hey, Moe!",
			},
			"Curly": {
				"rank":   "3",
				"saying": "Nyuk, nyuk, nyuk",
			},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sections, err := ParseData(tt.data)
			assert.Nil(t, err)
			assert.Equal(t, tt.want, sections)
		})
	}
}

func TestDataProvider_GameList(t *testing.T) {
	tests := []struct {
		name      string
		filename  string
		want      []string
		wantError bool
	}{
		{"Normal file",
			filepath.Join("testdata", "aisleriot"),
			[]string{"spider", "freecell", "canfield", "klondike"},
			false},
		{"Different .ini",
			filepath.Join("testdata", "stooges.ini"),
			nil,
			false},
		{"Non-existent file",
			filepath.Join("testdata", "non-existent.ini"),
			nil,
			true},
		{"Malformed file",
			filepath.Join("testdata", "bogus.ini"),
			nil,
			true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pdp, err := NewDataProvider(tt.filename)
			if tt.wantError {
				assert.NotNil(t, err, fmt.Sprint(err))
			} else {
				actual := pdp.GameList()
				expected := tt.want
				assert.Equal(t, expected, actual)
			}
		})
	}
}
