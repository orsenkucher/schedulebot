package route

import (
	"github.com/orsenkucher/schedulebot/root"
)

// BuildOSTree is used to build tree using OS File system
func BuildOSTree() *Tree {
	t := &Tree{Name: "root"}
	w := root.OSWalker{Root: root.Rootdir}
	w.Walk(bindToTree(t))
	return t
}

func bindToTree(tr *Tree) root.WalkFunc {
	return func(path, name string) root.WalkFunc {
		child := tr.MakeChild(name)
		return bindToTree(child)
	}
}
