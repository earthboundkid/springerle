package prompt

import (
	"fmt"
	"io"
	"os"

	"golang.org/x/term"
)

var console *term.Terminal

func Init() (stop func()) {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	console = term.NewTerminal(struct {
		io.Reader
		io.Writer
	}{os.Stdin, os.Stdout}, ": ")

	return func() { term.Restore(int(os.Stdin.Fd()), oldState) }
}

/*
Project name [My Project]:
File name [my-project]:
Class name [MyProject]:
Skip this (y/n) [y]:
*/
func Prompt(msg string, def string) string {
	if def != "" {
		for {
			fmt.Fprintf(console, "%s [%s]", msg, def)
			line, _ := console.ReadLine()
			switch line {
			case "":
				return def
			default:
				return line
			}
		}
	}
	for {
		fmt.Fprint(console, msg)
		line, _ := console.ReadLine()
		switch line {
		case "":
			fmt.Fprintln(console, "Must answer prompt.")
		default:
			return line
		}
	}
}

func YN(msg string, def bool) bool {
	defS := "n"
	if def {
		defS = "y"
	}
	for {
		fmt.Fprintf(console, "%s (y/n) [%s]", msg, defS)
		line, _ := console.ReadLine()
		switch line {
		case "Y", "y", "Yes", "yes", "true":
			return true
		case "N", "n", "No", "no", "false":
			return false
		case "":
			return def
		default:
			fmt.Fprintln(console, "Answer yes, no, or press return to accept the default.")
		}
	}
}
