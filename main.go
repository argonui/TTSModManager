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
	moddir       = flag.String("moddir", "testdata/simple", "a directory containing tts mod configs")
	rev          = flag.Bool("reverse", false, "Instead of building a json from file structure, build file structure from json.")
	writeToSrc   = flag.Bool("writesrc", false, "When unbundling Lua, save the included 'require' files to the src/ directory.")
	modfile      = flag.String("modfile", "", "where to read from when reversing.")
	objstatesdir = flag.String("objdir", "", "If non-empty, don't build/reverse a full mod, only the object state array")
	objfile      = flag.String("objfile", "", "if building only object state list, output to this filename")
)

const (
	luasrcSubdir   = "src"
	xmlsrcSubdir   = "xml"
	modsettingsDir = "modsettings"
	objectsSubdir  = "objects"
)

func main() {
	flag.Parse()

	lua := file.NewTextOpsMulti(
		[]string{path.Join(*moddir, luasrcSubdir), path.Join(*moddir, objectsSubdir)},
		path.Join(*moddir, objectsSubdir),
	)
	xml := file.NewTextOpsMulti(
		[]string{path.Join(*moddir, xmlsrcSubdir), path.Join(*moddir, objectsSubdir)},
		path.Join(*moddir, objectsSubdir),
	)
	xmlSrc := file.NewTextOps(path.Join(*moddir, xmlsrcSubdir))
	luaSrc := file.NewTextOps(path.Join(*moddir, luasrcSubdir))
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
			XMLWriter:         xml,
			ObjWriter:         objs,
			ObjDirCreeator:    objdir,
			RootWrite:         rootops,
		}
		if *writeToSrc {
			r.LuaSrcWriter = luaSrc
			r.XMLSrcWriter = xmlSrc
		}
		err = r.Write(raw)
		if err != nil {
			log.Fatalf("reverse.Write(<%s>) failed : %v", *modfile, err)
		}
		return
	}
	if *modfile == "" {
		*modfile = path.Join(*moddir, "output.json")
	}

	basename := path.Base(*modfile)
	outputOps := file.NewJSONOps(path.Dir(*modfile))
	// only create an object state list, not an entire mod
	if *objstatesdir != "" {
		lua = file.NewLuaOpsMulti(
			[]string{path.Join(*moddir, textSubdir), path.Join(*moddir, *objstatesdir)},
			path.Join(*moddir, *objstatesdir),
		)
		objs = file.NewJSONOps(path.Join(*moddir, *objstatesdir))
		objdir = file.NewDirOps(path.Join(*moddir, *objstatesdir))
		outputOps = file.NewJSONOps(path.Dir(path.Join(*moddir, *objfile)))
		basename = path.Base(*objfile)
	}
	m := &mod.Mod{

		Lua:           lua,
		XML:           xml,
		Modsettings:   ms,
		Objs:          objs,
		Objdirs:       objdir,
		RootRead:      rootops,
		RootWrite:     outputOps,
		OnlyObjStates: (*objstatesdir != ""),
	}
	err := m.GenerateFromConfig()
	if err != nil {
		fmt.Printf("generateMod(<config>) : %v\n", err)
		return
	}
	err = m.Print(basename)
	if err != nil {
		log.Fatalf("printMod(...) : %v", err)
	}
}

// prepForReverse creates the expected subdirectories in config path
func prepForReverse(cPath, modfile string) (types.J, error) {
	subDirs := []string{luasrcSubdir, modsettingsDir, objectsSubdir, xmlsrcSubdir}

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
