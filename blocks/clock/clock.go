/*
clock is a package implementing a clock display on i3bar
*/
package clock

import (
	"time"

	"clark/colors"
	"clark/protocol"
)

func currentTime() string {
	return time.Now().Format("Mon 2-Jan-2006 15:04")
}

// Run will write the current time once per second
func Run(defaultBlock *protocol.Block, in <-chan *protocol.Click, out chan<- *protocol.Block) {
	for {
		select {
		case <-time.After(time.Second):
			block := protocol.Block(*defaultBlock)
			block.FullText = currentTime()
			block.Color = colors.Green
			out <- &block
		case <-in:
			continue
		}
	}
}
