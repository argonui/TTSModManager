package bundler

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

const fullrawlua string = `
-- Bundled by luabundle {"version":"1.6.0"}
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
end)(nil)
__bundle_register("__root", function(require, _LOADED, __bundle_register, __bundle_modules)
require("core/AgendaDeck")
end)
__bundle_register("core/AgendaDeck", function(require, _LOADED, __bundle_register, __bundle_modules)
MIN_VALUE = -99
MAX_VALUE = 999

function onload(saved_data)
    light_mode = false
    val = 0

    if saved_data ~= "" then
        local loaded_data = JSON.decode(saved_data)
        light_mode = loaded_data[1]
        val = loaded_data[2]
    end

    createAll()
end

function updateSave()
    local data_to_save = {light_mode, val}
    saved_data = JSON.encode(data_to_save)
    self.script_state = saved_data
end

function createAll()
    s_color = {0.5, 0.5, 0.5, 95}

    if light_mode then
        f_color = {1,1,1,95}
    else
        f_color = {0,0,0,100}
    end



    self.createButton({
      label=tostring(val),
      click_function="add_subtract",
      function_owner=self,
      position={0,0.05,0},
      height=600,
      width=1000,
      alignment = 3,
      scale={x=1.5, y=1.5, z=1.5},
      font_size=600,
      font_color=f_color,
      color={0,0,0,0}
      })




    if light_mode then
        lightButtonText = "[ Set dark ]"
    else
        lightButtonText = "[ Set light ]"
    end

end

function removeAll()
    self.removeInput(0)
    self.removeInput(1)
    self.removeButton(0)
    self.removeButton(1)
    self.removeButton(2)
end

function reloadAll()
    removeAll()
    createAll()

    updateSave()
end

function swap_fcolor(_obj, _color, alt_click)
    light_mode = not light_mode
    reloadAll()
end

function swap_align(_obj, _color, alt_click)
    center_mode = not center_mode
    reloadAll()
end

function editName(_obj, _string, value)
    self.setName(value)
    setTooltips()
end

function add_subtract(_obj, _color, alt_click)
    mod = alt_click and -1 or 1
    new_value = math.min(math.max(val + mod, MIN_VALUE), MAX_VALUE)
    if val ~= new_value then
        val = new_value
      updateVal()
        updateSave()
    end
end

function updateVal()

    self.editButton({
        index = 0,
        label = tostring(val),

        })
end

function reset_val()
    val = 0
    updateVal()
    updateSave()
end

function setTooltips()
    self.editInput({
        index = 0,
        value = self.getName(),
        tooltip = ttText
        })
    self.editButton({
        index = 0,
        value = tostring(val),
        tooltip = ttText
        })
end

function null()
end

function keepSample(_obj, _string, value)
    reloadAll()
end

end)
return __bundle_require("__root")
`

type fakeLuaReader struct {
	fs map[string]string
}

func (f *fakeLuaReader) EncodeFromFile(s string) (string, error) {
	if _, ok := f.fs[s]; !ok {
		return "", fmt.Errorf("fake file <%s> not found", s)
	}
	return f.fs[s], nil
}

func TestUnbundle(t *testing.T) {
	got, err := Unbundle(fullrawlua)
	if err != nil {
		t.Fatalf("expected no err, got %v", err)
	}
	want := `require("core/AgendaDeck")`
	if want != got {
		t.Errorf("want <%s>, got <%s>\n", want, got)
	}
}

func TestUnmultiline(t *testing.T) {
	raw := `
__bundle_register("__root", function(require, _LOADED, __bundle_register, __bundle_modules)
require("core/AgendaDeck")
var a = '2'
require("core/AgendaDeck")
end)
`
	got, err := Unbundle(raw)
	if err != nil {
		t.Fatalf("expected no err, got %v", err)
	}
	want := `require("core/AgendaDeck")
var a = '2'
require("core/AgendaDeck")`
	if want != got {
		t.Errorf("want <%s>, got <%s>\n", want, got)
	}
}

func TestSmartBundle(t *testing.T) {
	fr := &fakeLuaReader{
		fs: map[string]string{
			"core/AgendaDeck.ttslua": `MIN_VALUE = -99
MAX_VALUE = 999

function onload(saved_data)
    light_mode = false
    val = 0

    if saved_data ~= "" then
        local loaded_data = JSON.decode(saved_data)
        light_mode = loaded_data[1]
        val = loaded_data[2]
    end

    createAll()
end

function updateSave()
    local data_to_save = {light_mode, val}
    saved_data = JSON.encode(data_to_save)
    self.script_state = saved_data
end

function createAll()
    s_color = {0.5, 0.5, 0.5, 95}

    if light_mode then
        f_color = {1,1,1,95}
    else
        f_color = {0,0,0,100}
    end



    self.createButton({
      label=tostring(val),
      click_function="add_subtract",
      function_owner=self,
      position={0,0.05,0},
      height=600,
      width=1000,
      alignment = 3,
      scale={x=1.5, y=1.5, z=1.5},
      font_size=600,
      font_color=f_color,
      color={0,0,0,0}
      })




    if light_mode then
        lightButtonText = "[ Set dark ]"
    else
        lightButtonText = "[ Set light ]"
    end

end

function removeAll()
    self.removeInput(0)
    self.removeInput(1)
    self.removeButton(0)
    self.removeButton(1)
    self.removeButton(2)
end

function reloadAll()
    removeAll()
    createAll()

    updateSave()
end

function swap_fcolor(_obj, _color, alt_click)
    light_mode = not light_mode
    reloadAll()
end

function swap_align(_obj, _color, alt_click)
    center_mode = not center_mode
    reloadAll()
end

function editName(_obj, _string, value)
    self.setName(value)
    setTooltips()
end

function add_subtract(_obj, _color, alt_click)
    mod = alt_click and -1 or 1
    new_value = math.min(math.max(val + mod, MIN_VALUE), MAX_VALUE)
    if val ~= new_value then
        val = new_value
      updateVal()
        updateSave()
    end
end

function updateVal()

    self.editButton({
        index = 0,
        label = tostring(val),

        })
end

function reset_val()
    val = 0
    updateVal()
    updateSave()
end

function setTooltips()
    self.editInput({
        index = 0,
        value = self.getName(),
        tooltip = ttText
        })
    self.editButton({
        index = 0,
        value = tostring(val),
        tooltip = ttText
        })
end

function null()
end

function keepSample(_obj, _string, value)
    reloadAll()
end
`,
		},
	}
	input := `require("core/AgendaDeck")`

	got, err := Bundle(input, fr)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	want := fullrawlua

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("want != got:\n%v\n", diff)
	}
}

func TestFailedUnbundle(t *testing.T) {
	rawlua := `  __bundle_register("core/AgendaDeck", function(require, _LOADED, __bundle_register, __bundle_modules)
	  MIN_VALUE = -99
	  MAX_VALUE = 999
`
	_, err := Unbundle(rawlua)
	if err == nil {
		t.Error("expected err, got no err")
	}
}

func TestNonBundled(t *testing.T) {
	rawlua := `
	  MIN_VALUE = -99
	  MAX_VALUE = 999
`
	got, err := Unbundle(rawlua)
	if err != nil {
		t.Fatalf("expected no err, got %v", err)
	}
	want := rawlua
	if want != got {
		t.Errorf("want <%s>, got <%s>\n", want, got)
	}
}
