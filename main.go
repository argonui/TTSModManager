package main

import (
	file "ModCreator/file"
	"ModCreator/mod"
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
)

const (
	textSubdir     = "src"
	modsettingsDir = "modsettings"
	objectsSubdir  = "objects"
)

func main() {
	flag.Parse()

	lua := file.NewLuaOpsMulti(
		[]string{path.Join(*moddir, textSubdir), path.Join(*moddir, objectsSubdir)},
		path.Join(*moddir, objectsSubdir),
	)
	ms := file.NewJSONOps(path.Join(*moddir, modsettingsDir))
	objs := file.NewJSONOps(path.Join(*moddir, objectsSubdir))
	objdir := file.NewDirOps(path.Join(*moddir, objectsSubdir))
	rootops := file.NewJSONOps(*moddir)

	if *rev {
		raw, err := prepForReverse(*moddir, *modfile)
		if err != nil {
			log.Fatalf("prepForReverse (%s) failed : %v", *modfile, err)
		}
		r := mod.Reverser{
			ModSettingsWriter: ms,
			LuaWriter:         lua,
			ObjWriter:         objs,
			ObjDirCreeator:    objdir,
			RootWrite:         rootops,
		}
		err = r.Write(raw)
		if err != nil {
			log.Fatalf("reverse.Write(<%s>) failed : %v", *modfile, err)
		}
		return
	}

	m := &mod.Mod{
		Lua:         lua,
		Modsettings: ms,
		Objs:        objs,
		Objdirs:     objdir,
		RootRead:    rootops,
		RootWrite:   rootops,
	}
	err := m.GenerateFromConfig()
	if err != nil {
		fmt.Printf("generateMod(<config>) : %v\n", err)
		return
	}
	err = m.Print()
	if err != nil {
		log.Fatalf("printMod(...) : %v", err)
	}
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
