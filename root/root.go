// Package root represents root domain
package root

// Schedule represents domain specific schedule
type Schedule struct {
	Name    string
	Events  []string
	Timings []int
	Spin    []int
}

// SubEvent represents subscription event
type SubEvent struct {
	SubID  int64
	Action SubAction
}

// SubAction represents user action
type SubAction int

// Add is when user Success
// Del is when user unsubbed
const (
	_ SubAction = iota
	Add
	Del
)

// MPW is total minutes in week
const MPW = 7 * 60 * 24
