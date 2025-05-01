# st8 — a fast, minimal status bar for DWM and friends

`st8` is a clone of [`slstatus`](https://tools.suckless.org/slstatus/) rewritten in Go.
It aims to be simpler to configure, easier to extend, and includes optional desktop notifications.

## Format Syntax

The `st8` status bar uses a simple template format string to render the output. Each dynamic element is written as:

```
${function_name[:argument]}
```

### Rules:
- `${...}` denotes a placeholder to be replaced with live system data.
- `function_name` is the name of a supported function (e.g. `netspeed_rx`, `username`, `battery_perc`).
- An optional `:argument` can be passed to the function (e.g. network interface like `eth0`).
- Plain text (like `hello`, `up:`, `down:`, `-`) is preserved as-is.

### Example:

```txt
hello up:${netspeed_rx:eth0} down:${netspeed_tx:eth0} - ${username}
```

This might render as:

```txt
hello up:103KB/s down:240KB/s - alice
```

### Notes:
- If a function fails or is unavailable, it returns an `<error>`.
- Some functions require at least two runs, e.g. `netspeed_*` must calculate its speed over time.

## Functions

| function          | description                          | argument (example)              |
| ----------------- | ------------------------------------ | ------------------------------- |
| battery_perc      | battery percentage                   | battery name (BAT0)             |
| battery_remaining | battery name (BAT0)                  |
| battery_state     | battery charging state               | battery name (BAT0)             |
| cat               | read arbitrary file                  | path (/home/foo/packages.txt)   |
| cpu_freq          | cpu frequency in MHz                 |                                 |
| cpu_perc          | cpu usage in percent                 |                                 |
| datetime          | date and time                        | format string (%F %T)           |
| disk_free         | free disk space in GB                | mountpoint path (/)             |
| disk_perc         | disk usage in percent                | mountpoint path (/)             |
| disk_total        | total disk space in GB               | mountpoint path (/)             |
| disk_used         | used disk space in GB                | mountpoint path (/)             |
| entropy           | available entropy                    |                                 |
| gid               | GID of current user                  |                                 |
| hostname          | hostname                             |                                 |
| ipv4              | IPv4 address                         | interface name (eth0)           |
| ipv6              | IPv6 address                         | interface name (eth0)           |
| kernel_release    | `uname -r`                           |                                 |
| load_avg          | load average                         |                                 |
| netspeed_rx       | receive network speed                | interface name (wlan0)          |
| netspeed_tx       | transfer network speed               | interface name (wlan0)          |
| num_files         | number of files in a directory  path | directory (/home/foo/Inbox/cur) |
| ram_free          | free memory in GB                    |                                 |
| ram_perc          | memory usage in percent              |                                 |
| ram_total         | total memory size in GB              |                                 |
| ram_used          | used memory in GB                    |                                 |
| run_command       | custom shell command                 | command (echo foo)              |
| swap_free         | free swap in GB                      |                                 |
| swap_perc         | swap usage in percent                |                                 |
| swap_total        | total swap size in GB                |                                 |
| swap_used         | used swap in GB                      |                                 |
| temp              | temperature in degree celsius        | sensor-filter (thermal_zone0)   |
|                   |                                      | none uses first occurring       |
| uid               | UID of current user                  |                                 |
| up                | interface is running                 | interface name (eth0)           |
| uptime            | system uptime                        |                                 |
| username          | username of current user             |                                 |
| wifi_essid        | WiFi ESSID                           | interface name (wlan0)          |
| wifi_perc         | WiFi signal in percent               | interface name (wlan0)          |

### Not implemented (yet)

| function                 | description                 | argument (example)        |
| ------------------------ | --------------------------- | ------------------------- |
| vol_perc                 | OSS/ALSA volume in percent  |                           |
| caps/num lock indicators | format string (c?n?)        | see keyboard_indicators.c |
| keymap                   | layout (variant) of current |                           |

If you want to add a function, look at the other implementations in `component/` and add the function in `component/interface.go`.

## License

This project uses zlib-license.