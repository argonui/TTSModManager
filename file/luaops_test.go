package file

import (
	"fmt"
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
