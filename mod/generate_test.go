package mod

import (
	"ModCreator/tests"
	"ModCreator/types"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestGenerate(t *testing.T) {
	for _, tc := range []struct {
		name             string
		inputRoot        types.J
		inputModSettings map[string]interface{}
		inputObjs        map[string]types.J
		inputLuaSrc      map[string]string
		inputObjTexts    map[string]string
		flags            map[string]interface{}
		want             map[string]interface{}
	}{
		{
			name: "LuaScript",
			inputRoot: map[string]interface{}{
				"LuaScript": "require(\"core/Global\")",
			},
			inputLuaSrc: map[string]string{
				"core/Global.ttslua": "var a = 42",
			},
			want: map[string]interface{}{
				"LuaScript":      "-- Bundled by luabundle {\"version\":\"1.6.0\"}\nlocal __bundle_require, __bundle_loaded, __bundle_register, __bundle_modules = (function(superRequire)\n\tlocal loadingPlaceholder = {[{}] = true}\n\n\tlocal register\n\tlocal modules = {}\n\n\tlocal require\n\tlocal loaded = {}\n\n\tregister = function(name, body)\n\t\tif not modules[name] then\n\t\t\tmodules[name] = body\n\t\tend\n\tend\n\n\trequire = function(name)\n\t\tlocal loadedModule = loaded[name]\n\n\t\tif loadedModule then\n\t\t\tif loadedModule == loadingPlaceholder then\n\t\t\t\treturn nil\n\t\t\tend\n\t\telse\n\t\t\tif not modules[name] then\n\t\t\t\tif not superRequire then\n\t\t\t\t\tlocal identifier = type(name) == 'string' and '\\\"' .. name .. '\\\"' or tostring(name)\n\t\t\t\t\terror('Tried to require ' .. identifier .. ', but no such module has been registered')\n\t\t\t\telse\n\t\t\t\t\treturn superRequire(name)\n\t\t\t\tend\n\t\t\tend\n\n\t\t\tloaded[name] = loadingPlaceholder\n\t\t\tloadedModule = modules[name](require, loaded, register, modules)\n\t\t\tloaded[name] = loadedModule\n\t\tend\n\n\t\treturn loadedModule\n\tend\n\n\treturn require, loaded, register, modules\nend)(nil)\n__bundle_register(\"__root\", function(require, _LOADED, __bundle_register, __bundle_modules)\nrequire(\"core/Global\")\nend)\n__bundle_register(\"core/Global\", function(require, _LOADED, __bundle_register, __bundle_modules)\nvar a = 42\nend)\nreturn __bundle_require(\"__root\")",
				"CameraStates":   nil,
				"ComponentTags":  map[string]interface{}{},
				"MusicPlayer":    map[string]interface{}{},
				"Sky":            "",
				"TabStates":      map[string]interface{}{},
				"Note":           "",
				"ObjectStates":   []interface{}{},
				"SaveName":       "",
				"Table":          "",
				"LuaScriptState": "",
				"SnapPoints":     nil,
				"XmlUI":          "",
				"Turns":          map[string]interface{}{},
				"VersionNumber":  "",
				"GameMode":       "",
				"GameType":       "",
				"Hands":          map[string]interface{}{},
				"Grid":           map[string]interface{}{},
				"Lighting":       map[string]interface{}{},
				"GameComplexity": "",
				"Decals":         nil,
				"CustomUIAssets": nil,
				"DecalPallet":    nil,
			},
		},
		{
			name: "Object State recursive",
			inputRoot: map[string]interface{}{
				"SaveName":           "cool mod",
				"ObjectStates_order": []interface{}{"parent"},
			},
			inputLuaSrc: map[string]string{
				"parent/eda22b/childstate2.ttslua": "var foo = 42\nvar foo = 42\nvar foo = 42\nvar foo = 42\nvar foo = 42\nvar foo = 42\nvar foo = 42\nvar foo = 42\n",
			},
			inputObjs: map[string]types.J{
				"parent.json": map[string]interface{}{
					"GUID": "parent",
					"States_path": map[string]interface{}{
						"2": "eda22b",
					},
					"ContainedObjects_path": "parent",
				},
				"parent/eda22b.json": map[string]interface{}{
					"GUID":                   "eda22b",
					"Autoraise":              true,
					"ContainedObjects_path":  "eda22b",
					"ContainedObjects_order": []string{"childstate2"},
				},
				"parent/eda22b/childstate2.json": map[string]interface{}{
					"Description":    "child of state 2",
					"GUID":           "childstate2",
					"LuaScript_path": "parent/eda22b/childstate2.ttslua",
				},
			},
			want: map[string]interface{}{
				"SaveName": "cool mod",
				"ObjectStates": []any{
					map[string]any{
						"GUID": "parent",
						"States": map[string]interface{}{
							"2": map[string]interface{}{
								"Autoraise": true,
								"GUID":      "eda22b",
								"ContainedObjects": []any{
									map[string]any{
										"Description": "child of state 2",
										"GUID":        "childstate2",
										"LuaScript":   "var foo = 42\nvar foo = 42\nvar foo = 42\nvar foo = 42\nvar foo = 42\nvar foo = 42\nvar foo = 42\nvar foo = 42\n",
									},
								},
							},
						},
					},
				},
				"CameraStates":   nil,
				"ComponentTags":  map[string]interface{}{},
				"MusicPlayer":    map[string]interface{}{},
				"Sky":            "",
				"TabStates":      map[string]interface{}{},
				"Note":           "",
				"Table":          "",
				"LuaScript":      "",
				"LuaScriptState": "",
				"SnapPoints":     nil,
				"XmlUI":          "",
				"Turns":          map[string]interface{}{},
				"VersionNumber":  "",
				"GameMode":       "",
				"GameType":       "",
				"Hands":          map[string]interface{}{},
				"Grid":           map[string]interface{}{},
				"Lighting":       map[string]interface{}{},
				"GameComplexity": "",
				"Decals":         nil,
				"CustomUIAssets": nil,
				"DecalPallet":    nil,
			},
		},
		{
			name: "Object State and contained object generation",
			inputRoot: map[string]interface{}{
				"SaveName":           "cool mod",
				"ObjectStates_order": []interface{}{"parent"},
			},
			inputLuaSrc: map[string]string{},
			inputObjs: map[string]types.J{
				"parent.json": map[string]interface{}{
					"GUID": "parent",
					"States_path": map[string]interface{}{
						"2": "state2",
					},
					"ContainedObjects_path":  "parent",
					"ContainedObjects_order": []string{"co123"},
				},
				"parent/state2.json": map[string]interface{}{
					"GUID":      "state2",
					"Autoraise": true,
				},
				"parent/co123.json": map[string]interface{}{
					"Description": "contained object",
					"GUID":        "co123",
				},
			},
			want: map[string]interface{}{
				"SaveName": "cool mod",
				"ObjectStates": []any{
					map[string]any{
						"GUID": "parent",
						"States": map[string]interface{}{
							"2": map[string]interface{}{
								"Autoraise": true,
								"GUID":      "state2",
							},
						},
						"ContainedObjects": []any{
							map[string]any{
								"GUID":        "co123",
								"Description": "contained object",
							},
						},
					},
				},
				"CameraStates":   nil,
				"ComponentTags":  map[string]interface{}{},
				"MusicPlayer":    map[string]interface{}{},
				"Sky":            "",
				"TabStates":      map[string]interface{}{},
				"Note":           "",
				"Table":          "",
				"LuaScript":      "",
				"LuaScriptState": "",
				"SnapPoints":     nil,
				"XmlUI":          "",
				"Turns":          map[string]interface{}{},
				"VersionNumber":  "",
				"GameMode":       "",
				"GameType":       "",
				"Hands":          map[string]interface{}{},
				"Grid":           map[string]interface{}{},
				"Lighting":       map[string]interface{}{},
				"GameComplexity": "",
				"Decals":         nil,
				"CustomUIAssets": nil,
				"DecalPallet":    nil,
			},
		},
		{
			name: "Saved Object (Simple)",
			inputObjs: map[string]types.J{
        "test123.json": map[string]interface{}{
					"GUID":        "test123",
					"Description": "A test object",
        },
    	},
			flags: map[string]interface{}{
				"OnlyObjStates": true,
				"SavedObj": true,
			},
			want: map[string]interface{}{
				"SaveName":       "",
				"Date":           "",
				"VersionNumber":  "",
				"GameMode":       "",
				"GameType":       "",
				"GameComplexity": "",
				"Tags":           []any{},
				"Gravity":        0.5,
				"PlayArea":       0.5,
				"Table":          "",
				"Sky":            "",
				"Note":           "",
				"TabStates":      map[string]interface{}{},
				"LuaScript":      "",
				"LuaScriptState": "",
				"XmlUI":          "",
				"ObjectStates": []interface{}{
					map[string]interface{}{
						"GUID":        "test123",
						"Description": "A test object",
					},
				},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			rootff := &tests.FakeFiles{
				Data: map[string]types.J{
					"config.json": tc.inputRoot,
				},
			}
			luaReadff := &tests.FakeFiles{
				Fs: tc.inputLuaSrc,
			}
			msff := &tests.FakeFiles{}
			objs := &tests.FakeFiles{
				Data: tc.inputObjs,
				Fs:   tc.inputObjTexts,
			}
			m := Mod{
				RootRead:    rootff,
				RootWrite:   rootff,
				Lua:         luaReadff,
				Modsettings: msff,
				Objs:        objs,
				Objdirs:     objs,
			}
			if OnlyObjStatesFlag, ok := tc.flags["OnlyObjStates"]; ok && OnlyObjStatesFlag == true {
				m.OnlyObjStates = "test123.json"
			}
			if savedObjFlag, ok := tc.flags["SavedObj"]; ok && savedObjFlag == true {
				m.SavedObj = true
			}
			err := m.GenerateFromConfig()
			if err != nil {
				t.Fatalf("Error reading config %v", err)
			}
			err = m.Print("output.json")
			if err != nil {
				t.Fatalf("Error printing config %v", err)
			}
			got, err := rootff.ReadObj("output.json")
			if err != nil {
				t.Fatalf("Error reading output: %v", err)
			}

			ignoreUnpredictable := func(k string, v interface{}) bool {
				if k == "Date" || k == "EpochTime" {
					return true
				}

				return false
			}
			if diff := cmp.Diff(tc.want, got, cmpopts.IgnoreMapEntries(ignoreUnpredictable)); diff != "" {
				t.Errorf("want != got:\n%v\n", diff)
			}
		})
	}
}

// TestGenerateBrokenPointer ensures that a present *_path key naming a file that
// cannot be read fails loudly rather than silently producing an empty value.
func TestGenerateBrokenPointer(t *testing.T) {
	rootff := &tests.FakeFiles{
		Data: map[string]types.J{
			"config.json": map[string]interface{}{
				// points at a lua file that was never written to disk
				"LuaScriptState_path": "missing/does-not-exist.luascriptstate",
			},
		},
	}
	m := Mod{
		RootRead:    rootff,
		RootWrite:   rootff,
		Lua:         &tests.FakeFiles{Fs: map[string]string{}},
		Modsettings: &tests.FakeFiles{},
		Objs:        &tests.FakeFiles{},
		Objdirs:     &tests.FakeFiles{},
	}
	err := m.GenerateFromConfig()
	if err == nil {
		t.Fatalf("expected error for broken *_path pointer, got nil")
	}
	if !strings.Contains(err.Error(), "missing/does-not-exist.luascriptstate") {
		t.Errorf("error should name the missing file, got: %v", err)
	}
}

// TestGenerateAbsentOptionalKey ensures that an optional key that is simply
// absent from config (no *_path and no inline value) stays lenient and does not
// error.
func TestGenerateAbsentOptionalKey(t *testing.T) {
	rootff := &tests.FakeFiles{
		Data: map[string]types.J{
			"config.json": map[string]interface{}{
				"SaveName": "a mod with no externalized optional fields",
			},
		},
	}
	m := Mod{
		RootRead:    rootff,
		RootWrite:   rootff,
		Lua:         &tests.FakeFiles{Fs: map[string]string{}},
		Modsettings: &tests.FakeFiles{},
		Objs:        &tests.FakeFiles{},
		Objdirs:     &tests.FakeFiles{},
	}
	if err := m.GenerateFromConfig(); err != nil {
		t.Fatalf("absent optional keys should not error, got: %v", err)
	}
	if got := m.Data["SaveName"]; got != "a mod with no externalized optional fields" {
		t.Errorf("SaveName not preserved, got %v", got)
	}
	// An absent optional string key is still filled with its zero value.
	if got, ok := m.Data["LuaScriptState"]; !ok || got != "" {
		t.Errorf("absent optional key LuaScriptState should be zero-valued, got %v (present=%v)", got, ok)
	}
}
