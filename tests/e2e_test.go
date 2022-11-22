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
		denyList := []string{
			"bundled_core", // currently depends on src file that test doesn't know about
			"bundled_lua",  // currently depends on src file that test doesn't know about
		}

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
				LuaSrcWriter:      objsAndLua,
				ObjWriter:         objsAndLua,
				ObjDirCreeator:    objsAndLua,
				RootWrite:         finalOutput,
			}
			err = r.Write(j)
			if err != nil {
				t.Fatalf("Error reversing : %v", err)
			}

			objsAndLua.DebugFileNames(t.Logf)
			finalOutput.DebugFileNames(t.Logf)
			reversedConfig, _ := finalOutput.ReadObj("config.json")
			t.Logf("%v\n", reversedConfig)

			m := &mod.Mod{
				Lua:         objsAndLua,
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
			// bundler.AnalyzeBundle(want["LuaScript"].(string), t.Logf)
			// bundler.AnalyzeBundle(got["LuaScript"].(string), t.Logf)
			wantBundles, err := bundler.UnbundleAll(want["LuaScript"].(string))
			if err != nil {
				t.Fatalf("unbundle want : %v", err)
			}
			gotBundles, err := bundler.UnbundleAll(got["LuaScript"].(string))
			if err != nil {
				t.Fatalf("unbundle got : %v", err)
			}
			if diff := cmp.Diff(mapOfKeys(wantBundles), mapOfKeys(gotBundles)); diff != "" {
				t.Errorf("want != got:\n%v\n", diff)
			}
			delete(want, "LuaScript")
			delete(got, "LuaScript")

			if diff := cmp.Diff(want, got, cmpopts.IgnoreMapEntries(ignoreUnpredictable)); diff != "" {
				t.Errorf("want != got:\n%v\n", diff)
			}
		})

	}
}

func mapOfKeys(m map[string]string) map[string]interface{} {
	r := map[string]interface{}{}
	for k := range m {
		r[k] = true
	}
	return r
}
