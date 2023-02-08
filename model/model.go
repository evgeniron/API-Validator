package model

import (
	"fmt"
	"net/http"

	"github.com/evgeniron/API-Validator/utils"
)

type Field struct {
	Name     string   `json:"name"`
	Types    []string `json:"types"`
	Required bool     `json:"required"`
}

type Endpoint struct {
	Path        string  `json:"path"`
	Method      string  `json:"method"`
	QueryParams []Field `json:"query_params"`
	Headers     []Field `json:"headers"`
	Body        []Field `json:"body"`
}

type Store interface {
	Insert(key string, model *Endpoint) error
	Get(model string) (*Endpoint, error)
}

func Decode(r *http.Request) ([]Endpoint, error) {
	var endpoints []Endpoint

	err := utils.Decode(r, endpoints)
	if err != nil {
		return nil, err
	}

	return endpoints, nil
}

func StoreModels(db Store, models []Endpoint) error {
	for _, model := range models {
		err := StoreModel(db, &model)
		if err != nil {
			return err
		}
	}
	return nil
}

func generateKey(path, method string) string {
	return fmt.Sprintf("%s-%s", path, method)
}

func StoreModel(db Store, model *Endpoint) error {
	key := generateKey(model.Path, model.Method)
	return db.Insert(key, model)
}

func GetModel(db Store, path, method string) (interface{}, error) {
	key := generateKey(path, method)
	return db.Get(key)
}
