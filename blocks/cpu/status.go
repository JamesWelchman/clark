package cpu

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

const path = "/proc/stat"

type cpuTimeStat struct {
	Cpu     string
	User    float64
	Nice    float64
	System  float64
	Idle    float64
	IoWait  float64
	Irq     float64
	SoftIrq float64
}

func (c *cpuTimeStat) Total() float64 {
	return c.User + c.Nice + c.System + c.IoWait + c.Irq + c.SoftIrq
}

// getCpuLines will read /proc/stat
// grep ^cpu /proc/stat
func getCpuLines(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
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

func getCpuTimeStat(cpuLine string) (*cpuTimeStat, error) {
	fields := strings.Fields(cpuLine)
	cpuTimes := &cpuTimeStat{}

	for i, f := range fields {
		if i == 0 {
			cpuTimes.Cpu = f
			continue
		}

		val, err := strconv.ParseFloat(f, 64)
		if err != nil {
			return nil, err
		}

		switch i {
		case 1:
			cpuTimes.User = val
		case 2:
			cpuTimes.Nice = val
		case 3:
			cpuTimes.System = val
		case 4:
			cpuTimes.Idle = val
		case 5:
			cpuTimes.IoWait = val
		case 6:
			cpuTimes.Irq = val
		case 7:
			cpuTimes.SoftIrq = val
		default:
			break
		}
	}

	return cpuTimes, nil
}

func computeLoad(c1, c2 *cpuTimeStat) float64 {
	t1Total := c1.Total()
	t2Total := c2.Total()

	// totalDiff is the amount of time in Jiffies
	// that the core has been processing
	totalDiff := t2Total - t1Total

	// total time in Jiffies the core
	// has been idle
	idleDiff := c2.Idle - c1.Idle

	// Return the percentage active
	per := totalDiff / (totalDiff + idleDiff)
	return per * 100
}

func buildAllCpuStats() ([]*cpuTimeStat, error) {
	cpuLines, err := getCpuLines(path)
	if err != nil {
		return nil, err
	}

	var cpuTimes []*cpuTimeStat
	for _, cpuLine := range cpuLines {
		cts, err := getCpuTimeStat(cpuLine)
		if err != nil {
			return nil, err
		}

		cpuTimes = append(cpuTimes, cts)
	}

	return cpuTimes, nil
}

type cpuStatus struct {
	cpuTimes []*cpuTimeStat
}

func NewCpuStatus() (*cpuStatus, error) {
	cpuTimes, err := buildAllCpuStats()
	if err != nil {
		return nil, err
	}

	return &cpuStatus{
		cpuTimes: cpuTimes,
	}, nil
}

func (self *cpuStatus) GetLoads() ([]float64, error) {
	cpuTimes, err := buildAllCpuStats()
	if err != nil {
		return nil, err
	}

	var loads []float64
	for i := 0; i < len(cpuTimes); i++ {
		load := computeLoad(self.cpuTimes[i], cpuTimes[i])
		loads = append(loads, load)
	}

	self.cpuTimes = cpuTimes

	return loads, nil
}
