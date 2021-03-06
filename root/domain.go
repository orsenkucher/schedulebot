// Package root represents root domain
package root

import (
	"path/filepath"
	"strings"
)

// Schedule represents domain specific schedule
type Schedule struct {
	Name    string
	Events  []string
	Timings []int
	Spin    []int
}

// Subscriber is subscribed to at least one Schedule
type Subscriber int64

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

// Rootdir path to the root directory
const Rootdir = "data"

// SchFile path to file with old schedule
const SchFile = "fbclient/sch.json"

// PopExt removes file extension from file name
func PopExt(name string) string {
	return strings.TrimSuffix(name, filepath.Ext(name))
}
