/*
blocks is the parent package to all block implementations.
*/
package blocks

import (
	"clark/protocol"
)

// RunFunc is the function signature which must be exported by individual packages.
// Implementations of RunFunc should not return and must do event looping themselves.
// Each instance of RunFunc is expected to take ownership and may freely modify the
// the variables passed to it.
type RunFunc func(*protocol.Block, <-chan *protocol.Click, chan<- *protocol.Block) error

// Block is the structure used by clark/conf/conf.go to specify
// exactly one block on i3bar.
type Block struct {
	// Name and Instance uniquely identify this block
	Name     string
	Instance string

	// Run is called in it's own goroutine and should not return.
	// Each package must implement it's own version.
	Run RunFunc
}
