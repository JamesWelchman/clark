package battery

import (
	"io/ioutil"
	"strconv"
	"strings"
)

const path = "/sys/class/power_supply/BAT0"

func getStatus() (string, error) {
	status, err := ioutil.ReadFile(path + "/status")
	if err != nil {
		return "", err
	}
	return strings.TrimSuffix(string(status[:]), "\n"), nil
}

func getFullCharge() (float64, error) {
	raw, err := ioutil.ReadFile(path + "/charge_full")
	if err != nil {
		return 0.0, err
	}
	chargeString := strings.TrimSuffix(string(raw[:]), "\n")
	fullCharge, err := strconv.ParseFloat(chargeString, 64)
	if err != nil {
		return 0.0, err
	}

	return fullCharge, nil
}

func getCurrentCharge() (float64, error) {
	raw, err := ioutil.ReadFile(path + "/charge_now")
	if err != nil {
		return 0.0, err
	}
	chargeString := strings.TrimSuffix(string(raw[:]), "\n")
	chargeNow, err := strconv.ParseFloat(chargeString, 64)
	if err != nil {
		return 0.0, err
	}

	return chargeNow, nil
}
