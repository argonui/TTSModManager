package bundler

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

const (
	fullrawlua string = `-- Bundled by luabundle {"version":"1.6.0"}
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
return __bundle_require("__root")`
	fullrawlua2 string = "-- Bundled by luabundle {\"version\":\"1.6.0\"}\r\nlocal __bundle_require, __bundle_loaded, __bundle_register, __bundle_modules = (function(superRequire)\r\n\tlocal loadingPlaceholder = {[{}] = true}\r\n\r\n\tlocal register\r\n\tlocal modules = {}\r\n\r\n\tlocal require\r\n\tlocal loaded = {}\r\n\r\n\tregister = function(name, body)\r\n\t\tif not modules[name] then\r\n\t\t\tmodules[name] = body\r\n\t\tend\r\n\tend\r\n\r\n\trequire = function(name)\r\n\t\tlocal loadedModule = loaded[name]\r\n\r\n\t\tif loadedModule then\r\n\t\t\tif loadedModule == loadingPlaceholder then\r\n\t\t\t\treturn nil\r\n\t\t\tend\r\n\t\telse\r\n\t\t\tif not modules[name] then\r\n\t\t\t\tif not superRequire then\r\n\t\t\t\t\tlocal identifier = type(name) == 'string' and '\\\"' .. name .. '\\\"' or tostring(name)\r\n\t\t\t\t\terror('Tried to require ' .. identifier .. ', but no such module has been registered')\r\n\t\t\t\telse\r\n\t\t\t\t\treturn superRequire(name)\r\n\t\t\t\tend\r\n\t\t\tend\r\n\r\n\t\t\tloaded[name] = loadingPlaceholder\r\n\t\t\tloadedModule = modules[name](require, loaded, register, modules)\r\n\t\t\tloaded[name] = loadedModule\r\n\t\tend\r\n\r\n\t\treturn loadedModule\r\n\tend\r\n\r\n\treturn require, loaded, register, modules\r\nend)(nil)\r\n__bundle_register(\"__root\", function(require, _LOADED, __bundle_register, __bundle_modules)\r\nrequire(\"core/DataHelper\")\r\nend)\r\n__bundle_register(\"core/DataHelper\", function(require, _LOADED, __bundle_register, __bundle_modules)\r\n-- set true to enable debug logging\r\nDEBUG = false\r\n\r\nfunction log(message)\r\n  if DEBUG then\r\n  print(message)\r\n  end\r\nend\r\n\r\n--[[\r\nKnown locations and clues. We check this to determine if we should\r\natttempt to spawn clues, first we look for <LOCATION_NAME>_<GUID> and if\r\nwe find nothing we look for <LOCATION_NAME>\r\nformat is [location_guid -> clueCount]\r\n]]\r\nLOCATIONS_DATA_JSON = [[\r\n{\r\n  \"Study\": {\"type\": \"perPlayer\", \"value\": 2, \"clueSide\": \"back\"},\r\n  \"Study_670914\": {\"type\": \"perPlayer\", \"value\": 1, \"clueSide\": \"back\"},\r\n  \"Attic_377b20\": {\"type\": \"perPlayer\", \"value\": 1, \"clueSide\": \"back\"},\r\n  \"Attic\": {\"type\": \"perPlayer\", \"value\": 2, \"clueSide\": \"back\"},\r\n  \"Cellar_5d3bcc\": {\"type\": \"perPlayer\", \"value\": 1, \"clueSide\": \"back\"},\r\n  \"Cellar\": {\"type\": \"perPlayer\", \"value\": 2, \"clueSide\": \"back\"},\r\n  \"Bathroom\": {\"type\": \"perPlayer\", \"value\": 1, \"clueSide\": \"back\"},\r\n  \"Bedroom\": {\"type\": \"perPlayer\", \"value\": 1, \"clueSide\": \"back\"},\r\n  \"Far Above Your House\": {\"type\": \"perPlayer\", \"value\": 1, \"clueSide\": \"back\"},\r\n  \"Deep Below Your House\": {\"type\": \"perPlayer\", \"value\": 1, \"clueSide\": \"back\"},\r\n\r\n  \"Northside_86faac\": {\"type\": \"perPlayer\", \"value\": 1, \"clueSide\": \"back\"},\r\n  \"Northside\": {\"type\" : \"perPlayer\", \"value\": 2, \"clueSide\": \"back\"},\r\n  \"Graveyard\": {\"type\": \"perPlayer\", \"value\": 2, \"clueSide\": \"back\"},\r\n  \"Miskatonic University_cedb0a\": {\"type\": \"perPlayer\", \"value\": 1, \"clueSide\": \"back\"},\r\n  \"Miskatonic University\": {\"type\": \"perPlayer\", \"value\": 2, \"clueSide\": \"back\"},\r\n  \"Downtown_1aa7cb\": {\"type\": \"perPlayer\", \"value\": 2, \"clueSide\": \"back\"},\r\n  \"Downtown\": {\"type\": \"perPlayer\", \"value\": 1, \"clueSide\": \"back\"},\r\n  \"St. Mary's Hospital\": {\"type\": \"perPlayer\", \"value\": 1, \"clueSide\": \"back\"},\r\n  \"Easttown_88245c\": {\"type\": \"perPlayer\", \"value\": 2, \"clueSide\": \"back\"},\r\n  \"Easttown\": {\"type\": \"perPlayer\", \"value\": 1, \"clueSide\": \"back\"},\r\n  \"Southside\": {\"type\": \"perPlayer\", \"value\": 1, \"clueSide\": \"back\"},\r\n  \"Rivertown\": {\"type\": \"perPlayer\", \"value\": 1, \"clueSide\": \"back\"},\r\n  \"Your House_377b20\": {\"type\": \"perPlayer\", \"value\": 1, \"clueSide\": \"back\"},\r\n  \"Your House_b28633\": {\"type\": \"perPlayer\", \"value\": 1, \"clueSide\": \"back\"},\r\n\r\n  \"Ritual Site\": {\"type\": \"perPlayer\", \"value\": 2, \"clueSide\": \"back\"},\r\n  \"Arkham Woods_e8e04b\": {\"type\": \"perPlayer\", \"value\": 0, \"clueSide\": \"back\"},\r\n  \"Arkham Woods\": {\"type\": \"perPlayer\", \"value\": 1, \"clueSide\": \"back\"},\r\n\r\n  \"New Orleans_5ab18a\": {\"type\": \"perPlayer\", \"value\": 0, \"clueSide\": \"back\"},\r\n  \"New Orleans\": {\"type\": \"perPlayer\", \"value\": 1, \"clueSide\": \"back\"},\r\n  \"Riverside_ab9d69\": {\"type\": \"perPlayer\", \"value\": 0, \"clueSide\": \"back\"},\r\n  \"Riverside\": {\"type\": \"perPlayer\", \"value\": 1, \"clueSide\": \"back\"},\r\n  \"Wilderness_3c5ea8\": {\"type\": \"perPlayer\", \"value\": 0, \"clueSide\": \"back\"},\r\n  \"Wilderness\": {\"type\": \"perPlayer\", \"value\": 1, \"clueSide\": \"back\"},\r\n  \"Unhallowed Land_552a1d\": {\"type\":', 'p1', '0', '0', '0', 'm1', 'm1', 'm1', 'm2', 'm2', 'skull', 'skull', 'elder', 'red', 'blue' } },\r\n  normal = { token = { 'p1', '0', '0', 'm1', 'm1', 'm1', 'm2', 'm2', 'm3', 'm4', 'skull', 'skull', 'elder', 'red', 'blue' } },\r\n  hard = { token = { 'p1', '0', 'm1', 'm1', 'm2', 'm2', 'm3', 'm4', 'm5', 'm6', 'skull', 'skull', 'elder', 'red', 'blue' } },\r\n  expert = { token = { '0', 'm1', 'm1', 'm2', 'm3', 'm4', 'm5', 'm6', 'm7', 'm8', 'skull', 'skull', 'elder', 'red', 'blue' } }\r\n  },\r\n  ['Pokemon'] = {\r\n  easy = { token = { 'p1', 'p1', '0', '0', '0', 'm1', 'm1', 'm2', 'm3', 'skull', 'skull', 'tablet', 'elder', 'red', 'blue' } },\r\n  normal = { token = { 'p1', '0', '0', '0', 'm1', 'm2', 'm2', 'm3', 'm5', 'skull', 'skull', 'tablet', 'elder', 'red', 'blue' } },\r\n  hard = { token = { 'p1', '0', '0', 'm1', 'm2', 'm3', 'm3', 'm4', 'm6', 'skull', 'skull', 'tablet', 'elder', 'red', 'blue' } },\r\n  expert = { token = { '0', 'm1', 'm2', 'm2', 'm3', 'm3', 'm4', 'm4', 'm6', 'm8', 'skull', 'skull', 'tablet', 'elder', 'red', 'blue' } }\r\n  },\r\n  ['Safari'] = {\r\n  normal = { token = { 'p1', '0', '0', '0', 'm1', 'm2', 'm2', 'm3', 'm5', 'skull', 'skull', 'cultist', 'tablet', 'elder', 'red', 'blue' } },\r\n  hard = { token = { 'p1', '0', '0', 'm1', 'm2', 'm3', 'm3', 'm4', 'm6', 'skull', 'skull', 'cultist', 'tablet', 'elder', 'red', 'blue' } },\r\n  },\r\n  ['Cerulean'] = {\r\n  normal = { token = { 'p1', '0', '0', '0', 'm1', 'm2', 'm2', 'm3', 'm5', 'skull', 'skull', 'cultist', 'cultist', 'tablet', 'elder', 'red', 'blue' } },\r\n  hard = { token = { 'p1', '0', '0', 'm1', 'm2', 'm3', 'm3', 'm4', 'm6', 'skull', 'skull', 'cultist', 'cultist', 'tablet', 'elder', 'red', 'blue' } },\r\n  },\r\n  ['Erich Zann'] = {\r\n  easy = { token = { 'p1', '0', '0', 'm1', 'm1', 'm2', 'm2', 'm3', 'skull', 'skull', 'cultist', 'tablet', 'elder', 'red', 'blue' } },\r\n  normal = { token = { 'p1', '0', 'm1', 'm1', 'm2', 'm3', 'm3', 'm4', 'skull', 'skull', 'cultist', 'tablet', 'elder', 'red', 'blue' } },\r\n  hard = { token = { '0', 'm1', 'm2', 'm3', 'm4', 'm4', 'm5', 'm6', 'skull', 'skull', 'cultist', 'tablet', 'elder', 'red', 'blue' } },\r\n  expert = { token = { '0', 'm1', 'm2', 'm3', 'm4', 'm5', 'm6', 'm8', 'skull', 'skull', 'cultist', 'tablet', 'elder', 'red', 'blue' } }\r\n  },\r\n  ['Kaimonogatari'] = {\r\n  easy = { token = { 'p1', 'p1', '0', '0', '0', 'm1', 'm1', 'm2', 'm2', 'skull', 'skull', 'cultist', 'red', 'blue' } },\r\n  normal = { token = { 'p1', '0', '0', 'm1', 'm2', 'm2', 'm3', 'm3', 'm4', 'skull', 'skull', 'cultist', 'red', 'blue' } },\r\n  hard = { token = { '0', '0', '0', 'm1', 'm2', 'm2', 'm3', 'm4', 'm4', 'm5', 'skull', 'skull', 'cultist', 'red', 'blue' } },\r\n  expert = { token = { '0', '0', 'm1', 'm1', 'm2', 'm3', 'm4', 'm5', 'm6', 'm6', 'm8', 'skull', 'skull', 'cultist', 'red', 'blue' } }\r\n  },\r\n  ['Sleepy Hollow'] = {\r\n  normal = { token = { 'p1', 'p1', '0', '0', '0', 'm1', 'm1', 'm1', 'm2', 'm2', 'm3', 'm3', 'm4', 'm4', 'm5', 'm6', 'skull', 'skull', 'skull', 'cultist', 'tablet', 'elder', 'red', 'blue' } },\r\n  hard = { token = { 'p1', '0', '0', '0', 'm1', 'm1', 'm1', 'm2', 'm2', 'm3', 'm3', 'm4', 'm4', 'm5', 'm6', 'm8', 'skull', 'skull', 'skull', 'cultist', 'tablet', 'elder', 'red', 'blue' } },\r\n  },\r\n  ['Flesh'] = {\r\n  easy = { token = { 'p1', 'p1', '0', '0', '0', 'm1', 'm1', 'm1', 'm2', 'm3', 'skull', 'skull', 'cultist', 'tablet', 'tablet', 'red', 'blue' } },\r\n  normal = { token = { 'p1', '0', '0', 'm1', 'm1', 'm1', 'm2', 'm2', 'm3', 'm4', 'skull', 'skull', 'cultist', 'tablet', 'tablet', 'red', 'blue' } },\r\n  hard = { token = { '0', '0', 'm1', 'm1', 'm2', 'm3', 'm3', 'm4', 'm4', 'm6', 'skull', 'skull', 'cultist', 'tablet', 'tablet', 'red', 'blue' } },\r\n  },\r\n  ['Dark Matter'] = {\r\n  easy = { token = { 'p1', 'p1', '0', '0', '0', 'm1', 'm1', 'm2', 'm2', 'skull', 'skull', 'cultist', 'cultist', 'red', 'blue' } },\r\n  normal = { token = { 'p1', '0', '0', 'm1', 'm1', 'm1', 'm2', 'm2', 'm3', 'm4', 'skull', 'skull', 'cultist', 'cultist', 'red', 'blue' } },\r\n  hard = { token = { '0', '0', '0', 'm1', 'm1', 'm2', 'm2', 'm3', 'm3', 'm4', 'm5', 'skull', 'skull', 'cultist', 'cultist', 'red', 'blue' } },\r\n  expert = { token = { '0', 'm1', 'm2', 'm2', 'm3', 'm3', 'm4', 'm4', 'm5', 'm6', 'm8', 'skull', 'skull', 'cultist', 'cultist', 'red', 'blue' } }\r\n  },\r\n  ['Dont Starve'] = {\r\n  normal = { token = { 'p1', '0', 'm1', 'm1', 'm2', 'm2', 'm3', 'm3', 'm5', 'skull', 'skull', 'cultist', 'tablet', 'elder', 'red', 'blue' } },\r\n  hard = { token = { '0', 'm1', 'm1', 'm2', 'm2', 'm3', 'm3', 'm5', 'm7', 'skull', 'skull', 'cultist', 'tablet', 'elder', 'red', 'blue' } },\r\n  },\r\n  ['XXXX'] = {\r\n  easy = { token = { 'p1', 'p1', '0', '0', '0', 'm1', 'm1', 'm1', 'm2', 'm2', 'skull', 'skull', 'cultist', 'tablet', 'red', 'blue' } },\r\n  normal = { token = { 'p1', '0', '0', 'm1', 'm1', 'm1', 'm2', 'm2', 'm3', 'm4', 'skull', 'skull', 'cultist', 'tablet', 'red', 'blue' } },\r\n  hard = { token = { '0', '0', '0', 'm1', 'm1', 'm2', 'm2', 'm3', 'm3', 'm4', 'm5', 'skull', 'skull', 'cultist', 'tablet', 'red', 'blue' } },\r\n  expert = { token = { '0', 'm1', 'm1', 'm2', 'm2', 'm3', 'm3', 'm4', 'm4', 'm5', 'm6', 'm8', 'skull', 'skull', 'cultist', 'tablet', 'red', 'blue' } }\r\n  },\r\n\r\n}\r\n\r\nfunction onSave()\r\n  local globalState = JSON.encode(SPAWNED_PLAYER_CARD_GUIDS)\r\n  log('saving global state:  ' .. globalState)\r\n  self.script_state = globalState\r\nend\r\n\r\nfunction onload(save_state)\r\n  if save_state ~= '' then\r\n  log('loading global state:  ' .. save_state)\r\n  SPAWNED_PLAYER_CARD_GUIDS = JSON.decode(save_state)\r\n  else\r\n  SPAWNED_PLAYER_CARD_GUIDS = {}\r\n  end\r\nend\r\n\r\nfunction getSpawnedPlayerCardGuid(params)\r\n  local guid = params[1]\r\n  if SPAWNED_PLAYER_CARD_GUIDS == nil then\r\n  return nil\r\n  end\r\n  return SPAWNED_PLAYER_CARD_GUIDS[guid]\r\nend\r\n\r\nfunction setSpawnedPlayerCardGuid(params)\r\n  local guid = params[1]\r\n  local value = params[2]\r\n  if SPAWNED_PLAYER_CARD_GUIDS ~= nil then\r\n  SPAWNED_PLAYER_CARD_GUIDS[guid] = value\r\n  return true\r\n  end\r\n  return false\r\nend\r\n\r\nfunction checkHiddenCard(name)\r\n  for _, n in ipairs(HIDDEN_CARD_DATA) do\r\n    if name == n then\r\n      return true\r\n    end\r\n  end\r\n  return false\r\nend\r\n\r\nfunction updateHiddenCards(args)\r\n    local custom_data_helper = getObjectFromGUID(args[1])\r\n    local data_hiddenCards = custom_data_helper.getTable(\"HIDDEN_CARD_DATA\")\r\n    for k, v in ipairs(data_hiddenCards) do\r\n        table.insert(HIDDEN_CARD_DATA, v)\r\n    end\r\nend\r\n\r\nend)\r\nreturn __bundle_require(\"__root\")"

	fullrawlua3 string = "-- Bundled by luabundle {\"version\":\"1.6.0\"}\nlocal __bundle_require, __bundle_loaded, __bundle_register, __bundle_modules = (function(superRequire)\n\tlocal loadingPlaceholder = {[{}] = true}\n\n\tlocal register\n\tlocal modules = {}\n\n\tlocal require\n\tlocal loaded = {}\n\n\tregister = function(name, body)\n\t\tif not modules[name] then\n\t\t\tmodules[name] = body\n\t\tend\n\tend\n\n\trequire = function(name)\n\t\tlocal loadedModule = loaded[name]\n\n\t\tif loadedModule then\n\t\t\tif loadedModule == loadingPlaceholder then\n\t\t\t\treturn nil\n\t\t\tend\n\t\telse\n\t\t\tif not modules[name] then\n\t\t\t\tif not superRequire then\n\t\t\t\t\tlocal identifier = type(name) == 'string' and '\\\"' .. name .. '\\\"' or tostring(name)\n\t\t\t\t\terror('Tried to require ' .. identifier .. ', but no such module has been registered')\n\t\t\t\telse\n\t\t\t\t\treturn superRequire(name)\n\t\t\t\tend\n\t\t\tend\n\n\t\t\tloaded[name] = loadingPlaceholder\n\t\t\tloadedModule = modules[name](require, loaded, register, modules)\n\t\t\tloaded[name] = loadedModule\n\t\tend\n\n\t\treturn loadedModule\n\tend\n\n\treturn require, loaded, register, modules\nend)(nil)\n__bundle_register(\"__root\", function(require, _LOADED, __bundle_register, __bundle_modules)\n-- Customizable Cards: Alchemical Distillation\r\n-- by Chr1Z\r\ninformation = {\r\n    version = \"1.7\",\r\n    last_updated = \"12.10.2022\"\r\n}\r\n\r\n\r\n-- Color information for buttons\r\nboxSize = 40\r\n\r\n-- static values\r\nx_1         = -0.933\r\nx_offset    = 0.075\r\ny_visible   = 0.25\r\ny_invisible = -0.5\r\n\r\n-- z-values (lines on the sheet)\r\nposZ = {\r\n    -0.892,\r\n    -0.665,\r\n    -0.430,\r\n    -0.092,\r\n    0.142,\r\n    0.376,\r\n    0.815\r\n}\r\n\r\n-- box setup (amount of boxes per line and amount of marked boxes in that line)\r\nexistingBoxes = { 1, 1, 1, 1, 2, 4, 5 }\r\n\r\ninputBoxes = {}\r\n\r\n-- override 'marked boxes' for debugging ('all' or 'none')\r\nmarkDEBUG = \"\"\r\n\r\n-- save state when going into bags / decks\r\nfunction onDestroy() self.script_state = onSave() end\r\n\r\nfunction onSave() return JSON.encode({ markedBoxes, inputValues }) end\r\n\r\n-- Startup procedure\r\nfunction onLoad(saved_data)\r\n    if saved_data ~= \"\" and markDEBUG == \"\" then\r\n        local loaded_data = JSON.decode(saved_data)\r\n        markedBoxes = loaded_data[1]\r\n        inputValues = loaded_data[2]\r\n    else\r\n        markedBoxes = { 0, 0, 0, 0, 0, 0, 0, 0, 0, 0 }\r\n        inputValues = { \"\", \"\", \"\", \"\", \"\" }\r\n    end\r\n\r\n    makeData()\r\n    createButtonsAndBoxes()\r\n\r\n    self.addContextMenuItem(\"Reset Inputs\", function() updateState() end)\r\n    self.addContextMenuItem(\"Scale: normal\", function() self.setScale({ 1, 1, 1 }) end)\r\n    self.addContextMenuItem(\"Scale: double\", function() self.setScale({ 2, 1, 2 }) end)\r\n    self.addContextMenuItem(\"Scale: triple\", function() self.setScale({ 3, 1, 3 }) end)\r\nend\r\n\r\nfunction updateState(markedBoxesNew)\r\n    if markedBoxesNew then markedBoxes = markedBoxesNew end\r\n    makeData()\r\n    createButtonsAndBoxes()\r\nend\r\n\r\n-- create Data\r\nfunction makeData()\r\n    Data = {}\r\n    Data.checkbox = {}\r\n    Data.textbox = {}\r\n\r\n    -- repeat this for each entry (= line) in existingBoxes\r\n    local totalCount = 0\r\n    for i = 1, #existingBoxes do\r\n        -- repeat this for each checkbox per line\r\n        for j = 1, existingBoxes[i] do\r\n            totalCount                      = totalCount + 1\r\n            Data.checkbox[totalCount]       = {}\r\n            Data.checkbox[totalCount].pos   = {}\r\n            Data.checkbox[totalCount].pos.x = x_1 + j * x_offset\r\n            Data.checkbox[totalCount].pos.z = posZ[i]\r\n            Data.checkbox[totalCount].row   = i\r\n\r\n            if (markDEBUG == \"all\") or (markedBoxes[i] >= j and markDEBUG ~= \"none\") then\r\n                Data.checkbox[totalCount].pos.y = y_visible\r\n                Data.checkbox[totalCount].state = true\r\n            else\r\n                Data.checkbox[totalCount].pos.y = y_invisible\r\n                Data.checkbox[totalCount].state = false\r\n            end\r\n        end\r\n    end\r\n\r\n    -- repeat this for each entry (= line) in inputBoxes\r\n    local totalCount = 0\r\n    for i = 1, #inputBoxes do\r\n        -- repeat this for each textbox per line\r\n        for j = 1, inputBoxes[i] do\r\n            totalCount                     = totalCount + 1\r\n            Data.textbox[totalCount]       = {}\r\n            Data.textbox[totalCount].pos   = inputPos[totalCount]\r\n            Data.textbox[totalCount].width = inputWidth[totalCount]\r\n            Data.textbox[totalCount].value = inputValues[totalCount]\r\n        end\r\n    end\r\nend\r\n\r\n-- checks or unchecks the given box\r\nfunction click_checkbox(tableIndex)\r\n    local row = Data.checkbox[tableIndex].row\r\n\r\n    if Data.checkbox[tableIndex].state == true then\r\n        Data.checkbox[tableIndex].pos.y = y_invisible\r\n        Data.checkbox[tableIndex].state = false\r\n\r\n        markedBoxes[row] = markedBoxes[row] - 1\r\n    else\r\n        Data.checkbox[tableIndex].pos.y = y_visible\r\n        Data.checkbox[tableIndex].state = true\r\n\r\n        markedBoxes[row] = markedBoxes[row] + 1\r\n    end\r\n\r\n    self.editButton({\r\n        index = tableIndex - 1,\r\n        position = Data.checkbox[tableIndex].pos\r\n    })\r\nend\r\n\r\n-- updates saved value for given text box\r\nfunction click_textbox(i, value, selected)\r\n    if selected == false then\r\n        inputValues[i] = value\r\n    end\r\nend\r\n\r\nfunction createButtonsAndBoxes()\r\n    self.clearButtons()\r\n    self.clearInputs()\r\n\r\n    for i, box_data in ipairs(Data.checkbox) do\r\n        local funcName = \"checkbox\" .. i\r\n        local func = function() click_checkbox(i) end\r\n        self.setVar(funcName, func)\r\n\r\n        self.createButton({\r\n            click_function = funcName,\r\n            function_owner = self,\r\n            position       = box_data.pos,\r\n            height         = boxSize,\r\n            width          = boxSize,\r\n            font_size      = box_data.size,\r\n            scale          = { 1, 1, 1 },\r\n            color          = { 0, 0, 0 },\r\n            font_color     = { 0, 0, 0 }\r\n        })\r\n    end\r\n\r\n    for i, box_data in ipairs(Data.textbox) do\r\n        local funcName = \"textbox\" .. i\r\n        local func = function(_, _, val, sel) click_textbox(i, val, sel) end\r\n        self.setVar(funcName, func)\r\n\r\n        self.createInput({\r\n            input_function = funcName,\r\n            function_owner = self,\r\n            label          = \"Click to type\",\r\n            alignment      = 2,\r\n            position       = box_data.pos,\r\n            scale          = buttonScale,\r\n            width          = box_data.width,\r\n            height         = (inputFontsize * 1) + 24,\r\n            font_size      = inputFontsize,\r\n            color          = \"White\",\r\n            font_color     = buttonFontColor,\r\n            value          = box_data.value\r\n        })\r\n    end\r\nend\r\n\nend)\nreturn __bundle_require(\"__root\")"
)

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

func TestUnbundle2(t *testing.T) {
	got, err := Unbundle(fullrawlua2)
	if err != nil {
		t.Fatalf("expected no err, got %v", err)
	}
	want := `require("core/DataHelper")`
	if want != got {
		t.Errorf("want <%s>, got <%s>\n", want, got)
	}
}

func TestUnbundle3(t *testing.T) {
	got, err := Unbundle(fullrawlua3)
	if err != nil {
		t.Fatalf("expected no err, got %v", err)
	}

	want := "-- Customizable Cards: Alchemical Distillation\r\n-- by Chr1Z\r\ninformation = {\r\n    version = \"1.7\",\r\n    last_updated = \"12.10.2022\"\r\n}\r\n\r\n\r\n-- Color information for buttons\r\nboxSize = 40\r\n\r\n-- static values\r\nx_1         = -0.933\r\nx_offset    = 0.075\r\ny_visible   = 0.25\r\ny_invisible = -0.5\r\n\r\n-- z-values (lines on the sheet)\r\nposZ = {\r\n    -0.892,\r\n    -0.665,\r\n    -0.430,\r\n    -0.092,\r\n    0.142,\r\n    0.376,\r\n    0.815\r\n}\r\n\r\n-- box setup (amount of boxes per line and amount of marked boxes in that line)\r\nexistingBoxes = { 1, 1, 1, 1, 2, 4, 5 }\r\n\r\ninputBoxes = {}\r\n\r\n-- override 'marked boxes' for debugging ('all' or 'none')\r\nmarkDEBUG = \"\"\r\n\r\n-- save state when going into bags / decks\r\nfunction onDestroy() self.script_state = onSave() end\r\n\r\nfunction onSave() return JSON.encode({ markedBoxes, inputValues }) end\r\n\r\n-- Startup procedure\r\nfunction onLoad(saved_data)\r\n    if saved_data ~= \"\" and markDEBUG == \"\" then\r\n        local loaded_data = JSON.decode(saved_data)\r\n        markedBoxes = loaded_data[1]\r\n        inputValues = loaded_data[2]\r\n    else\r\n        markedBoxes = { 0, 0, 0, 0, 0, 0, 0, 0, 0, 0 }\r\n        inputValues = { \"\", \"\", \"\", \"\", \"\" }\r\n    end\r\n\r\n    makeData()\r\n    createButtonsAndBoxes()\r\n\r\n    self.addContextMenuItem(\"Reset Inputs\", function() updateState() end)\r\n    self.addContextMenuItem(\"Scale: normal\", function() self.setScale({ 1, 1, 1 }) end)\r\n    self.addContextMenuItem(\"Scale: double\", function() self.setScale({ 2, 1, 2 }) end)\r\n    self.addContextMenuItem(\"Scale: triple\", function() self.setScale({ 3, 1, 3 }) end)\r\nend\r\n\r\nfunction updateState(markedBoxesNew)\r\n    if markedBoxesNew then markedBoxes = markedBoxesNew end\r\n    makeData()\r\n    createButtonsAndBoxes()\r\nend\r\n\r\n-- create Data\r\nfunction makeData()\r\n    Data = {}\r\n    Data.checkbox = {}\r\n    Data.textbox = {}\r\n\r\n    -- repeat this for each entry (= line) in existingBoxes\r\n    local totalCount = 0\r\n    for i = 1, #existingBoxes do\r\n        -- repeat this for each checkbox per line\r\n        for j = 1, existingBoxes[i] do\r\n            totalCount                      = totalCount + 1\r\n            Data.checkbox[totalCount]       = {}\r\n            Data.checkbox[totalCount].pos   = {}\r\n            Data.checkbox[totalCount].pos.x = x_1 + j * x_offset\r\n            Data.checkbox[totalCount].pos.z = posZ[i]\r\n            Data.checkbox[totalCount].row   = i\r\n\r\n            if (markDEBUG == \"all\") or (markedBoxes[i] >= j and markDEBUG ~= \"none\") then\r\n                Data.checkbox[totalCount].pos.y = y_visible\r\n                Data.checkbox[totalCount].state = true\r\n            else\r\n                Data.checkbox[totalCount].pos.y = y_invisible\r\n                Data.checkbox[totalCount].state = false\r\n            end\r\n        end\r\n    end\r\n\r\n    -- repeat this for each entry (= line) in inputBoxes\r\n    local totalCount = 0\r\n    for i = 1, #inputBoxes do\r\n        -- repeat this for each textbox per line\r\n        for j = 1, inputBoxes[i] do\r\n            totalCount                     = totalCount + 1\r\n            Data.textbox[totalCount]       = {}\r\n            Data.textbox[totalCount].pos   = inputPos[totalCount]\r\n            Data.textbox[totalCount].width = inputWidth[totalCount]\r\n            Data.textbox[totalCount].value = inputValues[totalCount]\r\n        end\r\n    end\r\nend\r\n\r\n-- checks or unchecks the given box\r\nfunction click_checkbox(tableIndex)\r\n    local row = Data.checkbox[tableIndex].row\r\n\r\n    if Data.checkbox[tableIndex].state == true then\r\n        Data.checkbox[tableIndex].pos.y = y_invisible\r\n        Data.checkbox[tableIndex].state = false\r\n\r\n        markedBoxes[row] = markedBoxes[row] - 1\r\n    else\r\n        Data.checkbox[tableIndex].pos.y = y_visible\r\n        Data.checkbox[tableIndex].state = true\r\n\r\n        markedBoxes[row] = markedBoxes[row] + 1\r\n    end\r\n\r\n    self.editButton({\r\n        index = tableIndex - 1,\r\n        position = Data.checkbox[tableIndex].pos\r\n    })\r\nend\r\n\r\n-- updates saved value for given text box\r\nfunction click_textbox(i, value, selected)\r\n    if selected == false then\r\n        inputValues[i] = value\r\n    end\r\nend\r\n\r\nfunction createButtonsAndBoxes()\r\n    self.clearButtons()\r\n    self.clearInputs()\r\n\r\n    for i, box_data in ipairs(Data.checkbox) do\r\n        local funcName = \"checkbox\" .. i\r\n        local func = function() click_checkbox(i) end\r\n        self.setVar(funcName, func)\r\n\r\n        self.createButton({\r\n            click_function = funcName,\r\n            function_owner = self,\r\n            position       = box_data.pos,\r\n            height         = boxSize,\r\n            width          = boxSize,\r\n            font_size      = box_data.size,\r\n            scale          = { 1, 1, 1 },\r\n            color          = { 0, 0, 0 },\r\n            font_color     = { 0, 0, 0 }\r\n        })\r\n    end\r\n\r\n    for i, box_data in ipairs(Data.textbox) do\r\n        local funcName = \"textbox\" .. i\r\n        local func = function(_, _, val, sel) click_textbox(i, val, sel) end\r\n        self.setVar(funcName, func)\r\n\r\n        self.createInput({\r\n            input_function = funcName,\r\n            function_owner = self,\r\n            label          = \"Click to type\",\r\n            alignment      = 2,\r\n            position       = box_data.pos,\r\n            scale          = buttonScale,\r\n            width          = box_data.width,\r\n            height         = (inputFontsize * 1) + 24,\r\n            font_size      = inputFontsize,\r\n            color          = \"White\",\r\n            font_color     = buttonFontColor,\r\n            value          = box_data.value\r\n        })\r\n    end\r\nend"

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("want != got:\n%v\n", diff)
	}
}

func TestUnmultiline(t *testing.T) {
	raw := `
__bundle_register("__root", function(require, _LOADED, __bundle_register, __bundle_modules)
require("core/AgendaDeck")
var a = '2'
require("core/AgendaDeck")
end)
__bundle_register("core/AgendaDeck", function(require, _LOADED, __bundle_register, __bundle_modules)
var b = '3'
end)
return __bundle_require("__root")
`
	got, err := UnbundleAll(raw)
	if err != nil {
		t.Fatalf("expected no err, got %v", err)
	}
	want := map[string]string{
		Rootname: `require("core/AgendaDeck")
var a = '2'
require("core/AgendaDeck")`,
		"core/AgendaDeck": "var b = '3'",
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("want != got:\n%v\n", diff)
	}
}

// This is still being non-deterministic
func DisabledTestSmartBundle(t *testing.T) {
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

func TestBundleNoRequires(t *testing.T) {
	fr := &fakeLuaReader{
		fs: map[string]string{},
	}
	input := `var foo = 42`

	got, err := Bundle(input, fr)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	want := `var foo = 42`

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("want != got:\n%v\n", diff)
	}
}

func TestFailedUnbundle(t *testing.T) {
	rawlua := `  __bundle_register("__root", function(require, _LOADED, __bundle_register, __bundle_modules)
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
