package memory

import (
	"fmt"
	"time"

	"clark/colors"
	"clark/protocol"
)

func update(block *protocol.Block) error {
	free, total, err := getMemory()
	if err != nil {
		return err
	}
	used := total - free
	perc := (used / total) * 100
	txt := fmt.Sprintf("Mem %.1f GB / %.1f GB [%.2f%%]", used, total, perc)
	block.FullText = txt

	return nil
}

func Run(defaultBlock *protocol.Block, in <-chan *protocol.Click, out chan<- *protocol.Block) error {
	color := colors.Grey

	throttle := time.After(0)
	for {
		select {
		case <-throttle:
			block := protocol.Block(*defaultBlock)
			if err := update(&block); err != nil {
				err = fmt.Errorf("couldn't update memory :: %v", err)
				return err
			}
			block.Color = color
			out <- &block
			throttle = time.After(time.Second)

		case click := <-in:
			if click.Button != 1 {
				continue
			}

			if color == colors.Grey {
				color = colors.White
			} else {
				color = colors.Grey
			}

			throttle = time.After(0)
		}
	}
}
