package route

import (
	"io/ioutil"
	"path/filepath"
	"strings"
)

// TreeCreator is capable of creating route tree
// from filesystem or firestore structure or whatever
type TreeCreator interface {
	Create() *Tree
}

// LocalCreator creates route tree from local file system
type LocalCreator struct {
	Root string
}

// Create walks through files starting from lc.Root and builds route.Tree
func (lc LocalCreator) Create() *Tree {
	t := &Tree{Name: "root"}
	lc.fillTree(t, lc.Root)
	return t
}

func (lc LocalCreator) fillTree(t *Tree, path string) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		name := file.Name()
		if !file.IsDir() {
			name = strings.TrimSuffix(name, filepath.Ext(name))
		}
		c := t.makeChild(name)
		if file.IsDir() {
			lc.fillTree(c, filepath.Join(path, name))
		}
	}
}
