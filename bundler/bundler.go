package bundler

import (
	"fmt"
	"regexp"
)

const (
	luabundleFile string = "../luabundle/src/metadata/index.ts"
)

// Unbundle takes luacode and strips it down to the root sub function
func Unbundle(rawlua string) (string, error) {
	anyBundle := regexp.MustCompile(`__bundle_register`)
	isbundled := len(anyBundle.FindStringSubmatch(rawlua)) > 0

	root := regexp.MustCompile(`(?s)__bundle_register\("__root", function\(require, _LOADED, __bundle_register, __bundle_modules\)\n\s*(.*?)\n\s*end\)`)
	matches := root.FindStringSubmatch(rawlua)
	if len(matches) <= 1 {
		if isbundled {
			return "", fmt.Errorf("no regexp match")
		}
		return rawlua, nil
	}
	return matches[1], nil

}
