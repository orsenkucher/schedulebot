package route

import (
	"encoding/json"
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
	return t.Parent.String() + " ⫶ " + t.Name
}

// MakePath is used to create valid path to current node
func (t *Tree) MakePath() string {
	chain := t.chain()
	bytes, err := json.Marshal(chain)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

func (t *Tree) chain() []string {
	chain := []string{t.Name}
	if t.Parent == nil {
		return chain
	}
	return append(t.Parent.chain(), chain...)
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

// Find searches for route node by path
// Fist selector of the path have to be the name of current node
func (t *Tree) Find(path string) (*Tree, bool) {
	var chain []string
	err := json.Unmarshal([]byte(path), &chain)
	if err != nil {
		panic(err)
	}
	for _, name := range chain[1:] {
		var ok bool
		t, ok = t.Select(name)
		if !ok {
			return nil, false
		}
	}
	return t, true
}

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
	t0111path := t0111.MakePath()
	fmt.Println(t0111path)
	found, _ := t0.Find(t0111path)
	fmt.Println(found)
	fmt.Println(found == t0111)
	return t0
}
