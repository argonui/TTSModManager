package tests

import (
	"ModCreator/types"
	"encoding/json"
	"fmt"
	"path"
	"strings"
)

// FakeFiles allows for mocking of the file system to use in tests
type FakeFiles struct {
	Fs   map[string]string
	Data map[string]types.J
}

// NewFF helps you not forget to initialize maps
func NewFF() *FakeFiles {
	return &FakeFiles{
		Fs:   map[string]string{},
		Data: map[string]types.J{},
	}
}

// EncodeFromFile satisfies LuaReader
func (f *FakeFiles) EncodeFromFile(s string) (string, error) {
	if _, ok := f.Fs[s]; !ok {
		return "", fmt.Errorf("fake file <%s> not found", s)
	}
	return f.Fs[s], nil
}

// ReadObj satisfies JSONReader
func (f *FakeFiles) ReadObj(s string) (map[string]interface{}, error) {
	if _, ok := f.Data[s]; !ok {
		return nil, fmt.Errorf("fake file <%s> not found", s)
	}
	b, err := json.MarshalIndent(f.Data[s], "", "  ")
	if err != nil {
		return nil, err
	}
	var v map[string]interface{}
	err = json.Unmarshal(b, &v)
	if err != nil {
		return nil, err
	}
	return v, nil
}

// ReadObjArray satisfies JSONReader
func (f *FakeFiles) ReadObjArray(s string) ([]map[string]interface{}, error) {
	wrapper, ok := f.Data[s]
	if !ok {
		return nil, fmt.Errorf("fake file <%s> not found", s)
	}
	raw, ok := wrapper["testarray"]
	if !ok {
		return nil, fmt.Errorf("fake obj(%s) was not faked to have array in it", s)
	}
	val, ok := raw.([]map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("fake obj array in %s was not array, was %T", s, raw)
	}
	return val, nil
}

// WriteObj satisfies JSONWriter
func (f *FakeFiles) WriteObj(data map[string]interface{}, path string) error {
	f.Data[path] = data
	return nil
}

// WriteObjArray satisfies JSONWriter
func (f *FakeFiles) WriteObjArray(data []map[string]interface{}, path string) error {
	f.Data[path] = types.J{
		"testarray": data,
	}
	return nil
}

// WriteSavedObj satisfies JSONWriter and mimics the behavior of the real WriteSavedObj.
func (f *FakeFiles) WriteSavedObj(data map[string]interface{}, path string) error {
	savedObject := map[string]interface{}{
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
		"ObjectStates":   []map[string]interface{}{data},
	}

	f.Data[path] = savedObject
	return nil
}

// EncodeToFile satisfies LuaWriter
func (f *FakeFiles) EncodeToFile(script, file string) error {
	f.Fs[file] = script
	return nil
}

// CreateDir satisfies DirCreator
func (f *FakeFiles) CreateDir(a, b string) (string, error) {
	// return the "chosen" directory name for next folder
	return b, nil
}

// Clear satisfies DirCreator
func (f *FakeFiles) Clear() error {
	return nil
}

// ListFilesAndFolders satisfies DirExplorer
func (f *FakeFiles) ListFilesAndFolders(relpath string) ([]string, []string, error) {
	// ignore non json files. i don't think they Matter
	files := []string{}
	folders := []string{}
	for k := range f.Data {
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

// DebugFileNames lets you log the file structure
func (f *FakeFiles) DebugFileNames(log func(s string, args ...interface{})) {
	log("txt files:\n")
	for fn := range f.Fs {
		log("\t%s\n", fn)
	}
	log("json files:\n")
	for fn := range f.Data {
		log("\t%s\n", fn)
	}
}
