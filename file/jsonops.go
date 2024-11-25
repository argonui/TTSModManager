package file

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
)

// JSONOps implements the corresponding reader & writer interfaces
type JSONOps struct {
	basepath string
}

// JSONReader allows for arbitrary reads and encoding of json
type JSONReader interface {
	ReadObj(string) (map[string]interface{}, error)
	ReadObjArray(string) ([]map[string]interface{}, error)
}

// JSONWriter allows for arbitrary writes and encoding of json
type JSONWriter interface {
	WriteObj(map[string]interface{}, string) error
	WriteObjArray([]map[string]interface{}, string) error
	WriteSavedObj(map[string]interface{}, string) error
}

type SavedObject struct {
	SaveName       string
	Date           string
	VersionNumber  string
	GameMode       string
	GameType       string
	GameComplexity string
	Tags           []string
	Gravity        float64
	PlayArea       float64
	Table          string
	Sky            string
	Note           string
	TabStates      map[string]interface{}
	LuaScript      string
	LuaScriptState string
	XmlUI          string
	ObjectStates   []map[string]interface{} // This will hold the original `m`
}

// NewJSONOps initializes our object on a directory
func NewJSONOps(base string) *JSONOps {
	return &JSONOps{
		basepath: base,
	}
}

// ReadObj pulls a file from configs and encodes it as a string.
func (j *JSONOps) ReadObj(filename string) (map[string]interface{}, error) {
	b, err := j.pullRawFile(filename)
	if err != nil {
		return map[string]interface{}{}, err
	}
	var v map[string]interface{}
	json.Unmarshal(b, &v)
	return v, nil
}

// ReadObjArray pulls a file from configs and encodes it as a string.
func (j *JSONOps) ReadObjArray(filename string) ([]map[string]interface{}, error) {
	b, err := j.pullRawFile(filename)
	if err != nil {
		return []map[string]interface{}{}, err
	}
	var v []map[string]interface{}
	json.Unmarshal(b, &v)
	return v, nil

}
func (j *JSONOps) pullRawFile(filename string) ([]byte, error) {
	p := path.Join(j.basepath, filename)
	jFile, err := os.Open(p)
	if err != nil {
		return nil, fmt.Errorf("os.Open(%s): %v", p, err)
	}
	defer jFile.Close()

	return io.ReadAll(jFile)
}

// WriteObj writes a serialized json object to a file.
func (j *JSONOps) WriteObj(m map[string]interface{}, filename string) error {
	b, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	b = append(b, '\n')
	p := path.Join(j.basepath, filename)
	err = os.MkdirAll(path.Dir(p), 0750)
	if err != nil && !os.IsExist(err) {
		return fmt.Errorf("MkdirAll(%s): %v", path.Dir(p), err)
	}
	return os.WriteFile(p, b, 0644)
}

// WriteSavedObj writes a serialized JSON object to a file with the boilerplate for TTS saved objects.
// For additional information, see saved-object-feature.md in the project repository.
func (j *JSONOps) WriteSavedObj(m map[string]interface{}, filename string) error {
	var b []byte
	var err error

	// Create a SavedObject struct and embed `m` in the ObjectStates field
	savedObject := SavedObject{
		SaveName:       "",
		Date:           "",
		VersionNumber:  "",
		GameMode:       "",
		GameType:       "",
		GameComplexity: "",
		Tags:           []string{},
		Gravity:        0.5,
		PlayArea:       0.5,
		Table:          "",
		Sky:            "",
		Note:           "",
		TabStates:      map[string]interface{}{},
		LuaScript:      "",
		LuaScriptState: "",
		XmlUI:          "",
		ObjectStates:   []map[string]interface{}{m}, // Embed the object here
	}

	b, err = json.MarshalIndent(savedObject, "", "  ")
	if err != nil {
		return err
	}

	// Add an end-of-file newline
	b = append(b, '\n')

	// Ensure the path exists
	p := path.Join(j.basepath, filename)
	err = os.MkdirAll(path.Dir(p), 0750)
	if err != nil && !os.IsExist(err) {
		return fmt.Errorf("MkdirAll(%s): %v", path.Dir(p), err)
	}

	// Write the JSON data to the file
	return os.WriteFile(p, b, 0644)
}

// WriteObjArray writes an array of serialized json objects to a file.
func (j *JSONOps) WriteObjArray(m []map[string]interface{}, filename string) error {
	b, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	b = append(b, '\n')
	p := path.Join(j.basepath, filename)
	err = os.MkdirAll(path.Dir(p), 0750)
	if err != nil && !os.IsExist(err) {
		return fmt.Errorf("MkdirAll(%s): %v", path.Dir(p), err)
	}
	return os.WriteFile(p, b, 0644)
}

// ReadRawFile allows for anyone who needs to to read json without objects.
func ReadRawFile(filename string) (map[string]interface{}, error) {
	jFile, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("os.Open(%s): %v", filename, err)
	}
	defer jFile.Close()

	b, err := io.ReadAll(jFile)
	if err != nil {
		return nil, err
	}

	var v map[string]interface{}
	err = json.Unmarshal(b, &v)
	if err != nil {
		return nil, fmt.Errorf("Unmarshal(<%s>) : %v", filename, err)
	}
	return v, nil
}
