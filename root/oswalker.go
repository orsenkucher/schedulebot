package root

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

// OSWalker is used to walk through folders starting from OSWalker.Root,
// Applying WalkFunc to every file (see WalkFunc behavior below)
type OSWalker struct {
	Root string
}

// WalkFunc is applied to every file on walkfn's level
// WalkFunc emits new WalkFunc which would be applied to files one level deeper
type WalkFunc func(path string, file os.FileInfo) WalkFunc

// Walk walks through files invoking current level WalkFunc on every file on current depth level
// New WalkFunc from Old WalkFunc is applied to sublevel children
// Walk starts from OSWalker.Root
func (w OSWalker) Walk(walkfn WalkFunc) {
	w.walk(walkfn, w.Root)
}

func (w OSWalker) walk(walkfn WalkFunc, path string) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		fn := walkfn(path, file)
		if file.IsDir() {
			w.walk(fn, filepath.Join(path, file.Name()))
		}
	}
}
