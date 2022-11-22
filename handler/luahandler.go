package handler

import (
	"ModCreator/bundler"
	"ModCreator/file"
	"fmt"
)

// LuaHandler is needed because handling if a script should be written to src or to objects folder,
// and if it it long enough to be written to separate file at all has become
// burdensome. abstract into struct
type LuaHandler struct {
	ObjWriter file.TextWriter
	SrcWriter file.TextWriter
	Reader    file.TextReader
}

// HandleAction describes at the end of handling,
// a single value should be written to a single key
type HandleAction struct {
	Noop  bool
	Key   string
	Value string
}

// WhileReadingFromFile consolidates expected behavior of both objects and root
// json while reading lua from file.
func (lh *LuaHandler) WhileReadingFromFile(rawj map[string]interface{}) (HandleAction, error) {
	rawscript := ""
	if spraw, ok := rawj["LuaScript_path"]; ok {
		sp, ok := spraw.(string)
		if !ok {
			return HandleAction{}, fmt.Errorf("Expected LuaScript_path to be type string, was %T", spraw)
		}
		encoded, err := lh.Reader.EncodeFromFile(sp)
		if err != nil {
			return HandleAction{}, fmt.Errorf("l.EncodeFromFile(%s) : %v", sp, err)
		}
		rawscript = encoded
	}
	if sraw, ok := rawj["LuaScript"]; ok {
		s, ok := sraw.(string)
		if !ok {
			return HandleAction{}, fmt.Errorf("Expected LuaScript to be type string, was %T", sraw)
		}
		rawscript = s
	}
	bundled, err := bundler.Bundle(rawscript, lh.Reader)
	if err != nil {
		return HandleAction{}, fmt.Errorf("Bundle(%s): %v", rawscript, err)
	}
	if bundled == "" {
		return HandleAction{
			Noop: true,
		}, nil
	}

	return HandleAction{
		Key:   "LuaScript",
		Value: bundled,
		Noop:  false,
	}, nil
}

// WhileWritingToFile consolidates the logic and flow of conditionally writing
// lua to a file, and which file to write to
func (lh *LuaHandler) WhileWritingToFile(rawj map[string]interface{}, possiblefname string) (HandleAction, error) {
	rawscript, ok := rawj["LuaScript"]
	if !ok {
		return HandleAction{Noop: true}, nil
	}
	script, ok := rawscript.(string)
	if !ok {
		return HandleAction{}, fmt.Errorf("Value of content at LuaScript expected to be string, was %v, type %T",
			rawscript, rawscript)
	}

	allScripts, err := bundler.UnbundleAll(script)
	if err != nil {
		return HandleAction{}, fmt.Errorf("UnbundleAll(...): %v", err)
	}
	// root bundle is promised to exist
	rootscript, _ := allScripts[bundler.Rootname]
	returnAction := HandleAction{Noop: false}
	if len(rootscript) > 80 {
		err = lh.ObjWriter.EncodeToFile(rootscript, possiblefname)
		if err != nil {
			return HandleAction{}, fmt.Errorf("EncodeToFile(<root script>, %s): %v", possiblefname, err)
		}
		returnAction.Key = "LuaScript_path"
		returnAction.Value = possiblefname
	} else {
		returnAction.Key = "LuaScript"
		returnAction.Value = rootscript
	}
	delete(allScripts, bundler.Rootname)

	for k, script := range allScripts {
		fname := fmt.Sprintf("%s.ttslua", k)
		if lh.SrcWriter == nil {
			break
		}
		err := lh.SrcWriter.EncodeToFile(script, fname)
		if err != nil {
			return HandleAction{}, fmt.Errorf("SrcWriter.EncodeToFile(<>, %s): %v", fname, err)
		}
	}

	return returnAction, nil
}
