/*
bat implemnts functions for parsing the files in /sys/class/power_supply/BAT0
It furthermore has some utility functions for doing percentage calculations.
*/
package bat

import (
	"io/ioutil"
	"strconv"
	"strings"
)

const filePath = "/sys/class/power_supply/BAT0"

// GetStatus will read the status file.
// Possible returns are "Charging", "Discharging" and "Unknown"
// All file read errors are returned
func GetStatus() (string, error) {
	status, err := ioutil.ReadFile(filePath + "/status")
	if err != nil {
		return "", err
	}
	return strings.TrimSuffix(string(status[:]), "\n"), nil
}

func parseFloat(p []byte) (float64, error) {
	amount := strings.TrimSuffix(string(p[:]), "\n")

	num, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		return 0, err
	}

	return num, nil
}

func parseFloatFile(filePath string) (float64, error) {
	raw, err := ioutil.ReadFile(filePath)
	if err != nil {
		return 0, err
	}
	return parseFloat(raw)
}

// GetFullCharge will return the number given for
// full charge - as a float64
func GetFullCharge() (float64, error) {
	return parseFloatFile(filePath + "/charge_full")

}

// GetCurrentcharge will return the number given for
// current charge - as a float64
func GetCurrentCharge() (float64, error) {
	return parseFloatFile(filePath + "/charge_now")
}

// GetChargePercentage will get the current
// charge percentage - as a float64
func GetChargePercentage() (float64, error) {
	full, err := GetFullCharge()
	if err != nil {
		return 0, err
	}

	current, err := GetCurrentCharge()
	if err != nil {
		return 0, err
	}

	return (current / full) * 100, nil
}
