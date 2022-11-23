package mod

import (
	"ModCreator/tests"
	"ModCreator/types"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestReverse(t *testing.T) {
	for _, tc := range []struct {
		name            string
		input           map[string]interface{}
		wantRootConfig  map[string]interface{}
		wantModSettings map[string]types.J
		wantObjs        map[string]types.J
		wantObjTexts    map[string]string
		wantSrcTexts    map[string]string
	}{
		{
			name: "SnapPoints",
			input: map[string]interface{}{
				"SnapPoints": []interface{}{
					map[string]interface{}{
						"Position": map[string]interface{}{
							"x": float64(12.123456),
							"y": float64(22.123456),
							"z": float64(32.123456),
						},
					},
				},
			},
			wantRootConfig: map[string]interface{}{
				"SnapPoints": []interface{}{
					map[string]interface{}{
						"Position": map[string]interface{}{
							"x": float64(12.123),
							"y": float64(22.123),
							"z": float64(32.123),
						},
					},
				},
			},
			wantModSettings: map[string]types.J{},
			wantSrcTexts:    map[string]string{},
			wantObjs:        map[string]types.J{},
			wantObjTexts:    map[string]string{},
		},
		{
			name: "SnapPointsOwnFile",
			input: map[string]interface{}{
				"SnapPoints": []interface{}{
					map[string]interface{}{
						"Position": map[string]interface{}{
							"x": float64(12.123456),
							"y": float64(22.123456),
							"z": float64(32.123456),
						},
					},
					map[string]interface{}{
						"Position": map[string]interface{}{
							"x": float64(12.123456),
							"y": float64(22.123456),
							"z": float64(32.123456),
						},
					},
					map[string]interface{}{
						"Position": map[string]interface{}{
							"x": float64(12.123456),
							"y": float64(22.123456),
							"z": float64(32.123456),
						},
					},
					map[string]interface{}{
						"Position": map[string]interface{}{
							"x": float64(12.123456),
							"y": float64(22.123456),
							"z": float64(32.123456),
						},
					},
					map[string]interface{}{
						"Position": map[string]interface{}{
							"x": float64(12.123456),
							"y": float64(22.123456),
							"z": float64(32.123456),
						},
					},
				},
			},
			wantRootConfig: map[string]interface{}{
				"SnapPoints_path": "SnapPoints.json",
			},
			wantObjs: map[string]types.J{},
			wantModSettings: map[string]types.J{
				"SnapPoints.json": types.J{
					"testarray": []map[string]interface{}{ // implementation detail of fake files
						map[string]interface{}{
							"Position": types.J{
								"x": float64(12.123),
								"y": float64(22.123),
								"z": float64(32.123),
							},
						},
						map[string]interface{}{
							"Position": types.J{
								"x": float64(12.123),
								"y": float64(22.123),
								"z": float64(32.123),
							},
						},
						map[string]interface{}{
							"Position": types.J{
								"x": float64(12.123),
								"y": float64(22.123),
								"z": float64(32.123),
							},
						},
						map[string]interface{}{
							"Position": types.J{
								"x": float64(12.123),
								"y": float64(22.123),
								"z": float64(32.123),
							},
						},
						map[string]interface{}{
							"Position": types.J{
								"x": float64(12.123),
								"y": float64(22.123),
								"z": float64(32.123),
							},
						},
					},
				},
			},
			wantSrcTexts: map[string]string{},
			wantObjTexts: map[string]string{},
		},
		{
			name: "ShortLua",
			input: map[string]interface{}{
				"LuaScript": "var foo = 42",
			},
			wantRootConfig: map[string]interface{}{
				"LuaScript": "var foo = 42",
			},
			wantSrcTexts:    map[string]string{},
			wantModSettings: map[string]types.J{},
			wantObjs:        map[string]types.J{},
			wantObjTexts:    map[string]string{},
		},
		{
			name: "ShortLuaBundled",
			input: map[string]interface{}{
				"LuaScript": "-- Bundled by luabundle {\"version\":\"1.6.0\"}\nlocal __bundle_require, __bundle_loaded, __bundle_register, __bundle_modules = (function(superRequire)\n\tlocal loadingPlaceholder = {[{}] = true}\n\n\tlocal register\n\tlocal modules = {}\n\n\tlocal require\n\tlocal loaded = {}\n\n\tregister = function(name, body)\n\t\tif not modules[name] then\n\t\t\tmodules[name] = body\n\t\tend\n\tend\n\n\trequire = function(name)\n\t\tlocal loadedModule = loaded[name]\n\n\t\tif loadedModule then\n\t\t\tif loadedModule == loadingPlaceholder then\n\t\t\t\treturn nil\n\t\t\tend\n\t\telse\n\t\t\tif not modules[name] then\n\t\t\t\tif not superRequire then\n\t\t\t\t\tlocal identifier = type(name) == 'string' and '\\\"' .. name .. '\\\"' or tostring(name)\n\t\t\t\t\terror('Tried to require ' .. identifier .. ', but no such module has been registered')\n\t\t\t\telse\n\t\t\t\t\treturn superRequire(name)\n\t\t\t\tend\n\t\t\tend\n\n\t\t\tloaded[name] = loadingPlaceholder\n\t\t\tloadedModule = modules[name](require, loaded, register, modules)\n\t\t\tloaded[name] = loadedModule\n\t\tend\n\n\t\treturn loadedModule\n\tend\n\n\treturn require, loaded, register, modules\nend)(nil)\n__bundle_register(\"__root\", function(require, _LOADED, __bundle_register, __bundle_modules)\nrequire(\"playermat/SkillToken\")\nend)\n__bundle_register(\"playermat/SkillToken\", function(require, _LOADED, __bundle_register, __bundle_modules)\nMIN_VALUE = -99\r\nMAX_VALUE = 999\r\n\r\nfunction onload(saved_data)\r\n    light_mode = false\r\n    val = 0\r\n\r\n    if saved_data ~= \"\" then\r\n        local loaded_data = JSON.decode(saved_data)\r\n        light_mode = loaded_data[1]\r\n        val = loaded_data[2]\r\n    end\r\n\r\n    createAll()\r\nend\r\n\r\nfunction updateSave()\r\n    local data_to_save = {light_mode, val}\r\n    saved_data = JSON.encode(data_to_save)\r\n    self.script_state = saved_data\r\nend\r\n\r\nfunction createAll()\r\n    s_color = {0.5, 0.5, 0.5, 95}\r\n\r\n    if light_mode then\r\n        f_color = {1,1,1,95}\r\n    else\r\n        f_color = {0,0,0,100}\r\n    end\r\n\r\n\r\n\r\n    self.createButton({\r\n      label=tostring(val),\r\n      click_function=\"add_subtract\",\r\n      function_owner=self,\r\n      position={0,0.05,0},\r\n      height=600,\r\n      width=1000,\r\n      alignment = 3,\r\n      scale={x=1.5, y=1.5, z=1.5},\r\n      font_size=600,\r\n      font_color=f_color,\r\n      color={0,0,0,0}\r\n      })\r\n\r\n\r\n\r\n\r\n    if light_mode then\r\n        lightButtonText = \"[ Set dark ]\"\r\n    else\r\n        lightButtonText = \"[ Set light ]\"\r\n    end\r\n\r\nend\r\n\r\nfunction removeAll()\r\n    self.removeInput(0)\r\n    self.removeInput(1)\r\n    self.removeButton(0)\r\n    self.removeButton(1)\r\n    self.removeButton(2)\r\nend\r\n\r\nfunction reloadAll()\r\n    removeAll()\r\n    createAll()\r\n\r\n    updateSave()\r\nend\r\n\r\nfunction swap_fcolor(_obj, _color, alt_click)\r\n    light_mode = not light_mode\r\n    reloadAll()\r\nend\r\n\r\nfunction swap_align(_obj, _color, alt_click)\r\n    center_mode = not center_mode\r\n    reloadAll()\r\nend\r\n\r\nfunction editName(_obj, _string, value)\r\n    self.setName(value)\r\n    setTooltips()\r\nend\r\n\r\nfunction add_subtract(_obj, _color, alt_click)\r\n    mod = alt_click and -1 or 1\r\n    new_value = math.min(math.max(val + mod, MIN_VALUE), MAX_VALUE)\r\n    if val ~= new_value then\r\n        val = new_value\r\n      updateVal()\r\n        updateSave()\r\n    end\r\nend\r\n\r\nfunction updateVal()\r\n\r\n    self.editButton({\r\n        index = 0,\r\n        label = tostring(val),\r\n\r\n        })\r\nend\r\n\r\nfunction reset_val()\r\n    val = 0\r\n    updateVal()\r\n    updateSave()\r\nend\r\n\r\nfunction setTooltips()\r\n    self.editInput({\r\n        index = 0,\r\n        value = self.getName(),\r\n        tooltip = ttText\r\n        })\r\n    self.editButton({\r\n        index = 0,\r\n        value = tostring(val),\r\n        tooltip = ttText\r\n        })\r\nend\r\n\r\nfunction null()\r\nend\r\n\r\nfunction keepSample(_obj, _string, value)\r\n    reloadAll()\r\nend\r\n\nend)\nreturn __bundle_require(\"__root\")",
			},
			wantRootConfig: map[string]interface{}{
				"LuaScript": "require(\"playermat/SkillToken\")",
			},
			wantObjs:        map[string]types.J{},
			wantModSettings: map[string]types.J{},
			wantSrcTexts: map[string]string{
				"playermat/SkillToken.ttslua": "MIN_VALUE = -99\r\nMAX_VALUE = 999\r\n\r\nfunction onload(saved_data)\r\n    light_mode = false\r\n    val = 0\r\n\r\n    if saved_data ~= \"\" then\r\n        local loaded_data = JSON.decode(saved_data)\r\n        light_mode = loaded_data[1]\r\n        val = loaded_data[2]\r\n    end\r\n\r\n    createAll()\r\nend\r\n\r\nfunction updateSave()\r\n    local data_to_save = {light_mode, val}\r\n    saved_data = JSON.encode(data_to_save)\r\n    self.script_state = saved_data\r\nend\r\n\r\nfunction createAll()\r\n    s_color = {0.5, 0.5, 0.5, 95}\r\n\r\n    if light_mode then\r\n        f_color = {1,1,1,95}\r\n    else\r\n        f_color = {0,0,0,100}\r\n    end\r\n\r\n\r\n\r\n    self.createButton({\r\n      label=tostring(val),\r\n      click_function=\"add_subtract\",\r\n      function_owner=self,\r\n      position={0,0.05,0},\r\n      height=600,\r\n      width=1000,\r\n      alignment = 3,\r\n      scale={x=1.5, y=1.5, z=1.5},\r\n      font_size=600,\r\n      font_color=f_color,\r\n      color={0,0,0,0}\r\n      })\r\n\r\n\r\n\r\n\r\n    if light_mode then\r\n        lightButtonText = \"[ Set dark ]\"\r\n    else\r\n        lightButtonText = \"[ Set light ]\"\r\n    end\r\n\r\nend\r\n\r\nfunction removeAll()\r\n    self.removeInput(0)\r\n    self.removeInput(1)\r\n    self.removeButton(0)\r\n    self.removeButton(1)\r\n    self.removeButton(2)\r\nend\r\n\r\nfunction reloadAll()\r\n    removeAll()\r\n    createAll()\r\n\r\n    updateSave()\r\nend\r\n\r\nfunction swap_fcolor(_obj, _color, alt_click)\r\n    light_mode = not light_mode\r\n    reloadAll()\r\nend\r\n\r\nfunction swap_align(_obj, _color, alt_click)\r\n    center_mode = not center_mode\r\n    reloadAll()\r\nend\r\n\r\nfunction editName(_obj, _string, value)\r\n    self.setName(value)\r\n    setTooltips()\r\nend\r\n\r\nfunction add_subtract(_obj, _color, alt_click)\r\n    mod = alt_click and -1 or 1\r\n    new_value = math.min(math.max(val + mod, MIN_VALUE), MAX_VALUE)\r\n    if val ~= new_value then\r\n        val = new_value\r\n      updateVal()\r\n        updateSave()\r\n    end\r\nend\r\n\r\nfunction updateVal()\r\n\r\n    self.editButton({\r\n        index = 0,\r\n        label = tostring(val),\r\n\r\n        })\r\nend\r\n\r\nfunction reset_val()\r\n    val = 0\r\n    updateVal()\r\n    updateSave()\r\nend\r\n\r\nfunction setTooltips()\r\n    self.editInput({\r\n        index = 0,\r\n        value = self.getName(),\r\n        tooltip = ttText\r\n        })\r\n    self.editButton({\r\n        index = 0,\r\n        value = tostring(val),\r\n        tooltip = ttText\r\n        })\r\nend\r\n\r\nfunction null()\r\nend\r\n\r\nfunction keepSample(_obj, _string, value)\r\n    reloadAll()\r\nend",
			},
			wantObjTexts: map[string]string{},
		},
		{
			name: "LongLua",
			input: map[string]interface{}{
				"LuaScript": "var foo = 42\nvar foo = 42\nvar foo = 42\nvar foo = 42\nvar foo = 42\nvar foo = 42\nvar foo = 42\nvar foo = 42\n",
			},
			wantRootConfig: map[string]interface{}{
				"LuaScript_path": "LuaScript.ttslua",
			},
			wantObjs:        map[string]types.J{},
			wantSrcTexts:    map[string]string{},
			wantModSettings: map[string]types.J{},
			wantObjTexts: map[string]string{
				"LuaScript.ttslua": "var foo = 42\nvar foo = 42\nvar foo = 42\nvar foo = 42\nvar foo = 42\nvar foo = 42\nvar foo = 42\nvar foo = 42\n",
			},
		},
		{
			name: "LongLuaBundled",
			input: map[string]interface{}{
				"LuaScript": "-- Bundled by luabundle {\"version\":\"1.6.0\"}\nlocal __bundle_require, __bundle_loaded, __bundle_register, __bundle_modules = (function(superRequire)\n\tlocal loadingPlaceholder = {[{}] = true}\n\n\tlocal register\n\tlocal modules = {}\n\n\tlocal require\n\tlocal loaded = {}\n\n\tregister = function(name, body)\n\t\tif not modules[name] then\n\t\t\tmodules[name] = body\n\t\tend\n\tend\n\n\trequire = function(name)\n\t\tlocal loadedModule = loaded[name]\n\n\t\tif loadedModule then\n\t\t\tif loadedModule == loadingPlaceholder then\n\t\t\t\treturn nil\n\t\t\tend\n\t\telse\n\t\t\tif not modules[name] then\n\t\t\t\tif not superRequire then\n\t\t\t\t\tlocal identifier = type(name) == 'string' and '\\\"' .. name .. '\\\"' or tostring(name)\n\t\t\t\t\terror('Tried to require ' .. identifier .. ', but no such module has been registered')\n\t\t\t\telse\n\t\t\t\t\treturn superRequire(name)\n\t\t\t\tend\n\t\t\tend\n\n\t\t\tloaded[name] = loadingPlaceholder\n\t\t\tloadedModule = modules[name](require, loaded, register, modules)\n\t\t\tloaded[name] = loadedModule\n\t\tend\n\n\t\treturn loadedModule\n\tend\n\n\treturn require, loaded, register, modules\nend)(nil)\n__bundle_register(\"__root\", function(require, _LOADED, __bundle_register, __bundle_modules)\nrequire(\"playermat/SkillToken\")\nvar foo = 42\nvar foo = 42\nvar foo = 42\nvar foo = 42\nvar foo = 42\nvar foo = 42\nvar foo = 42\nend)\n__bundle_register(\"playermat/SkillToken\", function(require, _LOADED, __bundle_register, __bundle_modules)\nMIN_VALUE = -99\r\nMAX_VALUE = 999\r\n\r\nfunction onload(saved_data)\r\n    light_mode = false\r\n    val = 0\r\n\r\n    if saved_data ~= \"\" then\r\n        local loaded_data = JSON.decode(saved_data)\r\n        light_mode = loaded_data[1]\r\n        val = loaded_data[2]\r\n    end\r\n\r\n    createAll()\r\nend\r\n\r\nfunction updateSave()\r\n    local data_to_save = {light_mode, val}\r\n    saved_data = JSON.encode(data_to_save)\r\n    self.script_state = saved_data\r\nend\r\n\r\nfunction createAll()\r\n    s_color = {0.5, 0.5, 0.5, 95}\r\n\r\n    if light_mode then\r\n        f_color = {1,1,1,95}\r\n    else\r\n        f_color = {0,0,0,100}\r\n    end\r\n\r\n\r\n\r\n    self.createButton({\r\n      label=tostring(val),\r\n      click_function=\"add_subtract\",\r\n      function_owner=self,\r\n      position={0,0.05,0},\r\n      height=600,\r\n      width=1000,\r\n      alignment = 3,\r\n      scale={x=1.5, y=1.5, z=1.5},\r\n      font_size=600,\r\n      font_color=f_color,\r\n      color={0,0,0,0}\r\n      })\r\n\r\n\r\n\r\n\r\n    if light_mode then\r\n        lightButtonText = \"[ Set dark ]\"\r\n    else\r\n        lightButtonText = \"[ Set light ]\"\r\n    end\r\n\r\nend\r\n\r\nfunction removeAll()\r\n    self.removeInput(0)\r\n    self.removeInput(1)\r\n    self.removeButton(0)\r\n    self.removeButton(1)\r\n    self.removeButton(2)\r\nend\r\n\r\nfunction reloadAll()\r\n    removeAll()\r\n    createAll()\r\n\r\n    updateSave()\r\nend\r\n\r\nfunction swap_fcolor(_obj, _color, alt_click)\r\n    light_mode = not light_mode\r\n    reloadAll()\r\nend\r\n\r\nfunction swap_align(_obj, _color, alt_click)\r\n    center_mode = not center_mode\r\n    reloadAll()\r\nend\r\n\r\nfunction editName(_obj, _string, value)\r\n    self.setName(value)\r\n    setTooltips()\r\nend\r\n\r\nfunction add_subtract(_obj, _color, alt_click)\r\n    mod = alt_click and -1 or 1\r\n    new_value = math.min(math.max(val + mod, MIN_VALUE), MAX_VALUE)\r\n    if val ~= new_value then\r\n        val = new_value\r\n      updateVal()\r\n        updateSave()\r\n    end\r\nend\r\n\r\nfunction updateVal()\r\n\r\n    self.editButton({\r\n        index = 0,\r\n        label = tostring(val),\r\n\r\n        })\r\nend\r\n\r\nfunction reset_val()\r\n    val = 0\r\n    updateVal()\r\n    updateSave()\r\nend\r\n\r\nfunction setTooltips()\r\n    self.editInput({\r\n        index = 0,\r\n        value = self.getName(),\r\n        tooltip = ttText\r\n        })\r\n    self.editButton({\r\n        index = 0,\r\n        value = tostring(val),\r\n        tooltip = ttText\r\n        })\r\nend\r\n\r\nfunction null()\r\nend\r\n\r\nfunction keepSample(_obj, _string, value)\r\n    reloadAll()\r\nend\r\n\nend)\nreturn __bundle_require(\"__root\")",
			},
			wantRootConfig: map[string]interface{}{
				"LuaScript_path": "LuaScript.ttslua",
			},
			wantSrcTexts: map[string]string{
				"playermat/SkillToken.ttslua": "MIN_VALUE = -99\r\nMAX_VALUE = 999\r\n\r\nfunction onload(saved_data)\r\n    light_mode = false\r\n    val = 0\r\n\r\n    if saved_data ~= \"\" then\r\n        local loaded_data = JSON.decode(saved_data)\r\n        light_mode = loaded_data[1]\r\n        val = loaded_data[2]\r\n    end\r\n\r\n    createAll()\r\nend\r\n\r\nfunction updateSave()\r\n    local data_to_save = {light_mode, val}\r\n    saved_data = JSON.encode(data_to_save)\r\n    self.script_state = saved_data\r\nend\r\n\r\nfunction createAll()\r\n    s_color = {0.5, 0.5, 0.5, 95}\r\n\r\n    if light_mode then\r\n        f_color = {1,1,1,95}\r\n    else\r\n        f_color = {0,0,0,100}\r\n    end\r\n\r\n\r\n\r\n    self.createButton({\r\n      label=tostring(val),\r\n      click_function=\"add_subtract\",\r\n      function_owner=self,\r\n      position={0,0.05,0},\r\n      height=600,\r\n      width=1000,\r\n      alignment = 3,\r\n      scale={x=1.5, y=1.5, z=1.5},\r\n      font_size=600,\r\n      font_color=f_color,\r\n      color={0,0,0,0}\r\n      })\r\n\r\n\r\n\r\n\r\n    if light_mode then\r\n        lightButtonText = \"[ Set dark ]\"\r\n    else\r\n        lightButtonText = \"[ Set light ]\"\r\n    end\r\n\r\nend\r\n\r\nfunction removeAll()\r\n    self.removeInput(0)\r\n    self.removeInput(1)\r\n    self.removeButton(0)\r\n    self.removeButton(1)\r\n    self.removeButton(2)\r\nend\r\n\r\nfunction reloadAll()\r\n    removeAll()\r\n    createAll()\r\n\r\n    updateSave()\r\nend\r\n\r\nfunction swap_fcolor(_obj, _color, alt_click)\r\n    light_mode = not light_mode\r\n    reloadAll()\r\nend\r\n\r\nfunction swap_align(_obj, _color, alt_click)\r\n    center_mode = not center_mode\r\n    reloadAll()\r\nend\r\n\r\nfunction editName(_obj, _string, value)\r\n    self.setName(value)\r\n    setTooltips()\r\nend\r\n\r\nfunction add_subtract(_obj, _color, alt_click)\r\n    mod = alt_click and -1 or 1\r\n    new_value = math.min(math.max(val + mod, MIN_VALUE), MAX_VALUE)\r\n    if val ~= new_value then\r\n        val = new_value\r\n      updateVal()\r\n        updateSave()\r\n    end\r\nend\r\n\r\nfunction updateVal()\r\n\r\n    self.editButton({\r\n        index = 0,\r\n        label = tostring(val),\r\n\r\n        })\r\nend\r\n\r\nfunction reset_val()\r\n    val = 0\r\n    updateVal()\r\n    updateSave()\r\nend\r\n\r\nfunction setTooltips()\r\n    self.editInput({\r\n        index = 0,\r\n        value = self.getName(),\r\n        tooltip = ttText\r\n        })\r\n    self.editButton({\r\n        index = 0,\r\n        value = tostring(val),\r\n        tooltip = ttText\r\n        })\r\nend\r\n\r\nfunction null()\r\nend\r\n\r\nfunction keepSample(_obj, _string, value)\r\n    reloadAll()\r\nend",
			},
			wantObjs:        map[string]types.J{},
			wantModSettings: map[string]types.J{},
			wantObjTexts: map[string]string{
				"LuaScript.ttslua": "require(\"playermat/SkillToken\")\nvar foo = 42\nvar foo = 42\nvar foo = 42\nvar foo = 42\nvar foo = 42\nvar foo = 42\nvar foo = 42",
			},
		},
		{
			name: "Date",
			input: map[string]interface{}{
				"Date":      "11/4/2022 8:34:56 PM",
				"EpochTime": 1667619296,
			},
			wantRootConfig:  map[string]interface{}{},
			wantObjs:        map[string]types.J{},
			wantModSettings: map[string]types.J{},
			wantSrcTexts:    map[string]string{},
			wantObjTexts:    map[string]string{},
		},
		{
			name: "XML in sub objects",
			input: map[string]interface{}{
				"SaveName": "cool mod",
				"ObjectStates": []map[string]interface{}{
					{
						"GUID": "123",
						"XmlUI": `<!-- include foo/bar -->
<Button id="bar"/>
<Button id="bar2"/>
<!-- include foo/bar -->
`,
					},
				},
			},
			wantRootConfig: map[string]interface{}{
				"SaveName":           "cool mod",
				"ObjectStates_order": []interface{}{"123"},
			},
			wantModSettings: map[string]types.J{},
			wantObjs: map[string]types.J{
				"123.json": map[string]interface{}{
					"GUID":  "123",
					"XmlUI": "<Include src=\"foo/bar\"/>\n",
				},
			},
			wantObjTexts: map[string]string{},
			wantSrcTexts: map[string]string{},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			finalOutput := tests.NewFF()
			modsettings := tests.NewFF()
			objsAndLua := tests.NewFF()
			srcTexts := tests.NewFF()
			r := Reverser{
				ModSettingsWriter: modsettings,
				LuaWriter:         objsAndLua,
				LuaSrcWriter:      srcTexts,
				XMLWriter:         srcTexts,
				ObjWriter:         objsAndLua,
				ObjDirCreeator:    objsAndLua,
				RootWrite:         finalOutput,
			}
			err := r.Write(tc.input)
			if err != nil {
				t.Fatalf("Error reversing : %v", err)
			}
			got, err := finalOutput.ReadObj("config.json")
			if err != nil {
				t.Fatalf("Error reading final config.json : %v", err)
			}
			if diff := cmp.Diff(tc.wantRootConfig, got); diff != "" {
				t.Errorf("want != got:\n%v\n", diff)
			}
			if diff := cmp.Diff(tc.wantModSettings, modsettings.Data); diff != "" {
				t.Errorf("want != got:\n%v\n", diff)
			}
			if diff := cmp.Diff(tc.wantObjTexts, objsAndLua.Fs); diff != "" {
				t.Errorf("want != got:\n%v\n", diff)
			}
			if diff := cmp.Diff(tc.wantSrcTexts, srcTexts.Fs); diff != "" {
				t.Errorf("want != got:\n%v\n", diff)
			}
			if diff := cmp.Diff(tc.wantObjs, objsAndLua.Data); diff != "" {
				t.Errorf("want != got:\n%v\n", diff)
			}
		})
	}

}
