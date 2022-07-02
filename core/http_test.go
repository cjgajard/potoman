package potoman

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestGetHttpVersion(t *testing.T) {
	for _, test := range []struct {
		content []byte
		exp     Attr
		isErr   bool
	}{
		{
			content: nil,
			exp:     Attr{"1", "0"},
			isErr:   false,
		},
		{
			content: []byte("HTTP1.1"),
			exp:     Attr{},
			isErr:   true,
		},
		{
			content: []byte("HTTP/2"),
			exp:     Attr{"2", "0"},
			isErr:   false,
		},
		{
			content: []byte("HTTP/1.1"),
			exp:     Attr{"1", "1"},
			isErr:   false,
		},
		{
			content: []byte("HTTP/1.99"),
			exp:     Attr{},
			isErr:   true,
		},
	} {
		v, err := getHttpVersion(test.content)
		if err != nil {
			if !test.isErr {
				t.Error(err)
			}
			break
		}
		if !cmp.Equal(test.exp, v) {
			t.Error(cmp.Diff(test.exp, v))
		}
	}
}
