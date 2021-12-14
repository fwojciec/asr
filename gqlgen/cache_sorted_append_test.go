package gqlgen

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

func TestSortedAppend(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name  string
		data  []string
		value string
		exp   []string
	}{
		{"add to nil", nil, "a", []string{"a"}},
		{"add to empty", []string{}, "a", []string{"a"}},
		{"add to single", []string{"b"}, "a", []string{"a", "b"}},
		{"add to the end", []string{"a"}, "b", []string{"a", "b"}},
		{"add to the middle", []string{"a", "c"}, "b", []string{"a", "b", "c"}},
		{"don't add duplicate", []string{"a", "b"}, "b", []string{"a", "b"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			res := sortedAppend(tc.data, tc.value)
			equals(t, tc.exp, res)
		})
	}
}

// equals fails the test if exp is not equal to act.
func equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}
