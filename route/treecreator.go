package route

import (
	"io/ioutil"
	"path/filepath"
)

// TreeCreator is capable of creating route tree
// from filesystem or firestore structure or whatever
type TreeCreator interface {
	Create() *Tree
}

// Rootdir path to the root directory
const Rootdir = "data"

// LocalCreator creates route tree from local file system
type LocalCreator struct {
	Root string
}

// Create walks through files starting from lc.Root and builds route.Tree
func (lc LocalCreator) Create() *Tree {
	t := &Tree{Name: "Root"}
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
		c := t.makeChild(name)
		if file.IsDir() {
			lc.fillTree(c, filepath.Join(path, name))
		}
	}

}
