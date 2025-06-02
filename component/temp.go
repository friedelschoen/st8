package component

import (
	"fmt"
	"strings"

	"github.com/friedelschoen/st8/notify"
	"github.com/shirou/gopsutil/v3/host"
)

// temp reads the first temperature sensor whose sensor key contains `sensorFilter`.
// Example input: "coretemp", "CPU", "acpitz", or "" to return the first available.
func Temperature(block *Block, args map[string]string, not *notify.Notification, cache *any) error {
	sensors, err := host.SensorsTemperatures()
	if err != nil {
		return fmt.Errorf("unable to get temperature: %w", err)
	}
	name := strings.ToLower(args["sensor"])

	for _, s := range sensors {
		if name == "" || strings.Contains(strings.ToLower(s.SensorKey), name) {
			block.Text = fmt.Sprintf("%.0f", s.Temperature)
			return nil
		}
	}

	return fmt.Errorf("no matching temperature sensor found")
}
