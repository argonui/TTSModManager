package objects

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"

	"ModCreator/tests"
	"ModCreator/types"

	"github.com/google/go-cmp/cmp"
)

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
				"GUID": "123456",
			},
		},
		{
			o: &objConfig{
				guid: "123456",
				data: types.J{
					"GUID": "123456",
					"XmlUI": `<Button id="bar"/>
					<Button id="bar"/>
					<Button id="bar"/>
					<Button id="bar"/>
					<Button id="bar"/>
					<Button id="bar"/>
					<Button id="bar"/>
					`,
				},
			},
			wantFilename: "123456.json",
			want: types.J{
				"GUID":       "123456",
				"XmlUI_path": "path/to/123456.xml",
			},
		},
		{
			o: &objConfig{
				guid: "123457",
				data: types.J{
					"GUID": "123457",
					"XmlUI": `<!-- include buttons -->
					<Button id="bar"/>
					<Button id="bar"/>
					<Button id="bar"/>
					<Button id="bar"/>
					<Button id="bar"/>
					<Button id="bar"/>
					<Button id="bar"/>
<!-- include buttons -->
`,
				},
			},
			wantFilename: "123457.json",
			want: types.J{
				"GUID":  "123457",
				"XmlUI": "<Include src=\"buttons\"/>\n",
			},
		},
	} {
		t.Run(tc.wantFilename, func(t *testing.T) {
			ff := tests.NewFF()
			p := &Printer{
				Lua: ff,
				XML: ff,
				Dir: ff,
				J:   ff,
			}
			err := tc.o.printToFile("path/to/", p)
			if err != nil {
				t.Fatalf("printing %v, got %v", tc.o, err)
			}
			ff.DebugFileNames(t.Logf)
			if diff := cmp.Diff(tc.want, ff.Data["path/to/"+tc.wantFilename]); diff != "" {
				t.Errorf("want != got:\n%v\n", diff)
			}
		})

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
		l := &tests.FakeFiles{
			Fs: map[string]string{
				"core/AgendaDeck.ttslua": "var a = 42;",
			},
		}
		got, err := tc.o.print(l, l)
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
		file    string
		content string
	}
	type jsonContent struct {
		file    string
		content types.J
	}
	for _, tc := range []struct {
		o             *objConfig
		folder        string
		wantObjs      []jsonContent
		wantLss       fileContent
		wantLssAsJson jsonContent
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
						"GUID": "123456",
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
				subObjOrder: []string{
					"coolobj.1237770", "1237771", "1237772",
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
						"GUID":                   "123777",
						"ContainedObjects_path":  "123777",
						"ContainedObjects_order": []string{"coolobj.1237770", "1237771", "1237772"},
					},
				},
				{
					file: "bar/123777/coolobj.1237770.json",
					content: types.J{
						"GUID":     "1237770",
						"Nickname": "coolobj",
					},
				},
				{
					file: "bar/123777/1237771.json",
					content: types.J{
						"GUID": "1237771",
					},
				},
				{
					file: "bar/123777/1237772.json",
					content: types.J{
						"GUID": "1237772",
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
					},
				},
			},
			wantLss: fileContent{
				file:    "foo/123456.luascriptstate",
				content: "fav color = green fav color = green fav color = green fav color = green fav color = green",
			},
		},
		{
			o: &objConfig{
				guid: "123456",
				data: types.J{
					"LuaScriptState": "{\"acknowledgedUpgradeVersions\":[],\"optionPanel\":{\"cardLanguage\":\"en\",\"changePlayAreaImage\":false,\"playAreaConnectionColor\":{\"a\":1,\"b\":0.4,\"g\":0.4,\"r\":0.4},\"useResourceCounters\":\"disabled\"}}",
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
					},
				},
			},
			wantLssAsJson: jsonContent{
				file: "foo/123456.luascriptstate",
				content: types.J{
					"acknowledgedUpgradeVersions": []any{},
					"optionPanel": map[string]any{
						"cardLanguage": string("en"),
						"changePlayAreaImage": bool(false),
						"playAreaConnectionColor": map[string]any{
							"a": float64(1),
							"b": float64(0.4),
							"g": float64(0.4),
							"r": float64(0.4),
						},
						"useResourceCounters": string("disabled"),
					},
				},
			},
		},
	} {
		ff := tests.NewFF()
		p := &Printer{
			Lua: ff,
			Dir: ff,
			J:   ff,
		}
		err := tc.o.printToFile(tc.folder, p)
		if err != nil {
			t.Errorf("printing %v, got %v", tc.o, err)
		}

		for _, wantFJ := range tc.wantObjs {
			got, ok := ff.Data[wantFJ.file]
			if !ok {
				log.Printf("%v\n", ff.Data)
				t.Fatalf("Wanted file %s, not found", wantFJ.file)
			}
			if diff := cmp.Diff(wantFJ.content, got); diff != "" {
				t.Errorf("want != got:\n%v\n", diff)
			}
		}

		// compare lua script state
		if tc.wantLss.file != "" {
			got, ok := ff.Fs[tc.wantLss.file]
			if !ok {
				ff.DebugFileNames(t.Logf)
				t.Errorf("wanted luascript state %s, didn't find", tc.wantLss.file)
			}
			if diff := cmp.Diff(tc.wantLss.content, got); diff != "" {
				t.Errorf("want != got:\n%v\n", diff)
			}
		}

		// compare lua script state that should output as json
		if tc.wantLssAsJson.file != "" {
			got, ok := ff.Data[tc.wantLssAsJson.file]
			if !ok {
				ff.DebugFileNames(t.Logf)
				t.Errorf("wanted luascript state %s, didn't find", tc.wantLssAsJson.file)
			}
			if diff := cmp.Diff(tc.wantLssAsJson.content, got); diff != "" {
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
			want: "OccultInvocation!!!!.010509",
		}, {
				data: types.J{
				"Nickname": "Verschwörung der Äxte!",
				"Name":     "Card",
			},
			guid: "010511",
			want: "VerschwörungderÄxte!.010511",
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
		objs       []map[string]interface{}
		wants      []wantFile
		wantsOrder []string
	}{
		{
			objs: []map[string]interface{}{
				{"GUID": "123456"},
				{"GUID": "123457"},
			},
			wants: []wantFile{
				{
					name:    "123456.json",
					content: types.J{"GUID": "123456"},
				}, {
					name:    "123457.json",
					content: types.J{"GUID": "123457"},
				},
			},
			wantsOrder: []string{
				"123456", "123457",
			},
		},
	} {
		ff := tests.NewFF()
		p := &Printer{
			Lua: ff,
			Dir: ff,
			J:   ff,
		}
		gotOrder, err := p.PrintObjectStates("", tc.objs)
		if err != nil {
			t.Fatalf("error not expected %v", err)
		}
		if diff := cmp.Diff(tc.wantsOrder, gotOrder); diff != "" {
			t.Errorf("want != got:\n%v\n", diff)
		}
		for _, w := range tc.wants {
			got, ok := ff.Data[w.name]
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
	ff := tests.NewFF()
	for _, tc := range []struct {
		root       map[string]*objConfig
		orderInput []string
		want       types.ObjArray
	}{
		{
			root: map[string]*objConfig{
				"123": &objConfig{
					data: types.J{"GUID": "123"},
				},
				"121": &objConfig{
					data: types.J{"GUID": "121"},
				},
				"122": &objConfig{
					data: types.J{"GUID": "122"},
				},
			},
			orderInput: []string{
				"121", "122", "123",
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
		got, err := db.print(ff, ff, tc.orderInput)
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

	ff := &tests.FakeFiles{
		Data: map[string]types.J{},
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
		ff.Data[tc.name] = tc.input
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
