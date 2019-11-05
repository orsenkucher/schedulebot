package route

import (
	"fmt"
)

// Tree is a node (and a tree) of a tree data structure.
// Used to navigate user to schedule he's interested in.
type Tree struct {
	Name     string
	Children []*Tree
	Parent   *Tree
}

func (t *Tree) String() string {
	if t.Parent == nil {
		return t.Name
	}
	return t.Parent.String() + "." + t.Name
}

func (t *Tree) makeChild(name string) *Tree {
	child := &Tree{
		Parent: t,
		Name:   name,
	}
	t.Children = append(t.Children, child)
	return child
}

// Select searches for child by its name
// returns found child and success flag
func (t *Tree) Select(childName string) (*Tree, bool) {
	for _, c := range t.Children {
		if c.Name == childName {
			return c, true
		}
	}
	return nil, false
}

// // Parent returns parent of this node and success flag
// func (t *Tree) Parent() (*Tree, bool) {
// 	if t.parent != nil {
// 		return t.parent, true
// 	}
// 	return nil, false
// }

// func (t *Tree) Child

// Routes are possible routes
var Routes = makeRoutes()

func makeRoutes() *Tree {
	t0 := &Tree{Name: "КНУ"}
	t01 := t0.makeChild("Мехмат")
	t02 := t0.makeChild("Фізфак")
	fmt.Println(t02)
	t011 := t01.makeChild("1 курс")
	t012 := t01.makeChild("2 курс")
	fmt.Println(t012)
	t0111 := t011.makeChild("1 група")
	fmt.Println(t0111)
	return t0
}
