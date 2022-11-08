package tests

import (
	"ModCreator/file"
	"ModCreator/mod"
	"ModCreator/types"
	"fmt"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

var (
	expectedStr       = []string{"SaveName", "Date", "VersionNumber", "GameMode", "GameType", "GameComplexity", "Table", "Sky", "Note", "LuaScript", "LuaScriptState", "XmlUI"}
	expectedObj       = []string{"TabStates", "MusicPlayer", "Grid", "Lighting", "Hands", "ComponentTags", "Turns"}
	expectedObjArr    = []string{"CameraStates", "DecalPallet", "CustomUIAssets", "SnapPoints", "Decals"}
	expectedObjStates = "ObjectStates"
)

type fakeFiles struct {
	fs   map[string]string
	data map[string]types.J
}

func newFF() *fakeFiles {
	return &fakeFiles{
		fs:   map[string]string{},
		data: map[string]types.J{},
	}
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
func (f *fakeFiles) ListFilesAndFolders(relpath string) ([]string, []string, error) {
	// ignore non json files. i don't think they Matter
	files := []string{}
	folders := []string{}
	for k := range f.data {
		if strings.HasPrefix(k, relpath) {
			left := k
			if relpath != "" {
				left = strings.Replace(k, relpath+"/", "", 1)
			}
			if strings.Contains(left, "/") {
				// this is a folder not a file
				folders = append(folders, path.Join(relpath, strings.Split(left, "/")[0]))
			} else {
				files = append(files, path.Join(relpath, left))
			}
		}
	}
	return files, folders, nil
}

func (f *fakeFiles) debugFileNames(log func(s string, args ...interface{})) {
	log("txt files:\n")
	for fn := range f.fs {
		log("\t%s\n", fn)
	}
	log("json files:\n")
	for fn := range f.data {
		log("\t%s\n", fn)
	}
}

func TestAllReverseThenBuild(t *testing.T) {
	paths, err := filepath.Glob(filepath.Join("testdata", "e2e", "*.json"))
	if err != nil {
		t.Fatal(err)
	}
	for _, path := range paths {
		_, filename := filepath.Split(path)
		testname := filename[:len(filename)-len(filepath.Ext(path))]
		t.Run(testname, func(t *testing.T) {
			j, err := file.ReadRawFile(path)
			if err != nil {
				t.Fatalf("Error parsing %s : %v", path, err)
			}
			modsettings := newFF()
			finalOutput := newFF()
			objsAndLua := newFF()

			r := mod.Reverser{
				ModSettingsWriter: modsettings,
				LuaWriter:         objsAndLua,
				ObjWriter:         objsAndLua,
				ObjDirCreeator:    objsAndLua,
				RootWrite:         finalOutput,
			}
			err = r.Write(j)
			if err != nil {
				t.Fatalf("Error reversing : %v", err)
			}

			objsAndLua.debugFileNames(t.Logf)
			finalOutput.debugFileNames(t.Logf)
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
			err = m.Print()
			if err != nil {
				t.Fatalf("printMod(...) : %v", err)
			}
			got, err := finalOutput.ReadObj("output.json")
			if err != nil {
				t.Fatalf("output.json not parsed : %v", err)
			}
			want := j

			if diff := cmp.Diff(want, got); diff != "" {
				t.Errorf("want != got:\n%v\n", diff)
			}
		})

	}
}
