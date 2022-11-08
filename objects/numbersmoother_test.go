package objects

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestPosition(t *testing.T) {
	j := map[string]interface{}{
		"posX": 42.41,
		"posY": 1.5,
		"posZ": -11.36,
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
