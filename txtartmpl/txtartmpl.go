package txtartmpl

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime/debug"
	"sort"
	"strings"
	"text/template"

	"github.com/Songmu/prompter"
	"github.com/carlmjohnson/flagext"
	"github.com/mitchellh/go-wordwrap"
	"golang.org/x/tools/txtar"
)

const AppName = "springerle"

// CLI runs the springerle command line application. The application name (os.Args[0])
// should not be passed to CLI. The returned error contains the CLI's the exit code.
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
	fl.Func("context", "`JSON` object to use as template context", app.setTmplCtx)
	fl.StringVar(&app.dumpCtx, "dump-context", "", "`path` to load/save context produced by user input")
	fl.StringVar(&app.leftD, "left-delim", "{{", "left `delimiter` to use when parsing template")
	fl.StringVar(&app.rightD, "right-delim", "}}", "right `delimiter` to use when parsing template")

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
	dstPath       string
	dryRun        bool
	leftD, rightD string
	src           io.ReadCloser
	tmplCtx       map[string]interface{}
	dumpCtx       string
	*log.Logger
}

func (app *appEnv) setusage(fl *flag.FlagSet) {
	fl.Usage = func() {
		version := "(unknown)"
		if i, ok := debug.ReadBuildInfo(); ok {
			version = i.Main.Version
		}
		s := fmt.Sprintf(
			`springerle %s - create simple projects with the txtar format and Go templates.

Usage:

	springerle [options] <project file or URL>

Project files are Go templates processed as txtar files. The preamble to the txtar file is used as a series of prompts for creating the template context. Each line should be formated as "key: User prompt question? default value" with colon and question mark used as delimiters. Lines beginning with # or without a colon are ignored. If the default value is "y" or "n", the prompt will be treated as a boolean. Prompt lines may use templates directives, e.g. to transform a prior prompt value into a default or skip irrelevant prompts, but premable template directives must be valid at the line level. That is, there can be no multiline blocks in the preamble.

To templatize files that contain other templates, set -left-delim and -right-delim options to something not used in the template.

In addition to the default Go template functions, templates can use the following functions.

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
			version,
			sortFuncMapNames(stringFuncMap()),
			sortFuncMapNames(filepathFuncMap()),
			sortFuncMapNames(timeFuncMap()),
			sortFuncMapNames(xStringFuncMap()),
			sortFuncMapNames(wordWrapFuncMap()),
		)
		fmt.Fprint(fl.Output(), wordwrap.WrapString(s, 79))
		fl.PrintDefaults()
		fmt.Fprintln(fl.Output())
	}
}

func (app *appEnv) setTmplCtx(s string) error {
	app.tmplCtx = make(map[string]interface{})
	return json.Unmarshal([]byte(s), &app.tmplCtx)
}

func sortFuncMapNames(m template.FuncMap) string {
	ss := make([]string, 0, len(m))
	for k := range m {
		ss = append(ss, k)
	}
	sort.Strings(ss)
	return strings.Join(ss, " ")
}

func (app *appEnv) Exec() (err error) {
	var buf bytes.Buffer

	if _, err = io.Copy(&buf, app.src); err != nil {
		return err
	}
	// check template validity
	t := template.New("").
		Option("missingkey=error").
		Delims(app.leftD, app.rightD).
		Funcs(stringFuncMap()).
		Funcs(filepathFuncMap()).
		Funcs(timeFuncMap()).
		Funcs(xStringFuncMap()).
		Funcs(wordWrapFuncMap())
	if t, err = t.Parse(buf.String()); err != nil {
		return fmt.Errorf("could not parse input as template: %w", err)
	}
	// read preamble by line, make up a Question context map
	ar := txtar.Parse(buf.Bytes())
	tctx, err := app.TemplateContextFrom(ar.Comment)
	if err != nil {
		return err
	}
	// feed src through template.Template
	buf.Reset()
	if err = t.Execute(&buf, tctx); err != nil {
		return err
	}

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
		name := filepath.Clean(f.Name)
		if strings.HasPrefix(name, "../") {
			return fmt.Errorf("won't write unsafe file name to disk: %q", f.Name)
		}
		name = filepath.FromSlash(filepath.Join(app.dstPath, name))
		if err := os.MkdirAll(filepath.Dir(name), 0o777); err != nil {
			return err
		}
		app.Printf("writing %q", name)
		var perm os.FileMode = 0o666
		if filepath.Ext(name) == ".sh" {
			perm = 0o777
		}
		if err := os.WriteFile(name, f.Data, perm); err != nil {
			return err
		}
	}

	return err
}

func (app *appEnv) dumpContext(tctx map[string]interface{}) {
	if app.dumpCtx == "" {
		return
	}
	b, err := json.Marshal(tctx)
	if err != nil {
		panic(err)
	}
	if err = os.WriteFile(app.dumpCtx, b, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "problem dumping context file: %v\n", err)
	}
}

func (app *appEnv) TemplateContextFrom(b []byte) (map[string]interface{}, error) {
	if app.tmplCtx != nil {
		return app.tmplCtx, nil
	}
	m := make(map[string]interface{})
	if app.dumpCtx != "" {
		if b, err := os.ReadFile(app.dumpCtx); err == nil {
			_ = json.Unmarshal(b, &m)
		}
	}
	t := template.New("").
		Delims(app.leftD, app.rightD).
		Funcs(stringFuncMap()).
		Funcs(filepathFuncMap()).
		Funcs(timeFuncMap()).
		Funcs(xStringFuncMap()).
		Funcs(wordWrapFuncMap())
	s := bufio.NewScanner(bytes.NewReader(b))
	for s.Scan() {
		if err := processLine(t, s.Text(), m); err != nil {
			return nil, err
		}
		app.dumpContext(m)
	}

	return m, s.Err()
}

func processLine(t *template.Template, line string, m map[string]interface{}) error {
	var (
		k, q, v string
		i       int
	)
	if strings.HasPrefix(line, "#") {
		return nil
	}
	if strings.Contains(line, "{"+"{") {
		var buf strings.Builder
		t, err := t.Parse(line)
		if err != nil {
			return fmt.Errorf("could not parse preliminary prompt as template: %w", err)
		}
		t.Execute(&buf, m)
		line = buf.String()
	}
	if strings.HasPrefix(line, "#") {
		return nil
	}
	if i = strings.IndexByte(line, ':'); i == -1 {
		return nil
	}
	k = strings.TrimSpace(line[:i])
	q = k
	line = line[i+1:]

	if i = strings.IndexByte(line, '?'); i != -1 {
		q = strings.TrimSpace(line[:i+1])
		line = line[i+1:]
	}
	v = strings.TrimSpace(line)

	if def, ok := m[k]; ok {
		if defb, ok := def.(bool); ok {
			m[k] = prompter.YN(q, defb)
			return nil
		}
		if defs, ok := def.(string); ok {
			m[k] = prompter.Prompt(q, defs)
			return nil
		}
	}

	if l := strings.ToLower(v); l == "y" || l == "n" {
		m[k] = prompter.YN(q, l == "y")
		return nil
	}

	m[k] = prompter.Prompt(q, v)
	return nil
}
