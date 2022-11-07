package reverse

import (
	"ModCreator/bundler"
	"ModCreator/file"
	"ModCreator/objects"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path"
)

// Reverser holds interfaces and configs for the reversing process
type Reverser struct {
	ModSettingsWriter file.JSONWriter
	LuaWriter         file.LuaWriter
	ObjWriter         file.JSONWriter
	ObjDirCreeator    file.DirCreator

	// expectations around types
	StringType   []string
	ObjType      []string
	ObjArrayType []string

	Root string
}

// Write executes the main purpose of the reverse library:
// to take a json object and create a file struture which mimics it.
func (r *Reverser) Write(raw map[string]interface{}) error {
	pathExt := "_path"
	for _, strKey := range r.StringType {
		rawVal, ok := raw[strKey]
		if !ok {
			log.Printf("expected string value in key %s, key not found\n", strKey)
			continue
		}
		strVal, ok := rawVal.(string)
		if !ok {
			return fmt.Errorf("expected string value in key %s, got %v", strKey, rawVal)
		}
		ext := ".luascriptstate"
		if strKey == "LuaScript" {
			ext = ".ttslua"

			unbundled, err := bundler.Unbundle(strVal)
			if err != nil {
				return fmt.Errorf("bundler.Unbundle(script from <%s>)\n: %v", r.Root, err)
			}
			strVal = unbundled
		}
		// decide if creating a separte file is worth it
		if len(strVal) < 80 {
			continue
		}

		createdFile := strKey + ext

		err := r.LuaWriter.EncodeToFile(strVal, createdFile)
		if err != nil {
			return fmt.Errorf("lua.EncodeToFile(<value>, %s) : %v", createdFile, err)
		}
		raw[strKey+pathExt] = createdFile
		delete(raw, strKey)
	}

	for _, objKey := range r.ObjType {
		rawVal, ok := raw[objKey]
		if ok {
			objVal, ok := rawVal.(map[string]interface{})
			if !ok {
				return fmt.Errorf("expected json object value in key %s, got %v", objKey, rawVal)
			}

			// decide if creating a separate file is worth it
			if len(fmt.Sprint(objVal)) < 100 {
				continue
			}

			createdFile := objKey + ".json"
			err := r.ModSettingsWriter.WriteObj(objVal, createdFile)
			if err != nil {
				return fmt.Errorf("j.WriteObj(<>, %s) : %v", createdFile, err)
			}
			raw[objKey+pathExt] = createdFile
			delete(raw, objKey)
		}
	}

	for _, objKey := range r.ObjArrayType {
		rawVal, ok := raw[objKey]
		if ok {
			arr, err := convertToObjArray(rawVal)
			if err != nil {
				return fmt.Errorf("mismatch expectations in key %s : %v", objKey, err)
			}

			// decide if creating a separate file is worth it
			if len(fmt.Sprint(arr)) < 200 {
				continue
			}

			createdFile := objKey + ".json"
			err = r.ModSettingsWriter.WriteObjArray(arr, createdFile)
			if err != nil {
				return fmt.Errorf("j.WriteObjArray(<>, %s) : %v", createdFile, err)
			}
			raw[objKey+pathExt] = createdFile
			delete(raw, objKey)
		}
	}

	if rawObjs, ok := raw["ObjectStates"]; ok {
		objStates, err := convertToObjArray(rawObjs)
		if err != nil {
			return fmt.Errorf("mismatch type expectations for ObjectStates : %v", err)
		}
		order, err := objects.PrintObjectStates("", r.LuaWriter, r.ObjWriter, r.ObjDirCreeator, objStates)
		if err != nil {
			return err
		}
		raw["ObjectStates_order"] = order
		delete(raw, "ObjectStates")
	}

	// write all that's Left
	err := writeJSON(raw, path.Join(r.Root, "config.json"))
	return err
}

func convertToObjArray(v interface{}) ([]map[string]interface{}, error) {
	arr := []map[string]interface{}{}

	rawArr, ok := v.([]interface{})
	if !ok {
		return nil, fmt.Errorf("%v is not an array", v)
	}

	for _, rv := range rawArr {
		objVal, ok := rv.(map[string]interface{})
		if !ok {
			if rv == nil {
				// if for some reason an array has nil object, just skip
				continue
			}
			return nil, fmt.Errorf("expected type json object, got %v", objVal)
		}
		arr = append(arr, objVal)
	}
	return arr, nil
}

func writeJSON(raw map[string]interface{}, filename string) error {
	b, err := json.MarshalIndent(raw, "", "  ")
	if err != nil {
		return fmt.Errorf("json.MarshalIndent() : %v", err)
	}
	return ioutil.WriteFile(filename, b, 0644)
}
