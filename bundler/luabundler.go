package bundler

import (
	"ModCreator/file"
	"fmt"
	"regexp"
	"sort"
	"strings"
)

const (
	metaprefix string = `-- Bundled by luabundle {"version":"1.6.0"}
local __bundle_require, __bundle_loaded, __bundle_register, __bundle_modules = (function(superRequire)
	local loadingPlaceholder = {[{}] = true}

	local register
	local modules = {}

	local require
	local loaded = {}

	register = function(name, body)
		if not modules[name] then
			modules[name] = body
		end
	end

	require = function(name)
		local loadedModule = loaded[name]

		if loadedModule then
			if loadedModule == loadingPlaceholder then
				return nil
			end
		else
			if not modules[name] then
				if not superRequire then
					local identifier = type(name) == 'string' and '\"' .. name .. '\"' or tostring(name)
					error('Tried to require ' .. identifier .. ', but no such module has been registered')
				else
					return superRequire(name)
				end
			end

			loaded[name] = loadingPlaceholder
			loadedModule = modules[name](require, loaded, register, modules)
			loaded[name] = loadedModule
		end

		return loadedModule
	end

	return require, loaded, register, modules
end)(nil)`
	metasuffix string = `return __bundle_require("__root")`

	funcprefixReplace string = `SRC_LOCATION`
	funcprefix        string = `__bundle_register("SRC_LOCATION", function(require, _LOADED, __bundle_register, __bundle_modules)`
	funcsuffix        string = `end)`
)

// Rootname is the bundle name for the batch of raw lua
var (
	Rootname string = `__root`
)

// IsBundled keeps regex bundling logic to this file
func IsBundled(rawlua string) bool {
	anyBundle := regexp.MustCompile(`__bundle_register`)
	if len(anyBundle.FindStringSubmatch(rawlua)) > 0 {
		return true
	}
	return false
}

// AnalyzeBundle exists to help test functions disect bundles
func AnalyzeBundle(rawlua string, log func(s string, a ...interface{})) {
	if !IsBundled(rawlua) {
		log("script is not bundled\n")
		return
	}
	results, err := UnbundleAll(rawlua)
	if err != nil {
		log("Couldn't unbundle to analyze: %v", err)
		return
	}
	for k, v := range results {
		log("\tmodule %s: %v\n", k, len(v))
	}
}

// UnbundleAll takes luacode generates all bundlenames and bundles
func UnbundleAll(rawlua string) (map[string]string, error) {
	if !IsBundled(rawlua) {
		return map[string]string{Rootname: rawlua}, nil
	}
	newRootInd := regexp.MustCompile(`__bundle_require\(".*"\)`).FindStringIndex(rawlua)
	if newRootInd != nil {
		Rootname = rawlua[newRootInd[0]+18 : newRootInd[1]-2]
	}
	scripts := map[string]string{}
	r, err := findNextBundledScript(rawlua)
	for r.leftover != "" {
		if err != nil {
			return nil, fmt.Errorf("findNextBundledScript(%s): %v", rawlua, err)
		}
		scripts[r.name] = r.body
		r, err = findNextBundledScript(r.leftover)
	}
	if _, ok := scripts[Rootname]; !ok {
		return nil, fmt.Errorf("Failed to find root bundle")
	}

	return scripts, nil
}

type result struct {
	name, body, leftover string
}

func findNextBundledScript(rawlua string) (result, error) {
	root := regexp.MustCompile(`(?s)__bundle_register\("(.*?)", function\(require, _LOADED, __bundle_register, __bundle_modules\)[\r\n\s]+(.*?)[\r\n ]+end\)[\n\r]+(return __bundle_require\(\"__root\"\)|__bundle_register)+`)
	m := root.FindStringSubmatchIndex(rawlua)
	if m == nil {
		return result{}, nil
	}
	if len(m) != 8 {
		return result{}, fmt.Errorf("Expected 8 indices, got %v", m)
	}
	// first 2 ints are indices of entire match
	// second 2 ints are indices of script name
	// third 2 are script body
	// fourth 2 are suffix needed to match
	return result{
		name:     rawlua[m[2]:m[3]],
		body:     rawlua[m[4]:m[5]],
		leftover: rawlua[m[6]:],
	}, nil
}

// Unbundle extracts the root bundle per
func Unbundle(rawlua string) (string, error) {
	srcmap, err := UnbundleAll(rawlua)
	if err != nil {
		return "", err
	}
	rt, ok := srcmap[Rootname]
	if !ok {
		return "", fmt.Errorf("Rootname not found in unbundled map")
	}
	return rt, nil
}

// Bundle grabs all dependencies and creates a single luascript
func Bundle(rawlua string, l file.TextReader) (string, error) {
	if IsBundled(rawlua) {
		return rawlua, nil
	}
	reqs := map[string]string{
		Rootname: rawlua,
	}
	todo := []string{Rootname}
	for len(todo) > 0 {
		fname := todo[0]
		todo = todo[1:] // pop first element off

		scriptToInvestigate := reqs[fname]
		reqsToLoad, err := getAllReqValues(scriptToInvestigate)
		if err != nil {
			return "", fmt.Errorf("for %s getAllReqValues(%s): %v", fname, scriptToInvestigate, err)
		}
		sort.Slice(reqsToLoad, func(i int, j int) bool {
			return reqsToLoad[i] < reqsToLoad[j]
		})
		for _, r := range reqsToLoad {
			val, err := l.EncodeFromFile(r + ".ttslua")
			if err != nil {
				return "", fmt.Errorf("EncodeFromFile(%s) : %v", r, err)
			}
			reqs[r] = val
		}
		todo = append(todo, reqsToLoad...)
	}
	if len(reqs) == 1 {
		// if there were no requires to load in, no need to bundle
		return rawlua, nil
	}

	bundlestr := metaprefix + "\n"

	sortedReqKeys := []string{}
	for k := range reqs {
		sortedReqKeys = append(sortedReqKeys, k)
	}
	sort.Strings(sortedReqKeys)

	for _, k := range sortedReqKeys {
		v := reqs[k]
		bundlestr += strings.Replace(funcprefix, funcprefixReplace, k, 1) + "\n"
		bundlestr += v + "\n"
		bundlestr += funcsuffix + "\n"
	}

	bundlestr += metasuffix

	return bundlestr, nil
}

func getAllReqValues(lua string) ([]string, error) {
	rsxp := regexp.MustCompile(`require\((\\)?\"[-a-zA-Z0-9/._@]+(\\)?\"\)`)
	reqs := rsxp.FindAllString(lua, -1)

	fnames := []string{}
	for _, req := range reqs {
		filexp := regexp.MustCompile(`require\(\\?"([-a-zA-Z0-9/._@]+)\\?"\)`)
		matches := filexp.FindSubmatch([]byte(req))
		if len(matches) != 2 {
			return nil, fmt.Errorf("regex error parsing requirement (%s)", req)
		}
		f := matches[1]
		fnames = append(fnames, string(f))
	}
	return fnames, nil
}
