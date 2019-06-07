package wifi

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/jameswelchman/clark/colors"
	"github.com/jameswelchman/clark/protocol"

	"github.com/jameswelchman/clark/blocks/wifi/wifibytes"
)

const (
	device     = "wlp2s0"
	textLength = len("down[9999.99 kbs] up[99.99]")
)

type runDetails struct {
	DefaultBlock *protocol.Block
	ClickChannel <-chan *protocol.Click
	BlockChannel chan<- *protocol.Block
	Down         float64
	Up           float64
	Color        string
	err          error
}

func (r *runDetails) ToggleColor() {
	if r.Color == colors.Grey {
		r.Color = colors.White
	} else {
		r.Color = colors.Grey
	}
}

func (r *runDetails) SendConnected() {
	block := protocol.Block(*r.DefaultBlock)

	// Set the text
	text := fmt.Sprintf("down[%.2f kbs] up[%.2f kbs]", r.Down, r.Up)
	if shortLength := textLength - len(text); shortLength > 0 {
		text += strings.Repeat(" ", shortLength)
	}

	block.FullText = text
	block.MinWidth = textLength
	block.Color = r.Color
	r.BlockChannel <- &block
}

func (r *runDetails) SendNotConnected() {
	block := protocol.Block(*r.DefaultBlock)
	block.FullText = "No Connection"
	block.Color = colors.Red
	r.BlockChannel <- &block
}

func (r *runDetails) RecentTraffic() bool {
	return r.Up != 0 || r.Down != 0
}

type stateFn func(*runDetails, *wifibytes.Client) stateFn

func notConnected(r *runDetails, c *wifibytes.Client) stateFn {
	r.SendNotConnected()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		var err error
		select {
		case <-ticker.C:
			r.Down, r.Up, err = c.GetKilobitsPerSecond()
			if err != nil {
				r.err = fmt.Errorf("couldn't get speed :: %v", err)
				return nil
			}

			if r.RecentTraffic() {
				return connected
			}

		case <-r.ClickChannel:
			continue
		}
	}
}

func connected(r *runDetails, c *wifibytes.Client) stateFn {
	r.SendConnected()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		var err error
		select {
		case <-ticker.C:
			r.Down, r.Up, err = c.GetKilobitsPerSecond()
			if err != nil {
				r.err = fmt.Errorf("couldn't get speeds :: %v", err)
				return nil
			}
			r.SendConnected()

			if !r.RecentTraffic() {
				return testConnection
			}

		case click := <-r.ClickChannel:
			if click.Button != 1 {
				continue
			}
			r.ToggleColor()
			r.SendConnected()
		}
	}
}

func testConnection(r *runDetails, c *wifibytes.Client) stateFn {
	errCh := make(chan error)
	go func() {
		cmd := exec.Command("ping", "-c", "1", "8.8.8.8")
		errCh <- cmd.Run()
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for n := 0; n < 20; n++ {
		var err error

		select {
		case <-ticker.C:
			r.Down, r.Up, err = c.GetKilobitsPerSecond()
			if err != nil {
				r.err = fmt.Errorf("couldn't get speeds :: %v", err)
				return nil
			}

			if r.RecentTraffic() {
				return connected
			}

		case click := <-r.ClickChannel:
			if click.Button != 1 {
				continue
			}
			r.ToggleColor()
			r.SendConnected()
		case err = <-errCh:
			if err != nil {
				return notConnected
			}
			// err == nil implies ping exited with status 0
			// This only happens if it was succesful
			return connected
		}
	}

	// Still no bytes after twenty seconds.
	// Ping has not exited after twenty seconds
	// Assume we've disconnected
	return notConnected
}

func Run(defaultBlock *protocol.Block, in <-chan *protocol.Click, out chan<- *protocol.Block) error {
	c, err := wifibytes.NewClient(10, device)
	if err != nil {
		err = fmt.Errorf("couldn't create client :: %v", err)
		return err
	}

	r := &runDetails{
		DefaultBlock: defaultBlock,
		ClickChannel: in,
		BlockChannel: out,
		Color:        colors.Grey,
	}

	state := notConnected
	for {
		state = state(r, c)
		if r.err != nil {
			return r.err
		}
	}
}
