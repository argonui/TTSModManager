package objects

import (
	"ModCreator/bundler"
	"ModCreator/file"
	"log"
	"path"
	"regexp"
	"strings"

	"fmt"
)

type j map[string]interface{}
type objArray []map[string]interface{}

type objConfig struct {
	guid               string
	data               j
	luascriptPath      string
	luascriptstatePath string
	subObjDir          string
	subObj             []*objConfig
}

func (o *objConfig) parseFromFile(filepath string, j file.JSONReader) error {
	d, err := j.ReadObj(filepath)
	if err != nil {
		return fmt.Errorf("ReadObj(%s): %v", filepath, err)
	}
	o.data = d

	dguid, ok := o.data["GUID"]
	if !ok {
		return fmt.Errorf("object at (%s) doesn't have a GUID field", filepath)
	}
	guid, ok := dguid.(string)
	if !ok {
		return fmt.Errorf("object at (%s) doesn't have a string GUID (%s)", filepath, o.data["GUID"])
	}
	o.guid = guid

	// TODO nead ability to read from script folder
	tryParseIntoStr(&o.data, "LuaScript_path", &o.luascriptPath)
	tryParseIntoStr(&o.data, "LuaScriptState_path", &o.luascriptstatePath)
	tryParseIntoStr(&o.data, "ContainedObjects_path", &o.subObjDir)

	return nil
}

func tryParseIntoStr(m *j, k string, dest *string) {
	if raw, ok := (*m)[k]; ok {
		if str, ok := raw.(string); ok {
			*dest = str
			delete((*m), k)
		}
	}
}

func (o *objConfig) parseFromJSON(data map[string]interface{}) error {
	o.data = data
	dguid, ok := o.data["GUID"]
	if !ok {
		return fmt.Errorf("object (%v) doesn't have a GUID field", data)
	}
	guid, ok := dguid.(string)
	if !ok {
		return fmt.Errorf("object (%v) doesn't have a string GUID (%s)", dguid, o.data["GUID"])
	}
	o.guid = guid
	o.subObj = []*objConfig{}

	tryParseIntoStr(&o.data, "LuaScript_path", &o.luascriptPath)
	tryParseIntoStr(&o.data, "LuaScriptState_path", &o.luascriptstatePath)
	tryParseIntoStr(&o.data, "ContainedObjects_path", &o.subObjDir)

	if rawObjs, ok := o.data["ContainedObjects"]; ok {
		rawArr, ok := rawObjs.([]interface{})
		if !ok {
			return fmt.Errorf("type mismatch in ContainedObjects : %v", rawArr)
		}
		for _, rawSubO := range rawArr {
			subO, ok := rawSubO.(map[string]interface{})
			if !ok {
				return fmt.Errorf("type mismatch in ContainedObjects : %v", rawSubO)
			}
			so := objConfig{}
			if err := so.parseFromJSON(subO); err != nil {
				return fmt.Errorf("printing sub object of %s : %v", o.guid, err)
			}
			o.subObj = append(o.subObj, &so)
		}
		delete(o.data, "ContainedObjects")
	}
	return nil
}

func (o *objConfig) print(l file.LuaReader) (j, error) {
	if o.luascriptPath != "" {
		encoded, err := l.EncodeFromFile(o.luascriptPath)
		if err != nil {
			return j{}, fmt.Errorf("l.EncodeFromFile(%s) : %v", o.luascriptPath, err)
		}
		bundleReqs, err := bundler.Bundle(encoded, l)
		if err != nil {
			return nil, fmt.Errorf("Bundle(%s) : %v", encoded, err)
		}
		o.data["LuaScript"] = bundleReqs
	}
	if o.luascriptstatePath != "" {
		encoded, err := l.EncodeFromFile(o.luascriptstatePath)
		if err != nil {
			return j{}, fmt.Errorf("l.EncodeFromFile(%s) : %v", o.luascriptstatePath, err)
		}
		o.data["LuaScriptState"] = encoded
	}

	if s, ok := o.data["LuaScript"]; ok {
		if str, ok := s.(string); ok && str != "" {
			bundleReqs, err := bundler.Bundle(str, l)
			if err != nil {
				return nil, fmt.Errorf("Bundle(%s) : %v", str, err)
			}
			o.data["LuaScript"] = bundleReqs
		}
	}

	subs := []j{}
	for _, sub := range o.subObj {
		printed, err := sub.print(l)
		if err != nil {
			return nil, err
		}
		subs = append(subs, printed)
	}
	if len(subs) > 0 {
		o.data["ContainedObjects"] = subs
	}
	return o.data, nil
}

func (o *objConfig) printToFile(filepath string, l file.LuaWriter, j file.JSONWriter, dir file.DirCreator) error {
	// maybe convert LuaScript or LuaScriptState
	if rawscript, ok := o.data["LuaScript"]; ok {
		if script, ok := rawscript.(string); ok {
			script, err := bundler.Unbundle(script)
			if err != nil {
				return fmt.Errorf("bundler.Unbundle(script from <%s>)\n: %v", o.guid, err)
			}
			if len(script) > 80 {
				createdFile := path.Join(filepath, o.getAGoodFileName()+".ttslua")
				o.data["LuaScript_path"] = createdFile
				if err := l.EncodeToFile(script, createdFile); err != nil {
					return fmt.Errorf("EncodeToFile(<obj %s>)", o.guid)
				}
				delete(o.data, "LuaScript")
			} else {
				// put the unbundled bit back in
				o.data["LuaScript"] = script
			}
		}
	}
	if rawscript, ok := o.data["LuaScriptState"]; ok {
		if script, ok := rawscript.(string); ok {
			if len(script) > 80 {
				createdFile := path.Join(filepath, o.getAGoodFileName()+".luascriptstate")
				o.data["LuaScriptState_path"] = createdFile
				if err := l.EncodeToFile(script, createdFile); err != nil {
					return fmt.Errorf("EncodeToFile(<obj %s>)", o.guid)
				}
				delete(o.data, "LuaScriptState")
			}
		}
	}

	// recurse if need be
	if o.subObj != nil && len(o.subObj) > 0 {
		subdirName, err := dir.CreateDir(filepath, o.getAGoodFileName())
		if err != nil {
			return fmt.Errorf("<%v>.CreateDir(%s, %s) : %v", o.guid, filepath, o.getAGoodFileName(), err)
		}
		o.data["ContainedObjects_path"] = subdirName
		o.subObjDir = subdirName
		for _, subo := range o.subObj {
			err = subo.printToFile(path.Join(filepath, subdirName), l, j, dir)
			if err != nil {
				return fmt.Errorf("printing file %s: %v", path.Join(filepath, subdirName), err)
			}
		}
	}

	// print self
	fname := path.Join(filepath, o.getAGoodFileName()+".json")
	return j.WriteObj(o.data, fname)
}

func (o *objConfig) getAGoodFileName() string {
	moreUUID := o.guid
	if o.subObjDir != "" {
		moreUUID = o.subObjDir
	}
	// only let alphanumberic, _, -, be put into names
	reg := regexp.MustCompile("[^a-zA-Z0-9_-]+")

	keyname, err := o.tryGetNonEmptyStr("Nickname")
	if err != nil {
		keyname, err = o.tryGetNonEmptyStr("Name")
	}
	if err != nil {
		return moreUUID
	}

	n := reg.ReplaceAllString(keyname, "")
	return n + "." + moreUUID
}

func (o *objConfig) tryGetNonEmptyStr(key string) (string, error) {
	rawname, ok := o.data[key]
	if !ok {
		return "", fmt.Errorf("no key %s", key)
	}
	name, ok := rawname.(string)
	if !ok {
		return "", fmt.Errorf("key %s is not string", key)
	}
	if name == "" {
		return "", fmt.Errorf("key %s is blank", key)
	}
	return name, nil
}

type db struct {
	root []*objConfig

	all map[string]*objConfig

	j   file.JSONReader
	dir file.DirExplorer
}

func (d *db) addObj(o, parent *objConfig) error {
	if parent == nil {
		d.root = append(d.root, o)
	} else {
		parent.subObj = append(parent.subObj, o)
	}
	if _, ok := d.all[o.guid]; ok {
		log.Printf("Found duplicate guid %s\n", o.guid)
	} else {
		d.all[o.guid] = o
	}
	return nil
}

func (d *db) print(l file.LuaReader) (objArray, error) {
	var oa objArray
	for _, o := range d.root {
		printed, err := o.print(l)
		if err != nil {
			return objArray{}, fmt.Errorf("obj (%s) did not print : %v", o.guid, err)
		}
		oa = append(oa, printed)
	}
	return oa, nil
}

// ParseAllObjectStates looks at a folder and creates a json map from it.
// It assumes that folder names under the 'objects' directory are valid guids
// of existing Objects.
// like:
// objects/
// --foo.json (guid=1234)
// --bar.json (guid=888)
// --888/
//    --baz.json (guid=999) << this is a child of bar.json
func ParseAllObjectStates(l file.LuaReader, j file.JSONReader, dir file.DirExplorer) ([]map[string]interface{}, error) {
	d := db{
		j:   j,
		dir: dir,
		all: map[string]*objConfig{},
	}
	err := d.parseFromFolder("", nil)
	if err != nil {
		return []map[string]interface{}{}, fmt.Errorf("parseFolder(%s): %v", "<root>", err)
	}
	return d.print(l)
}

func (d *db) parseFromFolder(relpath string, parent *objConfig) error {
	filenames, folnames, err := d.dir.ListFilesAndFolders(relpath)
	if err != nil {
		return fmt.Errorf("ListFilesAndFolders(%s) : %v", relpath, err)
	}

	whoseFolder := map[string]*objConfig{}
	for _, file := range filenames {
		if !strings.HasSuffix(file, ".json") {
			// expect luascriptstate files and ttslua files to be stored alongside
			continue
		}
		o, err := d.parseFromFile(file, parent)
		if err != nil {
			return fmt.Errorf("parseFromFile(%s, %v): %v", file, parent, err)
		}
		if o.subObjDir != "" {
			whoseFolder[o.subObjDir] = o
		}
	}
	for _, folder := range folnames {
		o, ok := whoseFolder[path.Base(folder)]
		if !ok {
			return fmt.Errorf("found folder %s without a peer object who claims it", folder)
		}
		if err := d.parseFromFolder(folder, o); err != nil {
			return fmt.Errorf("parseFromFolder(%s): %v", folder, err)
		}
	}
	return nil
}

func (d *db) parseFromFile(relpath string, parent *objConfig) (*objConfig, error) {
	var o objConfig
	err := o.parseFromFile(relpath, d.j)
	if err != nil {
		return nil, fmt.Errorf("parseFromFile(%s) : %v", relpath, err)
	}

	return &o, d.addObj(&o, parent)
}

// PrintObjectStates takes a list of json objects and prints them in the
// expected format outlined by ParseAllObjectStates
func PrintObjectStates(root string, f file.LuaWriter, j file.JSONWriter, dir file.DirCreator, objs []map[string]interface{}) error {
	for _, rootObj := range objs {
		oc := objConfig{}
		err := oc.parseFromJSON(rootObj)
		if err != nil {
			return err
		}
		err = oc.printToFile(root, f, j, dir)
		if err != nil {
			return err
		}
	}
	return nil
}
