package memory

import (
	"bufio"
	"errors"
	"os"
	"strconv"
	"strings"
)

const path = "/proc/meminfo"

func getLines(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	if err != nil {
		return nil, err
	}

	outPut := []string{"", ""}

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "MemTotal") {
			outPut[0] = line
		}

		if strings.HasPrefix(line, "MemAvailable") {
			outPut[1] = line
		}

		if len(outPut[0]) > 0 && len(outPut[1]) > 0 {
			return outPut, nil
		}
	}

	if err = scanner.Err(); err != nil {
		return nil, err
	}

	return nil, errors.New("couldn't find MemTotal and MemAvaiilable")
}

func parseMemTotal(memTotal string) (float64, error) {
	fields := strings.Split(memTotal, ":")
	if len(fields) != 2 {
		return 0, errors.New("couldn't parse MemTotal")
	}

	fields[1] = strings.TrimLeft(fields[1], "\t ")

	memory := strings.Split(fields[1], " ")
	if len(fields) != 2 {
		return 0, errors.New("field MemTotal missing units")
	}

	if memory[1] != "kB" {
		return 0, errors.New("field MemTotal not in kB")
	}

	numKB, err := strconv.ParseFloat(memory[0], 64)
	if err != nil {
		return 0, errors.New("field MemTotal couldn't parse float")
	}

	return (numKB / 1048576), nil
}

func parseMemAvailable(memAvailable string) (float64, error) {
	fields := strings.Split(memAvailable, ":")
	if len(fields) != 2 {
		return 0, errors.New("couldn't parse MemAvailable")
	}

	fields[1] = strings.TrimLeft(fields[1], "\t ")

	memory := strings.Split(fields[1], " ")
	if len(fields) != 2 {
		return 0, errors.New("field MemAvailable missing units")
	}

	if memory[1] != "kB" {
		return 0, errors.New("field MemAvailable not in kB")
	}

	numKB, err := strconv.ParseFloat(memory[0], 64)
	if err != nil {
		return 0, errors.New("field MemAvailable couldn't parse float")
	}

	return (numKB / 1048576), nil
}

func getMemory() (float64, float64, error) {
	memInfo, err := getLines(path)
	if err != nil {
		return 0, 0, err
	}

	memTotal, err := parseMemTotal(memInfo[0])
	if err != nil {
		return 0, 0, err
	}

	memAvailable, err := parseMemTotal(memInfo[1])
	if err != nil {
		return 0, 0, err
	}

	return memAvailable, memTotal, nil
}
