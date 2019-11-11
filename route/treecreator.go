package route

import (
	"os"

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
	return func(_ string, file os.FileInfo) root.WalkFunc {
		name := file.Name()
		if !file.IsDir() {
			name = root.PopExt(name)
		}
		child := tr.MakeChild(name)
		return bindToTree(child)
	}
}
