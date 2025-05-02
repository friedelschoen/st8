package component

import (
	"fmt"

	"github.com/friedelschoen/st8/notify"
)

type Component func(arg string, not *notify.Notification) (res string, err error)

var Functions = map[string]Component{
	"battery_state":     BatteryState,
	"battery_perc":      BatteryPercentage,
	"battery_remaining": BatteryRemaining,
	"cat":               ReadFile,
	"cpu_freq":          CPUFrequency,
	"cpu_perc":          CPUPercentage,
	"datetime":          Datetime,
	"disk_free":         DiskFree,
	"disk_perc":         DiskPercentage,
	"disk_total":        DiskTotal,
	"disk_used":         DiskUsed,
	"entropy":           EntropyAvailable,
	"gid":               Gid,
	"hostname":          Hostname,
	"ipv4":              IPv4,
	"ipv6":              IPv6,
	"kernel_release":    KernelRelease,
	"load_avg":          LoadAverage,
	"netspeed_rx":       NetspeedRx,
	"netspeed_tx":       NetspeedTx,
	"notify_appname":    NotifyAppName,
	"notify_appicon":    NotifyAppIcon,
	"notify_summary":    NotifySummary,
	"notify_body":       NotifyBody,
	"notify_actions":    NotifyActions,
	"period_command":    PeriodCommand,
	"num_files":         NumFiles,
	"ram_free":          RamFree,
	"ram_perc":          RamPercentage,
	"ram_total":         RamTotal,
	"ram_used":          RamUsed,
	"run_command":       RunCommand,
	"swap_free":         SwapFree,
	"swap_perc":         SwapPercentage,
	"swap_total":        SwapTotal,
	"swap_used":         SwapUsed,
	"temp":              Temperature,
	"uid":               Uid,
	"up":                Up,
	"uptime":            Uptime,
	"username":          Username,
	"wifi_essid":        WifiESSID,
	"wifi_perc":         WifiPerc,
}

// fmtHuman formats bytes to a human-readable string, e.g. "1.4 GiB"
func fmtHuman(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := unit, 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
