package file

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
)

// DirCreator abstracts folder creation
type DirCreator interface {
	CreateDir(relpath string, suggestion string) (string, error)
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
	err := os.Mkdir(path.Join(d.base, relpath, suggestion), 0644)
	tries := 0
	if os.IsExist(err) {
		return dirname, nil
	}
	for err != nil && tries < 100 {
		log.Printf("error creating %s, trying again\n%v\n", path.Join(d.base, relpath, suggestion), err)
		tries++
		dirname = fmt.Sprintf("%s_%v", suggestion, tries)
		err = os.Mkdir(path.Join(d.base, relpath, dirname), 0644)
	}
	if tries >= 100 {
		return "", fmt.Errorf("could not find sutible name for sub directory based on suggestion %s; %v", suggestion, err)
	}
	return dirname, nil
}

// ListFilesAndFolders allows for file exploration. returns relateive file or folder names
func (d *DirOps) ListFilesAndFolders(relpath string) ([]string, []string, error) {
	p := path.Join(d.base, relpath)
	files, err := ioutil.ReadDir(p)
	if err != nil {
		return nil, nil, fmt.Errorf("ioutil.ReadDir(%s) : %v", p, err)
	}
	fnames := make([]string, 0)
	folnames := make([]string, 0)
	for _, file := range files {
		if file.IsDir() {
			folnames = append(folnames, path.Join(relpath, file.Name()))
			continue
		}
		fnames = append(fnames, path.Join(relpath, file.Name()))
	}
	return fnames, folnames, nil
}
