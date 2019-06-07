/*
memory implemnts a function for reading /proc/meminfo.
This contains data of how much virtual memory is used/free etc.

  memAvailable, memTotal, err := memory.GetMemory()

memAvailable and memTotal are float64 and the units are GigaBytes.
*/
package memory

import (
	"bufio"
	"errors"
	"os"
	"strconv"
	"strings"
)

const filePath = "/proc/meminfo"

// GetMemory will read /proc/meminfo and return the available
// and used virtual memory.
func GetMemory() (memAvailable, memTotal float64, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, 0, err
	}
	defer file.Close()

	// First loop over the file line by line
	// The lines we want start "MemToal" and "MemAvailable"
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	var memTotalStr string
	var memAvailStr string
	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "MemTotal") {
			memTotalStr = line
		}

		if strings.HasPrefix(line, "MemAvailable") {
			memAvailStr = line
		}

		if len(memAvailStr) > 0 && len(memTotalStr) > 0 {
			goto PARSE_MEMORY_LINES
		}
	}

	// We didn't manage to find the line
	return 0, 0, errors.New("missing lines in /proc/meminfo")

PARSE_MEMORY_LINES:
	// We have a line like "MemTotal:  1619651 kB"
	// just pull the numbers out.
	memTotalStr, err = parseLine(memTotalStr)
	if err != nil {
		return 0, 0, err
	}
	memAvailStr, err = parseLine(memAvailStr)
	if err != nil {
		return 0, 0, err
	}

	// Parse these numbers into floats
	memTotal, err = strconv.ParseFloat(memTotalStr, 64)
	if err != nil {
		return 0, 0, err
	}

	memAvailable, err = strconv.ParseFloat(memAvailStr, 64)
	if err != nil {
		return 0, 0, err
	}

	return
}

func parseLine(line string) (string, error) {
	// We have a line like "MemTotal:  1619651 kB"
	// First split on the colon
	fields := strings.Split(line, ":")
	if len(fields) != 2 {
		return "", errors.New("couldn't parse memory line")
	}

	amount := strings.TrimLeft(fields[1], "\t ")

	// Next split on the space between the amout and the units
	fields = strings.Split(amount, ":")
	if fields[1] != "kB" {
		return "", errors.New("memory not in kB")
	}

	return fields[0], nil
}
