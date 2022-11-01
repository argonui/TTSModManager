package objects

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
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

func TestObjPrinting(t *testing.T) {
	for _, tc := range []struct {
		o    *objConfig
		want j
	}{
		{
			o: &objConfig{
				guid: "123456",
				data: j{
					"GUID": "123456",
				},
			},
			want: j{
				"GUID": "123456",
			},
		},
	} {
		l := &fakeLuaReader{
			fs: map[string]string{
				"core/AgendaDeck.ttslua": "var a = 42;",
			},
		}
		got, err := tc.o.print(l)
		if err != nil {
			t.Errorf("printing %v, got %v", tc.o, err)
		}
		if diff := cmp.Diff(tc.want, got); diff != "" {
			t.Errorf("want != got:\n%v\n", diff)
		}
	}
}

func TestName(t *testing.T) {
	for _, tc := range []struct {
		data            j
		guid, subObjDir string
		want            string
	}{
		{
			data: j{
				"Nickname": "Occult Invocation",
				"Name":     "Card",
			},
			guid: "010509",
			want: "OccultInvocation.010509",
		}, {
			data: j{
				"Nickname": "Occult Invocation!!!!~~()()",
				"Name":     "Card",
			},
			guid: "010509",
			want: "OccultInvocation.010509",
		}, {
			data: j{
				"Name": "Card",
			},
			guid: "1233",
			want: "Card.1233",
		}, {
			data: j{
				"Nickname": nil,
				"Name":     "Card",
			},
			guid: "1235",
			want: "Card.1235",
		}, {
			data: j{
				"Nickname": "",
				"Name":     "Card",
			},
			guid: "1234",
			want: "Card.1234",
		}, {
			data: j{
				"Nickname": "",
				"Name":     nil,
			},
			guid: "010509",
			want: "010509",
		}, {
			data: j{
				"Nickname": nil,
				"Name":     nil,
			},
			guid: "010509",
			want: "010509",
		}, {
			data: j{},
			guid: "010509",
			want: "010509",
		}, {
			data: j{
				"Nickname": "",
				"Name":     "",
			},
			guid: "010509",
			want: "010509",
		}, {
			data: j{
				"Nickname": "",
				"Name":     "",
			},
			subObjDir: "010509_2",
			guid:      "010509",
			want:      "010509_2",
		}, {
			data: j{
				"Nickname": "foo",
				"Name":     "",
			},
			subObjDir: "010509_2",
			guid:      "010509",
			want:      "foo.010509_2",
		},
	} {
		o := objConfig{
			data:      tc.data,
			guid:      tc.guid,
			subObjDir: tc.subObjDir,
		}
		got := o.getAGoodFileName()
		if tc.want != got {
			t.Errorf("want <%s> got <%s>", tc.want, got)
		}
	}

}
