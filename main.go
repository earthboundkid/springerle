package main

import (
	"os"

	"github.com/carlmjohnson/exitcode"
	"github.com/carlmjohnson/springerle/txtartmpl"
)

func main() {
	exitcode.Exit(txtartmpl.CLI(os.Args[1:]))
}
