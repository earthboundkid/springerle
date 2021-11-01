# Springerle [![GoDoc](https://godoc.org/github.com/carlmjohnson/springerle?status.svg)](https://godoc.org/github.com/carlmjohnson/springerle) [![Go Report Card](https://goreportcard.com/badge/github.com/carlmjohnson/springerle)](https://goreportcard.com/report/github.com/carlmjohnson/springerle)

Springerle are a kind of German prestamped cookie. Springerle is a command line tool for creating simple prestamped project files with the txtar format and Go templates. Inspired by [Cookiecutter](https://cookiecutter.readthedocs.io/) and [JTree Stamp](https://github.com/publicdomaincompany/jtree/tree/master/langs/stamp).

![](images/springerle.jpeg)

## Installation

First install [Go](http://golang.org).

If you just want to install the binary to your current directory and don't care about the source code, run

```bash
GOBIN=$(pwd) go install github.com/carlmjohnson/springerle@latest
```

## Screenshots

```
$ springerle -h
springerle v0.21.4 - create simple projects with the txtar format and Go templates.

Usage:

    springerle [options] <project file or URL>

Project files are Go templates processed as txtar files. The preamble to the
txtar file is used as prompts for creating the template context. Each line
should be formated as "key: User prompt question? default value" with colon and
question mark used as delimiters. Lines beginning with # or without a colon are
ignored. If the default value is "y" or "n", the prompt will be treated as a
boolean.

To templatize files that contain other templates, set -left-delim and
-right-delim options to something not used in the template.

Project files are Go templates processed as txtar files. The preamble to the
txtar file is used as a series of prompts for creating the template context.
Each line should be formated as "key: User prompt question? default value" with
colon and question mark used as delimiters. Lines beginning with # or without a
colon are ignored. If the default value is "y" or "n", the prompt will be
treated as a boolean. Prompt lines may use templates directives, e.g. to
transform a prior prompt value into a default or skip irrelevant prompts, but
premable template directives must be valid at the line level. That is, there
can be no multiline blocks in the preamble.

To templatize files that contain other templates, set -left-delim and
-right-delim options to something not used in the template.

In addition to the default Go template functions, templates can use the
functions listed below. In order to avoid name clashes, the added function
names follow a specific pattern: they combine their original package and
function names using no punctuation and only lowercase letters. E.g.,
strings.LastIndexByte becomes stringslastindexbyte.

From package strings:

stringscompare stringscontains stringscontainsany stringscontainsrune
stringscount stringsequalfold stringsfields stringsfieldsfunc stringshasprefix
stringshassuffix stringsindex stringsindexany stringsindexbyte
stringsindexfunc stringsindexrune stringsjoin stringslastindex
stringslastindexany stringslastindexbyte stringslastindexfunc stringsmap
stringsrepeat stringsreplace stringsreplaceall stringssplit stringssplitafter
stringssplitaftern stringssplitn stringstitle stringstolower
stringstolowerspecial stringstotitle stringstotitlespecial stringstoupper
stringstoupperspecial stringstovalidutf8 stringstrim stringstrimfunc
stringstrimleft stringstrimleftfunc stringstrimprefix stringstrimright
stringstrimrightfunc stringstrimspace stringstrimsuffix

From package path/filepath:

filepathabs filepathbase filepathclean filepathdir filepathext
filepathfromslash filepathisabs filepathjoin filepathmatch filepathrel
filepathsplit filepathsplitlist filepathtoslash filepathvolumename

From package time:

timedate timenow timeparse timeparseduration

From github.com/huandu/xstrings:

xstringscenter xstringscount xstringsdelete xstringsexpandtabs
xstringsfirstrunetolower xstringsfirstrunetoupper xstringsinsert
xstringslastpartition xstringsleftjustify xstringslen xstringspartition
xstringsreverse xstringsrightjustify xstringsrunewidth xstringsscrub
xstringsshuffle xstringsshufflesource xstringsslice xstringssqueeze
xstringssuccessor xstringsswapcase xstringstocamelcase xstringstokebabcase
xstringstosnakecase xstringstranslate xstringswidth xstringswordcount
xstringswordsplit

From github.com/mitchellh/go-wordwrap

wordwrapwrapstring wrapstring

The 'wordwrap' package is a slight exception to the rules for added function
names. Both 'wordwrapwrapstring' and 'wrapstring' are aliases to the same
function.

Options:
  -context JSON
        JSON object to use as template context
  -dest path
        destination path (default ".")
  -dry-run
        dry run output only (output txtar to stdout)
  -dump-context path
        path to load/save context produced by user input
  -left-delim delimiter
        left delimiter to use when parsing template (default "{{")
  -right-delim delimiter
        right delimiter to use when parsing template (default "}}")
  -verbose
        log debug output (default silent)
```
