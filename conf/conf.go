/*
conf is where global configuration may be set. See individual declarations
for details. This file was inspired by suckless config.h - full credit to them!
*/
package conf

import (
	"github.com/jameswelchman/clark/blocks"
	"github.com/jameswelchman/clark/blocks/battery"
	"github.com/jameswelchman/clark/blocks/clock"
	"github.com/jameswelchman/clark/blocks/cpu"
	"github.com/jameswelchman/clark/blocks/memory"
	"github.com/jameswelchman/clark/blocks/wifi"
	"github.com/jameswelchman/clark/colors"
	"github.com/jameswelchman/clark/protocol"
)

// Header is a string which must be the first line sent to i3bar as
// specified by the i3bar protocol.
const Header = `{"version": 1, "stop_signal": 10, "cont_signal": 12, "click_events": true}` + "\n"

// DefaultBlockJson is what we populate block entries with before
// we have any data from the corresponding block run function.
const DefaultBlockJson = `{"full_text": "no data"}`

// ErrorBlock is used to populate the block in i3bar when we fail
// to marshal to json the variable sent to us by a block run function.
const ErrorBlock = `{"full_text": "error"}`

// DefaultBlock is an instance of protocol.Block which specifies
// global defaults. All instances of running Blocks are given
// their own unique copy of the variable.
var DefaultBlock = protocol.Block{
	Color:     colors.Grey,
	MinWidth:  5,
	Align:     "right",
	Urgent:    false,
	Separator: true,
	Markup:    "none",
}

// ErrorThrottleBlock is used when a package is writing too many
// updates so the writer can't keep up. See clark/clarkio for details.
var ErrorThrottleBlock = protocol.Block{
	Color:     colors.Red,
	MinWidth:  5,
	Align:     "right",
	Urgent:    true,
	Separator: true,
	Markup:    "none",
	FullText:  "ERROR - block writing too many updates",
}

// NewBlock creates a new instance of protocol.Block with some global
// defaults already set. The global defaults are specified by editing
// the DefaultBlock variable.
func NewBlock() *protocol.Block {
	block := protocol.Block(DefaultBlock)
	return &block
}

// AllBlocks is an array of all blocks which we wish to include in our final
// output to i3bar. Furthermore this where we define the order in which they
// will appear in the final display. Adding and removing entries to this array
// is required/sufficient to activate/deactive a particular block.
var AllBlocks = [...]*blocks.Block{
	&blocks.Block{
		Name:     "clock",
		Instance: "1",
		Run:      clock.Run,
	},
	&blocks.Block{
		Name:     "wifi",
		Instance: "1",
		Run:      wifi.Run,
	},
	&blocks.Block{
		Name:     "memory",
		Instance: "1",
		Run:      memory.Run,
	},
	&blocks.Block{
		Name:     "cpu",
		Instance: "1",
		Run:      cpu.Run,
	},
	&blocks.Block{
		Name:     "battery",
		Instance: "1",
		Run:      battery.Run,
	},
}
