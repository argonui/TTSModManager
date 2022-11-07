package objects

import (
	"ModCreator/types"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestFileToJson(t *testing.T) {
	for _, tc := range []struct {
		jsonfiles map[string]types.J
		txtfiles  map[string]string
		want      types.J
	}{
		{
			jsonfiles: map[string]types.J{
				"foo/cool.123456.json": types.J{
					"GUID": "123456",
				},
			},
			want: types.J{
				"GUID": "123456",
			},
		}, {
			jsonfiles: map[string]types.J{
				"foo/cool.123456.json": types.J{
					"GUID":                  "123457",
					"Nickname":              "cool",
					"ContainedObjects_path": "cool.123456",
				},
				"foo/cool.123456/1.json": types.J{
					"GUID":          "1",
					"tts_mod_order": int64(1),
				},
				"foo/cool.123456/2.json": types.J{
					"tts_mod_order": int64(0),
					"GUID":          "2",
				},
			},
			want: types.J{
				"GUID":     "123457",
				"Nickname": "cool",
				"ContainedObjects": []types.J{
					{"GUID": "2"},
					{"GUID": "1"},
				},
			},
		},
		{
			jsonfiles: map[string]types.J{
				"foo/cool.123456.json": types.J{
					"GUID":                  "123457",
					"Nickname":              "cool",
					"ContainedObjects_path": "cool.123456",
				},
				"foo/cool.123456/1.json": types.J{
					"GUID":                "1",
					"tts_mod_order":       int64(1),
					"LuaScriptState_path": "foo/cool.123456/1.luascriptstate",
				},
				"foo/cool.123456/2.json": types.J{
					"tts_mod_order": int64(0),
					"GUID":          "2",
				},
			},
			txtfiles: map[string]string{
				"foo/cool.123456/1.luascriptstate": "state = up",
			},
			want: types.J{
				"GUID":     "123457",
				"Nickname": "cool",
				"ContainedObjects": []types.J{
					{"GUID": "2"},
					{"GUID": "1",
						"LuaScriptState": "state = up",
					},
				},
			},
		},
	} {
		ff := &fakeFiles{
			fs:   tc.txtfiles,
			data: tc.jsonfiles,
		}
		o := objConfig{}
		err := o.parseFromFile("foo/cool.123456.json", ff, ff)
		if err != nil {
			t.Fatalf("Expected no error parsing from %s: got %v", "foo/cool.123456.json", err)
		}

		got, err := o.print(ff)
		if err != nil {
			t.Errorf("Expected no error printing %s: got %v", o.guid, err)
		}
		if diff := cmp.Diff(tc.want, got); diff != "" {
			t.Errorf("want != got:\n%v\n", diff)
		}

	}

}
