package txtartmpl_test

import (
	"strings"
	"testing"

	"github.com/carlmjohnson/exitcode"
	"github.com/carlmjohnson/springerle/txtartmpl"
	"github.com/matryer/is"
)

func TestCLI(t *testing.T) {
	cases := map[string]struct {
		in   string
		code int
	}{
		"help":          {"-h", 0},
		"longhelp":      {"--help", 0},
		"bad-context":   {`-context }`, 1},
		"blank-context": {`-context {}`, 0},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			is := is.New(t)
			err := txtartmpl.CLI(strings.Fields(tc.in))
			is.Equal(exitcode.Get(err), tc.code) // exit code must match expectation
		})
	}
}
