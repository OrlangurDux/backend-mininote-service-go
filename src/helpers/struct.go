package helpers

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// normalizeFieldName converts field names like "category_id" and "categoryid" to a common base form.
func normalizeFieldName(name string) string {
	// Remove underscores and convert to lower case for a case-insensitive comparison
	return strings.ToLower(strings.ReplaceAll(name, "_", ""))
}

// SetField -> set value from value variable as key
func SetField(obj interface{}, name string, value interface{}) error {
	structValue := reflect.ValueOf(obj).Elem()

	// Normalize the name for comparison
	name = normalizeFieldName(name)

	// Find the field by normalized name
	var structFieldValue reflect.Value
	for i := 0; i < structValue.NumField(); i++ {
		field := structValue.Type().Field(i)
		if normalizeFieldName(field.Name) == name {
			structFieldValue = structValue.FieldByName(field.Name)
			break
		}
	}

	if !structFieldValue.IsValid() {
		return fmt.Errorf("no such field: %s in obj", name)
	}

	if !structFieldValue.CanSet() {
		return fmt.Errorf("cannot set %s field value", name)
	}

	structFieldType := structFieldValue.Type()
	val := reflect.ValueOf(value)

	// Check if the field type is time.Time and the value is a string
	if structFieldType == reflect.TypeOf(time.Time{}) && val.Kind() == reflect.String {
		// Parse the string to time.Time
		parsedTime, err := time.Parse(time.RFC3339, val.String())
		if err != nil {
			return fmt.Errorf("error parsing time: %s", err)
		}
		structFieldValue.Set(reflect.ValueOf(parsedTime))
	} else if strings.Contains(name, "id") && structFieldType == reflect.TypeOf(primitive.ObjectID{}) && val.Kind() == reflect.String {
		// Convert the string to ObjectID
		objectID, err := primitive.ObjectIDFromHex(val.String())
		if err != nil {
			return fmt.Errorf("error converting string to ObjectID: %s", err)
		}
		structFieldValue.Set(reflect.ValueOf(objectID))
	} else if strings.Contains(name, "active") && val.Kind() == reflect.String {
		bActive, err := strconv.ParseBool(val.String())
		if err != nil {
			return fmt.Errorf("error converting string to boolean: %s", err)
		}
		structFieldValue.Set(reflect.ValueOf(bActive))
	} else if structFieldType != val.Type() {
		return fmt.Errorf("provided value type didn't match obj field type name: %s", name)
	} else {
		structFieldValue.Set(val)
	}

	return nil
}
