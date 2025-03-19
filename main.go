package main

import (
	"os"

	"github.com/carlmjohnson/exitcode"
	"github.com/earthboundkid/springerle/v2/txtartmpl"
)

func main() {
	exitcode.Exit(txtartmpl.CLI(os.Args[1:]))
}
