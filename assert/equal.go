package assert

import (
	"fmt"
	"testing"
)

func EqualErrorf[T comparable](t *testing.T, expect, actual T, format string, args ...any) {
	if expect != actual {
		t.Errorf("\nexpect:%v\nactual:%v\nmsg:%s\n", expect, actual,
			fmt.Sprintf(format, args...))
	}
}

func EqualFatalf[T comparable](t *testing.T, expect, actual T, format string, args ...any) {
	if expect != actual {
		t.Fatalf("\nexpect:%v\nactual:%v\nmsg:%s\n", expect, actual,
			fmt.Sprintf(format, args...))
	}
}
