/*
clarkio implements the bottom layer for both sending and receiving messages
to i3bar.
*/
package clarkio

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"clark/conf"
	"clark/protocol"
)

// ReadClicks implements an event loop and will never return.
// We expect our reader to produce an infinite JSON stream as specified by the
// i3bar protocol. Each time we manage to unmarshal a Click into a protocol.Click
// struct we attempt to send said Click to the relevant channel. The relevant
// channel is found by examining Click.Name and Click.Instance. i.e
//     key := click.Name + "_" + click.Instance
// If this key is found in our map then we write the click event to this channel.
func ReadClicks(reader io.Reader, clickChannels map[string]chan<- *protocol.Click) {
	decoder := json.NewDecoder(reader)

	// Drain the initial '[' token
	err := drainArrayStart(decoder)
	if err != nil {
		/* Already logged */
		return
	}

	for {
		click := &protocol.Click{}

		// This will block until a complete JSON object is avaliable
		// on our reader.
		err = decoder.Decode(click)
		if err != nil {
			logError("failed to decode from stdin", err)
			continue
		}

		key := click.Name + "_" + click.Instance
		in, ok := clickChannels[key]
		if !ok {
			logName("couldn't find %s in channel map", key)
			continue
		}

		in <- click
	}
}

// drainArrayStart takes a json.Decoder and drains the first
// token. We return nil if the first token is '[' - otherwise
// we return an error.
func drainArrayStart(decoder *json.Decoder) error {
	// Take the first [
	tk, err := decoder.Token()
	if err != nil {
		logError("couldn't read from stdin", err)
		return err
	}

	delim, ok := tk.(json.Delim)
	if !ok {
		msg := "first token not a delim"
		log(msg)
		return errors.New(msg)
	}
	// Make sure the value of the delim is '['
	if delim != json.Delim('[') {
		msg := "first token not a ["
		log(msg)
		return errors.New(msg)
	}

	return nil
}

// WriteBlocks implements an event loop and will never return.
// We listen on all the channels in the blockChannels variable
// and write to the writer when updates are required.
func WriteBlocks(writer io.Writer, blockChannels []<-chan *protocol.Block) {

	// lineState is a buffer where we store previously
	// marshaled json, as we probably need to write it
	// for a second time in the near future. Each element of the
	// slice is a byte slice of a marshaled protocol.Block.
	var lineState [][]byte

	// Set every entry to the inital "no data"
	defaultJsonBlock := []byte(conf.DefaultBlockJson)
	for i := 0; i < len(conf.AllBlocks); i++ {
		lineState = append(lineState, defaultJsonBlock)
	}

	// Write the header
	writeWithError(writer, []byte(conf.Header))

	// Buffer our status line updates
	buffer := bufio.NewWriterSize(writer, 2048)

	// Start the infinite array
	writeWithError(buffer, []byte("[\n"))

	// The value of 100 milliseconds has been arrived at
	// by trial and error. It -seems- to keep CPU usage
	// very low and also keeps the bar responsive to user
	// interaction like clicks events.
	numUpdates := 0
	ticker := time.Tick(100 * time.Millisecond)
	for {
		select {
		case <-ticker:

			numUpdates = checkUpdate(lineState, blockChannels)
			if numUpdates > 0 {
				writeStatusLine(buffer, lineState)

				err := buffer.Flush()
				if err != nil {
					logError("failed to flush", err)
				}

				numUpdates = 0
			}
		}
	}
}

func writeWithError(w io.Writer, p []byte) {
	_, err := w.Write(p)
	if err != nil {
		logError("failed to write", err)
	}
}

var (
	openSquare  = []byte("[\n")
	comma       = []byte(",\n")
	closeSquare = []byte("\n],\n")
)

// writeStatusLine will write a complete JSON array to the writer.
// This complete array includes the blocks - one blocks per array element.
// lineState holds the serialized JSON - again one JSON object per element.
func writeStatusLine(writer io.Writer, lineState [][]byte) {
	writeWithError(writer, openSquare)

	for index, block := range lineState {
		if index != 0 {
			writeWithError(writer, comma)
		}
		writeWithError(writer, block)
	}

	writeWithError(writer, closeSquare)
}

// checkUpdate will poll all the Block channels and try to process the updates.
// We return the number of Blocks which have changes in the lineState variable.
func checkUpdate(lineState [][]byte, blockChannels []<-chan *protocol.Block) (numUpdates int) {

	for index, blockChan := range blockChannels {

		newBlock, ok := drainQueue(blockChan)
		if !ok {
			// Failed to drain the queue for this block
			msg := "couldn't drain queue for %s"
			logName(msg, conf.AllBlocks[index].Name)

			// Copy the error throttling block
			// We set a name/instance for it below
			n := protocol.Block(conf.ErrorThrottleBlock)
			newBlock = &n
		}
		if newBlock == nil {
			// No updates for this block
			continue
		}

		newBlock.Name = conf.AllBlocks[index].Name
		newBlock.Instance = conf.AllBlocks[index].Instance

		encodedBlock, err := json.Marshal(newBlock)
		if err != nil {
			logError("failed to marshal block", err)
			lineState[index] = []byte(conf.ErrorBlock)
		} else {
			lineState[index] = encodedBlock
		}
		numUpdates++
	}

	return numUpdates
}

// drainQueue examine a specific Block Channel. If there is no waiting
// data then return nil. Otherwise attempt to drain the channel, conflate
// all entries except the newest and return the protocol.Block instance.
func drainQueue(blockChan <-chan *protocol.Block) (*protocol.Block, bool) {
	// We use len() rather than select - this is because repeatadly calling
	// select on many channels causes the goroutine to start sleeping.
	if len(blockChan) == 0 {
		return nil, true
	}

	// Try to drain the channel - conflate if newer versions exist.
	var newBlock *protocol.Block
	for n := 0; len(blockChan) > 0; n++ {
		newBlock = <-blockChan

		// This is a safeguard for a misbehaving package.
		// We don't want to spend too much time reading the
		// queue of a single block.
		if n == 3 {
			return nil, false
		}
	}

	return newBlock, true
}

// Utility logging functions to write errors
// to stderr
func log(msg string) {
	fmt.Fprintln(os.Stderr, msg)
}

func logError(msg string, err error) {
	log(fmt.Sprintf(msg+" :: %v", err))
}

func logName(msg string, name string) {
	log(fmt.Sprintf(msg, name))
}
