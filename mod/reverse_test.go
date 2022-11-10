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
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			finalOutput := tests.NewFF()
			modsettings := tests.NewFF()
			objsAndLua := tests.NewFF()
			r := Reverser{
				ModSettingsWriter: modsettings,
				LuaWriter:         objsAndLua,
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
		})
	}

}
