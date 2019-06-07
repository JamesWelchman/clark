package cpu

import (
	"fmt"
	"strings"
	"time"

	"clark/colors"
	cpuClient "clark/pkg/cpu"
	"clark/protocol"
)

const (
	// Length of cpu [25.00]
	shortLength = len("cpu [25.00]")
	longLength  = len("cpu [25.00] | cpu [25.00] | cpu [25.00] | cpu [25.00]")
)

type runInfo struct {
	displayAll bool
	color      string
	loads      []float64
	client     *cpuClient.Client
}

func (r *runInfo) Update() error {
	var err error
	r.loads, err = r.client.GetLoads()
	if err != nil {
		return fmt.Errorf("couldn't get cpu loads :: %v", err)
	}

	return nil
}

func (r *runInfo) BuildBlock(defaultBlock *protocol.Block) *protocol.Block {
	block := protocol.Block(*defaultBlock)

	if r.displayAll {
		for i, f := range r.loads {
			if i != 0 {
				block.FullText += " |"
			}

			block.FullText += fmt.Sprintf("cpu%d [%.2f]", i, f)
		}
		if textLength := longLength - len(block.FullText); textLength > 0 {
			block.FullText += strings.Repeat(" ", textLength)
		}
	} else {
		block.FullText = fmt.Sprintf("cpu [%.2f]", r.loads[0])
		if textLength := shortLength - len(block.FullText); textLength > 0 {
			block.FullText += strings.Repeat(" ", textLength)
		}
	}

	block.Color = r.color

	return &block
}

func (r *runInfo) ToggleColor() {
	if r.color == colors.Grey {
		r.color = colors.White
	} else {
		r.color = colors.Grey
	}
}

func Run(defaultBlock *protocol.Block, in <-chan *protocol.Click, out chan<- *protocol.Block) error {
	client, err := cpuClient.NewClient()
	if err != nil {
		return fmt.Errorf("couldn't get cpu loads :: %v", err)
	}

	run := runInfo{
		client: client,
		color:  colors.Grey,
	}

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			err = run.Update()
			if err != nil {
				return err
			}

			block := run.BuildBlock(defaultBlock)
			out <- block
		case click := <-in:
			if click.Button == 1 {
				run.ToggleColor()
			} else if click.Button == 3 {
				run.displayAll = !run.displayAll
			} else {
				// Don't send an update if nothing changed
				continue
			}

			block := run.BuildBlock(defaultBlock)
			out <- block
		}
	}
}
