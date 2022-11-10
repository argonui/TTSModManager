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
	basepaths        []string
	writeBasepath    string
	readFileToBytes  func(string) ([]byte, error)
	writeBytesToFile func(string, []byte) error
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
	return NewLuaOpsMulti([]string{base}, base)
}

// NewLuaOpsMulti allows for luascript to be read from multiple directories
func NewLuaOpsMulti(readDirs []string, writeDir string) *LuaOps {
	return &LuaOps{
		basepaths:     readDirs,
		writeBasepath: writeDir,
		readFileToBytes: func(s string) ([]byte, error) {
			sFile, err := os.Open(s)
			if err != nil {
				return nil, fmt.Errorf("os.Open(%s): %v", s, err)
			}
			defer sFile.Close()

			b, err := ioutil.ReadAll(sFile)
			if err != nil {
				return nil, fmt.Errorf("ReadAll(%s): %v", s, err)
			}
			if l := len(b); l > 0 {
				if b[l-1] == '\n' {
					b = b[0 : l-1]
				}
			}
			return b, nil
		},
		writeBytesToFile: func(p string, b []byte) error {
			b = append(b, '\n')
			return os.WriteFile(p, b, 0644)
		},
	}
}

// EncodeFromFile pulls a file from configs and encodes it as a string.
func (l *LuaOps) EncodeFromFile(filename string) (string, error) {
	for _, base := range l.basepaths {
		p := path.Join(base, filename)
		b, err := l.readFileToBytes(p)
		if err != nil {
			continue
		}
		s := string(b)
		return s, nil
	}
	return "", fmt.Errorf("%s not found among any known paths", filename)
}

// EncodeToFile takes a single string and decodes escape characters; writes it.
func (l *LuaOps) EncodeToFile(script, file string) error {
	p := path.Join(l.writeBasepath, file)
	return l.writeBytesToFile(p, []byte(script))
}
