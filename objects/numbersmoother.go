package objects

import (
	"math"
)

var (
	// because arrays can't be constant
	posRounded = []string{"posX", "posY", "posZ"}
	rotRounded = []string{"rotX", "rotY", "rotZ"}
	colorRounded = []string{"b", "g", "r", "a"}
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
				obj[key] = roundFloat(fl, 3)
			}
		}
	}
	for _, key := range rotRounded {
		if val, ok := obj[key]; ok {
			if fl, ok := val.(float64); ok {
				obj[key] = math.Abs(math.Mod(roundFloat(fl, 0), 360))
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
