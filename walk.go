package goff

import (
	"io/fs"
	"log/slog"
	"path/filepath"
)

func Walk(root string, quickabort bool) FSOList {
	var ret FSOList
	err := filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			slog.Warn("goff.Walk err", "quickabort", quickabort, "path", path, "err", err)
			if quickabort {
				return err
			}
		}
		ret = append(ret, FSObject{Path: path, StatInfo: info})
		return err
	})
	if err != nil {
		slog.Warn("goff.Walk filepath.Walk error", "err", err)
	}
	return ret
}
