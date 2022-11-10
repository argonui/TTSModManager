package types

import "fmt"

// J is json
type J map[string]interface{}

// ObjArray is json array
type ObjArray []map[string]interface{}

// ConvertToObjArray handles type conversion
func ConvertToObjArray(v interface{}) ([]map[string]interface{}, error) {
	if easyconvert, ok := v.([]map[string]interface{}); ok {
		return easyconvert, nil
	}
	arr := []map[string]interface{}{}

	rawArr, ok := v.([]interface{})
	if !ok {
		return nil, fmt.Errorf("%v is not an array, is %T", v, v)
	}

	for _, rv := range rawArr {
		objVal, ok := rv.(map[string]interface{})
		if !ok {
			if rv == nil {
				// if for some reason an array has nil object, just skip
				continue
			}
			return nil, fmt.Errorf("expected type json object, got %v", objVal)
		}
		arr = append(arr, objVal)
	}
	return arr, nil
}
