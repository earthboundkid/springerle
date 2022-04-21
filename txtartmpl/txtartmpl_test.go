package txtartmpl_test

import (
	"strings"
	"testing"

	"github.com/carlmjohnson/be"
	"github.com/carlmjohnson/exitcode"
	"github.com/carlmjohnson/springerle/txtartmpl"
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
			err := txtartmpl.CLI(strings.Fields(tc.in))
			be.Equal(t, tc.code, exitcode.Get(err))
		})
	}
}
