package bundler

import (
	"ModCreator/file"
	"fmt"
	"regexp"
	"strings"
)

// BundleXML converts <Include ... >'s into full xml
func BundleXML(rawxml string, xr file.TextReader) (string, error) {
	lines := strings.Split(rawxml, "\n")
	final := []string{}
	inctag := regexp.MustCompile(`(?m)^(.*)<Include src="(.*)"/>`)

	for _, v := range lines {
		if !inctag.Match([]byte(v)) {
			final = append(final, v)
			continue
		}
		subs := inctag.FindSubmatch([]byte(v))
		indent := string(subs[1])
		name := string(subs[2])

		replacement := fmt.Sprintf("%s<!-- include %s -->", indent, name)
		final = append(final, replacement)

		fname := name
		if !strings.HasSuffix(name, ".xml") {
			fname = name + ".xml"
		}
		incXMLRaw, err := xr.EncodeFromFile(fname)
		if err != nil {
			return "", fmt.Errorf("EncodeFromFile(%s): %v", fname, err)
		}
		incXMLBundled, err := BundleXML(incXMLRaw, xr)
		if err != nil {
			return "", fmt.Errorf("BundleXML(<%s>): %v", fname, err)
		}
		final = append(final, indentString(incXMLBundled, indent))

		final = append(final, replacement)

	}
	return strings.Join(final, "\n"), nil
}

func indentString(s string, indent string) string {
	lines := strings.Split(s, "\n")
	final := []string{}
	for _, v := range lines {
		final = append(final, fmt.Sprintf("%s%s", indent, v))
	}

	return strings.Join(final, "\n")
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

			insertLine := fmt.Sprintf("%s<Include src=\"%s\"/>", indent, stack[0].name)
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
