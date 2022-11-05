package objects

import (
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
		"posX": 42.4102,
		"posY": 1.4999,
		"posZ": -11.3552,
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("want != got:\n%v\n", diff)
	}
}

func TestDegree(t *testing.T) {
	j := map[string]interface{}{
		"rotX": float64(90),
		"rotY": 89.83327,
		"rotZ": float64(0),
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
