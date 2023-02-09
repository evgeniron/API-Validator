package model

import (
	"fmt"
	"net/http"

	"github.com/evgeniron/API-Validator/utils"
)

type Store interface {
	Insert(key string, model interface{}) error
	Get(model string) (interface{}, error)
}
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

type FieldModel struct {
	Types    []string
	Requried bool
}

type EndpointModel struct {
	Path        string                `json:"path"`
	Method      string                `json:"method"`
	QueryParams map[string]FieldModel `json:"query_params"`
	Headers     map[string]FieldModel `json:"headers"`
	Body        map[string]FieldModel `json:"body"`
}

func Decode(r *http.Request) ([]Endpoint, error) {
	var endpoints []Endpoint

	err := utils.Decode(r, &endpoints)
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

	endpointModel := NewModel().
		WithPath(model.Path).
		WithMethod(model.Method).
		WithQueryParams(model.QueryParams).
		WithHeaders(model.Headers).
		WithBody(model.Body)

	return db.Insert(key, endpointModel)
}

func NewModel() *EndpointModel {
	return &EndpointModel{}
}

func (em *EndpointModel) WithPath(path string) *EndpointModel {
	if em == nil {
		return nil
	}

	em.Path = path
	return em
}

func (em *EndpointModel) WithMethod(method string) *EndpointModel {
	if em == nil {
		return nil
	}

	em.Method = method
	return em
}

func (em *EndpointModel) WithQueryParams(queryParams []Field) *EndpointModel {
	if em == nil {
		return nil
	}

	em.QueryParams = rawFeildToFieldModel(queryParams)
	return em
}

func (em *EndpointModel) WithHeaders(headers []Field) *EndpointModel {
	if em == nil {
		return nil
	}

	em.Headers = rawFeildToFieldModel(headers)
	return em
}

func (em *EndpointModel) WithBody(body []Field) *EndpointModel {
	if em == nil {
		return nil
	}

	em.Body = rawFeildToFieldModel(body)
	return em
}

func GetModel(db Store, path, method string) (*EndpointModel, error) {
	key := generateKey(path, method)
	record, err := db.Get(key)
	if err != nil {
		return nil, err
	}

	endpointModel, ok := record.(*EndpointModel)
	if !ok {
		return nil, fmt.Errorf("incorrect type - expecting EndpointModel type")
	}

	return endpointModel, nil
}

func rawFeildToFieldModel(fields []Field) map[string]FieldModel {
	fieldModel := make(map[string]FieldModel, len(fields))
	for _, field := range fields {
		types := make([]string, len(field.Types))
		copy(types, field.Types)
		fieldModel[field.Name] = FieldModel{
			Types:    field.Types,
			Requried: field.Required,
		}
	}
	return fieldModel
}
