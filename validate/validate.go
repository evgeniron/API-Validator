package validate

import (
	"fmt"
	"net/http"

	"github.com/evgeniron/API-Validator/model"
	"github.com/evgeniron/API-Validator/utils"
)

const (
	ErrMismatchType         = "value type mismatch"
	ErrMissingRequiredField = "missing required field"
	ErrUnrecognizedField    = "unrecognized field"
)

type Field struct {
	Name  string
	Value interface{}
}

type Endpoint struct {
	Path        string  `json:"path"`
	Method      string  `json:"method"`
	QueryParams []Field `json:"query_params"`
	Headers     []Field `json:"headers"`
	Body        []Field `json:"body"`
}

type ValidationError struct {
	FieldName    string
	ErrorType    string
	ExpectedType string
	ErrorValue   interface{}
}

type ValidationReport struct {
	Path        string
	Method      string
	QueryParams []ValidationError
	Headers     []ValidationError
	Body        []ValidationError
}

func NewValidationReport() *ValidationReport {
	return &ValidationReport{}
}

func (vr *ValidationReport) WithPath(path string) *ValidationReport {
	if vr == nil {
		return nil
	}

	vr.Path = path
	return vr
}

func (vr *ValidationReport) WithMethod(method string) *ValidationReport {
	if vr == nil {
		return nil
	}

	vr.Method = method
	return vr
}

func (vr *ValidationReport) WithQueryParams(queryParams []ValidationError) *ValidationReport {
	if vr == nil {
		return nil
	}

	vr.QueryParams = queryParams
	return vr
}

func (vr *ValidationReport) WithHeaders(headers []ValidationError) *ValidationReport {
	if vr == nil {
		return nil
	}

	vr.Headers = headers
	return vr
}

func (vr *ValidationReport) WithBody(body []ValidationError) *ValidationReport {
	if vr == nil {
		return nil
	}

	vr.Body = body
	return vr
}

func ParseEndpoint(r *http.Request) (*Endpoint, error) {
	var endpoint Endpoint

	err := utils.Decode(r, endpoint)
	if err != nil {
		return nil, err
	}

	return &endpoint, nil
}

func ValidateReport(endpoint *Endpoint, model *model.EndpointModel) (*ValidationReport, error) {
	validationReport := NewValidationReport().
		WithPath(endpoint.Path).
		WithMethod(endpoint.Method).
		WithQueryParams(validateFields(endpoint.QueryParams, model.QueryParams)).
		WithHeaders(validateFields(endpoint.Headers, model.Headers)).
		WithBody(validateFields(endpoint.Body, model.Body))

	if validationReport == nil {
		return nil, fmt.Errorf("failed to construct report")
	}
	return validationReport, nil
}

func validateFieldType(inputField Field, types []string) []ValidationError {
	var validationErrors []ValidationError
	for _, filterType := range types {
		validatorFunc, validatorExists := Validators[filterType]
		if !validatorExists {
			// log error here
			continue
		}

		if !validatorFunc(inputField.Value) {
			validationErrors = append(validationErrors, ValidationError{inputField.Name, ErrMismatchType, filterType, inputField.Value})
		}
	}

	return validationErrors
}

func validateField(inputField Field, expectedFieldModels map[string]model.FieldModel) []ValidationError {
	var validationErrors []ValidationError

	expectedFieldModel, fieldModelExists := expectedFieldModels[inputField.Name]
	if !fieldModelExists {
		validationErrors = append(validationErrors, ValidationError{inputField.Name, ErrUnrecognizedField, "", nil})
		return validationErrors
	}

	validationErrors = append(validationErrors, validateFieldType(inputField, expectedFieldModel.Types)...)

	return validationErrors
}

func validateRequiredFields(expectedFieldModels map[string]model.FieldModel, existingFields map[string]struct{}) []ValidationError {
	var validationErrors []ValidationError
	for fieldName, fieldModel := range expectedFieldModels {
		if !fieldModel.Required {
			continue
		}

		_, fieldExists := existingFields[fieldName]
		if !fieldExists {
			validationErrors = append(validationErrors, ValidationError{fieldName, ErrMissingRequiredField, "", nil})
		}
	}
	return validationErrors
}

func validateFields(inputFields []Field, expectedFieldModels map[string]model.FieldModel) []ValidationError {
	var validationErrors []ValidationError

	existingFields := make(map[string]struct{}, len(inputFields))
	for _, inputField := range inputFields {
		existingFields[inputField.Name] = struct{}{}
		validationErrors = append(validationErrors, validateField(inputField, expectedFieldModels)...)
	}

	validationErrors = append(validationErrors, validateRequiredFields(expectedFieldModels, existingFields)...)
	return validationErrors
}
