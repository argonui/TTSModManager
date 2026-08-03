package file

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

// Set of allowed file extensions for the safety check before clearing the objects folder
var allowedExtensions = map[string]struct{}{
	".json":           {},
	".gmnotes":        {},
	".luascriptstate": {},
	".ttslua":         {},
	".xml":            {},
}

// managedMarker is the hidden ownership marker (sentinel) this tool drops into any
// directory it manages. Its presence is proof that the directory was created by
// TTSModManager, which is what makes deleting the directory's contents safe. See the
// "Clear and its safety check" section of AGENTS.md.
const managedMarker = ".ttsmm-managed"

// DirCreator abstracts folder creation
type DirCreator interface {
	CreateDir(relpath string, suggestion string) (string, error)
	Clear() error
}

// DirExplorer allows files and folders to be enumerated
type DirExplorer interface {
	// ListFilesAndFolders returns files, folders, err with names sharing prefix of relpath
	ListFilesAndFolders(relpath string) ([]string, []string, error)
}

// DirOps abstracts away folder creation and other future folder oprations
type DirOps struct {
	base string
}

// NewDirOps allows for abstraction of creation of a directory operator
func NewDirOps(p string) *DirOps {
	return &DirOps{
		base: p,
	}
}

// CreateDir allows objects to abstract creation of sub directories without knowning the root path of the machine
func (d *DirOps) CreateDir(relpath, suggestion string) (string, error) {
	dirname := suggestion
	err := os.Mkdir(path.Join(d.base, relpath, suggestion), 0755)
	tries := 0
	if os.IsExist(err) {
		return dirname, nil
	}
	for err != nil && tries < 100 {
		log.Printf("error creating %s, trying again\n%v\n", path.Join(d.base, relpath, suggestion), err)
		tries++
		dirname = fmt.Sprintf("%s_%v", suggestion, tries)
		err = os.Mkdir(path.Join(d.base, relpath, dirname), 0755)
	}
	if tries >= 100 {
		return "", fmt.Errorf("could not find sutible name for sub directory based on suggestion %s; %v", suggestion, err)
	}
	return dirname, nil
}

// preClearCheck walks the directory and ensures all files have an allowed extension.
// It returns an error if a file with a disallowed extension is found.
func (d *DirOps) preClearCheck() error {
	walkErr := filepath.Walk(d.base, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		ext := filepath.Ext(info.Name())
		if _, isAllowed := allowedExtensions[ext]; !isAllowed {
			return fmt.Errorf("unsafe file type found: %s", path)
		}
		return nil
	})

	// If the directory doesn't exist, Walk returns an error. We treat this as "safe".
	if os.IsNotExist(walkErr) {
		return nil
	}
	return walkErr
}

// pathGuard refuses to operate on paths that are almost certainly a mistargeted
// --moddir rather than a real objects directory: the filesystem root, the user's home
// directory (or any ancestor of it), or any suspiciously shallow path. Deleting any of
// these would destroy unrelated user data.
func pathGuard(target string) error {
	abs, err := filepath.Abs(target)
	if err != nil {
		return fmt.Errorf("could not resolve absolute path of %s: %w", target, err)
	}
	abs = filepath.Clean(abs)

	// Refuse the filesystem root.
	if abs == filepath.Clean(string(filepath.Separator)) {
		return fmt.Errorf("refusing to clear filesystem root %q", abs)
	}

	// Refuse the home directory and any ancestor of it.
	if home, herr := os.UserHomeDir(); herr == nil && home != "" {
		home = filepath.Clean(home)
		if abs == home {
			return fmt.Errorf("refusing to clear home directory %q", abs)
		}
		// An ancestor of home (e.g. /home, / on some systems) is even more dangerous.
		// abs is a strict ancestor when the path from abs down to home never climbs out.
		if rel, rerr := filepath.Rel(abs, home); rerr == nil && rel != "." && !startsWithParent(rel) {
			return fmt.Errorf("refusing to clear %q, an ancestor of your home directory %q", abs, home)
		}
	}

	// Refuse suspiciously shallow paths: a single segment below the root such as
	// "/objects" is far more likely a typo than a real mod tree.
	trimmed := strings.Trim(abs, string(filepath.Separator))
	if trimmed == "" {
		return fmt.Errorf("refusing to clear filesystem root %q", abs)
	}
	if !strings.ContainsRune(trimmed, filepath.Separator) {
		return fmt.Errorf("refusing to clear suspiciously shallow path %q; pass a --moddir at least two levels deep", abs)
	}

	return nil
}

// startsWithParent reports whether a filepath.Rel result indicates the target is an
// ancestor of home (the relative path from target to home does not need to climb out).
func startsWithParent(rel string) bool {
	return rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator))
}

// isClearable decides whether it is safe to delete the contents of d.base. Deletion is
// allowed when the directory:
//   - does not exist, or
//   - is empty, or
//   - already contains our ownership marker (proof this tool created it), or
//   - passes the legacy extension allowlist (backward compatibility for trees written
//     by versions predating the marker).
//
// Anything else - a non-empty directory of unrecognized files with no marker, such as a
// mistargeted $HOME - is refused.
func (d *DirOps) isClearable() error {
	entries, err := os.ReadDir(d.base)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // absent: nothing to destroy
		}
		return fmt.Errorf("reading directory %s: %w", d.base, err)
	}

	if len(entries) == 0 {
		return nil // empty: nothing to destroy
	}

	// Marker present: this tool created the directory, so clearing is safe.
	for _, e := range entries {
		if !e.IsDir() && e.Name() == managedMarker {
			return nil
		}
	}

	// Backward compatibility: a directory that predates the marker but contains only
	// recognized content is still ours to clear (and gains a marker going forward).
	if err := d.preClearCheck(); err != nil {
		return fmt.Errorf(
			"%s is not empty, has no %s ownership marker, and contains unrecognized content: %w; "+
				"refusing to delete it in case --moddir is pointed at the wrong place. "+
				"If this really is a mod objects directory, remove the offending file or the whole directory by hand and re-run",
			d.base, managedMarker, err)
	}
	return nil
}

// writeMarker drops the ownership marker into d.base so future Clear calls recognize the
// directory as one this tool manages.
func (d *DirOps) writeMarker() error {
	marker := filepath.Join(d.base, managedMarker)
	content := []byte("This directory is managed by TTSModManager. It may be deleted and recreated by reverse mode.\n")
	if err := os.WriteFile(marker, content, 0644); err != nil {
		return fmt.Errorf("error writing ownership marker %s: %w", marker, err)
	}
	return nil
}

// Clear removes all contents from the base directory and recreates it.
func (d *DirOps) Clear() error {
	log.Println("Performing safety check...")
	startTime := time.Now()

	// Path guard: refuse obviously dangerous targets regardless of contents.
	if err := pathGuard(d.base); err != nil {
		return fmt.Errorf("pre-clear safety check failed, operation aborted: %w", err)
	}

	// Ownership guard: refuse to delete a directory this tool does not appear to own.
	if err := d.isClearable(); err != nil {
		return fmt.Errorf("pre-clear safety check failed, operation aborted: %w", err)
	}

	duration := time.Since(startTime)
	log.Printf("Safety check passed in %v. Proceeding with clear.", duration)

	// Remove the directory and all its contents
	if err := os.RemoveAll(d.base); err != nil {
		return fmt.Errorf("error clearing directory %s: %w", d.base, err)
	}

	// Recreate the empty directory
	if err := os.MkdirAll(d.base, 0755); err != nil {
		return fmt.Errorf("error recreating directory %s: %w", d.base, err)
	}

	// Drop the ownership marker so this directory is recognized as ours next time.
	if err := d.writeMarker(); err != nil {
		return err
	}

	log.Printf("Cleared and recreated directory: %s", d.base)
	return nil
}

// ListFilesAndFolders allows for file exploration. returns relateive file or folder names
func (d *DirOps) ListFilesAndFolders(relpath string) ([]string, []string, error) {
	p := filepath.Join(d.base, relpath)
	entries, err := os.ReadDir(p)
	if err != nil {
		return nil, nil, fmt.Errorf("os.ReadDir(%s + %s) : %v", d.base, relpath, err)
	}

	var fnames, folnames []string
	for _, entry := range entries {
		if entry.IsDir() {
			folnames = append(folnames, filepath.Join(relpath, entry.Name()))
		} else {
			// Hide the ownership marker: it is an internal file/ bookkeeping detail and
			// must not be seen as mod content by the rest of the pipeline.
			if entry.Name() == managedMarker {
				continue
			}
			fnames = append(fnames, filepath.Join(relpath, entry.Name()))
		}
	}
	return fnames, folnames, nil
}
