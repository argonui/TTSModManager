package bundler

import (
	"fmt"
	"regexp"
	"strings"
)

// BundleXML converts <Include ... >'s into full xml
func BundleXML() error {
	return nil
}

// UnbundleAllXML converts a bundled xml file to mapping of filenames to
// contents
func UnbundleAllXML(rawxml string) (map[string]string, error) {
	type inc struct {
		name  string
		start int
	}
	store := map[string]string{}
	inctag := regexp.MustCompile(`(?m)^(.*?)<!-- include (.*) -->`)
	stack := []inc{}
	xmlarray := strings.Split(rawxml, "\n")

	for ln := 0; ln < len(xmlarray); ln++ {
		val := xmlarray[ln]
		if !inctag.Match([]byte(val)) {
			continue
		}
		submatches := inctag.FindSubmatch([]byte(val))
		key := string(submatches[2])
		if len(stack) > 0 && stack[0].name == key {
			// found end of include, process it
			indent := string(submatches[1])
			indentedvals := xmlarray[stack[0].start+1 : ln]
			store[stack[0].name] = unindentAndJoin(indentedvals, indent)

			insertLine := fmt.Sprintf("%s<Include src=\"%s\">", indent, stack[0].name)
			// remove the include from xmlarray
			tmp := append(xmlarray[:stack[0].start], insertLine)
			xmlarray = append(tmp, xmlarray[ln+1:]...)
			ln = stack[0].start
			stack = stack[1:]
		} else {
			stack = append([]inc{inc{name: key, start: ln}}, stack...)
		}
	}
	if len(stack) != 0 {
		return nil, fmt.Errorf("Bundled xml left after finished reading file: %v", stack)
	}
	store[Rootname] = unindentAndJoin(xmlarray, "")
	return store, nil
}

func unindentAndJoin(raw []string, indent string) string {
	ret := []string{}
	for _, v := range raw {
		ret = append(ret, strings.Replace(v, indent, "", 1))
	}
	return strings.Join(ret, "\n")
}
