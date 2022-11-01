package file

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

const (
	expectedSuffix = ".ttslua"
)

// LuaOps allows for arbitrary reads and writes of luascript
type LuaOps struct {
	basepath        string
	readFileToBytes func(string) ([]byte, error)
}

// LuaReader serves to describe all ways to read luascripts
type LuaReader interface {
	EncodeFromFile(string) (string, error)
}

// LuaWriter serves to describe all ways to write luascripts
type LuaWriter interface {
	EncodeToFile(script, file string) error
}

// NewLuaOps initializes our object on a directory
func NewLuaOps(base string) *LuaOps {
	return &LuaOps{
		basepath: base,
		readFileToBytes: func(s string) ([]byte, error) {
			sFile, err := os.Open(s)
			if err != nil {
				return nil, fmt.Errorf("os.Open(%s): %v", s, err)
			}
			defer sFile.Close()

			return ioutil.ReadAll(sFile)
		},
	}
}

// EncodeFromFile pulls a file from configs and encodes it as a string.
func (l *LuaOps) EncodeFromFile(filename string) (string, error) {
	p := path.Join(l.basepath, filename)
	b, err := l.readFileToBytes(p)
	if err != nil {
		return "", err
	}
	s := string(b)

	return s, nil
}

// EncodeToFile takes a single string and decodes escape characters; writes it.
func (l *LuaOps) EncodeToFile(script, file string) error {
	p := path.Join(l.basepath, file)
	return os.WriteFile(p, []byte(script), 0644)
}
