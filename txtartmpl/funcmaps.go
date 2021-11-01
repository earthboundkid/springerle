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
		"xstringscenter":           xstrings.Center,
		"xstringscount":            xstrings.Count,
		"xstringsdelete":           xstrings.Delete,
		"xstringsexpandtabs":       xstrings.ExpandTabs,
		"xstringsfirstrunetolower": xstrings.FirstRuneToLower,
		"xstringsfirstrunetoupper": xstrings.FirstRuneToUpper,
		"xstringsinsert":           xstrings.Insert,
		"xstringslastpartition": func(str, sep string) [3]string {
			f, p, l := xstrings.LastPartition(str, sep)
			return [...]string{f, p, l}
		},
		"xstringsleftjustify": xstrings.LeftJustify,
		"xstringslen":         xstrings.Len,
		"xstringspartition": func(str, sep string) [3]string {
			f, p, l := xstrings.Partition(str, sep)
			return [...]string{f, p, l}
		},
		"xstringsreverse":       xstrings.Reverse,
		"xstringsrightjustify":  xstrings.RightJustify,
		"xstringsrunewidth":     xstrings.RuneWidth,
		"xstringsscrub":         xstrings.Scrub,
		"xstringsshuffle":       xstrings.Shuffle,
		"xstringsshufflesource": xstrings.ShuffleSource,
		"xstringsslice":         xstrings.Slice,
		"xstringssqueeze":       xstrings.Squeeze,
		"xstringssuccessor":     xstrings.Successor,
		"xstringsswapcase":      xstrings.SwapCase,
		"xstringstocamelcase":   xstrings.ToCamelCase,
		"xstringstokebabcase":   xstrings.ToKebabCase,
		"xstringstosnakecase":   xstrings.ToSnakeCase,
		"xstringstranslate":     xstrings.Translate,
		"xstringswidth":         xstrings.Width,
		"xstringswordcount":     xstrings.WordCount,
		"xstringswordsplit":     xstrings.WordSplit,
	}
}

func stringFuncMap() map[string]interface{} {
	return map[string]interface{}{
		"stringscompare":        strings.Compare,
		"stringscontains":       strings.Contains,
		"stringscontainsany":    strings.ContainsAny,
		"stringscontainsrune":   strings.ContainsRune,
		"stringscount":          strings.Count,
		"stringsequalfold":      strings.EqualFold,
		"stringsfields":         strings.Fields,
		"stringsfieldsfunc":     strings.FieldsFunc,
		"stringshasprefix":      strings.HasPrefix,
		"stringshassuffix":      strings.HasSuffix,
		"stringsindex":          strings.Index,
		"stringsindexany":       strings.IndexAny,
		"stringsindexbyte":      strings.IndexByte,
		"stringsindexfunc":      strings.IndexFunc,
		"stringsindexrune":      strings.IndexRune,
		"stringsjoin":           strings.Join,
		"stringslastindex":      strings.LastIndex,
		"stringslastindexany":   strings.LastIndexAny,
		"stringslastindexbyte":  strings.LastIndexByte,
		"stringslastindexfunc":  strings.LastIndexFunc,
		"stringsmap":            strings.Map,
		"stringsrepeat":         strings.Repeat,
		"stringsreplace":        strings.Replace,
		"stringsreplaceall":     strings.ReplaceAll,
		"stringssplit":          strings.Split,
		"stringssplitafter":     strings.SplitAfter,
		"stringssplitaftern":    strings.SplitAfterN,
		"stringssplitn":         strings.SplitN,
		"stringstitle":          strings.Title,
		"stringstolower":        strings.ToLower,
		"stringstolowerspecial": strings.ToLowerSpecial,
		"stringstotitle":        strings.ToTitle,
		"stringstotitlespecial": strings.ToTitleSpecial,
		"stringstoupper":        strings.ToUpper,
		"stringstoupperspecial": strings.ToUpperSpecial,
		"stringstovalidutf8":    strings.ToValidUTF8,
		"stringstrim":           strings.Trim,
		"stringstrimfunc":       strings.TrimFunc,
		"stringstrimleft":       strings.TrimLeft,
		"stringstrimleftfunc":   strings.TrimLeftFunc,
		"stringstrimprefix":     strings.TrimPrefix,
		"stringstrimright":      strings.TrimRight,
		"stringstrimrightfunc":  strings.TrimRightFunc,
		"stringstrimspace":      strings.TrimSpace,
		"stringstrimsuffix":     strings.TrimSuffix,
	}
}

func filepathFuncMap() map[string]interface{} {
	return map[string]interface{}{
		"filepathabs":       filepath.Abs,
		"filepathbase":      filepath.Base,
		"filepathclean":     filepath.Clean,
		"filepathdir":       filepath.Dir,
		"filepathext":       filepath.Ext,
		"filepathjoin":      filepath.Join,
		"filepathfromslash": filepath.FromSlash,
		"filepathisabs":     filepath.IsAbs,
		"filepathmatch":     filepath.Match,
		"filepathrel":       filepath.Rel,
		"filepathsplit": func(path string) [2]string {
			head, tail := filepath.Split(path)
			return [...]string{head, tail}
		},
		"filepathsplitlist":  filepath.SplitList,
		"filepathtoslash":    filepath.ToSlash,
		"filepathvolumename": filepath.VolumeName,
	}
}

func timeFuncMap() map[string]interface{} {
	return map[string]interface{}{
		"timedate":          time.Date,
		"timenow":           time.Now,
		"timeparse":         time.Parse,
		"timeparseduration": time.ParseDuration,
	}
}

func wordWrapFuncMap() map[string]interface{} {
	return map[string]interface{}{
		"wordwrapwrapstring": wordwrap.WrapString,
		"wrapstring":         wordwrap.WrapString,
	}
}
