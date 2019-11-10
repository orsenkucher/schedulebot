package route

import (
	"io/ioutil"
	"path/filepath"
	"strings"
)

// TreeCreator is capable of creating route tree
// from filesystem or firestore structure or whatever
type TreeCreator interface {
	Create(MyFn)
}

// LocalCreator creates route tree from local file system
type LocalCreator struct {
	Root string
}

// type ChildMaker interface{

// }

// MyFn recursive fn
type MyFn func(path, name string) MyFn

// Create walks through files starting from lc.Root and builds route.Tree
func (lc LocalCreator) Create(fn MyFn) {
	lc.fillTree(fn, lc.Root)
}

func (lc LocalCreator) fillTree(fn MyFn, path string) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		name := file.Name()
		if !file.IsDir() {
			name = strings.TrimSuffix(name, filepath.Ext(name))
		}
		fn := fn(path, name)
		if file.IsDir() {
			lc.fillTree(fn, filepath.Join(path, name))
		}
	}
}
