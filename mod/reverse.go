package mod

import (
	"ModCreator/file"
	"ModCreator/handler"
	"ModCreator/objects"
	"ModCreator/types"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// Reverser holds interfaces and configs for the reversing process
type Reverser struct {
	ModSettingsWriter file.JSONWriter
	LuaWriter         file.TextWriter
	LuaSrcWriter      file.TextWriter
	XMLWriter         file.TextWriter
	XMLSrcWriter      file.TextWriter
	ObjWriter         file.JSONWriter
	ObjDirCreeator    file.DirCreator
	RootWrite         file.JSONWriter

	// If not empty: holds the entire filename (C:...) of the json to read
	OnlyObjState string
}

func (r *Reverser) writeOnlyObjStates(raw map[string]interface{}) error {
	printer := &objects.Printer{
		Lua:    r.LuaWriter,
		LuaSrc: r.LuaSrcWriter,
		J:      r.ObjWriter,
		Dir:    r.ObjDirCreeator,
	}
	arraywrap := []map[string]interface{}{raw}
	_, err := printer.PrintObjectStates("", arraywrap)
	if err != nil {
		return fmt.Errorf("PrintObjectStates('', <%v objects>): %v", len(arraywrap), err)
	}
	return nil
}

// Write executes the main purpose of the reverse library:
// to take a json object and create a file struture which mimics it.
func (r *Reverser) Write(raw map[string]interface{}) error {
	if r.OnlyObjState != "" {
		return r.writeOnlyObjStates(raw)
	}
	pathExt := "_path"

	for _, strKey := range ExpectedStr {
		rawVal, ok := raw[strKey]
		if !ok {
			continue
		}
		strVal, ok := rawVal.(string)
		if !ok {
			return fmt.Errorf("expected string value in key %s, got %v", strKey, rawVal)
		}
		if strKey == "LuaScript" || strKey == "XmlUI" {
			// let the LuaHAndler handle the complicated case
			continue
		}
		ext := ".luascriptstate"
		// decide if creating a separate file is worth it
		if len(strVal) < 80 {
			raw[strKey] = strVal
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

	lh := handler.NewLuaHandler()
	lh.DefaultWriter = r.LuaWriter
	lh.SrcWriter = r.LuaSrcWriter

	act, err := lh.WhileWritingToFile(raw, "LuaScript.ttslua")
	if err != nil {
		return fmt.Errorf("WhileWritingToFile(<>, root luascript): %v", err)
	}
	if !act.Noop {
		delete(raw, "LuaScript")
		raw[act.Key] = act.Value
	}

	xh := handler.NewXMLHandler()
	xh.DefaultWriter = r.XMLWriter
	xh.SrcWriter = r.XMLSrcWriter

	act, err = xh.WhileWritingToFile(raw, "Root.xml")
	if err != nil {
		return fmt.Errorf("WhileWritingToFile(<>, root luascript): %v", err)
	}
	if !act.Noop {
		delete(raw, "XmlUI")
		raw[act.Key] = act.Value
	}

	for _, objKey := range ExpectedObj {
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

	for _, objKey := range ExpectedObjArr {
		rawVal, ok := raw[objKey]
		if ok {
			arr, err := types.ConvertToObjArray(rawVal)
			if err != nil {
				return fmt.Errorf("mismatch expectations in key %s : %v", objKey, err)
			}
			if objKey == "SnapPoints" {
				smoothed, err := objects.SmoothSnapPoints(arr)
				if err != nil {
					return fmt.Errorf("SmoothSnapPoints(): %v", err)
				}
				arr = smoothed
			}
			// decide if creating a separate file is worth it
			if len(fmt.Sprint(arr)) < 200 {
				raw[objKey] = arr
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
		objStates, err := types.ConvertToObjArray(rawObjs)
		if err != nil {
			return fmt.Errorf("mismatch type expectations for ObjectStates : %v", err)
		}
		printer := &objects.Printer{
			Lua:    r.LuaWriter,
			LuaSrc: r.LuaSrcWriter,
			XML:    r.XMLWriter,
			XMLSrc: r.XMLSrcWriter,
			J:      r.ObjWriter,
			Dir:    r.ObjDirCreeator,
		}
		order, err := printer.PrintObjectStates("", objStates)
		if err != nil {
			return fmt.Errorf("PrintObjectStates('', <%v objects>): %v", len(objStates), err)
		}
		raw["ObjectStates_order"] = order
		delete(raw, "ObjectStates")
	}

	// Time of the original mod being stored on file is meaningless,
	// they will be filled in on generate
	delete(raw, DateKey)
	delete(raw, EpochKey)

	// write all that's Left
	err = r.RootWrite.WriteObj(raw, "config.json")
	if err != nil {
		return fmt.Errorf("WriteObj(<obj>, %s) : %v", "config.json", err)
	}
	return nil
}

func writeJSON(raw map[string]interface{}, filename string) error {
	b, err := json.MarshalIndent(raw, "", "  ")
	if err != nil {
		return fmt.Errorf("json.MarshalIndent() : %v", err)
	}
	return ioutil.WriteFile(filename, b, 0644)
}
