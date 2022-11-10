package objects

import (
	"ModCreator/tests"
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
					"GUID":                   "123457",
					"Nickname":               "cool",
					"ContainedObjects_path":  "cool.123456",
					"ContainedObjects_order": []interface{}{"2", "1"},
				},
				"foo/cool.123456/1.json": types.J{
					"GUID": "1",
				},
				"foo/cool.123456/2.json": types.J{
					"GUID": "2",
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
					"GUID":                   "123457",
					"Nickname":               "cool",
					"ContainedObjects_path":  "cool.123456",
					"ContainedObjects_order": []interface{}{"2", "1"},
				},
				"foo/cool.123456/1.json": types.J{
					"GUID":                "1",
					"LuaScriptState_path": "foo/cool.123456/1.luascriptstate",
				},
				"foo/cool.123456/2.json": types.J{
					"GUID": "2",
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
		ff := &tests.FakeFiles{
			Fs:   tc.txtfiles,
			Data: tc.jsonfiles,
		}
		o := objConfig{}
		err := o.parseFromFile("foo/cool.123456.json", ff)
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

func TestJsonToFiles(t *testing.T) {
	for _, tc := range []struct {
		wantJsonfiles map[string]types.J
		wantTxtfiles  map[string]string
		input         types.J
		relpath       string
	}{
		{
			relpath: "foo",
			input: types.J{
				"GUID":   "123",
				"foobar": "baz",
			},
			wantJsonfiles: map[string]types.J{
				"foo/123.json": types.J{
					"GUID":   "123",
					"foobar": "baz",
				},
			},
			wantTxtfiles: map[string]string{},
		}, {
			relpath: "objrounding",
			input: types.J{
				"GUID":   "123",
				"foobar": "baz",
				"AltLookAngle": types.J{
					"x": 0.123456789,
					"y": -0.123456789,
					"z": float64(370),
				},
				"AttachedSnapPoints": []map[string]interface{}{
					{
						"Position": types.J{
							"x": -1.82239926,
							"y": 0.100342259,
							"z": 0.6163123,
						},
						"Rotation": types.J{
							"x": 3.36023078e-7,
							"y": 0.008230378,
							"z": -2.29263165e-7,
						},
					},
					{
						"Rotation": types.J{
							"x": 3.36023078e-7,
							"y": 0.008230378,
							"z": -2.29263165e-7,
						},
					},
					{
						"Position": types.J{
							"x": -1.82239926,
							"y": 0.100342259,
							"z": 0.6163123,
						},
					},
				},
			},
			wantJsonfiles: map[string]types.J{
				"objrounding/123.json": types.J{
					"GUID":   "123",
					"foobar": "baz",
					"AltLookAngle": types.J{
						"x": float64(0),
						"y": float64(0),
						"z": float64(10),
					},
					"AttachedSnapPoints": []map[string]interface{}{
						{
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
						{
							"Rotation": types.J{
								"x": float64(0),
								"y": float64(0),
								"z": float64(0),
							},
						},
						{
							"Position": types.J{
								"x": -1.822,
								"y": 0.100,
								"z": 0.616,
							},
						},
					},
				},
			},
			wantTxtfiles: map[string]string{},
		}, {
			relpath: "foo2",
			input: types.J{
				"GUID":           "123",
				"foobar":         "baz",
				"LuaScriptState": "var two = 2;var two = 2;var two = 2;var two = 2;var two = 2;var two = 2;var two = 2;",
			},
			wantJsonfiles: map[string]types.J{
				"foo2/123.json": types.J{
					"GUID":                "123",
					"foobar":              "baz",
					"LuaScriptState_path": "foo2/123.luascriptstate",
				},
			},
			wantTxtfiles: map[string]string{
				"foo2/123.luascriptstate": "var two = 2;var two = 2;var two = 2;var two = 2;var two = 2;var two = 2;var two = 2;",
			},
		}, {
			relpath: "foo",
			input: types.J{
				"GUID":   "123",
				"foobar": "baz",
				"ContainedObjects": []interface{}{
					map[string]interface{}{
						"GUID":     "1231",
						"Nickname": "first class",
					},
					map[string]interface{}{
						"GUID":     "1232",
						"Nickname": "second class",
					},
				},
			},
			wantJsonfiles: map[string]types.J{
				"foo/123.json": types.J{
					"GUID":                   "123",
					"foobar":                 "baz",
					"ContainedObjects_path":  "123",
					"ContainedObjects_order": []string{"firstclass.1231", "secondclass.1232"},
				},
				"foo/123/firstclass.1231.json": types.J{
					"GUID":     "1231",
					"Nickname": "first class",
				},
				"foo/123/secondclass.1232.json": types.J{
					"GUID":     "1232",
					"Nickname": "second class",
				},
			},
			wantTxtfiles: map[string]string{},
		},
	} {
		ff := tests.NewFF()
		o := objConfig{}
		err := o.parseFromJSON(tc.input)
		if err != nil {
			t.Fatalf("parseFromJSON(%s): %v", tc.input, err)
		}
		err = o.printToFile(tc.relpath, ff, ff, ff)
		if err != nil {
			t.Fatalf("printToFile(%s): %v", o.getAGoodFileName(), err)
		}
		if diff := cmp.Diff(tc.wantJsonfiles, ff.Data); diff != "" {
			t.Errorf("want != got:\n%v\n", diff)
		}
		if diff := cmp.Diff(tc.wantTxtfiles, ff.Fs); diff != "" {
			t.Errorf("want != got:\n%v\n", diff)
		}
	}
}
