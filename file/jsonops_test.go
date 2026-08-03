package file

import (
	"os"
	"path"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestWriteObjReadObjRoundTrip(t *testing.T) {
	dir := t.TempDir()
	j := NewJSONOps(dir)

	want := map[string]interface{}{
		"SaveName": "cool mod",
		"Nested": map[string]interface{}{
			"count": float64(3),
			"flag":  true,
		},
	}

	if err := j.WriteObj(want, "config.json"); err != nil {
		t.Fatalf("WriteObj(): %v", err)
	}

	got, err := j.ReadObj("config.json")
	if err != nil {
		t.Fatalf("ReadObj(): %v", err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("want != got:\n%v\n", diff)
	}
}

func TestWriteObjArrayReadObjArrayRoundTrip(t *testing.T) {
	dir := t.TempDir()
	j := NewJSONOps(dir)

	want := []map[string]interface{}{
		{"GUID": "a", "n": float64(1)},
		{"GUID": "b", "n": float64(2)},
	}

	if err := j.WriteObjArray(want, "arr.json"); err != nil {
		t.Fatalf("WriteObjArray(): %v", err)
	}

	got, err := j.ReadObjArray("arr.json")
	if err != nil {
		t.Fatalf("ReadObjArray(): %v", err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("want != got:\n%v\n", diff)
	}
}

func TestReadObjMalformedJSONReturnsError(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(path.Join(dir, "bad.json"), []byte("{not valid json"), 0644); err != nil {
		t.Fatalf("setup WriteFile(): %v", err)
	}
	j := NewJSONOps(dir)

	got, err := j.ReadObj("bad.json")
	if err == nil {
		t.Fatalf("ReadObj() on malformed JSON: wanted error, got none (got=%v)", got)
	}
	// The real impl returns a non-nil empty map alongside the error so callers
	// that ignore the error don't panic on a nil map.
	if got == nil {
		t.Errorf("ReadObj() on malformed JSON: wanted non-nil empty map, got nil")
	}
}

func TestReadObjArrayMalformedJSONReturnsError(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(path.Join(dir, "bad.json"), []byte("[not valid json"), 0644); err != nil {
		t.Fatalf("setup WriteFile(): %v", err)
	}
	j := NewJSONOps(dir)

	got, err := j.ReadObjArray("bad.json")
	if err == nil {
		t.Fatalf("ReadObjArray() on malformed JSON: wanted error, got none (got=%v)", got)
	}
	if got == nil {
		t.Errorf("ReadObjArray() on malformed JSON: wanted non-nil empty slice, got nil")
	}
}

func TestReadObjMissingFileReturnsError(t *testing.T) {
	j := NewJSONOps(t.TempDir())

	got, err := j.ReadObj("does-not-exist.json")
	if err == nil {
		t.Fatalf("ReadObj() on missing file: wanted error, got none")
	}
	if got == nil {
		t.Errorf("ReadObj() on missing file: wanted non-nil empty map, got nil")
	}
}
