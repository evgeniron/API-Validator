package validate

import (
	"net/mail"
	"regexp"
	"strings"
	"time"
)

var Validators = map[string]func(interface{}) bool{
	"Int":        IntValidator,
	"Date":       DateValidator,
	"Boolean":    BooleanValidator,
	"List":       JSONValidator,
	"UUID":       UUIDValidator,
	"BearerAuth": BearerAuthValidator,
	"String":     StringValidator,
}

// EmailValidator validates an email address (RFC 5322)
func EmailValidator(value interface{}) bool {
	email, ok := value.(string)
	if !ok {
		return false
	}

	_, err := mail.ParseAddress(strings.TrimSpace(email))
	return err == nil
}

// IntValidator validates an integer value
func IntValidator(value interface{}) bool {
	_, ok := value.(int)
	return ok
}

// StringValidator validates a string value
func StringValidator(value interface{}) bool {
	_, ok := value.(string)
	return ok
}

// BooleanValidator validates a boolean value
func BooleanValidator(value interface{}) bool {
	_, ok := value.(bool)
	return ok
}

// JSONValidator validates a list of JSON objects
func JSONValidator(value interface{}) bool {
	_, ok := value.([]interface{})
	return ok
}

// DateValidator validates a date with format "dd-mm-yyyy"
func DateValidator(value interface{}) bool {
	date, ok := value.(string)
	if !ok {
		return false
	}

	_, err := time.Parse("02-01-2006", date)
	return err == nil
}

// UUIDValidator validates a UUID value
func UUIDValidator(value interface{}) bool {
	uuid, ok := value.(string)
	if !ok {
		return false
	}

	re := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
	return re.MatchString(strings.TrimSpace(uuid))
}

// BearerAuthValidator validates a bearer authentication token
func BearerAuthValidator(value interface{}) bool {
	bearerToken, ok := value.(string)
	if !ok {
		return false
	}

	re := regexp.MustCompile("^Bearer [a-zA-Z0-9\\-._~+/]+=*$")
	return re.MatchString(strings.TrimSpace(bearerToken))
}
