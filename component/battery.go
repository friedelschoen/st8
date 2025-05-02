package component

import (
	"fmt"
	"time"

	"github.com/distatus/battery"
	"github.com/friedelschoen/st8/notify"
)

func BatteryState(_ string, _ *notify.Notification, _ *any) (string, error) {
	bat, err := battery.Get(0)
	if err != nil {
		return "", fmt.Errorf("unable to get battery status: %w", err)
	}

	return bat.State.String(), nil
}

func BatteryPercentage(_ string, _ *notify.Notification, _ *any) (string, error) {
	bat, err := battery.Get(0)
	if err != nil {
		return "", fmt.Errorf("unable to get battery status: %w", err)
	}

	perc := bat.Current / bat.Full
	return fmt.Sprintf("%.0f%%", perc), nil
}

func BatteryRemaining(_ string, _ *notify.Notification, _ *any) (string, error) {
	bat, err := battery.Get(0)
	if err != nil {
		return "", fmt.Errorf("unable to get battery status: %w", err)
	}

	var hours float64
	switch bat.State.Raw {
	case battery.Charging:
		hours = (bat.Full - bat.Current) / bat.ChargeRate
	case battery.Discharging:
		hours = bat.Current / bat.ChargeRate
	}
	dur := time.Hour * time.Duration(hours)

	return dur.String(), nil
}
