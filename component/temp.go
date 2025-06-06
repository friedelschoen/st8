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

func Temperature(block *Block, args map[string]string, not *notify.Notification, cacheptr *any) error {
	if *cacheptr == nil {
		unit, _ := args["unit"]
		switch strings.ToLower(unit) {
		case "", "c", "celsius":
			*cacheptr = TempCelsius
		case "f", "fahrenheit":
			*cacheptr = TempFahrenheit
		case "k", "kelvin":
			*cacheptr = TempKelvin
		}
	}
	unit := (*cacheptr).(TempUnit)

	content, err := os.ReadFile(fmt.Sprintf("/sys/class/thermal/%s/temp", args["sensor"]))
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

	block.OnClick = func(evt ClickEvent) {
		switch unit {
		case TempCelsius:
			*cacheptr = TempFahrenheit
		case TempFahrenheit:
			*cacheptr = TempKelvin
		case TempKelvin:
			*cacheptr = TempCelsius
		}
	}
	return nil
}
