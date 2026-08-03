package file

import (
	"os"
	"path"
	"path/filepath"
	"sort"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCreateDirFresh(t *testing.T) {
	base := t.TempDir()
	d := NewDirOps(base)

	got, err := d.CreateDir("", "objects")
	if err != nil {
		t.Fatalf("CreateDir(): %v", err)
	}
	if got != "objects" {
		t.Errorf("CreateDir() name = %q, want %q", got, "objects")
	}
	if info, err := os.Stat(path.Join(base, "objects")); err != nil || !info.IsDir() {
		t.Errorf("expected directory objects to exist, stat err=%v", err)
	}
}

// TestCreateDirCollision covers the case where the suggested directory already
// exists: CreateDir reuses it and returns the suggestion unchanged.
func TestCreateDirCollision(t *testing.T) {
	base := t.TempDir()
	if err := os.Mkdir(path.Join(base, "objects"), 0755); err != nil {
		t.Fatalf("setup Mkdir(): %v", err)
	}
	d := NewDirOps(base)

	got, err := d.CreateDir("", "objects")
	if err != nil {
		t.Fatalf("CreateDir(): %v", err)
	}
	if got != "objects" {
		t.Errorf("CreateDir() on collision name = %q, want %q", got, "objects")
	}
}

// TestCreateDirSuffixLoopError exercises the retry/suffix loop: when the parent
// path is a file, every mkdir attempt fails with a non-IsExist error, so the
// loop exhausts its retries and returns an error.
func TestCreateDirSuffixLoopError(t *testing.T) {
	base := t.TempDir()
	// Make "parent" a file, so path.Join(base, "parent", suggestion) can never
	// be created and IsExist is never true.
	if err := os.WriteFile(path.Join(base, "parent"), []byte("x"), 0644); err != nil {
		t.Fatalf("setup WriteFile(): %v", err)
	}
	d := NewDirOps(base)

	_, err := d.CreateDir("parent", "child")
	if err == nil {
		t.Fatalf("CreateDir() with file parent: wanted error, got none")
	}
}

func TestPreClearCheckAllowedPasses(t *testing.T) {
	base := t.TempDir()
	for _, name := range []string{"a.json", "b.ttslua", "sub/c.xml"} {
		full := path.Join(base, name)
		if err := os.MkdirAll(filepath.Dir(full), 0755); err != nil {
			t.Fatalf("setup MkdirAll(): %v", err)
		}
		if err := os.WriteFile(full, []byte("x"), 0644); err != nil {
			t.Fatalf("setup WriteFile(): %v", err)
		}
	}
	d := NewDirOps(base)

	if err := d.preClearCheck(); err != nil {
		t.Errorf("preClearCheck() with allowed extensions: wanted nil, got %v", err)
	}
}

func TestPreClearCheckDisallowedRejected(t *testing.T) {
	base := t.TempDir()
	if err := os.WriteFile(path.Join(base, "ok.json"), []byte("x"), 0644); err != nil {
		t.Fatalf("setup WriteFile(): %v", err)
	}
	if err := os.WriteFile(path.Join(base, "danger.txt"), []byte("x"), 0644); err != nil {
		t.Fatalf("setup WriteFile(): %v", err)
	}
	d := NewDirOps(base)

	if err := d.preClearCheck(); err == nil {
		t.Errorf("preClearCheck() with disallowed extension: wanted error, got nil")
	}
}

func TestClearRemovesAndRecreates(t *testing.T) {
	base := t.TempDir()
	if err := os.WriteFile(path.Join(base, "keep.json"), []byte("x"), 0644); err != nil {
		t.Fatalf("setup WriteFile(): %v", err)
	}
	d := NewDirOps(base)

	if err := d.Clear(); err != nil {
		t.Fatalf("Clear(): %v", err)
	}

	info, err := os.Stat(base)
	if err != nil || !info.IsDir() {
		t.Fatalf("expected base to exist after Clear, stat err=%v", err)
	}
	entries, err := os.ReadDir(base)
	if err != nil {
		t.Fatalf("ReadDir() after Clear: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected empty directory after Clear, found %d entries", len(entries))
	}
}

func TestClearRefusesUnsafeContents(t *testing.T) {
	base := t.TempDir()
	if err := os.WriteFile(path.Join(base, "danger.txt"), []byte("x"), 0644); err != nil {
		t.Fatalf("setup WriteFile(): %v", err)
	}
	d := NewDirOps(base)

	if err := d.Clear(); err == nil {
		t.Fatalf("Clear() with unsafe file: wanted error, got none")
	}
	// The unsafe file must not have been deleted.
	if _, err := os.Stat(path.Join(base, "danger.txt")); err != nil {
		t.Errorf("Clear() must not delete when the safety check fails: %v", err)
	}
}

func TestListFilesAndFolders(t *testing.T) {
	base := t.TempDir()
	if err := os.WriteFile(path.Join(base, "a.json"), []byte("x"), 0644); err != nil {
		t.Fatalf("setup WriteFile(): %v", err)
	}
	if err := os.WriteFile(path.Join(base, "b.json"), []byte("x"), 0644); err != nil {
		t.Fatalf("setup WriteFile(): %v", err)
	}
	if err := os.Mkdir(path.Join(base, "sub"), 0755); err != nil {
		t.Fatalf("setup Mkdir(): %v", err)
	}
	d := NewDirOps(base)

	files, folders, err := d.ListFilesAndFolders("")
	if err != nil {
		t.Fatalf("ListFilesAndFolders(): %v", err)
	}
	sort.Strings(files)
	sort.Strings(folders)

	if diff := cmp.Diff([]string{"a.json", "b.json"}, files); diff != "" {
		t.Errorf("files want != got:\n%v\n", diff)
	}
	if diff := cmp.Diff([]string{"sub"}, folders); diff != "" {
		t.Errorf("folders want != got:\n%v\n", diff)
	}
}
