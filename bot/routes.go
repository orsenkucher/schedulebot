package bot

import (
	"errors"
	"fmt"
)

type routeTree struct {
	name     string
	parent   *routeTree
	children []*routeTree
}

func (t *routeTree) String() string {
	if t.parent == nil {
		return t.name
	}
	return t.parent.String() + "." + t.name
}

func (t *routeTree) doChild(name string) *routeTree {
	child := &routeTree{
		parent: t,
		name:   name,
	}
	t.children = append(t.children, child)
	return child
}

func (t *routeTree) Select(childName string) (*routeTree, error) {
	for _, c := range t.children {
		if c.name == childName {
			return c, nil
		}
	}
	return nil, errors.New("Child not found")
}

// Routes are possible routes
var Routes = makeRoutes()

func makeRoutes() *routeTree {
	t0 := &routeTree{name: "КНУ"}
	t01 := t0.doChild("Мехмат")
	t02 := t0.doChild("Фізфак")
	fmt.Println(t02)
	t011 := t01.doChild("1 курс")
	t012 := t01.doChild("2 курс")
	fmt.Println(t012)
	t0111 := t011.doChild("1 група")
	fmt.Println(t0111)
	return t0
}

type user int64

var currentRoutes = make(map[user]*routeTree)
