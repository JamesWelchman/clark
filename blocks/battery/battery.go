package battery

import (
	"fmt"
	"time"

	"clark/colors"
	"clark/pkg/bat"
	"clark/protocol"
)

type runInfo struct {
	color         string
	chargePercent float64
	status        string
}

func (r *runInfo) BuildBlock(defaultBlock *protocol.Block) *protocol.Block {
	block := protocol.Block(*defaultBlock)

	// Color
	block.Color = r.color
	if r.status == "Discharging" {
		block.Color = colors.White

		if r.chargePercent < 40 {
			block.Color = colors.Yellow
		}

		if r.chargePercent < 10 {
			block.Color = colors.Red
		}
	}

	// Text
	block.FullText = fmt.Sprintf("%s %.0f%%", r.status, r.chargePercent)
	return &block
}

func (r *runInfo) Update() error {
	var err error

	r.status, err = bat.GetStatus()
	if err != nil {
		return fmt.Errorf("couldn't get battery status :: %v", err)
	}

	r.chargePercent, err = bat.GetChargePercentage()
	if err != nil {
		return fmt.Errorf("couldn't get battery charge :: %v", err)
	}

	return nil
}

func (r *runInfo) ToggleColor() {
	if r.color == colors.Grey {
		r.color = colors.White
	} else {
		r.color = colors.Grey
	}
}

func Run(defaultBlock *protocol.Block, in <-chan *protocol.Click, out chan<- *protocol.Block) error {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	run := &runInfo{
		color: colors.Grey,
	}

	for {
		select {
		case <-ticker.C:
			err := run.Update()
			if err != nil {
				return err
			}

			block := run.BuildBlock(defaultBlock)
			out <- block

		case click := <-in:
			if click.Button != 1 {
				continue
			}

			run.ToggleColor()
			block := run.BuildBlock(defaultBlock)
			out <- block
		}
	}
}
