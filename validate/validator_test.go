package validate

import "testing"

func TestEmailValidator(t *testing.T) {
	tests := []struct {
		value    interface{}
		expected bool
	}{
		{"example@example.com", true},
		{"invalidemail", false},
		{"invalidemail@", false},
		{"invalidemail.com", false},
		{123, false},
	}

	for _, test := range tests {
		result := EmailValidator(test.value)
		if result != test.expected {
			t.Errorf("EmailValidator(%v) = %v, expected %v", test.value, result, test.expected)
		}
	}
}

func TestIntValidator(t *testing.T) {
	tests := []struct {
		value    interface{}
		expected bool
	}{
		{1, true},
		{"1", false},
		{true, false},
		{"", false},
	}

	for _, test := range tests {
		result := IntValidator(test.value)
		if result != test.expected {
			t.Errorf("IntValidator(%v) = %v, expected %v", test.value, result, test.expected)
		}
	}
}

func TestStringValidator(t *testing.T) {
	tests := []struct {
		value    interface{}
		expected bool
	}{
		{"string", true},
		{1, false},
		{true, false},
	}

	for _, test := range tests {
		result := StringValidator(test.value)
		if result != test.expected {
			t.Errorf("StringValidator(%v) = %v, expected %v", test.value, result, test.expected)
		}
	}
}

func TestBooleanValidator(t *testing.T) {
	tests := []struct {
		value    interface{}
		expected bool
	}{
		{true, true},
		{"true", false},
		{1, false},
	}

	for _, test := range tests {
		result := BooleanValidator(test.value)
		if result != test.expected {
			t.Errorf("BooleanValidator(%v) = %v, expected %v", test.value, result, test.expected)
		}
	}
}

func TestJSONValidator(t *testing.T) {
	tests := []struct {
		value    interface{}
		expected bool
	}{
		{[]interface{}{1, 2, 3}, true},
		{"[1, 2, 3]", false},
		{1, false},
		{
			[]interface{}{
				map[string]interface{}{"hello": "world1", "num": float64(1)},
				map[string]interface{}{"hello": "world2", "num": float64(2)},
			},
			true,
		},
	}

	for _, test := range tests {
		result := JSONValidator(test.value)
		if result != test.expected {
			t.Errorf("JSONValidator(%v) = %v, expected %v", test.value, result, test.expected)
		}
	}
}

func TestDateValidator(t *testing.T) {
	validDate := "02-01-2023"
	invalidDate := "31-02-2023"

	if !DateValidator(validDate) {
		t.Errorf("Expected valid date to return true, got false")
	}

	if DateValidator(invalidDate) {
		t.Errorf("Expected invalid date to return false, got true")
	}
}
