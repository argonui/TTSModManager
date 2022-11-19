package objects

import (
	"ModCreator/types"
	"math"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestSmooth(t *testing.T) {
	for _, tc := range []struct {
		name        string
		input, want map[string]interface{}
	}{
		{
			name: "Position",
			input: map[string]interface{}{
				"posX": 42.4101944,
				"posY": 1.49994421,
				"posZ": -11.3552332,
			},
			want: map[string]interface{}{
				"posX": 42.41,
				"posY": 1.5,
				"posZ": -11.355,
			},
		},
		{
			name: "Rotation",
			input: map[string]interface{}{
				"rotX": float64(90),
				"rotY": 89.83327,
				"rotZ": float64(-0.004),
			},
			want: map[string]interface{}{
				"rotX": float64(90),
				"rotY": float64(90),
				"rotZ": float64(0),
			},
		},
		{
			name: "Rotation with overflow",
			input: map[string]interface{}{
				"rotX": float64(370),
				"rotY": -89.83327,
				"rotZ": float64(-0.004),
			},
			want: map[string]interface{}{
				"rotX": float64(10),
				"rotY": float64(270),
				"rotZ": float64(0),
			},
		},
		{
			name: "Color",
			input: map[string]interface{}{
				"r": 0.42513,
				"g": 0.333333333,
				"b": 0.914525304,
				"a": float64(0.5),
			},
			want: map[string]interface{}{
				"r": 0.42513,
				"g": 0.33333,
				"b": 0.91453,
				"a": 0.5,
			},
		},
		{
			name: "Color Transparent",
			input: map[string]interface{}{
				"r": 1.0,
				"g": 1.0,
				"b": 1.0,
				"a": 0,
			},
			want: map[string]interface{}{
				"r": 1.0,
				"g": 1.0,
				"b": 1.0,
				"a": 0,
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := Smooth(tc.input)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("want != got:\n%v\n", diff)
			}
		})
	}
}

func TestNegativeZeroPos(t *testing.T) {
	input := map[string]interface{}{
		"posX": float64(0),
		"posY": float64(-0),
		"posZ": float64(-0.00001),
	}
	gotraw := Smooth(input)
	got, ok := gotraw.(map[string]interface{})
	if !ok {
		t.Fatalf("bad typing")
	}
	for _, k := range posRounded {
		raw, ok := got[k]
		if !ok {
			t.Fatalf("key %s not found", k)
		}
		fl, ok := raw.(float64)
		if !ok {
			t.Errorf("key %s is not float64; is %T", k, raw)
		}
		if math.Signbit(fl) {
			t.Errorf("%s returned negative", k)
		}
	}
}

func TestNegativeZeroRot(t *testing.T) {
	input := map[string]interface{}{
		"rotX": float64(0),
		"rotY": float64(-0),
		"rotZ": float64(-0.00001),
	}
	gotraw := Smooth(input)
	got, ok := gotraw.(map[string]interface{})
	if !ok {
		t.Fatalf("bad typing")
	}
	for _, k := range rotRounded {
		raw, ok := got[k]
		if !ok {
			t.Fatalf("key %s not found", k)
		}
		fl, ok := raw.(float64)
		if !ok {
			t.Errorf("key %s is not float64; is %T", k, raw)
		}
		if math.Signbit(fl) {
			t.Errorf("%s returned negative", k)
		}
	}
}

func TestArray(t *testing.T) {
	j := []interface{}{
		map[string]interface{}{
			"Position": types.J{
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
		map[string]interface{}{
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

func TestBadKeyInSnapPoints(t *testing.T) {
	j := []interface{}{
		map[string]interface{}{
			"Position": types.J{
				"x":      float64(-1.822),
				"y":      float64(0.100),
				"z":      float64(0.616),
				"foobar": false,
			},
			"Rotation": types.J{
				"x": float64(0),
				"y": float64(0),
				"z": float64(0),
			},
		},
	}
	_, err := SmoothSnapPoints(j)
	if err == nil {
		t.Fatal("SmoothSnapPoints(): wanted error")
	}
	if !strings.Contains(err.Error(), "foobar") {
		t.Errorf("Expected an error about the unexpected key foobar")
	}
}
