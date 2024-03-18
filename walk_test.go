package goff

import (
	"fmt"
	"math/rand/v2"
	"os"
	"path/filepath"
	"testing"
)

func setup() string {
	tmpRoot, err := os.MkdirTemp("", "goff")
	if err != nil {
		panic(err)
	}
	makeTestDirs(0700, 501, 20, tmpRoot, "d1", "d1/deep1", "d1/deep2", "d2")
	r := rand.New(rand.NewPCG(1, 2))
	makeRandTestFiles(r, 0600, 501, 20, 60, 80, tmpRoot, "d1/deep1/f1.txt", "d1/deep2/f2.txt", "d2/f3.txt")
	os.Symlink(filepath.Join(tmpRoot, "d2/f3.txt"), filepath.Join(tmpRoot, "d2/link"))
	return tmpRoot
}

func paths(l FSOList) []string {
	var ret []string
	for _, f := range l {
		ret = append(ret, f.Path)
	}
	return ret
}

func absFunc(tmpRoot string) func(string) string {
	return func(s string) string {
		return filepath.Join(tmpRoot, s)
	}
}

func setCmp(s1, s2 []string) bool {
	m1 := make(map[string]struct{})
	m2 := make(map[string]struct{})
	for _, s := range s1 {
		m1[s] = struct{}{}
	}
	for _, s := range s2 {
		m2[s] = struct{}{}
	}
	same := true
	for k := range m1 {
		if _, ok := m2[k]; !ok {
			fmt.Printf("> %s\n", k)
			same = false
		}
	}
	for k := range m2 {
		if _, ok := m1[k]; !ok {
			fmt.Printf("< %s\n", k)
			same = false
		}
	}
	return same
}

func TestWalk(t *testing.T) {
	tmpRoot := setup()
	t.Logf("tmpRoot created %s", tmpRoot)
	defer os.RemoveAll(tmpRoot)
	abs := absFunc(tmpRoot)
	all := Walk(tmpRoot, true)
	all_paths_want := []string{
		abs(""),
		abs("d1"),
		abs("d1/deep1"),
		abs("d1/deep1/f1.txt"),
		abs("d1/deep2"),
		abs("d1/deep2/f2.txt"),
		abs("d2"),
		abs("d2/f3.txt"),
		abs("d2/link"),
	}
	all_paths := paths(all)
	if !setCmp(all_paths_want, all_paths) {
		t.Errorf("all_paths not equal have: %v want: %v", all_paths, all_paths_want)
	}

	all_symlinks_want := []string{
		abs("d2/link"),
	}

	symlinks_paths := paths(all.KeepOnly(IsSymLink))

	if !setCmp(all_symlinks_want, symlinks_paths) {
		t.Errorf("symlinks_paths not equal have: %v want: %v", symlinks_paths, all_symlinks_want)
	}

	all_non_symlinks_want := []string{
		abs(""),
		abs("d1"),
		abs("d1/deep1"),
		abs("d1/deep1/f1.txt"),
		abs("d1/deep2"),
		abs("d1/deep2/f2.txt"),
		abs("d2"),
		abs("d2/f3.txt"),
	}

	non_symlinks_paths := paths(all.Remove(IsSymLink))
	if !setCmp(all_non_symlinks_want, non_symlinks_paths) {
		t.Errorf("non_symlinks_paths not equal have: %v want: %v", non_symlinks_paths, all_non_symlinks_want)
	}

	dirs_want := []string{
		abs(""),
		abs("d1"),
		abs("d1/deep1"),
		abs("d1/deep2"),
		abs("d2"),
	}

	dir_paths := paths(all.KeepOnly(IsDir))
	if !setCmp(dirs_want, dir_paths) {
		t.Errorf("dir_paths not equal have: %v want: %v", dir_paths, dirs_want)
	}

	regular_want := []string{
		abs("d1/deep1/f1.txt"),
		abs("d1/deep2/f2.txt"),
		abs("d2/f3.txt"),
	}
	regular_paths := paths(all.KeepOnly(IsRegular))
	if !setCmp(regular_want, regular_paths) {
		t.Errorf("regular_paths not equal have: %v want: %v", regular_paths, regular_want)
	}

}
