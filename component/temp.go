package component

import (
	"fmt"
	"strings"

	"github.com/shirou/gopsutil/v3/host"
)

// temp reads the first temperature sensor whose sensor key contains `sensorFilter`.
// Example input: "coretemp", "CPU", "acpitz", or "" to return the first available.
func Temperature(name string) (string, error) {
	sensors, err := host.SensorsTemperatures()
	if err != nil {
		return "", fmt.Errorf("unable to get temperature: %w", err)
	}
	name = strings.ToLower(name)

	for _, s := range sensors {
		if name == "" || strings.Contains(strings.ToLower(s.SensorKey), name) {
			return fmt.Sprintf("%.0f", s.Temperature), nil
		}
	}

	return "", fmt.Errorf("no matching temperature sensor found")
}
