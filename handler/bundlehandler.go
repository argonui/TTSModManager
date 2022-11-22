package handler

import (
	"ModCreator/bundler"
	"ModCreator/file"
	"fmt"
	"strings"
)

// Handler is needed because handling if a script should be written to src or to objects folder,
// and if it it long enough to be written to separate file at all has become
// burdensome. abstract into struct. Also for XML bundling
type Handler struct {
	DefaultWriter file.TextWriter
	SrcWriter     file.TextWriter
	Reader        file.TextReader

	key, keypath, extension string
	bundle                  func(string, file.TextReader) (string, error)
	unbundle                func(string) (map[string]string, error)
}

// NewLuaHandler fills in relevant info for lua bundling
func NewLuaHandler() *Handler {
	return &Handler{
		key:       "LuaScript",
		keypath:   "LuaScript_path",
		extension: ".ttslua",
		bundle:    bundler.Bundle,
		unbundle:  bundler.UnbundleAll,
	}
}

// NewXMLHandler fills in relevant info for lua bundling
func NewXMLHandler() *Handler {
	return &Handler{
		key:       "XmlUI",
		keypath:   "XmlUI_path",
		extension: ".xml",
		bundle:    bundler.BundleXML,
		unbundle:  bundler.UnbundleAllXML,
	}
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
func (h *Handler) WhileReadingFromFile(rawj map[string]interface{}) (HandleAction, error) {
	rawscript := ""
	if spraw, ok := rawj[h.keypath]; ok {
		sp, ok := spraw.(string)
		if !ok {
			return HandleAction{}, fmt.Errorf("Expected %s to be type string, was %T", h.keypath, spraw)
		}
		encoded, err := h.Reader.EncodeFromFile(sp)
		if err != nil {
			return HandleAction{}, fmt.Errorf("l.EncodeFromFile(%s) : %v", sp, err)
		}
		rawscript = encoded
	}
	if sraw, ok := rawj[h.key]; ok {
		s, ok := sraw.(string)
		if !ok {
			return HandleAction{}, fmt.Errorf("Expected %s to be type string, was %T", h.key, sraw)
		}
		rawscript = s
	}
	bundled, err := h.bundle(rawscript, h.Reader)
	if err != nil {
		return HandleAction{}, fmt.Errorf("Bundle(%s): %v", rawscript, err)
	}
	if bundled == "" {
		return HandleAction{
			Noop: true,
		}, nil
	}

	return HandleAction{
		Key:   h.key,
		Value: bundled,
		Noop:  false,
	}, nil
}

// WhileWritingToFile consolidates the logic and flow of conditionally writing
// lua to a file, and which file to write to
func (h *Handler) WhileWritingToFile(rawj map[string]interface{}, possiblefname string) (HandleAction, error) {
	rawscript, ok := rawj[h.key]
	if !ok {
		return HandleAction{Noop: true}, nil
	}
	script, ok := rawscript.(string)
	if !ok {
		return HandleAction{}, fmt.Errorf("Value of content at %s expected to be string, was %v, type %T",
			h.key, rawscript, rawscript)
	}

	allScripts, err := h.unbundle(script)
	if err != nil {
		return HandleAction{}, fmt.Errorf("UnbundleAll(...): %v", err)
	}
	// root bundle is promised to exist
	rootscript, _ := allScripts[bundler.Rootname]
	returnAction := HandleAction{Noop: false}
	if len(rootscript) > 80 {
		err = h.DefaultWriter.EncodeToFile(rootscript, possiblefname)
		if err != nil {
			return HandleAction{}, fmt.Errorf("EncodeToFile(<root script>, %s): %v", possiblefname, err)
		}
		returnAction.Key = h.keypath
		returnAction.Value = possiblefname
	} else {
		returnAction.Key = h.key
		returnAction.Value = rootscript
	}
	delete(allScripts, bundler.Rootname)

	for k, script := range allScripts {
		fname := k
		if !strings.HasSuffix(k, h.extension) {
			fname = fmt.Sprintf("%s%s", k, h.extension)
		}
		if h.SrcWriter == nil {
			break
		}
		err := h.SrcWriter.EncodeToFile(script, fname)
		if err != nil {
			return HandleAction{}, fmt.Errorf("SrcWriter.EncodeToFile(<>, %s): %v", fname, err)
		}
	}

	return returnAction, nil
}
