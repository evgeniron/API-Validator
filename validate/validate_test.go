package validate

import (
	"testing"

	"github.com/evgeniron/API-Validator/model"
	"github.com/stretchr/testify/assert"
)

func TestValidateFieldType(t *testing.T) {
	// Test cases
	tests := []struct {
		name       string
		inputField Field
		types      []string
		expected   []ValidationError
	}{
		{
			name:       "Value Type Mismatch",
			inputField: Field{Name: "age", Value: "25"},
			types:      []string{"Int"},
			expected:   []ValidationError{{FieldName: "age", ErrorType: ErrMismatchType, ExpectedType: "Int", ErrorValue: "25"}},
		},
		{
			name:       "No validation errors",
			inputField: Field{Name: "email", Value: "test@example.com"},
			types:      []string{"string", "email"},
			expected:   []ValidationError{},
		},
		{
			name:       "Boolean value mismatch type",
			inputField: Field{Name: "flag", Value: "True"},
			types:      []string{"Boolean"},
			expected:   []ValidationError{{FieldName: "flag", ErrorType: ErrMismatchType, ExpectedType: "Boolean", ErrorValue: "True"}},
		},
		{
			name:       "No validation errors for boolean",
			inputField: Field{Name: "flag", Value: true},
			types:      []string{"Boolean"},
			expected:   []ValidationError{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call validateFieldType function with given input field and types
			actual := validateFieldType(tt.inputField, tt.types)

			// Check if the result matches with expected output
			if len(actual) != len(tt.expected) {
				t.Errorf("Expected %v but got %v", tt.expected, actual)
			}

			for i := 0; i < len(actual); i++ {
				if actual[i] != tt.expected[i] {
					t.Errorf("Expected %v but got %v", tt.expected[i], actual[i])
				}
			}
		})
	}
}

// ValidateField, validates the correctness of the fields according the types.
// It does not validate required fields.
func TestValidateField(t *testing.T) {
	tests := []struct {
		name           string
		inputField     Field
		expectedModels map[string]model.FieldModel
		expected       []ValidationError
	}{
		{
			name:       "Value Type Mismatch",
			inputField: Field{Name: "age", Value: "25"},
			expectedModels: map[string]model.FieldModel{
				"age":    {Types: []string{"Int"}, Required: true},
				"gender": {Types: []string{"String"}, Required: false},
			},
			expected: []ValidationError{{FieldName: "age", ErrorType: ErrMismatchType, ExpectedType: "Int", ErrorValue: "25"}},
		},
		{
			name:       "Only Value Type Mismatch",
			inputField: Field{Name: "age", Value: "25"},
			expectedModels: map[string]model.FieldModel{
				"age":    {Types: []string{"Int"}, Required: true},
				"gender": {Types: []string{"String"}, Required: true},
			},
			expected: []ValidationError{{FieldName: "age", ErrorType: ErrMismatchType, ExpectedType: "Int", ErrorValue: "25"}},
		},
		{
			name:       "No validation error",
			inputField: Field{Name: "email", Value: "test@example.com"},
			expectedModels: map[string]model.FieldModel{
				"email": {Types: []string{"string", "email"}},
			},
			expected: []ValidationError{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call validateField function with given input field and expected models
			actual := validateField(tt.inputField, tt.expectedModels)

			// Check if the result matches with expected output
			if len(actual) != len(tt.expected) {
				t.Errorf("Expected %v but got %v", tt.expected, actual)
			}

			for i := 0; i < len(actual); i++ {
				if actual[i] != tt.expected[i] {
					t.Errorf("Expected %v but got %v", tt.expected[i], actual[i])
				}
			}
		})
	}
}

func TestValidateRequiredFields(t *testing.T) {
	tests := []struct {
		name           string
		fieldModels    map[string]model.FieldModel
		existingFields map[string]struct{}
		expected       []ValidationError
	}{
		{
			name: "All required fields are present",
			fieldModels: map[string]model.FieldModel{
				"field1": {Required: true},
				"field2": {Required: true},
				"field3": {Required: false},
			},
			existingFields: map[string]struct{}{
				"field1": {},
				"field2": {},
			},
			expected: []ValidationError{},
		},
		{
			name: "A required field is missing",
			fieldModels: map[string]model.FieldModel{
				"field1": {Types: []string{"Int"}, Required: true},
				"field2": {Types: []string{"Int"}, Required: true},
				"field3": {Types: []string{"Int"}, Required: false},
			},
			existingFields: map[string]struct{}{
				"field1": {},
			},
			expected: []ValidationError{
				{FieldName: "field2", ErrorType: ErrMissingRequiredField, ExpectedType: "", ErrorValue: nil},
			},
		},
		{
			name: "All required fields are missing",
			fieldModels: map[string]model.FieldModel{
				"field1": {Required: true},
				"field2": {Required: true},
				"field3": {Required: true},
			},
			existingFields: map[string]struct{}{},
			expected: []ValidationError{
				{FieldName: "field1", ErrorType: ErrMissingRequiredField, ExpectedType: "", ErrorValue: nil},
				{FieldName: "field2", ErrorType: ErrMissingRequiredField, ExpectedType: "", ErrorValue: nil},
				{FieldName: "field3", ErrorType: ErrMissingRequiredField, ExpectedType: "", ErrorValue: nil},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := validateRequiredFields(tt.fieldModels, tt.existingFields)

			// Check if the result matches with expected output
			if len(actual) != len(tt.expected) {
				t.Errorf("Expected %v but got %v", tt.expected, actual)
			}

			results := make(map[string]ValidationError, len(actual))
			for _, element := range actual {
				results[element.FieldName] = element
			}

			for _, element := range tt.expected {
				actualResult, ok := results[element.FieldName]
				if !ok {
					t.Errorf("Expected %v but it is missing", element)
				}
				if actualResult != element {
					t.Errorf("Expected %v but got %v", element, actualResult)
				}
			}
		})
	}
}

func TestValidateFields(t *testing.T) {
	tests := []struct {
		inputFields         []Field
		expectedFieldModels map[string]model.FieldModel
		expectedErrors      []ValidationError
	}{
		{
			inputFields: []Field{
				{Name: "field1", Value: "hello world"},
				{Name: "field2", Value: 42},
				{Name: "field3", Value: true},
			},
			expectedFieldModels: map[string]model.FieldModel{
				"field1": {Types: []string{"String"}, Required: true},
				"field2": {Types: []string{"Int"}, Required: true},
				"field3": {Types: []string{"Boolean"}, Required: true},
				"field4": {Types: []string{"UUID"}, Required: true},
			},
			expectedErrors: []ValidationError{
				{FieldName: "field4", ErrorType: ErrMissingRequiredField, ExpectedType: "", ErrorValue: nil},
			},
		},
		{
			inputFields: []Field{
				{Name: "field1", Value: 42},
				{Name: "field2", Value: true},
				{Name: "field3", Value: "hello world"},
			},
			expectedFieldModels: map[string]model.FieldModel{
				"field1": {Types: []string{"String"}, Required: true},
				"field2": {Types: []string{"Int"}, Required: true},
				"field3": {Types: []string{"Boolean"}, Required: true},
			},
			expectedErrors: []ValidationError{
				{FieldName: "field1", ErrorType: ErrMismatchType, ExpectedType: "String", ErrorValue: 42},
				{FieldName: "field2", ErrorType: ErrMismatchType, ExpectedType: "Int", ErrorValue: true},
				{FieldName: "field3", ErrorType: ErrMismatchType, ExpectedType: "Boolean", ErrorValue: "hello world"},
			},
		},
		{
			inputFields: []Field{
				{Name: "field1", Value: 42},
				{Name: "field2", Value: true},
				{Name: "field3", Value: "hello world"},
			},
			expectedFieldModels: map[string]model.FieldModel{
				"field1": {Types: []string{"Int"}, Required: true},
				"field2": {Types: []string{"Boolean"}, Required: true},
				"field3": {Types: []string{"String", "Email"}, Required: true},
			},
			expectedErrors: nil,
		},
	}

	for _, test := range tests {
		validationErrors := validateFields(test.inputFields, test.expectedFieldModels)
		assert.Equal(t, test.expectedErrors, validationErrors)
	}
}
