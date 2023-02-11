package main

import (
	"bufio"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/evgeniron/API-Validator/store"
	"github.com/evgeniron/API-Validator/validate"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServer(t *testing.T) {
	db, err := store.NewInMemoryDB()
	require.NoError(t, err)
	srv := NewServer(db)

	modelFile, err := os.Open("test_data/models.json")
	require.NoError(t, err)
	w := httptest.NewRecorder()

	reqPutModel, err := http.NewRequest(http.MethodPut, modelRoute, bufio.NewReader(modelFile))
	require.NoError(t, err)
	srv.ServeHTTP(w, reqPutModel)
	require.Equal(t, http.StatusOK, w.Result().StatusCode)
	w.Result().Body.Close()

	tests := []struct {
		name            string
		endpointReqPath string
		status          int
		isValid         bool
	}{
		{
			name:            "Valid endpoint",
			endpointReqPath: "test_data/endpointReqValid.json",
			status:          http.StatusOK,
			isValid:         true,
		},
		{
			name:            "Invalid endpoint",
			endpointReqPath: "test_data/endpointReqInvalid.json",
			status:          http.StatusOK,
			isValid:         false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			endpoint, err := os.Open(tt.endpointReqPath)
			require.NoError(t, err)

			reqModel, err := http.NewRequest(http.MethodPost, validationRoute, bufio.NewReader(endpoint))
			require.NoError(t, err)

			w := httptest.NewRecorder()
			srv.ServeHTTP(w, reqModel)
			resp := w.Result()
			defer resp.Body.Close()
			require.Equal(t, tt.status, resp.StatusCode)

			var result validate.ValidationReport
			decoder := json.NewDecoder(resp.Body)
			err = decoder.Decode(&result)
			require.NoError(t, err)
			require.Equal(t, tt.isValid, result.Valid)
		})
	}
}

func TestRespond(t *testing.T) {
	tests := []struct {
		status int
		data   interface{}
		result string
	}{
		{
			status: http.StatusOK,
			data:   map[string]string{"message": "Hello, World!"},
			result: `{"message":"Hello, World!"}`,
		},
		{
			status: http.StatusBadRequest,
			data:   map[string]string{"error": "Invalid request"},
			result: `{"error":"Invalid request"}`,
		},
		{
			status: http.StatusOK,
			data: validate.ValidationReport{
				Path:   "/users/info",
				Method: "GET",
				QueryParams: []validate.ValidationError{
					{FieldName: "QueryParam1", ErrorType: validate.ErrMismatchType, ExpectedType: "String", ErrorValue: 42},
					{FieldName: "QueryParam2", ErrorType: validate.ErrMismatchType, ExpectedType: "Int", ErrorValue: true},
					{FieldName: "QueryParam3", ErrorType: validate.ErrMismatchType, ExpectedType: "Boolean", ErrorValue: "hello world"},
				},
				Headers: []validate.ValidationError{
					{FieldName: "Header1", ErrorType: validate.ErrMismatchType, ExpectedType: "String", ErrorValue: 42},
					{FieldName: "Header2", ErrorType: validate.ErrMismatchType, ExpectedType: "Int", ErrorValue: true},
					{FieldName: "Header3", ErrorType: validate.ErrMismatchType, ExpectedType: "Boolean", ErrorValue: "hello world"},
				},
				Body: []validate.ValidationError{
					{FieldName: "Body1", ErrorType: validate.ErrMismatchType, ExpectedType: "String", ErrorValue: 42},
					{FieldName: "Body2", ErrorType: validate.ErrMismatchType, ExpectedType: "Int", ErrorValue: true},
					{FieldName: "Body3", ErrorType: validate.ErrMismatchType, ExpectedType: "Boolean", ErrorValue: "hello world"},
				},
				Valid: false,
			},
			result: `{"Path":"/users/info","Method":"GET","QueryParams":[{"FieldName":"QueryParam1","ErrorType":"value type mismatch","ExpectedType":"String","ErrorValue":42},{"FieldName":"QueryParam2","ErrorType":"value type mismatch","ExpectedType":"Int","ErrorValue":true},{"FieldName":"QueryParam3","ErrorType":"value type mismatch","ExpectedType":"Boolean","ErrorValue":"hello world"}],"Headers":[{"FieldName":"Header1","ErrorType":"value type mismatch","ExpectedType":"String","ErrorValue":42},{"FieldName":"Header2","ErrorType":"value type mismatch","ExpectedType":"Int","ErrorValue":true},{"FieldName":"Header3","ErrorType":"value type mismatch","ExpectedType":"Boolean","ErrorValue":"hello world"}],"Body":[{"FieldName":"Body1","ErrorType":"value type mismatch","ExpectedType":"String","ErrorValue":42},{"FieldName":"Body2","ErrorType":"value type mismatch","ExpectedType":"Int","ErrorValue":true},{"FieldName":"Body3","ErrorType":"value type mismatch","ExpectedType":"Boolean","ErrorValue":"hello world"}],"Valid":false}`,
			/*
							{
				    "Path": "/users/info",
				    "Method": "GET",
				    "QueryParams": [
				        {
				            "FieldName": "QueryParam1",
				            "ErrorType": "value type mismatch",
				            "ExpectedType": "String",
				            "ErrorValue": 42
				        },
				        {
				            "FieldName": "QueryParam2",
				            "ErrorType": "value type mismatch",
				            "ExpectedType": "Int",
				            "ErrorValue": true
				        },
				        {
				            "FieldName": "QueryParam3",
				            "ErrorType": "value type mismatch",
				            "ExpectedType": "Boolean",
				            "ErrorValue": "hello world"
				        }
				    ],
				    "Headers": [
				        {
				            "FieldName": "Header1",
				            "ErrorType": "value type mismatch",
				            "ExpectedType": "String",
				            "ErrorValue": 42
				        },
				        {
				            "FieldName": "Header2",
				            "ErrorType": "value type mismatch",
				            "ExpectedType": "Int",
				            "ErrorValue": true
				        },
				        {
				            "FieldName": "Header3",
				            "ErrorType": "value type mismatch",
				            "ExpectedType": "Boolean",
				            "ErrorValue": "hello world"
				        }
				    ],
				    "Body": [
				        {
				            "FieldName": "Body1",
				            "ErrorType": "value type mismatch",
				            "ExpectedType": "String",
				            "ErrorValue": 42
				        },
				        {
				            "FieldName": "Body2",
				            "ErrorType": "value type mismatch",
				            "ExpectedType": "Int",
				            "ErrorValue": true
				        },
				        {
				            "FieldName": "Body3",
				            "ErrorType": "value type mismatch",
				            "ExpectedType": "Boolean",
				            "ErrorValue": "hello world"
				        }
				    ]
				}*/
		},
	}

	for _, test := range tests {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		respond(w, r, test.status, test.data)
		assert.Equal(t, test.status, w.Code)
		assert.Equal(t, test.result, strings.Replace(w.Body.String(), "\n", "", -1))
	}
}

/*


 */
