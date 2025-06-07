package component

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/friedelschoen/st8/notify"
)

const (
	urgentAt = 0.15 // 15%
)

type Battery struct {
	Name    string
	Status  string
	Current float64
	Full    float64
	Rate    float64
}

func defaultBatteryName() (string, error) {
	supplies, err := os.ReadDir("/sys/class/power_supply")
	if err != nil {
		return "", err
	}
	if len(supplies) == 0 {
		return "", fmt.Errorf("no supplies found")
	}
	for _, entry := range supplies {
		if entry.Name() == "BAT0" {
			return entry.Name(), nil
		}

		contents, err := os.ReadFile(fmt.Sprintf("/sys/class/power_supply/%s/type", entry.Name()))
		if err != nil {
			return "", err
		}
		if strings.Contains(strings.ToLower(string(contents)), "battery") {
			return entry.Name(), nil
		}
	}
	return "", fmt.Errorf("no supplies found")
}

func GetBattery(name string) (*Battery, error) {
	var bat Battery
	bat.Name = name
	var errs []error
	status, err := os.ReadFile(fmt.Sprintf("/sys/class/power_supply/%s/status", name))
	errs = append(errs, err)
	rate, err := os.ReadFile(fmt.Sprintf("/sys/class/power_supply/%s/power_now", name))
	errs = append(errs, err)
	var now, full []byte
	if _, err := os.Stat(fmt.Sprintf("/sys/class/power_supply/%s/energy_now", name)); err == nil {
		now, err = os.ReadFile(fmt.Sprintf("/sys/class/power_supply/%s/energy_now", name))
		errs = append(errs, err)
		full, err = os.ReadFile(fmt.Sprintf("/sys/class/power_supply/%s/energy_full", name))
		errs = append(errs, err)
	} else {
		now, err = os.ReadFile(fmt.Sprintf("/sys/class/power_supply/%s/charge_now", name))
		errs = append(errs, err)
		full, err = os.ReadFile(fmt.Sprintf("/sys/class/power_supply/%s/charge_full", name))
		errs = append(errs, err)
	}
	if err := errors.Join(errs...); err != nil {
		return nil, err
	}

	bat.Status = strings.TrimSpace(string(status))
	bat.Current, err = strconv.ParseFloat(strings.TrimSpace(string(now)), 64)
	if err != nil {
		return nil, err
	}
	bat.Full, err = strconv.ParseFloat(strings.TrimSpace(string(full)), 64)
	if err != nil {
		return nil, err
	}
	bat.Rate, err = strconv.ParseFloat(strings.TrimSpace(string(rate)), 64)
	if err != nil {
		return nil, err
	}

	return &bat, nil
}

func BatteryState(args map[string]string, events *EventHandlers) (Component, error) {
	return func(block *Block, not *notify.Notification) error {
		name, ok := args["battery"]
		if !ok {
			var err error
			name, err = defaultBatteryName()
			if err != nil {
				return err
			}
		}
		bat, err := GetBattery(name)
		if err != nil {
			return fmt.Errorf("unable to read battery status: %w", err)
		}

		block.Urgent = bat.Current/bat.Full <= urgentAt
		block.Text = bat.Status
		return nil
	}, nil
}

func BatteryPercentage(args map[string]string, events *EventHandlers) (Component, error) {
	return func(block *Block, not *notify.Notification) error {
		name, ok := args["battery"]
		if !ok {
			var err error
			name, err = defaultBatteryName()
			if err != nil {
				return err
			}
		}
		bat, err := GetBattery(name)
		if err != nil {
			return fmt.Errorf("unable to read battery status: %w", err)
		}

		block.Urgent = bat.Current/bat.Full <= urgentAt
		block.Text = fmt.Sprintf("%.0f%%", (bat.Current/bat.Full)*100)
		return nil
	}, nil
}

func BatteryRemaining(args map[string]string, events *EventHandlers) (Component, error) {
	return func(block *Block, not *notify.Notification) error {
		name, ok := args["battery"]
		if !ok {
			var err error
			name, err = defaultBatteryName()
			if err != nil {
				return err
			}
		}
		bat, err := GetBattery(name)
		if err != nil {
			return fmt.Errorf("unable to read battery status: %w", err)
		}

		var hours float64
		switch bat.Status {
		case "Charging":
			hours = (bat.Full - bat.Current) / bat.Rate
		case "Discharging":
			hours = bat.Current / bat.Rate
		default:
			block.Text = ""
			return nil
		}

		block.Urgent = bat.Current/bat.Full <= urgentAt
		block.Text = (time.Duration(hours * float64(time.Hour))).Round(time.Minute).String()
		return nil
	}, nil
}
