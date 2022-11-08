package objects

import (
	"ModCreator/bundler"
	"ModCreator/file"
	. "ModCreator/types"
	"path"
	"regexp"
	"strings"

	"fmt"
)

type objConfig struct {
	guid               string
	data               J
	luascriptPath      string
	luascriptstatePath string
	gmnotesPath	   string
	subObjDir          string
	subObjOrder        []string // array of base filenames of subobjects
	subObj             []*objConfig
}

func (o *objConfig) parseFromFile(filepath string, j file.JSONReader) error {
	d, err := j.ReadObj(filepath)
	if err != nil {
		return fmt.Errorf("ReadObj(%s): %v", filepath, err)
	}
	err = o.parseFromJSON(d)
	if err != nil {
		return fmt.Errorf("<%s>.parseFromJSON(): %v", filepath, err)
	}
	if o.subObjDir != "" {
		for _, oname := range o.subObjOrder {
			subo := &objConfig{}
			relFilename := path.Join(path.Dir(filepath), o.subObjDir, fmt.Sprintf("%s.json", oname))

			err = subo.parseFromFile(relFilename, j)
			if err != nil {
				return fmt.Errorf("parseFromFile(%s): %v", relFilename, err)
			}
			o.subObj = append(o.subObj, subo)
		}
	}
	return nil
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
	o.subObjOrder = []string{}

	file.TryParseIntoStr(&o.data, "LuaScript_path", &o.luascriptPath)
	file.TryParseIntoStr(&o.data, "LuaScriptState_path", &o.luascriptstatePath)
        file.TryParseIntoStr(&o.data, "GMNotes_path", &o.gmnotesPath)
	file.TryParseIntoStr(&o.data, "ContainedObjects_path", &o.subObjDir)
	file.TryParseIntoStrArray(&o.data, "ContainedObjects_order", &o.subObjOrder)

	if trans, ok := o.data["Transform"]; ok {
		o.data["Transform"] = Smooth(trans)
	}

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
				return fmt.Errorf("parsing sub object of %s : %v", o.guid, err)
			}
			o.subObj = append(o.subObj, &so)
			o.subObjOrder = append(o.subObjOrder, so.getAGoodFileName())
		}
		delete(o.data, "ContainedObjects")
	}

	return nil
}

func (o *objConfig) print(l file.LuaReader) (J, error) {
	if o.luascriptPath != "" {
		encoded, err := l.EncodeFromFile(o.luascriptPath)
		if err != nil {
			return J{}, fmt.Errorf("l.EncodeFromFile(%s) : %v", o.luascriptPath, err)
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
			return J{}, fmt.Errorf("l.EncodeFromFile(%s) : %v", o.luascriptstatePath, err)
		}
		o.data["LuaScriptState"] = encoded
	}
	if o.gmnotesPath != "" {
		encoded, err := l.EncodeFromFile(o.gmnotesPath)
		if err != nil {
			return J{}, fmt.Errorf("l.EncodeFromFile(%s) : %v", o.gmnotesPath, err)
		}
		o.data["GMNotes"] = encoded
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

	subs := []J{}
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
	if rawscript, ok := o.data["GMNotes"]; ok {
		if script, ok := rawscript.(string); ok {
			if len(script) > 80 {
				createdFile := path.Join(filepath, o.getAGoodFileName()+".gmnotes")
				o.data["GMNotes_path"] = createdFile
				if err := l.EncodeToFile(script, createdFile); err != nil {
					return fmt.Errorf("EncodeToFile(<obj %s>)", o.guid)
				}
				delete(o.data, "GMNotes")
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
		if len(o.subObj) != len(o.subObjOrder) {
			return fmt.Errorf("subobj order not getting filled in on %s", o.getAGoodFileName())
		}
		o.data["ContainedObjects_order"] = o.subObjOrder
	}

	// print self
	fname := path.Join(filepath, o.getAGoodFileName()+".json")
	return j.WriteObj(o.data, fname)
}

func (o *objConfig) getAGoodFileName() string {
	// only let alphanumberic, _, -, be put into names
	reg := regexp.MustCompile("[^a-zA-Z0-9_-]+")

	keyname, err := o.tryGetNonEmptyStr("Nickname")
	if err != nil {
		keyname, err = o.tryGetNonEmptyStr("Name")
	}
	if err != nil {
		return o.guid
	}

	n := reg.ReplaceAllString(keyname, "")
	return n + "." + o.guid
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
	root map[string]*objConfig

	all map[string]*objConfig

	j   file.JSONReader
	dir file.DirExplorer
}

func (d *db) print(l file.LuaReader, order []string) (ObjArray, error) {
	var oa ObjArray
	if len(order) != len(d.root) {
		return nil, fmt.Errorf("expected order and db.root to have same length, %v != %v", len(order), len(d.root))
	}
	for _, nextGUID := range order {
		if _, ok := d.root[nextGUID]; !ok {
			return nil, fmt.Errorf("order expected %s, not found in db", nextGUID)
		}
		printed, err := d.root[nextGUID].print(l)
		if err != nil {
			return ObjArray{}, fmt.Errorf("obj (%s) did not print : %v", nextGUID, err)
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
func ParseAllObjectStates(l file.LuaReader, j file.JSONReader, dir file.DirExplorer, order []string) ([]map[string]interface{}, error) {
	d := db{
		j:    j,
		dir:  dir,
		all:  map[string]*objConfig{},
		root: map[string]*objConfig{},
	}
	err := d.parseFromFolder("", nil)
	if err != nil {
		return []map[string]interface{}{}, fmt.Errorf("parseFolder(%s): %v", "<root>", err)
	}
	return d.print(l, order)
}

func (d *db) parseFromFolder(relpath string, parent *objConfig) error {
	filenames, _, err := d.dir.ListFilesAndFolders(relpath)
	if err != nil {
		return fmt.Errorf("ListFilesAndFolders(%s) : %v", relpath, err)
	}

	for _, file := range filenames {
		if !strings.HasSuffix(file, ".json") {
			// expect luascriptstate, gmnotes, and ttslua files to be stored alongside
			continue
		}
		var o objConfig
		err := o.parseFromFile(file, d.j)
		if err != nil {
			return fmt.Errorf("parseFromFile(%s, %v): %v", file, parent, err)
		}
		d.root[o.getAGoodFileName()] = &o
	}

	return nil
}

// PrintObjectStates takes a list of json objects and prints them in the
// expected format outlined by ParseAllObjectStates
func PrintObjectStates(root string, f file.LuaWriter, j file.JSONWriter, dir file.DirCreator, objs []map[string]interface{}) ([]string, error) {
	order := []string{}
	for _, rootObj := range objs {
		oc := objConfig{}

		err := oc.parseFromJSON(rootObj)
		if err != nil {
			return nil, err
		}
		order = append(order, oc.getAGoodFileName())
		err = oc.printToFile(root, f, j, dir)
		if err != nil {
			return nil, err
		}
	}
	return order, nil
}
