package main

import (
	file "ModCreator/file"
	objects "ModCreator/objects"
	"ModCreator/reverse"
	"ModCreator/types"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
)

var (
	moddir  = flag.String("moddir", "testdata/simple", "a directory containing tts mod configs")
	rev     = flag.Bool("reverse", false, "Instead of building a json from file structure, build file structure from json.")
	modfile = flag.String("modfile", "", "where to read from when reversing.")

	expectedStr       = []string{"SaveName", "Date", "VersionNumber", "GameMode", "GameType", "GameComplexity", "Table", "Sky", "Note", "LuaScript", "LuaScriptState", "XmlUI"}
	expectedObj       = []string{"TabStates", "MusicPlayer", "Grid", "Lighting", "Hands", "ComponentTags", "Turns"}
	expectedObjArr    = []string{"CameraStates", "DecalPallet", "CustomUIAssets", "SnapPoints", "Decals"}
	expectedObjStates = "ObjectStates"
)

const (
	textSubdir     = "src"
	modsettingsDir = "modsettings"
	objectsSubdir  = "objects"
)

// Config is how users will specify their mod's configuration.
type Config struct {
	Raw types.J `json:"-"`
}

// Mod is used as the accurate representation of what gets printed when
// module creation is done
type Mod struct {
	Data types.J

	lua         file.LuaReader
	modsettings file.JSONReader
	objs        file.JSONReader
	objdirs     file.DirExplorer
}

func main() {
	flag.Parse()

	lua := file.NewLuaOpsMulti(
		[]string{path.Join(*moddir, textSubdir), path.Join(*moddir, objectsSubdir)},
		path.Join(*moddir, objectsSubdir),
	)
	ms := file.NewJSONOps(path.Join(*moddir, modsettingsDir))
	objs := file.NewJSONOps(path.Join(*moddir, objectsSubdir))
	objdir := file.NewDirOps(path.Join(*moddir, objectsSubdir))

	if *rev {
		raw, err := prepForReverse(*moddir, *modfile)
		if err != nil {
			log.Fatalf("prepForReverse (%s) failed : %v", *modfile, err)
		}
		r := reverse.Reverser{
			ModSettingsWriter: ms,
			LuaWriter:         lua,
			ObjWriter:         objs,
			ObjDirCreeator:    objdir,
			StringType:        expectedStr,
			ObjType:           expectedObj,
			ObjArrayType:      expectedObjArr,
			Root:              *moddir,
		}
		err = r.Write(raw)
		if err != nil {
			log.Fatalf("reverse.Write(<%s>) failed : %v", *modfile, err)
		}
		return
	}

	c, err := readConfig(*moddir)
	if err != nil {
		fmt.Printf("readConfig(%s) : %v\n", *moddir, err)
		return
	}

	m := &Mod{
		lua:         lua,
		modsettings: ms,
		objs:        objs,
		objdirs:     objdir,
	}
	err = m.generate(c)
	if err != nil {
		fmt.Printf("generateMod(<config>) : %v\n", err)
		return
	}
	err = printMod(*moddir, m)
	if err != nil {
		log.Fatalf("printMod(...) : %v", err)
	}
}

func readConfig(cPath string) (*Config, error) {
	// Open our jsonFile
	cFile, err := os.Open(path.Join(cPath, "config.json"))
	// if we os.Open returns an error then handle it
	if err != nil {
		return nil, fmt.Errorf("os.Open(%s): %v", path.Join(cPath, "config.json"), err)
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer cFile.Close()

	b, err := ioutil.ReadAll(cFile)
	if err != nil {
		return nil, fmt.Errorf("ioutil.Readall(%s) : %v", path.Join(cPath, "config.json"), err)
	}
	var c Config

	err = json.Unmarshal(b, &c.Raw)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal(%s) : %v", b, err)
	}
	return &c, nil
}

func (m *Mod) generate(c *Config) error {
	if c == nil {
		return fmt.Errorf("nil config")
	}

	m.Data = c.Raw

	plainObj := func(s string) (interface{}, error) {
		return m.modsettings.ReadObj(s)
	}
	objArray := func(s string) (interface{}, error) {
		return m.modsettings.ReadObjArray(s)
	}
	luaGet := func(s string) (interface{}, error) {
		return m.lua.EncodeFromFile(s)
	}

	ext := "_path"
	for _, stringbased := range expectedStr {
		tryPut(&m.Data, stringbased+ext, stringbased, luaGet)
	}

	for _, objbased := range expectedObj {
		tryPut(&m.Data, objbased+ext, objbased, plainObj)
	}

	for _, objarraybased := range expectedObjArr {
		tryPut(&m.Data, objarraybased+ext, objarraybased, objArray)
	}

	objOrder := []string{}
	err := file.ForceParseIntoStrArray(&m.Data, "ObjectStates_order", &objOrder)
	if err != nil {
		return fmt.Errorf("ForceConvertToStrArray gave %v", err)
	}

	allObjs, err := objects.ParseAllObjectStates(m.lua, m.objs, m.objdirs, objOrder)
	if err != nil {
		return fmt.Errorf("objects.ParseAllObjectStates(%s) : %v", "", err)
	}
	m.Data[expectedObjStates] = allObjs
	return nil
}

func tryPut(d *types.J, from, to string, fun func(string) (interface{}, error)) {
	if d == nil {
		log.Println("Nil objects")
		return
	}

	var o interface{}
	fromFile, ok := (*d)[from]
	if !ok {
		fromFile = ""
		if _, ok := (*d)[to]; ok {
			// if there is not special key, but there is existant key, don't replace anything.
			return
		}
	}
	filename, ok := fromFile.(string)
	if !ok {
		log.Printf("non string filename found: %s", fromFile)
		filename = ""
	}

	o, _ = fun(filename)
	// ignore error for now

	(*d)[to] = o
	delete((*d), from)
}

func printMod(p string, m *Mod) error {
	b, err := json.MarshalIndent(m.Data, "", "  ")
	if err != nil {
		return fmt.Errorf("json.MarshalIndent(<mod>) : %v", err)
	}

	return ioutil.WriteFile(path.Join(p, "output.json"), b, 0644)
}

// prepForReverse creates the expected subdirectories in config path
func prepForReverse(cPath, modfile string) (types.J, error) {
	subDirs := []string{textSubdir, modsettingsDir, objectsSubdir}

	for _, s := range subDirs {
		p := path.Join(cPath, s)
		if _, err := os.Stat(p); err == nil {
			// directory already exists
		} else if os.IsNotExist(err) {
			err = os.Mkdir(p, 0777)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, fmt.Errorf("Undefined error checking for subdirectory %s : %v", s, err)
		}
	}

	mFile, err := os.Open(modfile)
	if err != nil {
		return nil, fmt.Errorf("os.Open(%s) : %v", modfile, err)
	}

	defer mFile.Close()

	b, err := ioutil.ReadAll(mFile)
	if err != nil {
		return nil, err
	}
	var o types.J
	err = json.Unmarshal(b, &o)
	if err != nil {
		return nil, err
	}

	return o, nil
}
