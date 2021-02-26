package txtartmpl

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	"github.com/carlmjohnson/flagext"
	"github.com/manifoldco/promptui"
	"github.com/mitchellh/go-wordwrap"
	"golang.org/x/tools/txtar"
)

const AppName = "springerle"

func CLI(args []string) error {
	var app appEnv
	err := app.ParseArgs(args)
	if err != nil {
		return err
	}

	if err = app.Exec(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}
	return err
}

func (app *appEnv) ParseArgs(args []string) error {
	fl := flag.NewFlagSet(AppName, flag.ContinueOnError)
	app.Logger = log.New(nil, AppName+" ", log.LstdFlags)
	flagext.LoggerVar(
		fl, app.Logger, "verbose", flagext.LogVerbose, "log debug output")
	fl.StringVar(&app.dstPath, "dest", ".", "destination `path`")
	fl.BoolVar(&app.dryRun, "dry-run", false, "dry run output only (output txtar to stdout)")
	app.setusage(fl)
	if err := fl.Parse(args); err != nil {
		return err
	}
	if err := flagext.ParseEnv(fl, AppName); err != nil {
		return err
	}
	if err := flagext.MustHaveArgs(fl, 0, 1); err != nil {
		return err
	}
	src := flagext.FileOrURL(flagext.StdIO, nil)
	app.src = src
	if fl.NArg() > 0 {
		return src.Set(fl.Arg(0))
	}
	return nil
}

type appEnv struct {
	dstPath string
	dryRun  bool
	src     io.ReadCloser
	*log.Logger
}

func (app *appEnv) setusage(fl *flag.FlagSet) {
	fl.Usage = func() {
		fmt.Fprintf(
			fl.Output(),
			`springerle - create simple projects with the txtar format and Go templates.

Usage:

	springerle [options] <project file or URL>

Project files are Go templates processed as txtar files. The preamble to the
txtar file is used as prompts for creating the template context. Each line
should be formated as "key: User prompt question? default value" with colon and
question mark used as delimiters. Lines beginning with # or without a colon are
ignored. If the default value is "y" or "n", the prompt will be treated as a
boolean.

In addition to the default Go template functions, templates can use the
following functions.

From package strings:

%s

From package path/filepath:

%s

From package time:

%s

From github.com/huandu/xstrings:

%s

From github.com/mitchellh/go-wordwrap

%s

Options:
`,
			sortFuncMapNames(StringFuncMap()),
			sortFuncMapNames(FilepathFuncMap()),
			sortFuncMapNames(TimeFuncMap()),
			sortFuncMapNames(XStringFuncMap()),
			sortFuncMapNames(WordWrapFuncMap()),
		)
		fl.PrintDefaults()
		fmt.Fprintln(fl.Output())
	}
}

func sortFuncMapNames(m template.FuncMap) string {
	ss := make([]string, 0, len(m))
	for k := range m {
		ss = append(ss, k)
	}
	sort.Strings(ss)
	return wordwrap.WrapString(strings.Join(ss, " "), 80)
}

func (app *appEnv) Exec() (err error) {
	var buf bytes.Buffer

	if _, err = io.Copy(&buf, app.src); err != nil {
		return err
	}
	// check template validity
	t := template.New("").
		Funcs(StringFuncMap()).
		Funcs(FilepathFuncMap()).
		Funcs(TimeFuncMap()).
		Funcs(XStringFuncMap()).
		Funcs(WordWrapFuncMap())
	if t, err = t.Parse(buf.String()); err != nil {
		return fmt.Errorf("could not parse input as template: %w", err)
	}
	// read preamble by line, make up a Question context map
	ar := txtar.Parse(buf.Bytes())
	tctx, err := app.getTCtx(ar.Comment)
	if err != nil {
		return err
	}
	// feed src through template.Template
	buf.Reset()
	t.Execute(&buf, tctx)

	// make all the files
	ar = txtar.Parse(buf.Bytes())
	if app.dryRun {
		app.Printf("dry run for %q", app.dstPath)
		ar.Comment = nil
		s := string(txtar.Format(ar))
		fmt.Print(s)
		if !strings.HasSuffix(s, "\n") {
			fmt.Println()
		}
		return nil
	}

	for _, f := range ar.Files {
		name := filepath.FromSlash(filepath.Join(app.dstPath, filepath.Clean(f.Name)))
		if err := os.MkdirAll(filepath.Dir(name), 0o777); err != nil {
			return err
		}
		app.Printf("writing %q", name)
		if err := os.WriteFile(name, f.Data, 0o666); err != nil {
			return err
		}
	}

	return err
}

func (app *appEnv) getTCtx(b []byte) (map[string]string, error) {
	m := make(map[string]string)
	s := bufio.NewScanner(bytes.NewReader(b))
	for s.Scan() {
		line := s.Text()
		var (
			k, q, v string
			i       int
		)
		if strings.HasPrefix(line, "#") {
			continue
		}
		if i = strings.IndexByte(line, ':'); i == -1 {
			continue
		}
		k = strings.TrimSpace(line[:i])
		q = k
		line = line[i+1:]

		if i = strings.IndexByte(line, '?'); i != -1 {
			q = strings.TrimSpace(line[:i+1])
			line = line[i+1:]
		}
		v = strings.TrimSpace(line)
		validator := boolPrompt
		if v != "y" && v != "n" {
			validator = nil
		}
		prompt := promptui.Prompt{
			Label:    q,
			Default:  v,
			Validate: validator,
		}
		var err error
		m[k], err = prompt.Run()
		if err != nil {
			return nil, err
		}
	}
	return m, s.Err()
}

var boolVals = [...]string{"y", "yes", "no", "n"}

func boolPrompt(input string) error {
	for _, value := range boolVals {
		if strings.ToLower(input) == value {
			return nil
		}
	}
	return fmt.Errorf("Response should be one of the values %v", boolVals)
}
