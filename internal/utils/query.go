package utils

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"github.com/srpvpn/tensor-go-sdk/internal/errors"
)

// BuildQueryParams converts a struct to URL query parameters
// It uses struct tags to determine parameter names and handles omitempty
func BuildQueryParams(req interface{}) (url.Values, error) {
	if req == nil {
		return url.Values{}, nil
	}

	params := url.Values{}
	v := reflect.ValueOf(req)
	t := reflect.TypeOf(req)

	// Handle pointer types
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return params, nil
		}
		v = v.Elem()
		t = t.Elem()
	}

	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected struct, got %T", req)
	}

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		// Skip unexported fields
		if !field.CanInterface() {
			continue
		}

		// Get the JSON tag for parameter name
		tag := fieldType.Tag.Get("json")
		if tag == "" || tag == "-" {
			continue
		}

		// Parse tag options
		tagParts := strings.Split(tag, ",")
		paramName := tagParts[0]
		hasOmitEmpty := len(tagParts) > 1 && contains(tagParts[1:], "omitempty")

		// Skip if omitempty and field is zero value
		if hasOmitEmpty && isZeroValue(field) {
			continue
		}

		// Convert field value to string
		paramValue, err := fieldToString(field)
		if err != nil {
			return nil, fmt.Errorf("error converting field %s: %w", fieldType.Name, err)
		}

		if paramValue != "" {
			params.Add(paramName, paramValue)
		}
	}

	return params, nil
}

// ValidateWalletAddress validates a Solana wallet address
func ValidateWalletAddress(wallet string) error {
	if wallet == "" {
		return &errors.ValidationError{
			Field:   "wallet",
			Message: "wallet address is required",
		}
	}

	// Basic Solana address validation
	// Solana addresses are base58 encoded and typically 32-44 characters
	if len(wallet) < 32 || len(wallet) > 44 {
		return &errors.ValidationError{
			Field:   "wallet",
			Message: "invalid wallet address length",
		}
	}

	// Check for valid base58 characters
	validChars := "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
	for _, char := range wallet {
		if !strings.ContainsRune(validChars, char) {
			return &errors.ValidationError{
				Field:   "wallet",
				Message: "wallet address contains invalid characters",
			}
		}
	}

	return nil
}

// Helper functions

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func isZeroValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.String:
		return v.String() == ""
	case reflect.Slice, reflect.Map, reflect.Array:
		return v.Len() == 0
	case reflect.Ptr, reflect.Interface:
		return v.IsNil()
	default:
		return false
	}
}

func fieldToString(v reflect.Value) (string, error) {
	switch v.Kind() {
	case reflect.String:
		return v.String(), nil
	case reflect.Bool:
		return strconv.FormatBool(v.Bool()), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(v.Uint(), 10), nil
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(v.Float(), 'f', -1, 64), nil
	case reflect.Slice:
		if v.Type().Elem().Kind() == reflect.String {
			// Handle []string by joining with commas
			strs := make([]string, v.Len())
			for i := 0; i < v.Len(); i++ {
				strs[i] = v.Index(i).String()
			}
			return strings.Join(strs, ","), nil
		}
		return "", fmt.Errorf("unsupported slice type: %s", v.Type())
	case reflect.Ptr:
		if v.IsNil() {
			return "", nil
		}
		return fieldToString(v.Elem())
	default:
		return "", fmt.Errorf("unsupported field type: %s", v.Type())
	}
}
