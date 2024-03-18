package goff

import "os"

// FSObjects are basically caching the walk results
type FSObject struct {
	Path     string // path to the object
	StatInfo os.FileInfo
}

// FSOLists are to attach filters to the results of walk
type FSOList []FSObject

func (fl FSOList) KeepOnly(ff FilterFunc) FSOList {
	var ret FSOList
	for _, f := range fl {
		if ff(f.Path, f.StatInfo) {
			ret = append(ret, f)
		}
	}
	return ret
}

func (fl FSOList) Remove(ff FilterFunc) FSOList {
	var ret FSOList
	for _, f := range fl {
		if !ff(f.Path, f.StatInfo) {
			ret = append(ret, f)
		}
	}
	return ret
}
