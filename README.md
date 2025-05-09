# st8 - a fast, minimal status bar for DWM and friends

`st8` is a clone of [`slstatus`](https://tools.suckless.org/slstatus/) rewritten in Go.
It aims to be simpler to configure, easier to extend, and includes optional desktop notifications.

## Format Syntax

The `st8` status bar uses a simple template format string to render the output. Each dynamic element is written as:

```
${function_name[!padding][:argument]}
```

### Rules:

- `${...}` denotes a placeholder to be replaced with live system data.
- `function_name` is the name of a supported function (e.g. `netspeed_rx`, `username`, `battery_perc`).
- An optional `!padding` can be passed
- An optional `:argument` can be passed to the function (e.g. network interface like `eth0`).
- Plain text (like `hello`, `up:`, `down:`, `-`) is preserved as-is.

Padding is done like in C's `printf`:

- `-010` pads with `0`'s until a length of 10 chars, aligns the text left
- `5` pads with spaces until a length of 5 chars, aligns the text right
- `_7` pads with `_`'s until a length of 7 chars, aligns the text right

### Example:

```
hello up:${netspeed_rx!8:eth0} down:${netspeed_tx!8:eth0} - ${username}
```

This might render as:

```
hello up:103KB/s down:240KB/s - alice
```

### Notes:
- If a function fails or is unavailable, it returns an `<error>`.

## Functions

| function                  | description                          | argument (example)               |
| ------------------------- | ------------------------------------ | -------------------------------- |
| `battery_perc`            | battery percentage                   | battery name (BAT0)              |
| `battery_remaining`       | battery name (BAT0)                  |
| `battery_state`           | battery charging state               | battery name (BAT0)              |
| `cat`                     | read arbitrary file                  | path (/home/foo/packages.txt)    |
| `cpu_freq`                | cpu frequency in MHz                 |                                  |
| `cpu_perc`<sup>1</sup>    | cpu usage in percent                 |                                  |
| `datetime`                | date and time                        | format string (%F %T)            |
| `disk_free`               | free disk space in GB                | mountpoint path (/)              |
| `disk_perc`               | disk usage in percent                | mountpoint path (/)              |
| `disk_total`              | total disk space in GB               | mountpoint path (/)              |
| `disk_used`               | used disk space in GB                | mountpoint path (/)              |
| `entropy`                 | available entropy                    |                                  |
| `gid`                     | GID of current user                  |                                  |
| `hostname`                | hostname                             |                                  |
| `ipv4`                    | IPv4 address                         | interface name (eth0)            |
| `ipv6`                    | IPv6 address                         | interface name (eth0)            |
| `kernel_release`          | current kernel                       |                                  |
| `load_avg`                | load average                         | average over time in minutes (5) |
|                           |                                      | if empty `<1min> <5min> <15min>` |
| `netspeed_rx`<sup>1</sup> | receive network speed                | interface name (wlan0)           |
| `netspeed_tx`<sup>1</sup> | transfer network speed               | interface name (wlan0)           |
| `num_files`               | number of files in a directory  path | directory (/home/foo/Inbox/cur)  |
| `ram_free`                | free memory in GB                    |                                  |
| `ram_perc`                | memory usage in percent              |                                  |
| `ram_total`               | total memory size in GB              |                                  |
| `ram_used`                | used memory in GB                    |                                  |
| `run_command`             | runs shell command                   | command (echo foo)               |
| `swap_free`               | free swap in GB                      |                                  |
| `swap_perc`               | swap usage in percent                |                                  |
| `swap_total`              | total swap size in GB                |                                  |
| `swap_used`               | used swap in GB                      |                                  |
| `temp`                    | temperature in degree celsius        | sensor-filter (thermal_zone0)    |
|                           |                                      | if empty first occurring         |
| `uid`                     | UID of current user                  |                                  |
| `up`                      | interface is running                 | interface name (eth0)            |
| `uptime`                  | system uptime                        |                                  |
| `username`                | username of current user             |                                  |
| `wifi_essid`              | WiFi ESSID                           | interface name (wlan0)           |
| `wifi_perc`               | WiFi signal in percent               | interface name (wlan0)           |

New in `st8`:

| function                     | description                           | argument (example)               |
| ---------------------------- | ------------------------------------- | -------------------------------- |
| `notify_appname`<sup>2</sup> | senders name                          |                                  |
| `notify_appicon`<sup>2</sup> | senders icon                          |                                  |
| `notify_summary`<sup>2</sup> | notification summary                  |                                  |
| `notify_body`<sup>2</sup>    | notification content                  |                                  |
| `notify_actions`<sup>2</sup> | notification actions, comma-seperated |                                  |
| `period_command`             | runs command at an interval           | period,command (10m,checkupdate) |


<sup>1</sup>Some functions require at least two runs, e.g. `netspeed_*` must calculate its speed over time.

<sup>2</sup>Do only work in `notify.txt`

### Notes:
- `period_command` is non-blocking and is meant for heavy commands which does not have to be run every period, e.g. checking for updates.
  `run_command` is blocking and is meant for lightweight commands, e.g. `sv status sshd`

### Not implemented (yet)

| function              | description                 |
| --------------------- | --------------------------- |
| `vol_perc`            | OSS/ALSA volume in percent  |
| `keyboard_indicators` | caps/num lock indicators    |
| `keymap`              | layout (variant) of current |

If you want to add a function, look at the other implementations in `component/` and add the function in `component/interface.go`.

## Notifications

`st8` also implements a lightweight D-Bus notification daemon (`org.freedesktop.Notifications`).
You can send notifications from scripts using tools like `notify-send`:

```sh
notify-send "Build finished" "Your code compiled successfully."
```

These notifications will appear inline in your status bar, replacing the output temporarily.

## License

This project uses zlib-license.