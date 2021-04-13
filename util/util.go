package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

func DecodeGetParameters(query string, v interface{}) error {
	keyValuePairs := strings.Split(query, "&")
	var m = map[string]string{}
	for _, pair := range keyValuePairs {
		if strings.Index(pair, "=") > 0 {
			key := strings.Title(strings.Split(pair, "=")[0])
			value := strings.Split(pair, "=")[1]
			m[key] = value
		}
	}
	jsonStr, err := json.Marshal(m)
	if err != nil {
		return err
	}

	err = json.Unmarshal(jsonStr, &v)
	if err != nil {
		return err
	}

	return nil
}

func setField(obj interface{}, name string, value interface{}) error {
	structValue := reflect.ValueOf(obj).Elem()
	structFieldValue := structValue.FieldByName(name)

	if !structFieldValue.IsValid() {
		return fmt.Errorf("No such field: %s in obj", name)
	}

	if !structFieldValue.CanSet() {
		return fmt.Errorf("Cannot set %s field value", name)
	}

	structFieldType := structFieldValue.Type()
	val := reflect.ValueOf(value)
	if structFieldType != val.Type() {
		return errors.New("Provided value type didn't match obj field type")
	}

	structFieldValue.Set(val)
	return nil
}

func mapToStruct(s interface{}, m map[string]interface{}) error {
	for k, v := range m {
		err := setField(s, k, v)
		if err != nil {
			return err
		}
	}
	return nil
}

func CopyStruct(dst interface{}, src interface{}) error {
	// make src and dst into map
	// traverse the maps, if a value of key has nil value, set it to value of key of src, else skip.
	fmt.Println(dst)
	dstMap := structToMap(dst)
	srcMap := structToMap(src)
	for key, value := range dstMap {
		if reflect.TypeOf(value).Kind() == reflect.String {
			if len(value.(string)) == 0 {
				dstMap[key] = srcMap[key]
			}
		}
		if reflect.TypeOf(value).Kind() == reflect.Int32 {
			if value.(int32) == 0 {
				dstMap[key] = srcMap[key]
			}
		}
	}
	fmt.Println("finally:", dstMap)

	jsonStr, _ := json.Marshal(dstMap)

	if err := json.Unmarshal(jsonStr, &dst); err != nil {
		return err
	}
	return nil
}

func structToMap(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
    v := reflect.ValueOf(obj)
    data := make(map[string]interface{})
    if t.Kind() == reflect.Ptr {
    	t = t.Elem()
	}
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Name] = v.Field(i).Interface()
	}
	return data
}