package component

import (
	"fmt"
	"time"

	"github.com/distatus/battery"
	"github.com/friedelschoen/st8/notify"
)

func BatteryState(block *Block, args map[string]string, not *notify.Notification, cache *any) error {
	bat, err := battery.Get(0)
	if err != nil {
		return fmt.Errorf("unable to get battery status: %w", err)
	}

	block.Text = bat.State.String()
	return nil
}

func BatteryPercentage(block *Block, args map[string]string, not *notify.Notification, cache *any) error {
	bat, err := battery.Get(0)
	if err != nil {
		return fmt.Errorf("unable to get battery status: %w", err)
	}

	perc := bat.Current / bat.Full
	block.Text = fmt.Sprintf("%.0f%%", perc)
	return nil
}

func BatteryRemaining(block *Block, args map[string]string, not *notify.Notification, cache *any) error {
	bat, err := battery.Get(0)
	if err != nil {
		return fmt.Errorf("unable to get battery status: %w", err)
	}

	var hours float64
	switch bat.State.Raw {
	case battery.Charging:
		hours = (bat.Full - bat.Current) / bat.ChargeRate
	case battery.Discharging:
		hours = bat.Current / bat.ChargeRate
	}
	dur := time.Hour * time.Duration(hours)

	block.Text = dur.String()
	return nil
}
