package helpers

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/greatcloak/decimal"
	"gorm.io/gorm"
)

var (
	ErrMessageValueNotFound        = "value is not found"
	ErrMessageFieldWithTagNotFound = "field with the specified tag not found"
)

func MapToStruct(m map[string]any, s any) error {
	structType := reflect.TypeOf(s).Elem()
	structValue := reflect.ValueOf(s).Elem()

	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		fieldValue := structValue.FieldByName(field.Name)
		fieldName := field.Tag.Get("json")

		// Check if the map key exists with the original field name
		if value, ok := m[fieldName]; ok {
			if fieldValue.CanSet() {
				fieldValue.Set(reflect.ValueOf(value))
			} else {
				return fmt.Errorf("cannot set value for field %s", fieldName)
			}

			continue
		}

		// If the map key does not exist with the original field name,
		// check if it exists with the modified name (e.g., full_name -> FullName)
		modifiedFieldName := strings.ToLower(fieldName)
		if value, ok := m[modifiedFieldName]; ok {
			if fieldValue.CanSet() {
				fieldValue.Set(reflect.ValueOf(value))
			} else {
				return fmt.Errorf("cannot set value for field %s", fieldName)
			}
		}
	}

	return nil
}

func GetStructTagValue(s reflect.StructField, prefix string) string {
	return s.Tag.Get(prefix)
}

func StructToBytes(data any) []byte {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return jsonData
	}

	return jsonData
}

func StructToMap[T any](val any) (res map[string]T) {
	res = make(map[string]T)
	b, _ := json.Marshal(val)
	json.Unmarshal(b, &res)

	return
}

func ExtractStructToBytes(data any) []byte {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return []byte(fmt.Sprintf("%+v", data))
	}

	return jsonData
}

func ExtractStructToString(data any) string {
	jsonBytes := ExtractStructToBytes(data)
	return string(jsonBytes)
}

func SanitizeString(secret string, censoredLength int, censoredWith string) string {
	var censored string

	length := len(secret)

	if censoredLength > length {
		censoredLength = length
	}

	for i := 0; i < length-censoredLength; i++ {
		censored += censoredWith
	}

	sanitized := secret[0:censoredLength] + censored

	return sanitized
}

func SanitizeStruct(val reflect.Value) any {
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return val.Interface()
	}

	typ := val.Type()
	newStruct := reflect.New(typ).Elem()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		if newStruct.Field(i).CanSet() {
			newStruct.Field(i).Set(field)
		}

		sanitizeField := GetStructTagValue(fieldType, "sanitize")

		if sanitizeField == "on" && field.Kind() == reflect.String {
			originalValue := field.String()
			sanitizedValue := SanitizeString(originalValue, 3, "*")
			newStruct.Field(i).SetString(sanitizedValue)
		}
	}

	return newStruct.Interface()
}

func FormatReflectValue(val reflect.Value) string {
	var result string

	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return fmt.Sprintf("%s<nil>", val.Type().Name())
		}

		val = val.Elem()
	}

	typ := val.Type()

	switch val.Kind() {
	case reflect.Struct:
		if typ == reflect.TypeOf(time.Time{}) {
			v, ok := val.Interface().(time.Time)
			if !ok {
				return ""
			}

			result += v.String()

			return result
		}

		if typ == reflect.TypeOf(gorm.DeletedAt{}) {
			v, ok := val.Interface().(gorm.DeletedAt)
			if !ok {
				return ""
			}

			if v.Valid {
				result += v.Time.String()
			} else {
				result += "<nil>"
			}

			return result
		}

		if typ == reflect.TypeOf(decimal.Decimal{}) {
			v, ok := val.Interface().(decimal.Decimal)
			if !ok {
				return ""
			}

			result += v.String()

			return result
		}

		for i := 0; i < val.NumField(); i++ {
			field := val.Field(i)
			fieldType := typ.Field(i)
			fieldName := fieldType.Name

			if !field.CanInterface() {
				result += fmt.Sprintf("%s:<unexported> ", fieldName)
			}

			if fieldType.Type == reflect.TypeOf(json.RawMessage{}) {
				jsonData, ok := field.Interface().(json.RawMessage)
				if !ok {
					return ""
				}

				result += fmt.Sprintf("%s:%s ", fieldName, string(jsonData))

				continue
			}

			result += fmt.Sprintf("%s:%s ", fieldName, FormatReflectValue(field))
		}
	case reflect.Slice:
		result += "["

		for i := 0; i < val.Len(); i++ {
			element := val.Index(i)
			result += FormatReflectValue(element) + " "
		}

		result += "]"
	default:
		if val.CanInterface() {
			result = fmt.Sprintf("%v", val.Interface())
		} else {
			result = "<unexported>"
		}
	}

	return result
}

func StructIsEmpty(strct any) bool {
	v := reflect.ValueOf(strct)

	// Check if the value is a pointer and dereference it
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// Check if the struct is valid and of kind Struct
	if v.Kind() != reflect.Struct {
		return true // or handle this case as needed
	}

	// Iterate through all fields in the struct
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		// If any field is not zero, return false
		if field.IsValid() && !field.IsZero() {
			return false
		}
	}

	return true // All fields are zero values
}

func BytesToMap(data []byte) (map[string]any, error) {
	var result map[string]any

	err := json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Find field by tag.
func GetFieldByTag(entity interface{}, tagValue string, tagKey string) (interface{}, error) {
	v := reflect.ValueOf(entity)

	// Handle pointer structs
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return nil, errors.New("provided value is not a struct")
	}

	// Get the type of the struct
	t := v.Type()

	// Iterate over the struct fields
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get(tagKey)

		// Check if the tag contains the "column:<tagValue>"
		if tag == tagValue {
			// Return the value of the field
			return v.Field(i).Interface(), nil
		}
	}

	return nil, errors.New("value is not found")
}

// SetFieldByTag sets the value of a struct field identified by a tag.
func SetFieldByTag(entity interface{}, tagValue string, tagKey string, newValue interface{}) error {
	v := reflect.ValueOf(entity)

	// Handle pointer structs
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	} else {
		return errors.New("entity must be a pointer to a struct")
	}

	if v.Kind() != reflect.Struct {
		return errors.New("provided value is not a struct")
	}

	// Get the type of the struct
	t := v.Type()

	// Iterate over the struct fields
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get(tagKey)

		// Check if the tag matches the tagValue
		if tag == tagValue {
			fieldValue := v.Field(i)

			// Ensure the field is settable
			if !fieldValue.CanSet() {
				return errors.New("field is not settable")
			}

			// Convert newValue to the field's type if necessary
			newValueReflected := reflect.ValueOf(newValue)

			if newValueReflected.Type().ConvertibleTo(fieldValue.Type()) {
				fieldValue.Set(newValueReflected.Convert(fieldValue.Type()))
			} else {
				return errors.New(fmt.Sprintf("cannot assign value of type %s to field of type %s", newValueReflected.Type(), fieldValue.Type()))
			}

			return nil
		}
	}

	return errors.New(ErrMessageFieldWithTagNotFound)
}
