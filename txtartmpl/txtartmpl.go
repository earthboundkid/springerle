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
	"strings"
	"text/template"

	"github.com/carlmjohnson/flagext"
	"github.com/huandu/xstrings"
	"github.com/manifoldco/promptui"
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
	fl.Usage = app.usage(fl)
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

func (app *appEnv) usage(fl *flag.FlagSet) func() {
	return func() {
		fmt.Fprintf(fl.Output(), `springerle - a Go CLI application template cat clone

Usage:

	springerle [options]

TODO: write help, document template funcs

Options:
`)
		fl.PrintDefaults()
		fmt.Fprintln(fl.Output(), "")
	}
}

func (app *appEnv) Exec() (err error) {
	var buf bytes.Buffer

	if _, err = io.Copy(&buf, app.src); err != nil {
		return err
	}
	// check template validity
	t := template.New("").
		Funcs(XStringFuncMap()).
		Funcs(StringFuncMap()).
		Funcs(FilepathFuncMap())
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
		prompt := promptui.Prompt{
			Label:     q,
			Default:   v,
			IsConfirm: v == "y" || v == "n",
		}
		var err error
		m[k], err = prompt.Run()
		if err != nil {
			return nil, err
		}
	}
	return m, s.Err()
}

func XStringFuncMap() map[string]interface{} {
	return map[string]interface{}{
		"center":           xstrings.Center,
		"countpattern":     xstrings.Count,
		"delete":           xstrings.Delete,
		"expandtabs":       xstrings.ExpandTabs,
		"firstrunetolower": xstrings.FirstRuneToLower,
		"firstrunetoupper": xstrings.FirstRuneToUpper,
		"insert":           xstrings.Insert,
		"lastpartition": func(str, sep string) [3]string {
			f, p, l := xstrings.LastPartition(str, sep)
			return [...]string{f, p, l}
		},
		"leftjustify": xstrings.LeftJustify,
		"runelen":     xstrings.Len,
		"partition": func(str, sep string) [3]string {
			f, p, l := xstrings.Partition(str, sep)
			return [...]string{f, p, l}
		},
		"reverse":       xstrings.Reverse,
		"rightjustify":  xstrings.RightJustify,
		"runewidth":     xstrings.RuneWidth,
		"scrub":         xstrings.Scrub,
		"shuffle":       xstrings.Shuffle,
		"shufflesource": xstrings.ShuffleSource,
		"slice":         xstrings.Slice,
		"squeeze":       xstrings.Squeeze,
		"successor":     xstrings.Successor,
		"swapcase":      xstrings.SwapCase,
		"tocamelcase":   xstrings.ToCamelCase,
		"tokebabcase":   xstrings.ToKebabCase,
		"tosnakecase":   xstrings.ToSnakeCase,
		"translate":     xstrings.Translate,
		"width":         xstrings.Width,
		"wordcount":     xstrings.WordCount,
		"wordsplit":     xstrings.WordSplit,
	}
}

func StringFuncMap() map[string]interface{} {
	return map[string]interface{}{
		"compare":        strings.Compare,
		"contains":       strings.Contains,
		"containsany":    strings.ContainsAny,
		"containsrune":   strings.ContainsRune,
		"count":          strings.Count,
		"equalfold":      strings.EqualFold,
		"fields":         strings.Fields,
		"fieldsfunc":     strings.FieldsFunc,
		"hasprefix":      strings.HasPrefix,
		"hassuffix":      strings.HasSuffix,
		"index":          strings.Index,
		"indexany":       strings.IndexAny,
		"indexbyte":      strings.IndexByte,
		"indexfunc":      strings.IndexFunc,
		"indexrune":      strings.IndexRune,
		"join":           strings.Join,
		"lastindex":      strings.LastIndex,
		"lastindexany":   strings.LastIndexAny,
		"lastindexbyte":  strings.LastIndexByte,
		"lastindexfunc":  strings.LastIndexFunc,
		"map":            strings.Map,
		"repeat":         strings.Repeat,
		"replace":        strings.Replace,
		"replaceall":     strings.ReplaceAll,
		"split":          strings.Split,
		"splitafter":     strings.SplitAfter,
		"splitaftern":    strings.SplitAfterN,
		"splitn":         strings.SplitN,
		"title":          strings.Title,
		"tolower":        strings.ToLower,
		"tolowerspecial": strings.ToLowerSpecial,
		"totitle":        strings.ToTitle,
		"totitlespecial": strings.ToTitleSpecial,
		"toupper":        strings.ToUpper,
		"toupperspecial": strings.ToUpperSpecial,
		"tovalidutf8":    strings.ToValidUTF8,
		"trim":           strings.Trim,
		"trimfunc":       strings.TrimFunc,
		"trimleft":       strings.TrimLeft,
		"trimleftfunc":   strings.TrimLeftFunc,
		"trimprefix":     strings.TrimPrefix,
		"trimright":      strings.TrimRight,
		"trimrightfunc":  strings.TrimRightFunc,
		"trimspace":      strings.TrimSpace,
		"trimsuffix":     strings.TrimSuffix,
	}
}

func FilepathFuncMap() map[string]interface{} {
	return map[string]interface{}{
		"abs":          filepath.Abs,
		"base":         filepath.Base,
		"clean":        filepath.Clean,
		"dir":          filepath.Dir,
		"evalsymlinks": filepath.EvalSymlinks,
		"ext":          filepath.Ext,
		"fromslash":    filepath.FromSlash,
		"glob":         filepath.Glob,
		"hasprefix":    filepath.HasPrefix,
		"isabs":        filepath.IsAbs,
		"join":         filepath.Join,
		"match":        filepath.Match,
		"rel":          filepath.Rel,
		"split": func(path string) [2]string {
			head, tail := filepath.Split(path)
			return [...]string{head, tail}
		},
		"splitlist":  filepath.SplitList,
		"toslash":    filepath.ToSlash,
		"volumename": filepath.VolumeName,
	}
}
