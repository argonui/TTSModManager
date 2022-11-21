package objects

import (
	"ModCreator/types"
	"fmt"
	"math"
)

var (
	// because arrays can't be constant
	posRounded       = []string{"posX", "posY", "posZ"}
	rotRounded       = []string{"rotX", "rotY", "rotZ"}
	colorRounded     = []string{"b", "g", "r", "a"}
	scaleRounded     = []string{"scaleX", "scaleY", "scaleZ"}
	arbitraryRounded = []string{"x", "y", "z"}
)

func roundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

// Smooth will round numbers or degrees to reasonable precision
func Smooth(objraw interface{}) interface{} {
	obj, ok := objraw.(map[string]interface{})
	if !ok {
		return objraw
	}

	for _, key := range posRounded {
		if val, ok := obj[key]; ok {
			if fl, ok := val.(float64); ok {
				obj[key] = smoothPos(fl)
			}
		}
	}
	for _, key := range rotRounded {
		if val, ok := obj[key]; ok {
			if fl, ok := val.(float64); ok {
				obj[key] = smoothRot(fl)
			}
		}
	}
	for _, key := range colorRounded {
		if val, ok := obj[key]; ok {
			if fl, ok := val.(float64); ok {
				obj[key] = roundFloat(fl, 5)
			}
		}
	}
	for _, key := range scaleRounded {
		if val, ok := obj[key]; ok {
			if fl, ok := val.(float64); ok {
				obj[key] = roundFloat(fl, 2)
			}
		}
	}
	return obj
}

func smoothArbitrary(objraw interface{}, round func(float64) float64) (types.J, error) {
	smoothed := types.J{}
	obj, ok := objraw.(map[string]interface{})
	if !ok {
		obj, ok = objraw.(types.J)
		if !ok {
			return nil, fmt.Errorf("Couldn't convert object, was type %T", objraw)
		}
	}

	for _, key := range arbitraryRounded {
		if val, ok := obj[key]; ok {
			if fl, ok := val.(float64); ok {
				smoothed[key] = round(fl)
			} else {
				smoothed[key] = val
			}
		} else {
			return nil, fmt.Errorf("expected key %s to be exist", key)
		}
	}
	for k := range obj {
		if _, ok := smoothed[k]; !ok {
			return nil, fmt.Errorf("Key not smoothed: %s, expected to be a float64 at one of %v", k, arbitraryRounded)
		}
	}

	return smoothed, nil
}

// SmoothSnapPoints will consume an array of snap points and smooth all of them
func SmoothSnapPoints(rawsps interface{}) ([]map[string]interface{}, error) {
	arr, err := types.ConvertToObjArray(rawsps)
	if err != nil {
		return nil, fmt.Errorf("ConvertToObjArray(%v): %v", rawsps, err)
	}
	smth := []map[string]interface{}{}
	for _, sp := range arr {

		smthsp := types.J{}
		if pos, ok := sp["Position"]; ok {
			v, err := smoothArbitrary(pos, smoothPos)
			if err != nil {
				return nil, fmt.Errorf("smoothArbitrary(%v, positional): %v", pos, err)
			}
			smthsp["Position"] = v
		}
		if rot, ok := sp["Rotation"]; ok {
			v, err := smoothArbitrary(rot, smoothRot)
			if err != nil {
				return nil, fmt.Errorf("smoothArbitrary(%v, rotational): %v", rot, err)
			}
			smthsp["Rotation"] = v
		}
		if tags, ok := sp["Tags"]; ok {
			smthsp["Tags"] = tags
		}
		if len(sp) != len(smthsp) {
			return nil, fmt.Errorf("Unexpected key(s). full obj: %v", sp)
		}

		smth = append(smth, smthsp)
	}
	return smth, nil
}

// SmoothAngle smooths x,y,z per angle rules
func SmoothAngle(objraw interface{}) (interface{}, error) {
	return smoothArbitrary(objraw, smoothRot)
}

func smoothPos(f float64) float64 {
	r := roundFloat(f, 3)
	if r == float64(0) {
		// for some reason -0.000001 is being returned as "-0"
		return float64(0)
	}
	return r
}

func smoothRot(f float64) float64 {
	r := math.Mod(roundFloat(f, 0)+360, 360)
	if r == float64(0) {
		// for some reason -0.000001 is being returned as "-0"
		return float64(0)
	}
	return r
}
