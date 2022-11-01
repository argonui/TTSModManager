package bundler

import (
	"ModCreator/file"
	"fmt"
	"regexp"
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
	rootfuncname      string = `__root`
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

// Bundle grabs all dependencies and creates a single luascript
func Bundle(rawlua string, l file.LuaReader) (string, error) {
	reqs := map[string]string{
		rootfuncname: rawlua,
	}
	todo := []string{rootfuncname}
	for len(todo) > 0 {
		fname := todo[0]
		todo = todo[1:] // pop first element off

		scriptToInvestigate := reqs[fname]
		reqsToLoad, err := getAllReqValues(scriptToInvestigate)
		if err != nil {
			return "", fmt.Errorf("for %s getAllReqValues(%s): %v", fname, scriptToInvestigate, err)
		}
		for _, r := range reqsToLoad {
			val, err := l.EncodeFromFile(r + ".ttslua")
			if err != nil {
				return "", fmt.Errorf("EncodeFromFile(%s) : %v", r, err)
			}
			reqs[r] = val
		}
		todo = append(todo, reqsToLoad...)
	}

	bundlestr := "\n" + metaprefix + "\n"

	for k, v := range reqs {
		bundlestr += strings.Replace(funcprefix, funcprefixReplace, k, 1) + "\n"
		bundlestr += v + "\n"
		bundlestr += funcsuffix + "\n"
	}

	bundlestr += metasuffix + "\n"

	return bundlestr, nil
}

func getAllReqValues(lua string) ([]string, error) {
	rsxp := regexp.MustCompile(`(?m)^require\((\\)?\"[a-zA-Z0-9/]*(\\)?\"\)\s*$`)
	reqs := rsxp.FindAllString(lua, -1)

	fnames := []string{}
	for _, req := range reqs {
		filexp := regexp.MustCompile(`require\(\\?"([a-zA-Z0-9/]*)\\?"\)`)
		matches := filexp.FindSubmatch([]byte(req))
		if len(matches) != 2 {
			return nil, fmt.Errorf("regex error parsing requirement (%s)", req)
		}
		f := matches[1]
		fnames = append(fnames, string(f))
	}
	return fnames, nil
}
