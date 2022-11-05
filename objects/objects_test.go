package objects

import (
	"fmt"
	"log"
	"path"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type fakeFiles struct {
	fs   map[string]string
	data map[string]j
}

func (f *fakeFiles) EncodeFromFile(s string) (string, error) {
	if _, ok := f.fs[s]; !ok {
		return "", fmt.Errorf("fake file <%s> not found", s)
	}
	return f.fs[s], nil
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
	return a + b, nil
}

func TestObjPrintToFile(t *testing.T) {
	for _, tc := range []struct {
		o            *objConfig
		want         j
		wantFilename string
	}{
		{
			o: &objConfig{
				guid: "123456",
				data: j{
					"GUID": "123456",
				},
			},
			wantFilename: "123456.json",
			want: j{
				"GUID": "123456",
			},
		},
	} {
		ff := &fakeFiles{
			fs:   map[string]string{},
			data: map[string]j{},
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
		want j
	}{
		{
			o: &objConfig{
				guid: "123456",
				data: j{
					"GUID": "123456",
				},
			},
			want: j{
				"GUID": "123456",
			},
		}, {
			o: &objConfig{
				guid: "123456",
				data: j{
					"GUID": "123456",
				},
				subObjDir: "Foo.123456",
			},
			want: j{
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
	for _, tc := range []struct {
		o       *objConfig
		folder  string
		want    j
		wantLSS fileContent
	}{
		{
			o: &objConfig{
				guid: "123456",
				data: j{
					"GUID": "123456",
				},
			},
			folder: "foo",
			want: j{
				"GUID": "123456",
			},
		}, {
			o: &objConfig{
				guid: "123456",
				data: j{
					"LuaScriptState": "fav color = green",
					"GUID":           "123456",
				},
			},
			folder: "foo",
			want: j{
				"GUID":           "123456",
				"LuaScriptState": "fav color = green",
			},
			// want no LSS file because it's short
		}, {
			o: &objConfig{
				guid: "123456",
				data: j{
					"LuaScriptState": "fav color = green fav color = green fav color = green fav color = green fav color = green",
					"GUID":           "123456",
				},
			},
			folder: "foo",
			want: j{
				"GUID":                "123456",
				"LuaScriptState_path": "foo/123456.luascriptstate",
			},
			wantLSS: fileContent{
				file:    "foo/123456.luascriptstate",
				content: "fav color = green fav color = green fav color = green fav color = green fav color = green",
			},
		},
	} {
		ff := &fakeFiles{
			fs:   map[string]string{},
			data: map[string]j{},
		}
		err := tc.o.printToFile(tc.folder, ff, ff, ff)
		if err != nil {
			t.Errorf("printing %v, got %v", tc.o, err)
		}
		got, ok := ff.data[path.Join(tc.folder, tc.o.getAGoodFileName()+".json")]
		if !ok {
			log.Printf("%v\n", ff.data)
			t.Fatalf("data not found in fake files as expected")

		}
		if diff := cmp.Diff(tc.want, got); diff != "" {
			t.Errorf("want != got:\n%v\n", diff)
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
		data            j
		guid, subObjDir string
		want            string
	}{
		{
			data: j{
				"Nickname": "Occult Invocation",
				"Name":     "Card",
			},
			guid: "010509",
			want: "OccultInvocation.010509",
		}, {
			data: j{
				"Nickname": "Occult Invocation!!!!~~()()",
				"Name":     "Card",
			},
			guid: "010509",
			want: "OccultInvocation.010509",
		}, {
			data: j{
				"Name": "Card",
			},
			guid: "1233",
			want: "Card.1233",
		}, {
			data: j{
				"Nickname": nil,
				"Name":     "Card",
			},
			guid: "1235",
			want: "Card.1235",
		}, {
			data: j{
				"Nickname": "",
				"Name":     "Card",
			},
			guid: "1234",
			want: "Card.1234",
		}, {
			data: j{
				"Nickname": "",
				"Name":     nil,
			},
			guid: "010508",
			want: "010508",
		}, {
			data: j{
				"Nickname": nil,
				"Name":     nil,
			},
			guid: "010509",
			want: "010509",
		}, {
			data: j{},
			guid: "010507",
			want: "010507",
		}, {
			data: j{
				"Nickname": "",
				"Name":     "",
			},
			guid: "010506",
			want: "010506",
		}, {
			data: j{
				"Nickname": "",
				"Name":     "",
			},
			subObjDir: "010504_2",
			guid:      "010504",
			want:      "010504",
		}, {
			data: j{
				"Nickname": "foo",
				"Name":     "",
			},
			subObjDir: "010503_2",
			guid:      "010503",
			want:      "foo.010503",
		}, {
			data: j{
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
