package file

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
		return map[string]interface{}{}, fmt.Errorf("pullRawFile(%s): %v", filename, err)
	}
	var v map[string]interface{}
	err = json.Unmarshal(b, &v)
	if err != nil {
		return map[string]interface{}{}, fmt.Errorf("Unmarshal(%s): %v", filename, err)
	}
	return v, nil
}

// ReadObjArray pulls a file from configs and encodes it as a string.
func (j *JSONOps) ReadObjArray(filename string) ([]map[string]interface{}, error) {
	b, err := j.pullRawFile(filename)
	if err != nil {
		return nil, err
	}
	var v []map[string]interface{}
	err = json.Unmarshal(b, &v)
	if err != nil {
		return nil, fmt.Errorf("Unmarshal(): %v", err)
	}
	return v, nil

}
func (j *JSONOps) pullRawFile(filename string) ([]byte, error) {
	p := path.Join(j.basepath, filename)
	jFile, err := os.Open(p)
	if err != nil {
		return nil, fmt.Errorf("os.Open(%s): %v", p, err)
	}
	defer jFile.Close()

	return ioutil.ReadAll(jFile)
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
	return ioutil.WriteFile(p, b, 0644)
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
	return ioutil.WriteFile(p, b, 0644)
}

// ReadRawFile allows for anyone who needs to to read json without objects.
func ReadRawFile(filename string) (map[string]interface{}, error) {
	jFile, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("os.Open(%s): %v", filename, err)
	}
	defer jFile.Close()

	b, err := ioutil.ReadAll(jFile)
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
