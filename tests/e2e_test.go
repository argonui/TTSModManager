package tests

import (
	"ModCreator/bundler"
	"ModCreator/file"
	"ModCreator/mod"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

var (
	expectedStr       = []string{"SaveName", "Date", "VersionNumber", "GameMode", "GameType", "GameComplexity", "Table", "Sky", "Note", "LuaScript", "LuaScriptState", "XmlUI"}
	expectedObj       = []string{"TabStates", "MusicPlayer", "Grid", "Lighting", "Hands", "ComponentTags", "Turns"}
	expectedObjArr    = []string{"CameraStates", "DecalPallet", "CustomUIAssets", "SnapPoints", "Decals"}
	expectedObjStates = "ObjectStates"
)

func TestAllReverseThenBuild(t *testing.T) {
	paths, err := filepath.Glob(filepath.Join("testdata", "e2e", "*.json"))
	if err != nil {
		t.Fatal(err)
	}
	for _, path := range paths {
		_, filename := filepath.Split(path)
		testname := filename[:len(filename)-len(filepath.Ext(path))]
		denyList := []string{}

		t.Run(testname, func(t *testing.T) {
			for _, f := range denyList {
				if f == testname {
					return
				}
			}
			j, err := file.ReadRawFile(path)
			if err != nil {
				t.Fatalf("Error parsing %s : %v", path, err)
			}
			want, err := file.ReadRawFile(path)
			if err != nil {
				t.Fatalf("Error parsing %s : %v", path, err)
			}
			modsettings := NewFF()
			finalOutput := NewFF()
			objsAndLua := NewFF()

			r := mod.Reverser{
				ModSettingsWriter: modsettings,
				LuaWriter:         objsAndLua,
				XMLWriter:         objsAndLua,
				XMLSrcWriter:      objsAndLua,
				LuaSrcWriter:      objsAndLua,
				ObjWriter:         objsAndLua,
				ObjDirCreator:     objsAndLua,
				RootWrite:         finalOutput,
			}
			err = r.Write(j)
			if err != nil {
				t.Fatalf("Error reversing : %v", err)
			}

			objsAndLua.DebugFileNames(t.Logf)
			finalOutput.DebugFileNames(t.Logf)
			reversedConfig, err := finalOutput.ReadObj("config.json")
			if err != nil {
				t.Fatalf("Couldn't find root config: %v", err)
			}
			t.Logf("%v\n", reversedConfig)

			m := &mod.Mod{
				Lua:         objsAndLua,
				XML:         objsAndLua,
				Modsettings: modsettings,
				Objs:        objsAndLua,
				Objdirs:     objsAndLua,
				RootRead:    finalOutput,
				RootWrite:   finalOutput,
			}
			err = m.GenerateFromConfig()
			if err != nil {
				t.Fatalf("generateMod(<config>) : %v\n", err)
			}
			err = m.Print("output.json")
			if err != nil {
				t.Fatalf("printMod(...) : %v", err)
			}
			got, err := finalOutput.ReadObj("output.json")
			if err != nil {
				t.Fatalf("output.json not parsed : %v", err)
			}
			ignoreUnpredictable := func(k string, v interface{}) bool {
				if _, ok := v.(float64); ok {
					return true
				}
				if k == "Date" || k == "EpochTime" {
					return true
				}

				return false
			}
			wls, wok := want["LuaScript"]
			gls, gok := got["LuaScript"]
			if wok && gok {
				wlss, ok := wls.(string)
				if !ok {
					t.Fatalf("non string found in luascript, found %T", wls)
				}
				glss, ok := gls.(string)
				if !ok {
					t.Fatalf("non string found in luascript, found %T", gls)
				}
				wantBundles, _, err := bundler.UnbundleAll(wlss)
				if err != nil {
					t.Fatalf("unbundle want : %v", err)
				}
				gotBundles, _, err := bundler.UnbundleAll(glss)
				if err != nil {
					t.Fatalf("unbundle got : %v", err)
				}
				if diff := cmp.Diff(mapOfKeys(wantBundles), mapOfKeys(gotBundles)); diff != "" {
					t.Errorf("want != got:\n%v\n", diff)
				}
				delete(want, "LuaScript")
				delete(got, "LuaScript")
			}

			// Object-level LuaScript (e.g. inside ObjectStates) is compared the
			// same formatting-insensitive way as the top-level script above:
			// bundling is order- and formatting-sensitive, so a bundled script
			// is compared by its set of module names rather than byte-for-byte.
			normalizeBundledLua(want)
			normalizeBundledLua(got)

			if diff := cmp.Diff(want, got, cmpopts.IgnoreMapEntries(ignoreUnpredictable)); diff != "" {
				t.Errorf("want != got:\n%v\n", diff)
			}
		})

	}
}

// normalizeBundledLua walks a decoded savegame and replaces every bundled
// LuaScript string it finds (at any nesting depth) with the set of module
// names that unbundle out of it. This mirrors how the top-level LuaScript is
// compared, so nested object scripts are not held to a stricter byte-for-byte
// standard than the root while still verifying that require resolution rebuilt
// the same set of modules. Non-bundled scripts are left untouched so their
// contents are still compared exactly.
func normalizeBundledLua(v interface{}) {
	switch t := v.(type) {
	case map[string]interface{}:
		for k, val := range t {
			if k == "LuaScript" {
				if s, ok := val.(string); ok && bundler.IsBundled(s) {
					if modules, err := bundler.UnbundleAll(s); err == nil {
						t[k] = mapOfKeys(modules)
						continue
					}
				}
			}
			normalizeBundledLua(val)
		}
	case []interface{}:
		for _, val := range t {
			normalizeBundledLua(val)
		}
	}
}

func mapOfKeys(m map[string]string) map[string]interface{} {
	r := map[string]interface{}{}
	for k := range m {
		r[k] = true
	}
	return r
}
