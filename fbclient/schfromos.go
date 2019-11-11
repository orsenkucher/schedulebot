package fbclient

import (
	"os"
	"path/filepath"

	"github.com/orsenkucher/schedulebot/root"
	"github.com/orsenkucher/schedulebot/route"
)

// OSSch is a pair of (TreePath, OSPath) pointing to schedule file
type OSSch struct {
	TrPath string
	OSPath string
}

// BuildOSSchedule is used to build tree using OS File system
func BuildOSSchedule() []OSSch {
	var res []OSSch
	w := root.OSWalker{Root: root.Rootdir}
	t := route.BuildOSTree()
	w.Walk(bindToTree(&res, t))
	return res
}

func bindToTree(res *[]OSSch, t *route.Tree) root.WalkFunc {
	return func(path string, file os.FileInfo) root.WalkFunc {
		name := file.Name()
		child, ok := t.Select(root.PopExt(name))
		if !ok {
			panic("Must be OK")
		}
		if !file.IsDir() {
			*res = append(*res, OSSch{
				OSPath: filepath.Join(path, name),
				TrPath: child.MakePath(),
			})
		}
		return bindToTree(res, child)
	}
}
