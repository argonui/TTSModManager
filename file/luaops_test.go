package file

import (
	"fmt"
	"testing"
)

type fakeFiles struct {
	fs map[string][]byte
}

func (ff *fakeFiles) read(s string) ([]byte, error) {
	if val, ok := ff.fs[s]; ok {
		return val, nil
	}
	return nil, fmt.Errorf("fake file %s not found", s)
}

func TestRequirSimple(t *testing.T) {
	ff := &fakeFiles{
		fs: map[string][]byte{
			"src/base": []byte(`a + b = c`),
		},
	}
	l := &LuaOps{
		basepath:        "src",
		readFileToBytes: ff.read,
	}

	want := "a + b = c"
	got, err := l.EncodeFromFile("base")
	if err != nil {
		t.Errorf("encode error %v", err)
	}
	if want != got {
		t.Errorf("want <%s> got <%s>", want, got)
	}
}

func TestRequireOnce(t *testing.T) {
	ff := &fakeFiles{
		fs: map[string][]byte{
			"src/base":    []byte(`require(\"foo/bar\")`),
			"src/foo/bar": []byte(`a + b = c`),
		},
	}
	l := &LuaOps{
		basepath:        "src",
		readFileToBytes: ff.read,
	}

	want := `
a + b = c
`
	got, err := l.EncodeFromFile("base")
	if err != nil {
		t.Errorf("encode error %v", err)
	}
	if want != got {
		t.Errorf("want <%s> got <%s>", want, got)
	}
}

func TestRequireUnescaped(t *testing.T) {
	ff := &fakeFiles{
		fs: map[string][]byte{
			"src/base":    []byte(`require("foo/bar")`),
			"src/foo/bar": []byte(`a + b = c`),
		},
	}
	l := &LuaOps{
		basepath:        "src",
		readFileToBytes: ff.read,
	}

	want := `
a + b = c
`
	got, err := l.EncodeFromFile("base")
	if err != nil {
		t.Errorf("encode error %v", err)
	}
	if want != got {
		t.Errorf("want <%s> got <%s>", want, got)
	}
}

func TestRequireNested(t *testing.T) {
	ff := &fakeFiles{
		fs: map[string][]byte{
			"src/base": []byte(`require(\"foo/bar\")
require(\"foo/baz\")`),
			"src/foo/bar": []byte(`a + b = c`),
			"src/foo/baz": []byte(`a + b = d
require(\"util\")
var y = 55
`),
			"src/util": []byte(`var x = 42`),
		},
	}
	l := &LuaOps{
		basepath:        "src",
		readFileToBytes: ff.read,
	}

	want := `
a + b = c


a + b = d

var x = 42

var y = 55

`
	got, err := l.EncodeFromFile("base")
	if err != nil {
		t.Errorf("encode error %v", err)
	}
	if want != got {
		t.Errorf("want <%s> got <%s>", want, got)
	}
}
