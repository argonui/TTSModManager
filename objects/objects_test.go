package objects

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"

	"ModCreator/types"

	"github.com/google/go-cmp/cmp"
)

type fakeFiles struct {
	fs   map[string]string
	data map[string]types.J
}

func (f *fakeFiles) EncodeFromFile(s string) (string, error) {
	if _, ok := f.fs[s]; !ok {
		return "", fmt.Errorf("fake file <%s> not found", s)
	}
	return f.fs[s], nil
}
func (f *fakeFiles) ReadObj(s string) (map[string]interface{}, error) {
	if _, ok := f.data[s]; !ok {
		return nil, fmt.Errorf("fake file <%s> not found", s)
	}
	return f.data[s], nil
}
func (f *fakeFiles) ReadObjArray(s string) ([]map[string]interface{}, error) {
	return nil, fmt.Errorf("unimplemented")
}
func (f *fakeFiles) WriteObj(data map[string]interface{}, path string) error {
	f.data[path] = data
	return nil
}
func (f *fakeFiles) WriteObjArray(data []map[string]interface{}, path string) error {
	return fmt.Errorf("unimplemented")
}
func (f *fakeFiles) EncodeToFile(script, file string) error {
	f.fs[file] = script
	return nil
}
func (f *fakeFiles) CreateDir(a, b string) (string, error) {
	// return the "chosen" directory name for next folder
	return b, nil
}

func TestObjPrintToFile(t *testing.T) {
	for _, tc := range []struct {
		o            *objConfig
		want         types.J
		wantFilename string
	}{
		{
			o: &objConfig{
				guid: "123456",
				data: types.J{
					"GUID": "123456",
				},
			},
			wantFilename: "123456.json",
			want: types.J{
				"GUID":          "123456",
				"tts_mod_order": int64(0),
			},
		},
	} {
		ff := &fakeFiles{
			fs:   map[string]string{},
			data: map[string]types.J{},
		}
		err := tc.o.printToFile("path/to/", ff, ff, ff)
		if err != nil {
			t.Errorf("printing %v, got %v", tc.o, err)
		}
		if diff := cmp.Diff(tc.want, ff.data["path/to/"+tc.wantFilename]); diff != "" {
			t.Errorf("want != got:\n%v\n", diff)
		}
	}
}

func TestObjPrinting(t *testing.T) {
	for _, tc := range []struct {
		o    *objConfig
		want types.J
	}{
		{
			o: &objConfig{
				guid: "123456",
				data: types.J{
					"GUID": "123456",
				},
				subObj: []*objConfig{
					{
						guid: "1234561",
						data: types.J{
							"GUID": "1234561",
						},
					}, {
						guid: "1234562",
						data: types.J{
							"GUID": "1234562",
						},
					}, {
						guid: "1234563",
						data: types.J{
							"GUID": "1234563",
						},
					},
				},
			},
			want: types.J{
				"GUID": "123456",
				"ContainedObjects": []types.J{
					{"GUID": "1234561"},
					{"GUID": "1234562"},
					{"GUID": "1234563"},
				},
			},
		}, {
			o: &objConfig{
				guid: "123456",
				data: types.J{
					"GUID": "123456",
				},
				subObjDir: "Foo.123456",
			},
			want: types.J{
				"GUID": "123456",
			},
		},
	} {
		l := &fakeFiles{
			fs: map[string]string{
				"core/AgendaDeck.ttslua": "var a = 42;",
			},
		}
		got, err := tc.o.print(l)
		if err != nil {
			t.Errorf("printing %v, got %v", tc.o, err)
		}
		if diff := cmp.Diff(tc.want, got); diff != "" {
			t.Errorf("want != got:\n%v\n", diff)
		}
	}
}

func TestObjPrintingToFile(t *testing.T) {
	type fileContent struct {
		file, content string
	}
	type jsonContent struct {
		file    string
		content types.J
	}
	for _, tc := range []struct {
		o        *objConfig
		folder   string
		wantObjs []jsonContent
		wantLSS  fileContent
	}{
		{
			o: &objConfig{
				guid: "123456",
				data: types.J{
					"GUID": "123456",
				},
			},
			folder: "foo",
			wantObjs: []jsonContent{
				{
					file: "foo/123456.json",
					content: types.J{
						"GUID":          "123456",
						"tts_mod_order": int64(0),
					},
				},
			},
		},
		{
			o: &objConfig{
				guid: "123777",
				data: types.J{
					"GUID": "123777",
				},
				subObj: []*objConfig{
					&objConfig{
						guid: "1237770",
						data: types.J{
							"Nickname": "coolobj",
							"GUID":     "1237770",
						},
					},
					&objConfig{
						guid: "1237771",
						data: types.J{
							"GUID": "1237771",
						},
					},
					&objConfig{
						guid: "1237772",
						data: types.J{
							"GUID": "1237772",
						},
					},
				},
			},
			folder: "bar",
			wantObjs: []jsonContent{
				{
					file: "bar/123777.json",
					content: types.J{
						"GUID":                  "123777",
						"ContainedObjects_path": "123777",
						"tts_mod_order":         int64(0),
					},
				},
				{
					file: "bar/123777/coolobj.1237770.json",
					content: types.J{
						"GUID":          "1237770",
						"tts_mod_order": int64(0),
						"Nickname":      "coolobj",
					},
				},
				{
					file: "bar/123777/1237771.json",
					content: types.J{
						"tts_mod_order": int64(1),
						"GUID":          "1237771",
					},
				},
				{
					file: "bar/123777/1237772.json",
					content: types.J{
						"tts_mod_order": int64(2),
						"GUID":          "1237772",
					},
				},
			},
		},
		{
			o: &objConfig{
				guid: "123456",
				data: types.J{
					"LuaScriptState": "fav color = green",
					"GUID":           "123456",
				},
			},
			folder: "foo",
			wantObjs: []jsonContent{
				{
					file: "foo/123456.json",
					content: types.J{
						"GUID":           "123456",
						"LuaScriptState": "fav color = green",
						"tts_mod_order":  int64(0),
					},
				},
			},
			// want no LSS file because it's short
		},
		{
			o: &objConfig{
				guid: "123456",
				data: types.J{
					"LuaScriptState": "fav color = green fav color = green fav color = green fav color = green fav color = green",
					"GUID":           "123456",
				},
			},
			folder: "foo",
			wantObjs: []jsonContent{
				{
					file: "foo/123456.json",
					content: types.J{
						"GUID":                "123456",
						"LuaScriptState_path": "foo/123456.luascriptstate",
						"tts_mod_order":       int64(0),
					},
				},
			},
			wantLSS: fileContent{
				file:    "foo/123456.luascriptstate",
				content: "fav color = green fav color = green fav color = green fav color = green fav color = green",
			},
		},
	} {
		ff := &fakeFiles{
			fs:   map[string]string{},
			data: map[string]types.J{},
		}
		err := tc.o.printToFile(tc.folder, ff, ff, ff)
		if err != nil {
			t.Errorf("printing %v, got %v", tc.o, err)
		}

		for _, wantFJ := range tc.wantObjs {
			got, ok := ff.data[wantFJ.file]
			if !ok {
				log.Printf("%v\n", ff.data)
				t.Fatalf("Wanted file %s, not found", wantFJ.file)
			}
			if diff := cmp.Diff(wantFJ.content, got); diff != "" {
				t.Errorf("want != got:\n%v\n", diff)
			}
		}

		// compare lua script state
		if tc.wantLSS.file != "" {
			got, ok := ff.fs[tc.wantLSS.file]
			if !ok {
				log.Printf("fs: %v", ff.fs)
				log.Printf("data %v", ff.data)
				t.Errorf("wanted luascript state %s, didn't find", tc.wantLSS.file)
			}
			if diff := cmp.Diff(tc.wantLSS.content, got); diff != "" {
				t.Errorf("want != got:\n%v\n", diff)
			}
		}
	}
}

func TestName(t *testing.T) {
	for _, tc := range []struct {
		data            types.J
		guid, subObjDir string
		want            string
	}{
		{
			data: types.J{
				"Nickname": "Occult Invocation",
				"Name":     "Card",
			},
			guid: "010509",
			want: "OccultInvocation.010509",
		}, {
			data: types.J{
				"Nickname": "Occult Invocation!!!!~~()()",
				"Name":     "Card",
			},
			guid: "010509",
			want: "OccultInvocation.010509",
		}, {
			data: types.J{
				"Name": "Card",
			},
			guid: "1233",
			want: "Card.1233",
		}, {
			data: types.J{
				"Nickname": nil,
				"Name":     "Card",
			},
			guid: "1235",
			want: "Card.1235",
		}, {
			data: types.J{
				"Nickname": "",
				"Name":     "Card",
			},
			guid: "1234",
			want: "Card.1234",
		}, {
			data: types.J{
				"Nickname": "",
				"Name":     nil,
			},
			guid: "010508",
			want: "010508",
		}, {
			data: types.J{
				"Nickname": nil,
				"Name":     nil,
			},
			guid: "010509",
			want: "010509",
		}, {
			data: types.J{},
			guid: "010507",
			want: "010507",
		}, {
			data: types.J{
				"Nickname": "",
				"Name":     "",
			},
			guid: "010506",
			want: "010506",
		}, {
			data: types.J{
				"Nickname": "",
				"Name":     "",
			},
			subObjDir: "010504_2",
			guid:      "010504",
			want:      "010504",
		}, {
			data: types.J{
				"Nickname": "foo",
				"Name":     "",
			},
			subObjDir: "010503_2",
			guid:      "010503",
			want:      "foo.010503",
		}, {
			data: types.J{
				"Nickname": "foo",
				"Name":     "",
			},
			subObjDir: "foo.010502_2",
			guid:      "010502",
			want:      "foo.010502",
		},
	} {
		o := objConfig{
			data:      tc.data,
			guid:      tc.guid,
			subObjDir: tc.subObjDir,
		}
		got := o.getAGoodFileName()
		if tc.want != got {
			t.Errorf("want <%s> got <%s>", tc.want, got)
		}
	}

}

func TestPrintAllObjs(t *testing.T) {
	type wantFile struct {
		name    string
		content types.J
	}
	for _, tc := range []struct {
		objs  []map[string]interface{}
		wants []wantFile
	}{
		{
			objs: []map[string]interface{}{
				{"GUID": "123456"},
				{"GUID": "123457"},
			},
			wants: []wantFile{
				{
					name:    "123456.json",
					content: types.J{"GUID": "123456", "tts_mod_order": int64(0)},
				}, {
					name:    "123457.json",
					content: types.J{"GUID": "123457", "tts_mod_order": int64(1)},
				},
			},
		},
	} {
		ff := &fakeFiles{
			data: map[string]types.J{},
			fs:   map[string]string{},
		}
		err := PrintObjectStates("", ff, ff, ff, tc.objs)
		if err != nil {
			t.Fatalf("error not expected %v", err)
		}

		for _, w := range tc.wants {
			got, ok := ff.data[w.name]
			if !ok {
				t.Errorf("wanted filename %s not present in data", w.name)
			}
			if diff := cmp.Diff(w.content, got); diff != "" {
				t.Errorf("want != got:\n%v\n", diff)
			}
		}
	}

}

func TestDBPrint(t *testing.T) {
	ff := &fakeFiles{
		fs:   map[string]string{},
		data: map[string]types.J{},
	}
	for _, tc := range []struct {
		root []*objConfig
		want types.ObjArray
	}{
		{
			root: []*objConfig{
				&objConfig{
					data:  types.J{"GUID": "123"},
					order: int64(3),
				},
				&objConfig{
					data:  types.J{"GUID": "121"},
					order: int64(1),
				},
				&objConfig{
					data:  types.J{"GUID": "122"},
					order: int64(2),
				},
			},
			want: types.ObjArray{
				{"GUID": "121"},
				{"GUID": "122"},
				{"GUID": "123"},
			},
		},
	} {
		db := db{
			root: tc.root,
		}
		got, err := db.print(ff)
		if err != nil {
			t.Fatalf("got unexpected err %v", err)
		}
		if diff := cmp.Diff(tc.want, got); diff != "" {
			t.Errorf("want != got:\n%v\n", diff)
		}

	}
}

func jn(i int) json.Number {
	return json.Number(fmt.Sprint(i))
}

func TestParseFromFile(t *testing.T) {

	ff := &fakeFiles{
		data: map[string]types.J{},
	}
	for _, tc := range []struct {
		name  string
		input types.J
		want  objConfig
	}{
		{
			name: "mod order",
			input: types.J{
				"GUID": "123",
			},
			want: objConfig{
				guid: "123",
				data: types.J{"GUID": "123"},
			},
		},
	} {
		ff.data[tc.name] = tc.input
		o := objConfig{}
		err := o.parseFromFile(tc.name, ff)
		if err != nil {
			t.Fatalf("failed to preset data in %s\n", tc.name)
		}
		if diff := cmp.Diff(tc.want.data, o.data); diff != "" {
			t.Errorf("want != got:\n%v\n", diff)
		}
		if tc.want.guid != o.guid {
			t.Errorf("guid mismatch want %s got %s", tc.want.guid, o.guid)
		}
	}
}
