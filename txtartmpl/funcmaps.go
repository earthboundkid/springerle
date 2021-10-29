package txtartmpl

import (
	"path/filepath"
	"strings"
	"time"

	"github.com/huandu/xstrings"
	"github.com/mitchellh/go-wordwrap"
)

func xStringFuncMap() map[string]interface{} {
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

func stringFuncMap() map[string]interface{} {
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

func filepathFuncMap() map[string]interface{} {
	return map[string]interface{}{
		"abs":          filepath.Abs,
		"base":         filepath.Base,
		"clean":        filepath.Clean,
		"dir":          filepath.Dir,
		"ext":          filepath.Ext,
		"filepathjoin": filepath.Join,
		"filepathsplit": func(path string) [2]string {
			head, tail := filepath.Split(path)
			return [...]string{head, tail}
		},
		"fromslash":  filepath.FromSlash,
		"isabs":      filepath.IsAbs,
		"match":      filepath.Match,
		"rel":        filepath.Rel,
		"splitlist":  filepath.SplitList,
		"toslash":    filepath.ToSlash,
		"volumename": filepath.VolumeName,
	}
}

func timeFuncMap() map[string]interface{} {
	return map[string]interface{}{
		"parseduration": time.ParseDuration,
		"date":          time.Date,
		"now":           time.Now,
		"parse":         time.Parse,
	}
}

func wordWrapFuncMap() map[string]interface{} {
	return map[string]interface{}{
		"wrapstring": wordwrap.WrapString,
	}
}
