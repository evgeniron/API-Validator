package utils

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDecode(t *testing.T) {
	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	tests := []struct {
		name        string
		json        string
		expected    Person
		expectError bool
	}{
		{
			name:     "Valid JSON",
			json:     `{"name": "Gaia", "age": 2}`,
			expected: Person{Name: "Gaia", Age: 2},
		},
		{
			name:        "Invalid JSON",
			json:        `{"name": "Jonathan", "age": "invalid"}`,
			expected:    Person{Name: "Jonathan"},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, "/", bytes.NewBuffer([]byte(tt.json)))
			require.NoError(t, err)

			var p Person
			err = Decode(req, &p)
			if !tt.expectError {
				require.NoError(t, err)
			}
			if tt.expected.Name != p.Name {
				t.Errorf("expected name %q but got %q", tt.expected.Name, p.Name)
			}
			if tt.expected.Age != p.Age {
				t.Errorf("expected age %d but got %d", tt.expected.Age, p.Age)
			}
		})
	}
}
