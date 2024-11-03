package polycode

import (
	"encoding/json"
	"reflect"
)

func ToMap(data map[string][]string) map[string]any {
	// Create a map to hold the final map[string]interface{}
	result := make(map[string]interface{})

	// Iterate through the query parameters
	for key, value := range data {
		if len(value) == 1 {
			// If there's only one value for this key, store it as a string
			result[key] = value[0]
		} else {
			// If there are multiple values, store them as a []string
			result[key] = value
		}
	}
	return result
}

func GetTypeNameFromT[T any]() string {
	var zeroValue T
	typeName := reflect.TypeOf(zeroValue).Name()
	return typeName
}

func GetTypeName[T any](value T) string {
	t := reflect.TypeOf(value)

	// Handling for pointer types to get the base type
	if t.Kind() == reflect.Pointer || t.Kind() == reflect.Slice {
		t = t.Elem()
	}
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	// For named structs, this should now return the struct name
	//fmt.Println("The type name is:", t.Name())
	return t.Name()
}

func GetTypeNameWithPkg[T any](value T) (string, string) {
	t := reflect.TypeOf(value)

	if t.Kind() == reflect.Ptr || t.Kind() == reflect.Slice {
		panic("pointer | slice type not allowed")
	}
	//// Handling for pointer types to get the base type
	//if t.Kind() == reflect.Pointer || t.Kind() == reflect.Slice {
	//	t = t.Elem()
	//}

	// For named structs, this should now return the struct name
	//fmt.Println("The type name is:", t.Name())
	return t.PkgPath(), t.Name()
}

func IsPointer(data interface{}) bool {
	return reflect.TypeOf(data).Kind() == reflect.Ptr
}

func ConvertMap(m interface{}) interface{} {
	switch x := m.(type) {
	case map[interface{}]interface{}:
		converted := make(map[string]interface{})
		for k, v := range x {
			converted[k.(string)] = ConvertMap(v) // Recursively convert values
		}
		return converted
	case []interface{}:
		for i, v := range x {
			x[i] = ConvertMap(v)
		}
	}
	return m
}

func ConvertType(input any, output any) error {
	in, err := json.Marshal(input)
	if err != nil {
		return err
	}

	return json.Unmarshal(in, output)
}
