package main

import (
	file "ModCreator/file"
	"encoding/json"
	"flag"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	altModfile = flag.String("altmodfile", "", "where to read second mod from when comparing.")
)

func compareDelta(t *testing.T, filea, fileb string) error {
	a, err := file.ReadRawFile(filea)
	if err != nil {
		return err
	}
	b, err := file.ReadRawFile(fileb)
	if err != nil {
		return err
	}
	osKey := "ObjectStates"
	arawOS, ok := a[osKey]
	if !ok {
		return fmt.Errorf("Expected key %s in map", osKey)
	}
	asubOs, err := toObjArray(arawOS)
	if err != nil {
		return fmt.Errorf("cannot cast to obj array %v", err)
	}
	brawOS, ok := b[osKey]
	if !ok {
		return fmt.Errorf("Expected key %s in map", osKey)
	}
	bsubOs, err := toObjArray(brawOS)
	if err != nil {
		return fmt.Errorf("cannot cast to obj array %v", err)
	}
	err = compareObjArrays(t, "root", asubOs, bsubOs)
	if err != nil {
		t.Errorf("compareObjs(<>) : %v", err)
	}

	delete(a, osKey)
	delete(b, osKey)

	abytes, err := json.Marshal(a)
	if err != nil {
		return err
	}
	bbytes, err := json.Marshal(b)
	if err != nil {
		return err
	}
	require.JSONEq(t, string(abytes), string(bbytes))
	return nil
}

func toObjArray(i interface{}) ([]map[string]interface{}, error) {
	arr := []map[string]interface{}{}

	ir, ok := i.([]interface{})
	if !ok {
		return nil, fmt.Errorf("Could not cast input as array")
	}
	for _, rawo := range ir {
		o, ok := rawo.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("could not cast element to json obj: \n%v", rawo)
		}
		arr = append(arr, o)
	}
	return arr, nil
}

func compareObjArrays(t *testing.T, guid string, a, b []map[string]interface{}) error {
	if len(a) != len(b) {
		return fmt.Errorf("length mismatch %v vs %v", len(a), len(b))
	}
	am, err := convertToMetaMap(a)
	if err != nil {
		return err
	}
	bm, err := convertToMetaMap(b)
	if err != nil {
		return err
	}
	for k, av := range am {
		bv, ok := bm[k]
		if !ok {
			return fmt.Errorf("b doesn't have GUID %s", k)
		}
		err = compareObjs(t, k, av, bv)
		if err != nil {
			return fmt.Errorf("object %s found diff: %v", k, err)
		}
	}
	return nil
}

func convertToMetaMap(arr []map[string]interface{}) (map[string]map[string]interface{}, error) {
	m := map[string]map[string]interface{}{}
	for _, a := range arr {
		rawG, ok := a["GUID"]
		if !ok {
			return nil, fmt.Errorf("some object doesn't have a Guid: %v", a)
		}
		strG, ok := rawG.(string)
		if !ok {
			return nil, fmt.Errorf("some object doesn't have a string for a Guid: %v", rawG)
		}
		m[strG] = a
	}
	return m, nil
}

func compareObjs(t *testing.T, guid string, a, b map[string]interface{}) error {
	subKey := "ContainedObjects"

	aSub, aok := a[subKey]

	bSub, bok := b[subKey]

	if aok && bok {
		aArr, err := toObjArray(aSub)
		if err != nil {
			return err
		}
		bArr, err := toObjArray(bSub)
		if err != nil {
			return err
		}

		err = compareObjArrays(t, guid, aArr, bArr)
		if err != nil {
			return fmt.Errorf("subObjects of %s[ContainedObjects] have diff: %v", guid, err)
		}

		delete(a, subKey)
		delete(b, subKey)
	} else if !aok && !bok {
		// ignore, neither object has sub objects
	} else {
		return fmt.Errorf("in obj %s, one has sub-objects, the other does not", guid)
	}

	abytes, err := json.Marshal(a)
	if err != nil {
		return err
	}
	bbytes, err := json.Marshal(b)
	if err != nil {
		return err
	}
	require.JSONEq(t, string(abytes), string(bbytes))

	return nil
}

func TestDiff(t *testing.T) {
	if *altModfile == "" || *modfile == "" {
		// if run automatically, ignore this test
		return
	}
	err := compareDelta(t, *modfile, *altModfile)
	if err != nil {
		t.Errorf("compareDelta(%s,%s) : %v", *modfile, *altModfile, err)
	}
}
