# st8 Components

## `battery_state`
* **Arguments**: `battery` (e.g., `BAT0`)
* **Returns**: Current battery status (Charging, Discharging, Full, etc.)
* **Example Output**: `Charging`
* **Reads from**: `/sys/class/power_supply/<battery>/status`
* **Urgency**: Yes, when battery percentage <= 15%

## `battery_perc`
* **Arguments**: `battery`
* **Returns**: Battery percentage as integer
* **Example Output**: `87%`
* **Reads from**:
  * `/sys/class/power_supply/<battery>/energy_now` or `/charge_now`
  * `/sys/class/power_supply/<battery>/energy_full` or `/charge_full`
* **Urgency**: Yes, when battery percentage <= 15%

## `battery_remaining`
* **Arguments**: `battery`
* **Returns**: Remaining battery time (e.g., `1h32m`)
* **Example Output**: `2h45m`
* **Reads from**: Same as `battery_perc` + `/power_now`
* **Urgency**: Yes, when battery percentage <= 15%

## `cat`
* **Arguments**: `file`
* **Returns**: Contents of a file (trimmed)
* **Example Output**: `Hello, World!`
* **Reads from**: Any file specified by path

## `counter`
* **Arguments**: *(none)*
* **Returns**: Counter value
* **Example Output**: `3`
* **OnClick**: Increments the counter

## `cpu_perc`
* **Arguments**: *(none)*
* **Returns**: CPU usage percentage
* **Example Output**: `42`
* **Reads from**: `/proc/stat`

## `datetime`
* **Arguments**: `datefmt` (strftime format)
* **Returns**: Current date/time
* **Example Output**: `2025-06-06 14:12`

## `disk_free`, `disk_used`, `disk_total`, `disk_perc`
* **Arguments**: `path`
* **Returns**:
  * `disk_free`: free space (e.g., `12.4 GiB`)
  * `disk_used`: used space
  * `disk_total`: total available space
  * `disk_perc`: used percentage (e.g., `67`)
* **Reads from**: Statfs via syscall

## `entropy`
* **Arguments**: *(none)*
* **Returns**: Available system entropy
* **Example Output**: `2760`
* **Reads from**: `/proc/sys/kernel/random/entropy_avail`

## `gid`, `uid`, `username`
* **Arguments**: *(none)*
* **Returns**:
  * `gid`: current group ID
  * `uid`: current user ID
  * `username`: current username
* **Reads from**: system calls and `/etc/passwd`

## `hostname`
* **Arguments**: *(none)*
* **Returns**: System hostname

## `ipv4`, `ipv6`, `up`
* **Arguments**: `interface`
* **Returns**:
  * `ipv4`: list of IPv4 addresses
  * `ipv6`: list of IPv6 addresses
  * `up`: `up` or `down`
* **Reads from**: `net.InterfaceByName`, `Interface.Addrs()`

## `kernel_release`
* **Arguments**: *(none)*
* **Returns**: Kernel release string
* **Reads from**: `syscall.Uname`

## `load_avg`
* **Arguments**: `minutes` (1, 5, or 15)
* **Returns**: System load average
* **Example Output**: `0.32`
* **Reads from**: `/proc/loadavg`

## `netspeed_rx`, `netspeed_tx`
* **Arguments**: `interface`
* **Returns**: RX/TX speed (e.g., `1.5 MiB/s`)
* **Reads from**: `/proc/net/dev`

## `notify_appname`, `notify_appicon`, `notify_summary`, `notify_body`, `notify_actions`
* **Arguments**: *(none)*
* **Returns**: Field from most recent notification
* **Example Output**: `firefox`
* **Use**: Only within notification context

## `num_files`
* **Arguments**: `path`
* **Returns**: Number of files in a directory

## `ram_free`, `ram_used`, `ram_total`, `ram_perc`
* **Arguments**: *(none)*
* **Returns**: Memory usage info
  * `ram_free`: Available memory
  * `ram_used`: Used memory
  * `ram_total`: Total memory
  * `ram_perc`: Usage percent
* **Reads from**: `/proc/meminfo`

## `run_command`
* **Arguments**: `command`
* **Returns**: Result of running command
* **OnClick**: *(none)*

## `period_command`
* **Arguments**: `command`, `interval` (e.g., `5s`)
* **Returns**: Output of the command, updated periodically

## `swap_free`, `swap_used`, `swap_total`, `swap_perc`
* **Arguments**: *(none)*
* **Returns**: Swap memory info
* **Reads from**: `/proc/meminfo`

## `temp`
* **Arguments**: `sensor`, `unit` (c/f/k)
* **Returns**: Temperature value in selected unit
* **Reads from**: `/sys/class/thermal/<sensor>/temp`
* **OnClick**: Toggles unit between °C, °F, and K

## `uptime`
* **Arguments**: *(none)*
* **Returns**: System uptime
* **Reads from**: `/proc/loadavg` (first field misused as uptime)

## `wifi_essid`, `wifi_perc`
* **Arguments**: `interface`
* **Returns**:
  * `wifi_essid`: Connected Wi-Fi network SSID
  * `wifi_perc`: Signal strength percentage
* **Uses**: `github.com/mdlayher/wifi` (netlink API)
