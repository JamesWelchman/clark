package clarkio

import (
	"github.com/jameswelchman/clark/blocks"
	"github.com/jameswelchman/clark/conf"
	"github.com/jameswelchman/clark/protocol"
)

func RunBlock(run blocks.RunFunc, c <-chan *protocol.Click, b chan<- *protocol.Block) {
	for {
		err := run(conf.NewBlock(), c, b)
		logError("block stopped", err)
	}
}
