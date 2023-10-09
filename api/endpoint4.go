package api

import "strings"

func NewHeaderInfo(active, endpath, section string, path ...string) HeaderInfo {
	if len(endpath) > 50 {
		endpath = endpath[:50]
	}
	return HeaderInfo{
		Active:  active,
		Path:    strings.Join(path, " > "),
		EndPath: endpath,
		Section: section,
	}
}
