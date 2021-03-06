package route

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/orsenkucher/schedulebot/hash"
)

// Tree is a node (and a tree) of a tree data structure.
// Used to navigate user to schedule he's interested in.
type Tree struct {
	Name     string
	Children []*Tree
	Parent   *Tree
}

func (t *Tree) String() string {
	chain := t.chain()
	return strings.Join(chain, " ⫶ ")
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

// CalcHash64 calculates unique hash for current node path
func (t *Tree) CalcHash64() string {
	return hash.EncodeAsBase64(t.MakePath())
}

func (t *Tree) chain() []string {
	if t.Parent == nil {
		return nil
	}
	return append(t.Parent.chain(), t.Name)
}

// MakeChild creates child node with provided name
func (t *Tree) MakeChild(name string) *Tree {
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
		log.Println("route Tree.Find invalid path")
		return nil, false
	}
	for _, name := range chain {
		var ok bool
		t, ok = t.Select(name)
		if !ok {
			return nil, false
		}
	}
	return t, true
}

// Print prints this tree
func (t *Tree) Print() {
	t.print(1)
}

func (t *Tree) print(n int) {
	fmt.Println(t.Name) // ┊ ╰
	for _, child := range t.Children {
		fmt.Print(strings.Repeat("  ", n))
		fmt.Print("╰")
		child.print(n + 1)
	}
}

// GenerateUsersTree is pub
func GenerateUsersTree(nodes [][]string) *Tree {
	t := &Tree{Name: "Root"}
	for _, node := range nodes {
		t.addNode(node)
	}
	return t
}

func (t *Tree) addNode(path []string) {
	for _, node := range path {
		next, ok := t.Select(node)
		if !ok {
			t = t.MakeChild(node)
		} else {
			t = next
		}
	}
}

// Drop drops down to the nearest fork
func (t *Tree) Drop() *Tree {
	if t.Children == nil {
		return t
	}
	for t.Children != nil && len(t.Children) == 1 {
		t = t.Children[0]
	}
	if t.Children == nil {
		t = t.Parent
	}
	return t
}

// Jump drop to fork
func (t *Tree) Jump() *Tree {
	if t := t.Parent; t != nil {
		for t.Parent != nil && len(t.Children) == 1 {
			t = t.Parent
		}
		return t
	}
	return t
}
