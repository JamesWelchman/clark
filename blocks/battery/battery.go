package battery

import (
	"fmt"
	"time"

	"clark/colors"
	"clark/protocol"
)

func updateCharge(block *protocol.Block, chargeFull float64) error {
	// status
	status, err := getStatus()
	if err != nil {
		return fmt.Errorf("unable to get battery status %v\n", err)
	}

	// charge now
	chargeNow, err := getCurrentCharge()
	if err != nil {
		return fmt.Errorf("unable to get battery status %v\n", err)
	}

	// percentage
	chargePercentage := calcChargePercentage(chargeNow, chargeFull)

	// Color
	if status == "Discharging" {
		block.Color = colors.White

		if chargePercentage < 40 {
			block.Color = colors.Yellow
		}

		if chargePercentage < 10 {
			block.Color = colors.Red
		}
	}

	block.FullText = fmt.Sprintf("%s %d%%", status, chargePercentage)
	return nil
}

func Run(defaultBlock *protocol.Block, in <-chan *protocol.Click, out chan<- *protocol.Block) error {
	defaultColor := colors.Grey

	// Populate chargeFull
	chargeFull, err := getFullCharge()
	if err != nil {
		err = fmt.Errorf("unable to get full battery charge :: %v", err)
		return err
	}

	throttle := time.After(0)
	for {
		select {
		case <-throttle:
			block := protocol.Block(*defaultBlock)
			block.Color = defaultColor

			if err := updateCharge(&block, chargeFull); err != nil {
				return err
			}

			out <- &block
			throttle = time.After(time.Second)

		case click := <-in:
			if click.Button != 1 {
				continue
			}

			if defaultColor == colors.Grey {
				defaultColor = colors.White
			} else {
				defaultColor = colors.Grey
			}

			throttle = time.After(0)
		}
	}
}

func calcChargePercentage(now, full float64) int {
	percentage := (now / full) * 100
	return int(percentage)
}
