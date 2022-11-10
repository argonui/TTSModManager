package tests

import (
	"ModCreator/types"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestTestframework(t *testing.T) {
	for _, tc := range []struct {
		input       map[string]types.J
		path        string
		wantFiles   []string
		wantFolders []string
	}{
		{
			input: map[string]types.J{
				"foo.json": types.J{},
			},
			path:        "",
			wantFiles:   []string{"foo.json"},
			wantFolders: []string{},
		},
		{
			input: map[string]types.J{
				"foo.json":      types.J{},
				"foo/bar.json":  types.J{},
				"foo2/bar.json": types.J{},
			},
			path:        "",
			wantFiles:   []string{"foo.json"},
			wantFolders: []string{"foo", "foo2"},
		},
		{
			input: map[string]types.J{
				"foo/bar.json":     types.J{},
				"foo/bar/baz.json": types.J{},
			},
			path:        "foo",
			wantFiles:   []string{"foo/bar.json"},
			wantFolders: []string{"foo/bar"},
		},
	} {
		ff := FakeFiles{
			Data: tc.input,
		}
		gotFiles, gotFolders, err := ff.ListFilesAndFolders(tc.path)
		if err != nil {
			t.Fatalf("ListFilesAndFolders(%s): %v", tc.path, err)
		}
		asset := cmp.Transformer("tomap", func(in []string) map[string]bool {
			out := map[string]bool{}
			for _, i := range in {
				out[i] = true
			}
			return out
		})
		if diff := cmp.Diff(tc.wantFiles, gotFiles, asset); diff != "" {
			t.Errorf("want != got:\n%v\n", diff)
		}
		if diff := cmp.Diff(tc.wantFolders, gotFolders, asset); diff != "" {
			t.Errorf("want != got:\n%v\n", diff)
		}
	}

}
