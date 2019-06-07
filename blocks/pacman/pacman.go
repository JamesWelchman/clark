package pacman

import (
	"fmt"
	"time"

	"github.com/jameswelchman/clark/colors"
	pacmanClient "github.com/jameswelchman/clark/pkg/pacman"
	"github.com/jameswelchman/clark/protocol"
)

func sendZero(defaultBlock *protocol.Block, out chan<- *protocol.Block) {
	block := protocol.Block(*defaultBlock)

	block.Color = colors.Grey
	block.FullText = "[0]"

	out <- &block
}

func Run(defaultBlock *protocol.Block, in <-chan *protocol.Click, out chan<- *protocol.Block) error {
	sendZero(defaultBlock, out)

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	numUpdates := 0

	for {
		select {
		case <-ticker.C:
			newUpdates, err := pacmanClient.NumUpdates()
			if err != nil {
				return fmt.Errorf("couldn't get pacman updates :: %v", err)
			}

			// Don't update for no reason
			if newUpdates == numUpdates {
				continue
			}

			block := protocol.Block(*defaultBlock)
			block.Color = colors.Grey
			if newUpdates > 0 {
				block.Color = colors.Blue
			}
			block.FullText = fmt.Sprintf("[%d]", newUpdates)

			out <- &block
			numUpdates = newUpdates
		case <-in:
			continue
		}
	}
}
