package clarkio

import (
	"clark/blocks"
	"clark/conf"
	"clark/protocol"
)

func RunBlock(run blocks.RunFunc, c <-chan *protocol.Click, b chan<- *protocol.Block) {
	for {
		err := run(conf.NewBlock(), c, b)
		logError("block stopped", err)
	}
}
