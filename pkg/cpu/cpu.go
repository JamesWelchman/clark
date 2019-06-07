/*
cpu will examine /proc/stat to estimate how busy the CPU is.
/proc/stat holds data of how much time the CPU has been idle/busy etc. since uptime.

  client, err := cpu.NewClient()
  // .. handle error
  loads, err := client.GetLoads()

loads is a slice of how busy each CPU appears to be.
loads[0] will almost certainly be the average over all cpus.
All subsequent entries will be for each individual CPU.
*/
package cpu

import (
	"bufio"
	"errors"
	"io"
	"os"
	"strconv"
	"strings"
)

const filePath = "/proc/stat"

type (
	cpuTimeStat struct {
		cpu  string
		busy float64
		idle float64
	}

	// Client stores one read of the stat file.
	// It is used to estimate CPU load.
	Client struct {
		cpuReads []*cpuTimeStat
	}
)

// GetLoads() will perform a read of /proc/stat
// It will update the internal buffer with the information from said file.
// It returns how much time the cpus have been busy since the *last* call
// to get loads as a percentage.
// The order of the slice is the same as the order of the cpus in /proc/stat.
// The cpu total will almost certainly be in index 0, all subsequent entries
// are for the individual cores.
func (c *Client) GetLoads() ([]float64, error) {
	cpuReads, err := buildAllCpuTimeStats()
	if err != nil {
		return nil, err
	}

	var loads []float64
	for i, r := range c.cpuReads {
		load, err := computeLoad(r, cpuReads[i])
		if err != nil {
			return nil, err
		}

		loads = append(loads, load)

		c.cpuReads[i] = r
	}

	return loads, nil
}

// NewClient creates a new instance of the client
func NewClient() (*Client, error) {
	cpuReads, err := buildAllCpuTimeStats()
	if err != nil {
		return nil, err
	}

	return &Client{
		cpuReads: cpuReads,
	}, nil
}

func parseCpuLine(reader io.Reader) ([]string, error) {
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)

	var cpus []string
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "cpu") {
			cpus = append(cpus, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return cpus, nil
}

func getCpuLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return parseCpuLine(file)
}

func newCpuTimeStat(cpuLine string) (*cpuTimeStat, error) {
	fields := strings.Fields(cpuLine)
	c := &cpuTimeStat{}

	if len(fields) < 8 {
		return nil, errors.New("not enough fields in cpu line")
	}

	for i, f := range fields {
		if i == 0 {
			c.cpu = f
			continue
		}

		if i == 8 {
			break
		}

		val, err := strconv.ParseFloat(f, 64)
		if err != nil {
			return nil, err
		}

		// The idle field
		if i == 4 {
			c.idle = val
			continue
		}

		// All other fields are "busy"
		c.busy += val
	}

	return c, nil
}

func buildAllCpuTimeStats() ([]*cpuTimeStat, error) {
	var cpuReads []*cpuTimeStat

	cpuLines, err := getCpuLines(filePath)
	if err != nil {
		return nil, err
	}

	for _, line := range cpuLines {
		c, err := newCpuTimeStat(line)
		if err != nil {
			return nil, err
		}

		cpuReads = append(cpuReads, c)
	}

	return cpuReads, nil
}

func computeLoad(c1, c2 *cpuTimeStat) (float64, error) {
	busy := c2.busy - c1.busy
	idle := c2.idle - c1.idle

	ratio := busy / (busy + idle)
	return ratio * 100, nil
}
