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
	// Clear now drops the ownership marker into the recreated dir (#96),
	// so the only permitted entry is that marker.
	for _, e := range entries {
		if e.Name() != managedMarker {
			t.Errorf("expected only the ownership marker after Clear, found %q", e.Name())
		}
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

// --- ownership-marker + path-guard tests (#96) ---

func writeFile(t *testing.T, p, content string) {
	t.Helper()
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("could not write %s: %v", p, err)
	}
}

func markerExists(dir string) bool {
	_, err := os.Stat(filepath.Join(dir, managedMarker))
	return err == nil
}

// (a) An unmarked, non-empty directory containing a disallowed file is refused, and its
// contents are left untouched.
func TestClearRefusesUnmarkedForeignContent(t *testing.T) {
	dir := t.TempDir()
	victim := filepath.Join(dir, "important.txt")
	writeFile(t, victim, "the user's real data")

	d := NewDirOps(dir)
	if err := d.Clear(); err == nil {
		t.Fatalf("expected Clear to refuse an unmarked directory with foreign content, got nil")
	}

	if _, err := os.Stat(victim); err != nil {
		t.Fatalf("Clear deleted the user's file despite refusing: %v", err)
	}
}

// (b) A directory that already carries the ownership marker clears successfully, and the
// marker is present afterwards.
func TestClearAllowsMarkedDirectory(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, managedMarker), "managed")
	// Even foreign content is fine once the marker proves we own the directory.
	writeFile(t, filepath.Join(dir, "leftover.bin"), "stale")

	d := NewDirOps(dir)
	if err := d.Clear(); err != nil {
		t.Fatalf("expected Clear to succeed on a marked directory, got %v", err)
	}

	if _, err := os.Stat(filepath.Join(dir, "leftover.bin")); !os.IsNotExist(err) {
		t.Fatalf("expected old content to be removed, stat err = %v", err)
	}
	if !markerExists(dir) {
		t.Fatalf("expected marker to be re-written after Clear")
	}
}

// (c) An empty directory, and an absent directory, both clear and gain a marker.
func TestClearEmptyAndAbsentGetMarker(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		dir := t.TempDir()
		d := NewDirOps(dir)
		if err := d.Clear(); err != nil {
			t.Fatalf("expected Clear to succeed on empty directory, got %v", err)
		}
		if !markerExists(dir) {
			t.Fatalf("expected marker after clearing empty directory")
		}
	})

	t.Run("absent", func(t *testing.T) {
		dir := filepath.Join(t.TempDir(), "does-not-exist-yet")
		d := NewDirOps(dir)
		if err := d.Clear(); err != nil {
			t.Fatalf("expected Clear to succeed on absent directory, got %v", err)
		}
		if !markerExists(dir) {
			t.Fatalf("expected marker after creating absent directory")
		}
	})
}

// Backward compatibility: an unmarked directory whose files all pass the legacy
// extension allowlist still clears (and gains a marker going forward).
func TestClearAllowsLegacyRecognizedContent(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "obj.json"), "{}")
	writeFile(t, filepath.Join(dir, "script.ttslua"), "-- lua")

	d := NewDirOps(dir)
	if err := d.Clear(); err != nil {
		t.Fatalf("expected Clear to succeed on legacy recognized content, got %v", err)
	}
	if !markerExists(dir) {
		t.Fatalf("expected marker after clearing legacy directory")
	}
}

// (d) The path guard refuses the filesystem root and the home directory outright.
func TestPathGuardRefusesDangerousTargets(t *testing.T) {
	if err := pathGuard(string(filepath.Separator)); err == nil {
		t.Fatalf("expected path guard to refuse filesystem root")
	}

	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		t.Skip("no home directory available on this platform")
	}
	if err := pathGuard(home); err == nil {
		t.Fatalf("expected path guard to refuse home directory %q", home)
	}

	// Clear must also refuse the home directory, not just the raw guard.
	d := NewDirOps(home)
	if err := d.Clear(); err == nil {
		t.Fatalf("expected Clear to refuse home directory %q", home)
	}
}

// The path guard refuses suspiciously shallow single-segment paths.
func TestPathGuardRefusesShallowPaths(t *testing.T) {
	if err := pathGuard(string(filepath.Separator) + "objects"); err == nil {
		t.Fatalf("expected path guard to refuse a single-segment path")
	}
}
