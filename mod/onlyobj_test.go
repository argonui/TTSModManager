package mod

import (
	"encoding/json"
	"strings"
	"testing"

	"ModCreator/tests"
	"ModCreator/types"

	"github.com/google/go-cmp/cmp"
)

// TestObjectOnlyRoundTrip covers the --objin/--objout "downloadable content"
// feature (OnlyObjState / OnlyObjStates): a single object is decomposed into
// files by reverse and reassembled by generate. Reverse->generate must
// reproduce the original object, including its contained objects and states.
func TestObjectOnlyRoundTrip(t *testing.T) {
	orig := types.J{
		"GUID":           "aaa111",
		"Nickname":       "Rich Bag",
		"LuaScript":      "function onLoad()\n  print(\"hi\")\nend",
		"LuaScriptState": "{\"x\":1}",
		"ContainedObjects": []interface{}{
			types.J{"GUID": "child1", "Nickname": "Card A"},
			types.J{"GUID": "child2", "Nickname": "Card B"},
		},
		"States": types.J{
			"2": types.J{"GUID": "state2", "Nickname": "flipped"},
		},
	}

	ff := tests.NewFF()

	// Reverse (object-only): decompose the single object into files.
	r := Reverser{
		ModSettingsWriter: ff,
		LuaWriter:         ff,
		LuaSrcWriter:      ff,
		XMLWriter:         ff,
		XMLSrcWriter:      ff,
		ObjWriter:         ff,
		ObjDirCreator:     ff,
		RootWrite:         ff,
		OnlyObjState:      "ready.json",
	}
	if err := r.Write(jsonClone(t, orig)); err != nil {
		t.Fatalf("object-only reverse: %v", err)
	}

	// The root object is written as "<sanitized nickname>.<guid>.json" at the
	// top level; find it so we can point generate at it (mirrors how main.go
	// derives OnlyObjStates from the --objin basename).
	var rootFile string
	for k := range ff.Data {
		if strings.HasSuffix(k, ".json") && !strings.Contains(k, "/") {
			rootFile = k
		}
	}
	if rootFile == "" {
		t.Fatalf("reverse wrote no root object file; json files: %v", ffKeys(ff))
	}

	// Generate (object-only): reassemble the single object.
	m := &Mod{
		Lua:           ff,
		XML:           ff,
		Modsettings:   ff,
		Objs:          ff,
		Objdirs:       ff,
		RootRead:      ff,
		RootWrite:     ff,
		OnlyObjStates: rootFile,
	}
	if err := m.GenerateFromConfig(); err != nil {
		t.Fatalf("object-only generate: %v", err)
	}

	// Compare through a JSON round trip so slice/map element types
	// (e.g. []interface{} vs []types.J) don't produce spurious diffs.
	if diff := cmp.Diff(jsonNormalize(t, orig), jsonNormalize(t, m.Data)); diff != "" {
		t.Errorf("object-only round trip differs (-want +got):\n%s", diff)
	}
}

func jsonClone(t *testing.T, v interface{}) map[string]interface{} {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var out map[string]interface{}
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	return out
}

func jsonNormalize(t *testing.T, v interface{}) interface{} {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var out interface{}
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	return out
}

func ffKeys(ff *tests.FakeFiles) []string {
	ks := []string{}
	for k := range ff.Data {
		ks = append(ks, k)
	}
	return ks
}
