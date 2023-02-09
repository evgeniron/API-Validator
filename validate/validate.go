package validate

import (
	"fmt"
	"net/http"

	"github.com/evgeniron/API-Validator/model"
	"github.com/evgeniron/API-Validator/utils"
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

type ValidationReport struct {
	Path        string
	Method      string
	QueryParams []error
	Headers     []error
	Body        []error
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

func (vr *ValidationReport) WithQueryParams(queryParams []error) *ValidationReport {
	if vr == nil {
		return nil
	}

	vr.QueryParams = queryParams
	return vr
}

func (vr *ValidationReport) WithHeaders(headers []error) *ValidationReport {
	if vr == nil {
		return nil
	}

	vr.Headers = headers
	return vr
}

func (vr *ValidationReport) WithBody(body []error) *ValidationReport {
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

func validateFields(inputFields []Field, expectedFieldModels map[string]model.FieldModel) []error {
	var validationErrors []error

	existingFields := make(map[string]struct{})
	for _, inputField := range inputFields {
		existingFields[inputField.Name] = struct{}{}

		expectedFieldModel, fieldModelExists := expectedFieldModels[inputField.Name]
		if !fieldModelExists {
			validationErrors = append(validationErrors, fmt.Errorf("unrecognized field: %s", inputField.Name))
			continue
		}

		for _, filterType := range expectedFieldModel.Types {
			validatorFunc, validatorExists := Validators[filterType]
			if !validatorExists {
				validationErrors = append(validationErrors, fmt.Errorf("validator missing for type: %s", filterType))
				continue
			}

			if !validatorFunc(inputField.Value) {
				validationErrors = append(validationErrors, fmt.Errorf("value type mismatch for field: %s, expected: %s", inputField.Name, filterType))
			}
		}
	}

	validationErrors = append(validationErrors, validateRequiredFields(expectedFieldModels, existingFields)...)
	return validationErrors
}

func validateRequiredFields(expectedFieldModels map[string]model.FieldModel, existingFields map[string]struct{}) []error {
	var validationErrors []error
	for fieldName, fieldModel := range expectedFieldModels {
		if !fieldModel.Requried {
			continue
		}

		_, fieldExists := existingFields[fieldName]
		if !fieldExists {
			validationErrors = append(validationErrors, fmt.Errorf("required field missing: %s", fieldName))
		}
	}
	return validationErrors
}
