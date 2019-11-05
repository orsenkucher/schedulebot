package route

// TreeCreator is capable of creating route tree
// from filesystem or firestore structure or whatever
type TreeCreator interface {
	Create() *Tree
}
