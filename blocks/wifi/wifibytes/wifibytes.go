/*
wifi bytes is a package for examing how many bytes are being
transmitted over a wireless network interface.
*/
package wifibytes

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	path = "/proc/net/dev"
)

type (
	// singleRead holds data for a specific read
	// of /proc/net/dev. It holds the total number
	// of bytes sent over the interface from uptime.
	singleRead struct {
		bytesUp   int
		bytesDown int
		readTime  time.Time
	}

	Client struct {
		reads  []*singleRead
		pos    int
		device string
	}
)

func newSingleRead(device string) (*singleRead, error) {
	bytesDown, bytesUp, err := readNetDevFile(path, device)
	if err != nil {
		return nil, err
	}

	return &singleRead{
		readTime:  time.Now(),
		bytesUp:   bytesUp,
		bytesDown: bytesDown,
	}, nil
}

// GetKilobitsPerSecond will return the number of kilobits
// per second which are currently going over the network interface.
// The caller is expected to know how big the internal buffer and
// how long between calls are used.
func (c *Client) GetKilobitsPerSecond() (float64, float64, error) {
	r := c.reads[c.pos]
	down, up, err := r.getKilobitsSecond(c.device)
	if err != nil {
		return 0, 0, err
	}

	c.pos = (c.pos + 1) % len(c.reads)

	return down, up, nil
}

// NewClient returns a new client ready for reading
// numReads is how big the user wishes the buffer
// of previous data points to be.
func NewClient(numReads int, device string) (*Client, error) {
	var reads []*singleRead
	for i := 0; i < numReads; i++ {
		r, err := newSingleRead(device)
		if err != nil {
			return nil, err
		}

		reads = append(reads, r)
	}

	return &Client{
		reads:  reads,
		pos:    0,
		device: device,
	}, nil
}

// getKilobitsSecond will refresh the singleRead instance to
// the current value of bytes over the interface. It will
// return the computed kilobits per second since the last time
// getKilobitsSecond was called.
func (s *singleRead) getKilobitsSecond(device string) (float64, float64, error) {
	bytesDown, bytesUp, err := readNetDevFile(path, device)
	if err != nil {
		return 0, 0, err
	}

	readTime := time.Now()
	period := readTime.Sub(s.readTime).Seconds()

	// Compute how many bits over the wire
	changeDown := float64(bytesDown-s.bytesDown) * 8
	changeUp := float64(bytesUp-s.bytesUp) * 8

	// Compute how many kilobits over the wire
	changeDown = changeDown / 1000
	changeUp = changeUp / 1000

	// Compute kilobites per second
	down := changeDown / period
	up := changeUp / period

	// Set s to the current read
	s.readTime = readTime
	s.bytesUp = bytesUp
	s.bytesDown = bytesDown

	return down, up, nil
}

// readNetDevFile attempts to parse a file for number of
// bytes (up and down) for a given device.
func readNetDevFile(filePath, device string) (int, int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, 0, err
	}
	defer file.Close()
	return readNetworkBytes(file, device)
}

// readNetworkBytes attempts to read the total number of bytes for
// a given device from an io.Reader. It returns a triplet.
// bytesUp, bytesDown and thirdly any errors encountered.
func readNetworkBytes(reader io.Reader, device string) (int, int, error) {
	var err error
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, device) {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 10 {
			return 0, 0, fmt.Errorf("not enough fields in %s for %s",
				path, device)
		}

		bytesDown, err := strconv.Atoi(fields[1])
		if err != nil {
			return 0, 0, err
		}

		bytesUp, err := strconv.Atoi(fields[9])
		if err != nil {
			return 0, 0, err
		}

		return bytesDown, bytesUp, nil
	}

	if err = scanner.Err(); err != nil {
		return 0, 0, err
	}

	/* couldn't find line for this device */
	return 0, 0, fmt.Errorf("couldn't find line for %s", device)
}
