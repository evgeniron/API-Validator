package store

import "fmt"

/*
In real life application I'd use document database. Document databases have better schema flexibility,
better performance due to locality, and closer to data structures in applications. While Relational databases have better
support for joins, many to one and many to many relationships.

Our doucment is not deeply nested so it's not difficult to directly access it. Also, the whole document is required
so there is a performance advantage to this storage locality (single lookup, no joins required)
*/

type inMemory map[string]interface{}

type Store struct {
	db inMemory
}

func NewInMemoryDB() (*Store, error) {
	db := make(inMemory, 10)
	return &Store{db: db}, nil
}

func (s *Store) Insert(key string, model interface{}) error {
	if key == "" {
		return fmt.Errorf("empty key")
	}
	s.db[key] = model
	return nil
}

func (s *Store) Get(key string) (interface{}, error) {
	record, ok := s.db[key]
	if !ok {
		return nil, &RecordNotFoundError{Key: key}
	}
	return record, nil
}
