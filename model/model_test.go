package model

import (
	"bufio"
	"net/http"
	"os"
	"testing"

	"github.com/evgeniron/API-Validator/store"
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

func TestGetModel(t *testing.T) {
	tests := []struct {
		data        *EndpointModel
		expectedErr bool
	}{
		{
			data: &EndpointModel{
				"path",
				"POST",
				map[string]FieldModel{
					"testing1": {
						[]string{"1", "2"},
						true,
					},
				},
				map[string]FieldModel{
					"testing2": {
						[]string{"3", "4"},
						true,
					},
				},
				map[string]FieldModel{
					"testing3": {
						[]string{"5", "6"},
						true,
					},
				},
			},
			expectedErr: false,
		},
	}
	for _, tt := range tests {
		db, err := store.NewInMemoryDB()
		require.NoError(t, err)
		key := generateKey(tt.data.Path, tt.data.Method)
		err = db.Insert(key, tt.data)
		require.NoError(t, err)

		result, err := GetModel(db, tt.data.Path, tt.data.Method)
		require.NoError(t, err)

		require.Equal(t, tt.data, result)
	}
}
