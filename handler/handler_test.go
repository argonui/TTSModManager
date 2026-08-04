package handler

import (
	"ModCreator/tests"
	"strings"
	"testing"
)

// constructors returns both handlers so each behavior is asserted for Lua and
// XML, which share all of handler's logic and differ only in configuration.
func constructors() []struct {
	name    string
	h       func() *Handler
	key     string
	keypath string
	ext     string
} {
	return []struct {
		name    string
		h       func() *Handler
		key     string
		keypath string
		ext     string
	}{
		{"lua", NewLuaHandler, "LuaScript", "LuaScript_path", ".ttslua"},
		{"xml", NewXMLHandler, "XmlUI", "XmlUI_path", ".xml"},
	}
}

// TestWhileWritingToFileInlineVsFile covers the 80-char root-script boundary:
// a short root stays inline (returned under key), a long one is written to a
// file (returned under keypath).
func TestWhileWritingToFileInlineVsFile(t *testing.T) {
	for _, c := range constructors() {
		t.Run(c.name+"/short_stays_inline", func(t *testing.T) {
			short := "print('hi')"
			if len(short) > 80 {
				t.Fatalf("test fixture should be <= 80 chars, was %d", len(short))
			}
			def := tests.NewFF()
			h := c.h()
			h.DefaultWriter = def

			got, err := h.WhileWritingToFile(map[string]interface{}{c.key: short}, "some.file")
			if err != nil {
				t.Fatalf("WhileWritingToFile: %v", err)
			}
			if got.Noop {
				t.Fatalf("expected non-noop action")
			}
			if got.Key != c.key {
				t.Errorf("inline should return key %q, got %q", c.key, got.Key)
			}
			if got.Value != short {
				t.Errorf("inline value = %q, want %q", got.Value, short)
			}
			if len(def.Fs) != 0 {
				t.Errorf("nothing should have been written to a file, got %v", def.Fs)
			}
		})

		t.Run(c.name+"/long_written_to_file", func(t *testing.T) {
			long := strings.Repeat("a", 100)
			if len(long) <= 80 {
				t.Fatalf("test fixture should be > 80 chars, was %d", len(long))
			}
			def := tests.NewFF()
			h := c.h()
			h.DefaultWriter = def

			got, err := h.WhileWritingToFile(map[string]interface{}{c.key: long}, "some.file")
			if err != nil {
				t.Fatalf("WhileWritingToFile: %v", err)
			}
			if got.Noop {
				t.Fatalf("expected non-noop action")
			}
			if got.Key != c.keypath {
				t.Errorf("file should return keypath %q, got %q", c.keypath, got.Key)
			}
			if got.Value != "some.file" {
				t.Errorf("file value = %q, want the possiblefname %q", got.Value, "some.file")
			}
			if def.Fs["some.file"] != long {
				t.Errorf("file content = %q, want %q", def.Fs["some.file"], long)
			}
		})
	}
}

// TestWhileWritingToFileNoop covers a rawj that lacks the content key entirely.
func TestWhileWritingToFileNoop(t *testing.T) {
	for _, c := range constructors() {
		t.Run(c.name, func(t *testing.T) {
			def := tests.NewFF()
			h := c.h()
			h.DefaultWriter = def

			got, err := h.WhileWritingToFile(map[string]interface{}{}, "some.file")
			if err != nil {
				t.Fatalf("WhileWritingToFile: %v", err)
			}
			if !got.Noop {
				t.Errorf("missing key should be a noop, got %+v", got)
			}
			if len(def.Fs) != 0 {
				t.Errorf("nothing should have been written, got %v", def.Fs)
			}
		})
	}
}

// bundledWithModule produces a bundled script whose unbundling yields both a
// root module and one required module, so the SrcWriter branch is exercised.
// modName is returned without the handler extension.
func bundledWithModule(t *testing.T, h *Handler, name, ext string) (string, string) {
	t.Helper()
	reader := tests.NewFF()
	var root string
	switch ext {
	case ".ttslua":
		reader.Fs[name+".ttslua"] = "-- " + name + " module\nreturn 42\n"
		root = "require(\"" + name + "\")\nprint('root')"
	case ".xml":
		reader.Fs[name+".xml"] = "<Text>hi</Text>"
		root = "<Panel>\n  <Include src=\"" + name + "\"/>\n</Panel>"
	default:
		t.Fatalf("unexpected extension %q", ext)
	}
	bundled, err := h.bundle(root, reader)
	if err != nil {
		t.Fatalf("bundling fixture: %v", err)
	}
	// sanity: it must actually unbundle into more than just the root
	all, _, err := h.unbundle(bundled)
	if err != nil {
		t.Fatalf("unbundling fixture: %v", err)
	}
	if _, ok := all[name]; !ok {
		t.Fatalf("fixture did not contain required module %q, got keys %v", name, keys(all))
	}
	return bundled, name
}

func keys(m map[string]string) []string {
	out := []string{}
	for k := range m {
		out = append(out, k)
	}
	return out
}

// TestWhileWritingToFileSrcWriter covers that required modules are written to
// SrcWriter when it is non-nil and skipped when it is nil.
func TestWhileWritingToFileSrcWriter(t *testing.T) {
	for _, c := range constructors() {
		t.Run(c.name+"/nonnil_writes_modules", func(t *testing.T) {
			h := c.h()
			bundled, mod := bundledWithModule(t, h, "mymod", c.ext)

			def := tests.NewFF()
			src := tests.NewFF()
			h.DefaultWriter = def
			h.SrcWriter = src

			got, err := h.WhileWritingToFile(map[string]interface{}{c.key: bundled}, "root.file")
			if err != nil {
				t.Fatalf("WhileWritingToFile: %v", err)
			}
			if got.Noop {
				t.Fatalf("expected non-noop action")
			}
			wantFile := mod + c.ext
			if _, ok := src.Fs[wantFile]; !ok {
				t.Errorf("required module %q should have been written to SrcWriter, got %v", wantFile, src.Fs)
			}
		})

		t.Run(c.name+"/nil_skips_modules", func(t *testing.T) {
			h := c.h()
			bundled, mod := bundledWithModule(t, h, "mymod", c.ext)

			def := tests.NewFF()
			src := tests.NewFF()
			h.DefaultWriter = def
			h.SrcWriter = nil // --writesrc not set

			got, err := h.WhileWritingToFile(map[string]interface{}{c.key: bundled}, "root.file")
			if err != nil {
				t.Fatalf("WhileWritingToFile: %v", err)
			}
			if got.Noop {
				t.Fatalf("expected non-noop action")
			}
			if _, ok := src.Fs[mod+c.ext]; ok {
				t.Errorf("required module should NOT have been written when SrcWriter is nil")
			}
			if len(src.Fs) != 0 {
				t.Errorf("SrcWriter fake should be untouched, got %v", src.Fs)
			}
		})
	}
}

// TestWhileReadingFromFileInline covers content present inline under key.
func TestWhileReadingFromFileInline(t *testing.T) {
	for _, c := range constructors() {
		t.Run(c.name, func(t *testing.T) {
			content := "print('inline')"
			h := c.h()
			h.Reader = tests.NewFF()

			got, err := h.WhileReadingFromFile(map[string]interface{}{c.key: content})
			if err != nil {
				t.Fatalf("WhileReadingFromFile: %v", err)
			}
			if got.Noop {
				t.Fatalf("expected non-noop action")
			}
			if got.Key != c.key {
				t.Errorf("read should return key %q, got %q", c.key, got.Key)
			}
			if got.Value != content {
				t.Errorf("read value = %q, want %q", got.Value, content)
			}
		})
	}
}

// TestWhileReadingFromFilePath covers content stored in a *_path file, which is
// read through the injected reader.
func TestWhileReadingFromFilePath(t *testing.T) {
	for _, c := range constructors() {
		t.Run(c.name, func(t *testing.T) {
			content := "print('from file')"
			reader := tests.NewFF()
			reader.Fs["root"+c.ext] = content
			h := c.h()
			h.Reader = reader

			got, err := h.WhileReadingFromFile(map[string]interface{}{c.keypath: "root" + c.ext})
			if err != nil {
				t.Fatalf("WhileReadingFromFile: %v", err)
			}
			if got.Noop {
				t.Fatalf("expected non-noop action")
			}
			if got.Key != c.key {
				t.Errorf("read should return key %q, got %q", c.key, got.Key)
			}
			if got.Value != content {
				t.Errorf("read value = %q, want %q", got.Value, content)
			}
		})
	}
}

// TestWhileReadingFromFileNoop covers neither key nor keypath being present.
func TestWhileReadingFromFileNoop(t *testing.T) {
	for _, c := range constructors() {
		t.Run(c.name, func(t *testing.T) {
			h := c.h()
			h.Reader = tests.NewFF()

			got, err := h.WhileReadingFromFile(map[string]interface{}{})
			if err != nil {
				t.Fatalf("WhileReadingFromFile: %v", err)
			}
			if !got.Noop {
				t.Errorf("empty rawj should be a noop, got %+v", got)
			}
		})
	}
}
