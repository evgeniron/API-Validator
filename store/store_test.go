package store

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInsertAndGet(t *testing.T) {
	tests := []struct {
		key         string
		data        interface{}
		expectedErr bool
	}{
		{
			key:         "string",
			data:        "testing",
			expectedErr: false,
		},
		{
			key:         "int",
			data:        123,
			expectedErr: false,
		},
		{
			key:         "map",
			data:        map[string]int{"one": 1},
			expectedErr: false,
		},
		{
			key:         "struct",
			data:        struct{ name string }{name: "test"},
			expectedErr: false,
		},
		{
			key:         "",
			data:        "empty key",
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		db, err := NewInMemoryDB()
		require.NoError(t, err)

		err = db.Insert(tt.key, tt.data)
		if tt.expectedErr && err == nil {
			t.Fatalf("expected error, got: nil")
		}

		if !tt.expectedErr && err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		returnedData, err := db.Get(tt.key)
		if tt.expectedErr && err == nil {
			t.Fatalf("expected error, got: nil")
		}

		if !tt.expectedErr && err != nil {
			t.Fatalf("unexpected error, %v", err)
		}

		if !tt.expectedErr && !reflect.DeepEqual(tt.data, returnedData) {
			t.Fatalf("expected: %v, got: %v", tt.data, returnedData)
		}
	}
}
