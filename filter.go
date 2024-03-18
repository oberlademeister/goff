package goff

import "os"

type FilterFunc func(path string, StatInfo os.FileInfo) bool

func IsDir(path string, statInfo os.FileInfo) bool {
	return statInfo.Mode().IsDir()
}

func IsRegular(path string, statInfo os.FileInfo) bool {
	return statInfo.Mode().IsRegular()
}

func IsSymLink(path string, statInfo os.FileInfo) bool {
	return !(statInfo.Mode().Type()&os.ModeSymlink == 0)
}
