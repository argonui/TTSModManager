package file

import (
	"ModCreator/types"
	"fmt"
)

// TryParseIntoStr attempts to assign string
func TryParseIntoStr(m *types.J, k string, dest *string) {
	_ = ForceParseIntoStr(m, k, dest)
}

// ForceParseIntoStr will error if unsuccessful
func ForceParseIntoStr(m *types.J, k string, dest *string) error {
	if raw, ok := (*m)[k]; ok {
		if str, ok := raw.(string); ok {
			*dest = str
			delete((*m), k)
			return nil
		}
		return fmt.Errorf("key %s not convertable to string", k)
	}
	return fmt.Errorf("key %s not found", k)
}

// TryParseIntoStrArray will not error
func TryParseIntoStrArray(m *types.J, k string, dest *[]string) {
	_ = ForceParseIntoStrArray(m, k, dest)
}

// ForceParseIntoStrArray will error
func ForceParseIntoStrArray(m *types.J, k string, dest *[]string) error {
	raw, ok := (*m)[k]
	if !ok {
		return fmt.Errorf("key %s not found", k)
	}
	if rawstrarr, ok := raw.([]string); ok {
		*dest = rawstrarr
		delete((*m), k)
		return nil
	}

	rawarr, ok := raw.([]interface{})
	if !ok {
		return fmt.Errorf("key %s is not an array type; is %T", k, raw)
	}

	strarr := []string{}
	for i, val := range rawarr {
		str, ok := val.(string)
		if !ok {
			return fmt.Errorf("element %v of %s is not convertable to string. type %T", i, k, val)
		}
		strarr = append(strarr, str)
	}
	*dest = strarr
	delete((*m), k)

	return nil
}

func TryParseIntoStrMap(m *types.J, k string, dest *map[string]string) {
	_ = ForceParseIntoStrMap(m, k, dest)
}

func ForceParseIntoStrMap(m *types.J, k string, dest *map[string]string) error {
	raw, ok := (*m)[k]
	if !ok {
		return fmt.Errorf("key %s not found", k)
	}
	if rawmapstrings, ok := raw.(map[string]string); ok {
		*dest = rawmapstrings
		delete((*m), k)
		return nil
	}
	if rawmap, ok := raw.(map[string]any); ok {
		massaged := map[string]string{}
		for k, v := range rawmap {
			vstr, ok := v.(string)
			if !ok {
				return fmt.Errorf("m[%s] of type %T, not string", k, v)
			}
			massaged[k] = vstr
		}
		*dest = massaged
		delete((*m), k)
		return nil
	}
	return fmt.Errorf("key %s was not type map[string]string, was %T", k, raw)
}

// TryParseIntoInt will not throw error
func TryParseIntoInt(m *types.J, k string, dest *int64) {
	_ = ForceParseIntoInt(m, k, dest)
}

// ForceParseIntoInt will throw error
func ForceParseIntoInt(m *types.J, k string, dest *int64) error {
	if raw, ok := (*m)[k]; ok {
		if in, ok := raw.(int64); ok {
			*dest = in
			delete((*m), k)
			return nil
		}
		if fl, ok := raw.(float64); ok {
			*dest = int64(fl)
			delete((*m), k)
			return nil
		}
		return fmt.Errorf("key %s unable to be parsed as float64 or int64", k)
	}
	return fmt.Errorf("key %s not found", k)
}
