package route

// TreeRoot is the root node of route.Tree
type TreeRoot struct {
	Rootnode *Tree
	hashmap  map[string]string
}

// NewTreeRoot creates TreeRoot with provided node
func NewTreeRoot(rootnode *Tree) *TreeRoot {
	hashmap := make(map[string]string)
	tr := &TreeRoot{
		Rootnode: rootnode,
		hashmap:  *buildHashMap(rootnode, &hashmap)}
	return tr
}

func buildHashMap(from *Tree, hashmap *map[string]string) *map[string]string {
	(*hashmap)[from.CalcHash64()] = from.MakePath()
	for _, child := range from.Children {
		buildHashMap(child, hashmap)
	}
	return hashmap
}

// Find searches for route node by its hash value
func (tr *TreeRoot) Find(hash string) (*Tree, bool) {
	path, ok := tr.hashmap[hash]
	if !ok {
		return nil, false
	}
	node, ok := tr.Rootnode.Find(path)
	if !ok {
		return nil, false
	}
	return node, true
}
