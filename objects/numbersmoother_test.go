package objects

import (
	"ModCreator/types"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestPosition(t *testing.T) {
	j := map[string]interface{}{
		"posX": 42.4101944,
		"posY": 1.49994421,
		"posZ": -11.3552332,
	}
	got := Smooth(j)
	want := map[string]interface{}{
		"posX": 42.41,
		"posY": 1.5,
		"posZ": -11.355,
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("want != got:\n%v\n", diff)
	}
}

func TestDegree(t *testing.T) {
	j := map[string]interface{}{
		"rotX": float64(90),
		"rotY": 89.83327,
		"rotZ": float64(-0.004),
	}
	got := Smooth(j)
	want := map[string]interface{}{
		"rotX": float64(90),
		"rotY": float64(90),
		"rotZ": float64(0),
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("want != got:\n%v\n", diff)
	}
}

func TestDegreeAbs(t *testing.T) {
	j := map[string]interface{}{
		"rotX": float64(370),
		"rotY": -89.83327,
		"rotZ": float64(-0.004),
	}
	got := Smooth(j)
	want := map[string]interface{}{
		"rotX": float64(10),
		"rotY": float64(270),
		"rotZ": float64(0),
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("want != got:\n%v\n", diff)
	}
}

func TestColor(t *testing.T) {
	j := map[string]interface{}{
		"r": 0.42513,
		"g": 0.333333333,
		"b": 0.914525304,
		"a": float64(0.5),
	}
	got := Smooth(j)
	want := map[string]interface{}{
		"r": 0.42513,
		"g": 0.33333,
		"b": 0.91453,
		"a": 0.5,
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("want != got:\n%v\n", diff)
	}
}

func TestColorTransparent(t *testing.T) {
	j := map[string]interface{}{
		"r": 1.0,
		"g": 1.0,
		"b": 1.0,
		"a": 0,
	}
	got := Smooth(j)
	want := map[string]interface{}{
		"r": 1.0,
		"g": 1.0,
		"b": 1.0,
		"a": 0,
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("want != got:\n%v\n", diff)
	}
}

func TestArray(t *testing.T) {
	j := []map[string]interface{}{
		{"Position": types.J{
			"x": float64(-1.822),
			"y": float64(0.100),
			"z": float64(0.616),
		},
			"Rotation": types.J{
				"x": float64(0),
				"y": float64(0),
				"z": float64(0),
			},
		},
		{
			"Rotation": types.J{
				"x": float64(0),
				"y": float64(0),
				"z": float64(0),
			},
		},
	}
	got, err := SmoothSnapPoints(j)
	if err != nil {
		t.Fatalf("SmoothSnapPoints(): %v", err)
	}
	want := []map[string]interface{}{
		{"Position": types.J{
			"x": -1.822,
			"y": 0.100,
			"z": 0.616,
		},
			"Rotation": types.J{
				"x": float64(0),
				"y": float64(0),
				"z": float64(0),
			},
		},
		{
			"Rotation": types.J{
				"x": float64(0),
				"y": float64(0),
				"z": float64(0),
			},
		},
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("want != got:\n%v\n", diff)
	}
}
