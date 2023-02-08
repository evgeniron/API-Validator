package store

import (
	"fmt"

	"github.com/evgeniron/API-Validator/model"
)

type inMemory map[string]model.Endpoint

type Store struct {
	db inMemory
}

func NewDB() (*Store, error) {
	db := make(inMemory, 0)
	return &Store{db: db}, nil
}

func (s *Store) Insert(key string, model *model.Endpoint) error {
	s.db[key] = *model
	return nil
}

func (s *Store) Get(key string) (*model.Endpoint, error) {
	record, ok := s.db[key]
	if !ok {
		return nil, fmt.Errorf("record not found: %s", key)
	}
	return &record, nil
}
