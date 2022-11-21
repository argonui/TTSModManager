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
	ObjWriter file.LuaWriter
	SrcWriter file.LuaWriter
	Reader    file.LuaReader
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
