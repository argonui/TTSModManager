// Command moddiff compares two Tabletop Simulator mod files and prints the
// meaningful differences between them, ignoring values that are expected to
// vary between otherwise-equivalent savegames (timestamps and numeric jitter).
//
// It exits 0 when the two mods are equivalent and non-zero when they differ or
// when either file cannot be read.
//
// Usage:
//
//	moddiff -a path/to/first.json -b path/to/second.json
package main

import (
	file "ModCreator/file"
	"flag"
	"fmt"
	"os"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

var (
	modfilea = flag.String("a", "", "path to the first mod file to compare")
	modfileb = flag.String("b", "", "path to the second mod file to compare")
)

// differ accumulates the human-readable differences discovered while comparing
// two mods. It replaces the *testing.T that this logic used to lean on when it
// lived in compare_test.go.
type differ struct {
	diffs []string
}

func (d *differ) reportf(format string, args ...interface{}) {
	d.diffs = append(d.diffs, fmt.Sprintf(format, args...))
}

// ignoreUnpredictable reports whether a map entry should be skipped when
// diffing. Numeric values and the Date/EpochTime fields drift between
// otherwise-equivalent savegames, so comparing them only produces noise.
func ignoreUnpredictable(k string, v interface{}) bool {
	// Date and EpochTime are regenerated on every build by design.
	if k == "Date" || k == "EpochTime" {
		return true
	}

	return false
}

// approxFloats compares float64 values with a small absolute tolerance
// consistent with number smoothing (positions 3dp, scale 2dp, colors 5dp),
// rather than ignoring them outright.
var approxFloats = cmpopts.EquateApprox(0, 1e-4)

func compareDelta(d *differ, filea, fileb string) error {
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
	if err := compareObjArrays(d, asubOs, bsubOs); err != nil {
		d.reportf("compareObjs(<>) : %v", err)
	}

	delete(a, osKey)
	delete(b, osKey)

	if diff := cmp.Diff(a, b, cmpopts.IgnoreMapEntries(ignoreUnpredictable), approxFloats); diff != "" {
		d.reportf("want != got:\n%v\n", diff)
	}
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

func compareObjArrays(d *differ, a, b []map[string]interface{}) error {
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
		if err := compareObjs(d, k, av, bv); err != nil {
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

func compareObjs(d *differ, guid string, a, b map[string]interface{}) error {
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

		if err := compareObjArrays(d, aArr, bArr); err != nil {
			return fmt.Errorf("subObjects of %s[ContainedObjects] have diff: %v", guid, err)
		}

		delete(a, subKey)
		delete(b, subKey)
	} else if !aok && !bok {
		// ignore, neither object has sub objects
	} else {
		return fmt.Errorf("in obj %s, one has sub-objects, the other does not", guid)
	}

	if diff := cmp.Diff(a, b, cmpopts.IgnoreMapEntries(ignoreUnpredictable), approxFloats); diff != "" {
		d.reportf("want != got:\n%v\n", diff)
	}
	return nil
}

func main() {
	flag.Parse()

	if *modfilea == "" || *modfileb == "" {
		fmt.Fprintln(os.Stderr, "both -a and -b must be set to mod file paths")
		flag.Usage()
		os.Exit(2)
	}

	d := &differ{}
	if err := compareDelta(d, *modfilea, *modfileb); err != nil {
		fmt.Fprintf(os.Stderr, "compareDelta(%s,%s) : %v\n", *modfilea, *modfileb, err)
		os.Exit(1)
	}

	if len(d.diffs) > 0 {
		for _, diff := range d.diffs {
			fmt.Println(diff)
		}
		fmt.Fprintf(os.Stderr, "mods differ: found %d difference(s)\n", len(d.diffs))
		os.Exit(1)
	}

	fmt.Println("no differences")
}
