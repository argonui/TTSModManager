package objects

import (
	"ModCreator/file"
	"ModCreator/handler"
	. "ModCreator/types"
	"encoding/json"
	"fmt"
	"path"
	"regexp"
	"strings"
)

type objConfig struct {
	guid               string
	data               J
	luascriptstatePath string
	gmnotesPath        string
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

	file.TryParseIntoStr(&o.data, "LuaScriptState_path", &o.luascriptstatePath)
	file.TryParseIntoStr(&o.data, "GMNotes_path", &o.gmnotesPath)
	file.TryParseIntoStr(&o.data, "ContainedObjects_path", &o.subObjDir)
	file.TryParseIntoStrArray(&o.data, "ContainedObjects_order", &o.subObjOrder)

	for _, needSmoothing := range []string{"Transform", "ColorDiffuse"} {
		if v, ok := o.data[needSmoothing]; ok {
			o.data[needSmoothing] = Smooth(v)
		}
	}
	if v, ok := o.data["AltLookAngle"]; ok {
		vv, err := SmoothAngle(v)
		if err != nil {
			return fmt.Errorf("SmoothAngle(<%s>): %v", "AltLookAngle", err)
		}
		o.data["AltLookAngle"] = vv
	}
	if sp, ok := o.data["AttachedSnapPoints"]; ok {
		sm, err := SmoothSnapPoints(sp)
		if err != nil {
			return fmt.Errorf("SmoothSnapPoints(<%s>): %v", o.guid, err)
		}
		o.data["AttachedSnapPoints"] = sm
	}

	// apply smoothing to object states
	if states, ok := o.data["States"]; ok {
		statesMap, ok := states.(map[string]interface{})
		if !ok {
			return fmt.Errorf("type mismatch in States : %v", states)
		}
		for stateName, stateData := range statesMap {
			stateObj, ok := stateData.(map[string]interface{})
			if !ok {
				return fmt.Errorf("type mismatch in State %s : %v", stateName, stateData)
			}

			for _, needSmoothing := range []string{"Transform", "ColorDiffuse"} {
				if v, ok := stateObj[needSmoothing]; ok {
					stateObj[needSmoothing] = Smooth(v)
				}
			}

			if v, ok := stateObj["AltLookAngle"]; ok {
				vv, err := SmoothAngle(v)
				if err != nil {
					return fmt.Errorf("SmoothAngle(<%s>) in State %s: %v", "AltLookAngle", stateName, err)
				}
				stateObj["AltLookAngle"] = vv
			}

			if sp, ok := stateObj["AttachedSnapPoints"]; ok {
				sm, err := SmoothSnapPoints(sp)
				if err != nil {
					return fmt.Errorf("SmoothSnapPoints(<%s>) in State %s: %v", o.guid, stateName, err)
				}
				stateObj["AttachedSnapPoints"] = sm
			}

			statesMap[stateName] = stateObj
		}
		o.data["States"] = statesMap
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

func (o *objConfig) print(l, x file.TextReader) (J, error) {
	out := o.data

	lh := handler.NewLuaHandler()
	lh.Reader = l
	act, err := lh.WhileReadingFromFile(o.data)
	if err != nil {
		return nil, fmt.Errorf("WhileReadingFromFile(): %v", err)
	}
	if !act.Noop {
		delete(out, "LuaScript")
		delete(out, "LuaScript_path")
		out[act.Key] = act.Value
	}
	xh := handler.NewXMLHandler()
	xh.Reader = x
	act, err = xh.WhileReadingFromFile(o.data)
	if err != nil {
		return nil, fmt.Errorf("WhileReadingFromFile(): %v", err)
	}
	if !act.Noop {
		delete(out, "XmlUI")
		delete(out, "XmlUI_path")
		out[act.Key] = act.Value
	}

	if o.gmnotesPath != "" {
		encoded, err := l.EncodeFromFile(o.gmnotesPath)
		if err != nil {
			return J{}, fmt.Errorf("l.EncodeFromFile(%s) : %v", o.gmnotesPath, err)
		}
		out["GMNotes"] = encoded
	}
	if o.luascriptstatePath != "" {
		encoded, err := l.EncodeFromFile(o.luascriptstatePath)
		if err != nil {
			return J{}, fmt.Errorf("l.EncodeFromFile(%s) : %v", o.luascriptstatePath, err)
		}
		out["LuaScriptState"] = encoded
	}

	subs := []J{}
	for _, sub := range o.subObj {
		printed, err := sub.print(l, x)
		if err != nil {
			return nil, err
		}
		subs = append(subs, printed)
	}
	if len(subs) > 0 {
		out["ContainedObjects"] = subs
	}
	return out, nil
}

func (o *objConfig) printToFile(filepath string, p *Printer) error {
	out := o.data

	lh := handler.NewLuaHandler()
	lh.DefaultWriter = p.Lua
	lh.SrcWriter = p.LuaSrc

	maybeNeededFname := path.Join(filepath, o.getAGoodFileName()+".ttslua")
	act, err := lh.WhileWritingToFile(o.data, maybeNeededFname)
	if err != nil {
		return fmt.Errorf("WhileWritingToFile(<>, %s): %v", maybeNeededFname, err)
	}
	if !act.Noop {
		delete(out, "LuaScript")
		delete(out, "LuaScript_path")
		out[act.Key] = act.Value
	}

	xh := handler.NewXMLHandler()
	xh.DefaultWriter = p.XML
	xh.SrcWriter = p.XMLSrc
	maybeNeededFname = path.Join(filepath, o.getAGoodFileName()+".xml")
	act, err = xh.WhileWritingToFile(o.data, maybeNeededFname)
	if err != nil {
		return fmt.Errorf("WhileWritingToFile(<>, %s): %v", maybeNeededFname, err)
	}
	if !act.Noop {
		delete(out, "XmlUI")
		delete(out, "XmlUI_path")
		out[act.Key] = act.Value
	}

	if rawscript, ok := o.data["LuaScriptState"]; ok {
		if script, ok := rawscript.(string); ok {
			if len(script) > 80 {
				createdFile := path.Join(filepath, o.getAGoodFileName()+".luascriptstate")
				out["LuaScriptState_path"] = createdFile

				var jsonInterface map[string]interface{}
				err := json.Unmarshal([]byte(script), &jsonInterface)
				if err == nil {
					// if it's JSON, use the JSONWriter
					if err := p.J.WriteObj(jsonInterface, createdFile); err != nil {
						return fmt.Errorf("j.WriteObj(<obj %s>)", o.guid)
					}
				} else {
					// default to the regular text writing
					if err := p.Lua.EncodeToFile(script, createdFile); err != nil {
						return fmt.Errorf("EncodeToFile(<obj %s>)", o.guid)
					}
				}

				delete(out, "LuaScriptState")
			}
		}
	}
	if rawscript, ok := o.data["GMNotes"]; ok {
		if script, ok := rawscript.(string); ok {
			if len(script) > 80 {
				createdFile := path.Join(filepath, o.getAGoodFileName()+".gmnotes")
				o.data["GMNotes_path"] = createdFile
				if err := p.Lua.EncodeToFile(script, createdFile); err != nil {
					return fmt.Errorf("EncodeToFile(<obj %s>)", o.guid)
				}
				delete(o.data, "GMNotes")
			}
		}
	}

	// recurse if need be
	if o.subObj != nil && len(o.subObj) > 0 {
		subdirName, err := p.Dir.CreateDir(filepath, o.getAGoodFileName())
		if err != nil {
			return fmt.Errorf("<%v>.CreateDir(%s, %s) : %v", o.guid, filepath, o.getAGoodFileName(), err)
		}
		out["ContainedObjects_path"] = subdirName
		o.subObjDir = subdirName
		for _, subo := range o.subObj {
			err = subo.printToFile(path.Join(filepath, subdirName), p)
			if err != nil {
				return fmt.Errorf("printing file %s: %v", path.Join(filepath, subdirName), err)
			}
		}
		if len(o.subObj) != len(o.subObjOrder) {
			return fmt.Errorf("subobj order not getting filled in on %s", o.getAGoodFileName())
		}
		out["ContainedObjects_order"] = o.subObjOrder
	}

	// print self
	fname := path.Join(filepath, o.getAGoodFileName()+".json")
	return p.J.WriteObj(out, fname)
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

func (d *db) print(l file.TextReader, x file.TextReader, order []string) (ObjArray, error) {
	var oa ObjArray
	if len(order) != len(d.root) {
		return nil, fmt.Errorf("expected order (%v) and db.root (%v) to have same length", len(order), len(d.root))
	}
	for _, nextGUID := range order {
		if _, ok := d.root[nextGUID]; !ok {
			return nil, fmt.Errorf("order expected %s, not found in db <%v>", nextGUID, d.root)
		}
		printed, err := d.root[nextGUID].print(l, x)
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
func ParseAllObjectStates(l file.TextReader, x file.TextReader, j file.JSONReader, dir file.DirExplorer, order []string) ([]map[string]interface{}, error) {
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
	return d.print(l, x, order)
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

// Printer holds the info of which writer is which, because parameter lists are rough
type Printer struct {
	Lua    file.TextWriter
	LuaSrc file.TextWriter
	XML    file.TextWriter
	XMLSrc file.TextWriter
	J      file.JSONWriter
	Dir    file.DirCreator
}

// PrintObjectStates takes a list of json objects and prints them in the
// expected format outlined by ParseAllObjectStates
func (p *Printer) PrintObjectStates(root string, objs []map[string]interface{}) ([]string, error) {
	order := []string{}

	for _, rootObj := range objs {
		oc := objConfig{}

		err := oc.parseFromJSON(rootObj)
		if err != nil {
			return nil, err
		}
		order = append(order, oc.getAGoodFileName())
		err = oc.printToFile(root, p)
		if err != nil {
			return nil, err
		}
	}
	return order, nil
}
