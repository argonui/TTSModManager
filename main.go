package main

import (
	file "ModCreator/file"
	objects "ModCreator/objects"
	"ModCreator/reverse"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
)

var (
	config  = flag.String("moddir", "testdata/simple", "a directory containing tts mod configs")
	rev     = flag.Bool("reverse", false, "Instead of building a json from file structure, build file structure from json.")
	modfile = flag.String("modfile", "", "where to read from when reversing.")

	expectedStr       = []string{"SaveName", "Date", "VersionNumber", "GameMode", "GameType", "GameComplexity", "Table", "Sky", "Note", "LuaScript", "LuaScriptState", "XmlUI"}
	expectedObj       = []string{"TabStates", "MusicPlayer", "Grid", "Lighting", "Hands", "ComponentTags", "Turns"}
	expectedObjArr    = []string{"CameraStates", "DecalPallet", "CustomUIAssets", "SnapPoints"}
	expectedObjStates = "ObjectStates"
)

const (
	textSubdir    = "src"
	jsonSubdir    = "modsettings"
	objectsSubdir = "objects"
)

// Config is how users will specify their mod's configuration.
type Config struct {
	Raw Obj `json:"-"`
}

// Obj is a simpler way to refer to a json map.
type Obj map[string]interface{}

// ObjArray is a simple way to refer to an array of json maps
type ObjArray []map[string]interface{}

// Mod is used as the accurate representation of what gets printed when
// module creation is done
type Mod struct {
	Data Obj
}

func main() {
	flag.Parse()

	lua := file.NewLuaOps(path.Join(*config, textSubdir))
	j := file.NewJSONOps(path.Join(*config, jsonSubdir))

	if *rev {
		raw, err := prepForReverse(*config, *modfile)
		if err != nil {
			log.Fatalf("prepForReverse (%s) failed : %v", *modfile, err)
		}
		err = reverse.Write(raw, lua, j, *config, expectedStr, expectedObj, expectedObjArr)
		if err != nil {
			log.Fatalf("reverse.Write(<%s>) failed : %v", *modfile, err)
		}
		return
	}

	c, err := readConfig(*config)
	if err != nil {
		fmt.Printf("readConfig(%s) : %v\n", *config, err)
		return
	}

	m, err := generateMod(*config, lua, j, c)
	if err != nil {
		fmt.Printf("generateMod(<config>) : %v\n", err)
		return
	}
	err = printMod(*config, m)
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

func generateMod(p string, lua file.LuaReader, j file.JSONReader, c *Config) (*Mod, error) {
	if c == nil {
		return nil, fmt.Errorf("nil config")
	}
	var m Mod

	m.Data = c.Raw

	plainObj := func(s string) (interface{}, error) {
		return j.ReadObj(s)
	}
	objArray := func(s string) (interface{}, error) {
		return j.ReadObjArray(s)
	}
	luaGet := func(s string) (interface{}, error) {
		return lua.EncodeFromFile(s)
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

	allObjs, err := objects.ParseAllObjectStates(path.Join(p, objectsSubdir), lua)
	if err != nil {
		return nil, fmt.Errorf("objects.ParseAllObjectStates(%s) : %v", path.Join(p, objectsSubdir), err)
	}
	m.Data[expectedObjStates] = allObjs
	return &m, nil
}

func tryPut(d *Obj, from, to string, fun func(string) (interface{}, error)) {
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
func prepForReverse(cPath, modfile string) (Obj, error) {
	subDirs := []string{textSubdir, jsonSubdir, objectsSubdir}

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
	var o Obj
	err = json.Unmarshal(b, &o)
	if err != nil {
		return nil, err
	}

	return o, nil
}
