package goff

import (
	"math/rand/v2"
	"os"
	"path/filepath"
)

func makeTestDirs(mode os.FileMode, uid, gid int, rootdir string, elements ...string) error {
	for _, e := range elements {
		fp := filepath.Join(rootdir, e)
		err := os.Mkdir(fp, mode)
		if err != nil {
			return err
		}
		if uid != -1 || gid != -1 {
			err = os.Chown(fp, uid, gid)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func randStringRunes(r *rand.Rand, n int) string {
	letterRunes := []rune("0123456789abcdef")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[r.IntN(len(letterRunes))]
	}
	return string(b)
}

func makeRandTestFiles(r *rand.Rand, mode os.FileMode, uid, gid, minSize, maxSize int, rootdir string, elements ...string) error {
	for _, e := range elements {
		fp := filepath.Join(rootdir, e)
		size := r.IntN(maxSize-minSize) + minSize
		content := randStringRunes(r, size)
		err := os.WriteFile(fp, []byte(content), mode)
		if err != nil {
			return err
		}
		if uid != -1 || gid != -1 {
			err = os.Chown(fp, uid, gid)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
