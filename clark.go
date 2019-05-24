/*
clark is a daemon implementing the i3bar protocol.
*/
package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"clark/clarkio"
	"clark/conf"
	"clark/protocol"
)

func main() {
	// Hold our channels
	// clickChannels is a map where the keys are block name/instance
	// the values are the click channels which the individual packages
	// may listen on.
	clickChannels := map[string]chan<- *protocol.Click{}
	blockChannels := []<-chan *protocol.Block{}

	// Build the channels
	for _, block := range conf.AllBlocks {
		c := make(chan *protocol.Click, 4)
		b := make(chan *protocol.Block, 4)

		key := block.Name + "_" + block.Instance
		clickChannels[key] = c

		blockChannels = append(blockChannels, b)

		go clarkio.RunBlock(block.Run, c, b)
	}

	// Start listening on stdin
	go func() {
		clarkio.ReadClicks(os.Stdin, clickChannels)
		fmt.Fprintln(os.Stderr, "routine listening on stdin closed")
	}()

	// Start wrting on stdout
	go func() {
		clarkio.WriteBlocks(os.Stdout, blockChannels)
		fmt.Fprintln(os.Stderr, "routine writing on stdout closed")
	}()

	// Run until we get a SIGINT
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT)
	<-sigs
}
