package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func notify(title, message, urgency string) {
	var cmd *exec.Cmd
	if runtime.GOOS == "linux" {
		cmd = exec.Command("notify-send", "-u", urgency, "-t", "5000", title, message)
	} else {
		fmt.Printf("Unsupported OS for dunst: %s\n", runtime.GOOS)
		return
	}

	err := cmd.Run()
	if err != nil {
		log.Printf("Error sending notification: %v", err)
	}
}

// Auxiliary functions to help out with parsing the data from files
func readIntFromFile(filepath string) (int, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(strings.TrimSpace(string(data)))
}

func readFileToString(filepath string) (string, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// Manaing battery status with help from auxiliary functions
func checkBatteryStatus() {
	// The locations of the files which contain the data to battery may vary.
	// Hence, isolating them into separate variables for easy changes.
	current_enery_file := "/sys/class/power_supply/BAT0/energy_now"
	energy_full_file := "/sys/class/power_supply/BAT0/energy_full"
	charging_status_file := "/sys/class/power_supply/BAT0/status"

	chargeNow, err := readIntFromFile(current_enery_file)
	if err != nil {
		log.Println("Error reading energy_now:", err)
		return
	}
	chargeFull, err := readIntFromFile(energy_full_file)
	if err != nil {
		log.Println("Error reading energy_full:", err)
		return
	}

	status, err := readFileToString(charging_status_file)
	if err != nil {
		log.Println("Error reading status:", err)
		return
	}

	percent := int(float64(chargeNow) / float64(chargeFull) * 100)
	isCharging := strings.TrimSpace(status) == "Charging"

	// Debugging statement
	// notify("Reporting battery percentage", fmt.Sprintf("Battery level: %d%%", percent), "low")

	if percent <= 20 {
		notify("Battery Low", fmt.Sprintf("Battery level is %d%%. Please charge your device.",
			percent), "critical")
	}
	if isCharging && percent >= 90 {
		notify("Battery Charged Sufficiently.", fmt.Sprintf("Battery level has reached %d%%. Turn off charging.",
			percent), "normal")
	}

	if isCharging && percent >= 95 {
		notify("Battery Overflow.", fmt.Sprintf("Battery incurring damage"), "critical")
	}
}

func main() {
	for {
		checkBatteryStatus()
		time.Sleep(10 * time.Second)
	}
}
