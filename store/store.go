package store

import (
	"github.com/evgeniron/API-Validator/model"
)

type inMemory map[string]model.Endpoint

type Store struct {
	db inMemory
}

func NewInMemoryDB() (*Store, error) {
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
		return nil, &RecordNotFoundError{Key: key}
	}
	return &record, nil
}
