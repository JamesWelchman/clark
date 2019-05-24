package cpu

import (
	"fmt"
	"time"

	"clark/colors"
	"clark/protocol"
)

func updateFull(block *protocol.Block, loads []float64) {
	for i, f := range loads {
		if i != 0 {
			block.FullText += " |"
		}

		block.FullText += fmt.Sprintf("cpu%d [%.2f]", i, f)
	}
}

func updateSmall(block *protocol.Block, load float64) {
	block.FullText = fmt.Sprintf("cpu [%.2f]", load)
}

func Run(defaultBlock *protocol.Block, in <-chan *protocol.Click, out chan<- *protocol.Block) error {
	cpuStatus, err := NewCpuStatus()
	if err != nil {
		return fmt.Errorf("couldn't get cpu load :: %v", err)
	}

	displayAll := false
	color := colors.Grey

	throttle := time.After(0)
	for {
		select {
		case <-throttle:
			if cpuStatus == nil {
				continue
			}

			loads, err := cpuStatus.GetLoads()
			if err != nil {
				err = fmt.Errorf("couldn't get cpu load :: %v", err)
				return err
			}
			block := protocol.Block(*defaultBlock)

			if displayAll {
				updateFull(&block, loads[1:])
			} else {
				updateSmall(&block, loads[0])
			}

			block.Color = color

			out <- &block
			throttle = time.After(time.Second)

		case click := <-in:
			// Toggle the color between white and gray
			if click.Button == 1 {
				if color == colors.Grey {
					color = colors.White
				} else {
					color = colors.Grey
				}

				throttle = time.After(0)
				continue
			}

			// Toggle between per-core and average
			if click.Button == 3 {
				displayAll = !displayAll
				throttle = time.After(0)
				continue
			}
		}
	}
}
