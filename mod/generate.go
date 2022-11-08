package mod

import (
	"ModCreator/file"
	"ModCreator/objects"
	"ModCreator/types"
	"fmt"
	"log"
)

var (
	// ExpectedStr are op level keys of mods that expected to have string values
	ExpectedStr = []string{"SaveName", "Date", "VersionNumber", "GameMode", "GameType", "GameComplexity", "Table", "Sky", "Note", "LuaScript", "LuaScriptState", "XmlUI"}
	// ExpectedObj are op level keys of mods that expected to have json values
	ExpectedObj = []string{"TabStates", "MusicPlayer", "Grid", "Lighting", "Hands", "ComponentTags", "Turns"}
	// ExpectedObjArr are op level keys of mods that expected to have array of json values
	ExpectedObjArr = []string{"CameraStates", "DecalPallet", "CustomUIAssets", "SnapPoints", "Decals"}
	// ExpectedObjStates is the key which holds all objects in a mod
	ExpectedObjStates = "ObjectStates"
)

// Mod is used as the accurate representation of what gets printed when
// module creation is done
type Mod struct {
	Data types.J

	RootRead    file.JSONReader
	RootWrite   file.JSONWriter
	Lua         file.LuaReader
	Modsettings file.JSONReader
	Objs        file.JSONReader
	Objdirs     file.DirExplorer
}

// GenerateFromConfig uses RootRead for reading entire mod config
func (m *Mod) GenerateFromConfig() error {
	raw, err := m.RootRead.ReadObj("config.json")
	if err != nil {
		return fmt.Errorf("could not read a root config: %v", err)
	}
	return m.generate(raw)
}

func (m *Mod) generate(raw types.J) error {
	m.Data = raw

	plainObj := func(s string) (interface{}, error) {
		return m.Modsettings.ReadObj(s)
	}
	objArray := func(s string) (interface{}, error) {
		return m.Modsettings.ReadObjArray(s)
	}
	luaGet := func(s string) (interface{}, error) {
		return m.Lua.EncodeFromFile(s)
	}

	ext := "_path"
	for _, stringbased := range ExpectedStr {
		tryPut(&m.Data, stringbased+ext, stringbased, luaGet)
	}

	for _, objbased := range ExpectedObj {
		tryPut(&m.Data, objbased+ext, objbased, plainObj)
	}

	for _, objarraybased := range ExpectedObjArr {
		tryPut(&m.Data, objarraybased+ext, objarraybased, objArray)
	}

	objOrder := []string{}
	files, _, _ := m.Objdirs.ListFilesAndFolders("")
	hasObjects := len(files) > 0

	err := file.ForceParseIntoStrArray(&m.Data, "ObjectStates_order", &objOrder)
	if hasObjects && err != nil {
		return fmt.Errorf("Has Objects, but can't discern their order: %v", err)
	}

	allObjs, err := objects.ParseAllObjectStates(m.Lua, m.Objs, m.Objdirs, objOrder)
	if err != nil {
		return fmt.Errorf("objects.ParseAllObjectStates(%s) : %v", "", err)
	}
	m.Data[ExpectedObjStates] = allObjs
	return nil
}

// Print outputs internal representation of mod to json file with indents
func (m *Mod) Print() error {
	return m.RootWrite.WriteObj(m.Data, "output.json")
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

//
// func readConfig(cPath string) (*Config, error) {
// 	// Open our jsonFile
// 	cFile, err := os.Open(path.Join(cPath, "config.json"))
// 	// if we os.Open returns an error then handle it
// 	if err != nil {
// 		return nil, fmt.Errorf("os.Open(%s): %v", path.Join(cPath, "config.json"), err)
// 	}
// 	// defer the closing of our jsonFile so that we can parse it later on
// 	defer cFile.Close()
//
// 	b, err := ioutil.ReadAll(cFile)
// 	if err != nil {
// 		return nil, fmt.Errorf("ioutil.Readall(%s) : %v", path.Join(cPath, "config.json"), err)
// 	}
// 	var c Config
//
// 	err = json.Unmarshal(b, &c.Raw)
// 	if err != nil {
// 		return nil, fmt.Errorf("json.Unmarshal(%s) : %v", b, err)
// 	}
// 	return &c, nil
// }
