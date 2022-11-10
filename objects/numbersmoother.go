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
			}
		}
	}
	if len(smoothed) != len(obj) {
		return nil, fmt.Errorf("unexpected keys match in %v", obj)
	}

	return smoothed, nil
}

// SmoothSnapPoints will consume an array of snap points and smooth all of them
func SmoothSnapPoints(rawsps interface{}) ([]map[string]interface{}, error) {
	arr, ok := rawsps.([]map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("SnapPoints is expected to be an array of objects; was %T", rawsps)
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
		if len(sp) != len(smthsp) {
			return nil, fmt.Errorf("Unexpected key found in array of snap points. full obj: %v", sp)
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
	return roundFloat(f, 3)
}

func smoothRot(f float64) float64 {
	return math.Mod(roundFloat(f, 0)+360, 360)
}
