package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func Decode(r *http.Request, v interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return fmt.Errorf("error decoding request body: %w", err)
	}
	return nil
}
