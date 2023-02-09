package store

import (
	"testing"

	"github.com/evgeniron/API-Validator/model"
	"github.com/stretchr/testify/require"
)

func TestStore(t *testing.T) {
	db, err := NewInMemoryDB()
	require.NoError(t, err)

	err = db.Insert("test", &model.Endpoint{Path: "path/test"})
	require.NoError(t, err)

	result, err := db.Get("test")
	require.NoError(t, err)
	require.Equal(t, "path/test", result.Path)
}
