package file

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"regexp"
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
	ReplaceRequire(string) (string, error)
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

	return l.ReplaceRequire(s)
}

// ReplaceRequire will examine any luascript and recursively replace
// require statements with their contents
func (l *LuaOps) ReplaceRequire(script string) (string, error) {
	rsxp := regexp.MustCompile(`(?m)^require\((\\)?\"[a-zA-Z0-9/]*(\\)?\"\)\s*$`)

	notReqs := rsxp.Split(script, -1)

	reqs := rsxp.FindAllString(script, -1)

	if len(reqs)+1 != len(notReqs) {
		return "", fmt.Errorf("I've done something wrong with <%s>", script)
	}
	reqsExpanded := []string{}
	for _, req := range reqs {
		log.Printf("matching on <%s>\n", req)

		filexp := regexp.MustCompile(`require\(\\?"([a-zA-Z0-9/]*)\\?"\)`)
		matches := filexp.FindSubmatch([]byte(req))
		f := matches[1]
		exp, err := l.EncodeFromFile(string(f) + expectedSuffix)
		if err != nil {
			return "", fmt.Errorf("expanding require(%s): %v", f, err)
		}
		reqsExpanded = append(reqsExpanded, exp)
	}

	finalStr := ""
	for i := 0; i < len(reqs); i++ {
		finalStr += notReqs[i] + "\n"
		finalStr += reqsExpanded[i] + "\n"
	}
	finalStr += notReqs[len(notReqs)-1]

	return finalStr, nil
}

// EncodeToFile takes a single string and decodes escape characters; writes it.
func (l *LuaOps) EncodeToFile(script, file string) error {
	p := path.Join(l.basepath, file)
	return os.WriteFile(p, []byte(script), 0644)
}
