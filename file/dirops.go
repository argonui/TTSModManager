package file

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
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

// Clear removes all contents from the base directory and recreates it.
func (d *DirOps) Clear() error {
	log.Println("Performing safety check...")
	startTime := time.Now()

	if err := d.preClearCheck(); err != nil {
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
			fnames = append(fnames, filepath.Join(relpath, entry.Name()))
		}
	}
	return fnames, folnames, nil
}
