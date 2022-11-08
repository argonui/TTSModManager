package file

import (
	"ModCreator/types"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestStrArr(t *testing.T) {
	for _, tc := range []struct {
		name    string
		input   interface{}
		want    []string
		wantErr bool
	}{
		{
			name:    "simple",
			input:   []interface{}{"foo", "bar"},
			want:    []string{"foo", "bar"},
			wantErr: false,
		},
		{
			name:    "bool in array",
			input:   []interface{}{"foo", "bar", false},
			want:    []string{},
			wantErr: true,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			m := types.J{}
			m["key"] = tc.input
			got := []string{}
			err := ForceParseIntoStrArray(&m, "key", &got)
			if (err == nil) == tc.wantErr {
				t.Fatalf("wantErr %v got %v", tc.wantErr, err)
			}
			if tc.wantErr {
				return
			}
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("want != got:\n%v\n", diff)
			}
			if _, ok := m["key"]; ok {
				t.Error("failed to remove key from map")
			}
		})
	}
}

func TestStrArrNokey(t *testing.T) {
	m := types.J{}
	got := []string{}
	err := ForceParseIntoStrArray(&m, "notkey", &got)
	if err == nil {
		t.Error("Expected no key found error, got none")
	}
}
