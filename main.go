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
	"path/filepath"
)

var (
	moddir     = flag.String("moddir", "testdata/simple", "a directory containing tts mod configs")
	rev        = flag.Bool("reverse", false, "instead of building a json from file structure, build file structure from json.")
	writeToSrc = flag.Bool("writesrc", false, "when unbundling Lua, save the included 'require' files to the src/ directory.")
	modfile    = flag.String("modfile", "", "where to read from when reversing.")
	objin      = flag.String("objin", "", "if non-empty, don't build/reverse a full mod, only an object state array")
	objout     = flag.String("objout", "", "if building only object state list, output to this filename")
)

var (
	luasrcSubdir   = "src"
	xmlsrcSubdir   = "xml"
	modsettingsDir = "modsettings"
	objectsSubdir  = "objects"
)

func main() {
	flag.Parse()

	if (*objin == "") != (*objout == "") {
		log.Fatalln("Must set either both or neither of {objin,objout}.")
	}

	lua := file.NewTextOpsMulti(
		[]string{filepath.Join(*moddir, luasrcSubdir), filepath.Join(*moddir, objectsSubdir)},
		filepath.Join(*moddir, objectsSubdir),
	)
	xml := file.NewTextOpsMulti(
		[]string{filepath.Join(*moddir, xmlsrcSubdir), filepath.Join(*moddir, objectsSubdir)},
		filepath.Join(*moddir, objectsSubdir),
	)
	xmlSrc := file.NewTextOps(filepath.Join(*moddir, xmlsrcSubdir))
	luaSrc := file.NewTextOps(filepath.Join(*moddir, luasrcSubdir))
	ms := file.NewJSONOps(filepath.Join(*moddir, modsettingsDir))
	objs := file.NewJSONOps(filepath.Join(*moddir, objectsSubdir))
	objdir := file.NewDirOps(filepath.Join(*moddir, objectsSubdir))
	rootops := file.NewJSONOps(*moddir)

	basename := filepath.Base(*modfile)
	outputOps := file.NewJSONOps(filepath.Dir(*modfile))

	// handling for saved objects instead of a full savegame
	if *objin != "" {
		objs = file.NewJSONOps(filepath.Dir(*objin))
		objdir = file.NewDirOps(filepath.Dir(*objin))
		lua = file.NewTextOpsMulti(
			[]string{filepath.Join(*moddir, luasrcSubdir), filepath.Dir(*objin)},
			filepath.Dir(*objout),
		)
		xml = file.NewTextOpsMulti(
			[]string{filepath.Join(*moddir, xmlsrcSubdir), filepath.Dir(*objin)},
			filepath.Dir(*objout),
		)
		basename = filepath.Base(*objout)
		outputOps = file.NewJSONOps(filepath.Dir(*objout))
	}

	if *rev {
		if *objin != "" {
			*modfile = *objin
			objs = file.NewJSONOps(filepath.Dir(*objout))
		}
		raw, err := prepForReverse(*moddir, *modfile)
		if err != nil {
			log.Fatalf("prepForReverse (%s) failed : %v", *modfile, err)
		}

		r := mod.Reverser{
			ModSettingsWriter: ms,
			LuaWriter:         lua,
			XMLWriter:         xml,
			ObjWriter:         objs,
			ObjDirCreator:     objdir,
			RootWrite:         rootops,
			OnlyObjState:      *objin,
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
		*modfile = filepath.Join(*moddir, "output.json")
	}

	// setting this to empty instead of default return value (".") if not found
	OnlyObjStates := filepath.Base(*objin)
	if OnlyObjStates == "." {
		OnlyObjStates = ""
	}

	m := &mod.Mod{
		Lua:           lua,
		XML:           xml,
		Modsettings:   ms,
		Objs:          objs,
		Objdirs:       objdir,
		RootRead:      rootops,
		RootWrite:     outputOps,
		OnlyObjStates: OnlyObjStates,
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
		p := filepath.Join(cPath, s)
		if _, err := os.Stat(p); err == nil {
			// directory already exists
		} else if os.IsNotExist(err) {
			err = os.Mkdir(p, 0777)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, fmt.Errorf("undefined error checking for subdirectory %s : %v", s, err)
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
