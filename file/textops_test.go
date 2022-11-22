package file

import (
	"fmt"

	"testing"

	"github.com/google/go-cmp/cmp"

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

func (ff *fakeFiles) write(path string, val []byte) error {
	ff.fs[path] = val
	return nil
}

func TestRead(t *testing.T) {
	ff := &fakeFiles{
		fs: map[string][]byte{"foo/bar": []byte("bar = 2")},
	}
	l := &TextOps{
		basepaths:       []string{"foo"},
		readFileToBytes: ff.read,
	}

	want := "bar = 2"
	got, err := l.EncodeFromFile("bar")
	if err != nil {
		t.Errorf("EncodeFromFile(bar): %v", err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("want != got:\n%v\n", diff)
	}
}

func TestReadMulti(t *testing.T) {
	ff := &fakeFiles{
		fs: map[string][]byte{
			"foo/old":  []byte("not here"),
			"foo2/bar": []byte("bar = 2"),
		},
	}
	l := &TextOps{
		basepaths:       []string{"foo", "foo2"},
		readFileToBytes: ff.read,
	}

	want := "bar = 2"
	got, err := l.EncodeFromFile("bar")
	if err != nil {
		t.Fatalf("EncodeFromFile(bar): %v", err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("want != got:\n%v\n", diff)
	}
}

func TestReadErr(t *testing.T) {
	ff := &fakeFiles{
		fs: map[string][]byte{
			"foo/old":  []byte("not here"),
			"foo2/bar": []byte("bar = 2"),
		},
	}
	l := &TextOps{
		basepaths:       []string{"foo", "foo2"},
		readFileToBytes: ff.read,
	}

	_, err := l.EncodeFromFile("notfound")
	if err == nil {
		t.Fatalf("Wanted error, got none")
	}
}

func TestWrite(t *testing.T) {
	ff := &fakeFiles{
		fs: map[string][]byte{},
	}
	l := TextOps{
		writeBasepath:    "foo",
		writeBytesToFile: ff.write,
	}

	err := l.EncodeToFile("var bar = 2", "bar.lua")
	if err != nil {
		t.Fatalf("No error expected got %v", err)
	}

	want := []byte("var bar = 2")
	got, ok := ff.fs["foo/bar.lua"]
	if !ok {
		t.Errorf("fake files not found")
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("want != got:\n%v\n", diff)
	}
}
