package component

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/friedelschoen/st8/notify"
)

type TempUnit int

const (
	TempCelsius TempUnit = iota
	TempFahrenheit
	TempKelvin
)

func temperature(args map[string]string, events *EventHandlers) (Component, error) {
	sensor, ok := args["sensor"]
	if !ok {
		return nil, fmt.Errorf("missing argument: sensor")
	}

	unit := TempCelsius
	events.OnClick = func(evt ClickEvent) {
		switch unit {
		case TempCelsius:
			unit = TempFahrenheit
		case TempFahrenheit:
			unit = TempKelvin
		case TempKelvin:
			unit = TempCelsius
		}
	}

	return func(block *Block, not *notify.Notification) error {
		content, err := os.ReadFile(fmt.Sprintf("/sys/class/thermal/%s/temp", sensor))
		if err != nil {
			return err
		}
		cels1000, err := strconv.Atoi(strings.TrimSpace(string(content)))
		switch unit {
		case TempCelsius:
			block.Text = fmt.Sprintf("%.1f °C", float64(cels1000)/1000)
		case TempFahrenheit:
			// °F = (°C × 1.8) + 32
			block.Text = fmt.Sprintf("%.1f °F", (float64(cels1000)*1.8/1000)+32)
		case TempKelvin:
			block.Text = fmt.Sprintf("%.1f K", (float64(cels1000)/1000)-274.15)
		}
		return nil
	}, nil
}

func init() {
	Install("temp", temperature)
}
