package fbclient

import "github.com/orsenkucher/schedulebot/route"

// SchTree is like route.Tree, but with schpath pointed to schedule file
type SchTree struct {
	route.Tree
	SchPath string
}

// Flatten flattens tree skipping joints
func (st *SchTree) Flatten() []*SchTree {
	// var res []*SchTree
	return nil
}

func (st *SchTree) flatten(res *[]*SchTree) {
	if st.Children == nil {
		// *res = append(*res, st.Children...)
	} else {
		// for _, child := range st.Children
	}
}
