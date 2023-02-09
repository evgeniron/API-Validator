package model

import (
	"bufio"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStoreModels(t *testing.T) {
	tests := []struct {
		name string
		file string
	}{
		{
			name: "happy",
			file: "test_data/models.json",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// make test w and r
			file, err := os.Open(tt.file)
			require.NoError(t, err)
			r, err := http.NewRequest(http.MethodPut, "/", bufio.NewReader(file))
			require.NoError(t, err)

			models, err := Decode(r)
			require.NoError(t, err)

			require.Equal(t, 5, len(models))
			require.Equal(t, "UUID", models[0].QueryParams[1].Types[1])
		})
	}
}
