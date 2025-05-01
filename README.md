# DWM-status -- a clone of `slstatus` rewritting in Go

> _a better name will follow_

| implemented | function                                     | description                                 | argument (example)        |
| ----------- | -------------------------------------------- | ------------------------------------------- | ------------------------- |
| [x]         | battery_perc                                 | battery percentage                          | battery name (BAT0)       |
| [x]         | battery_remaining   battery remaining HH:MM  | battery name (BAT0)                         |
| [x]         | battery_state                                | battery charging state                      | battery name (BAT0)       |
| [x]         | cat                                          | read arbitrary file                         | path                      |
| [x]         | cpu_freq                                     | cpu frequency in MHz                        | NULL                      |
| [x]         | cpu_perc                                     | cpu usage in percent                        | NULL                      |
| [x]         | datetime                                     | date and time                               | format string (%F %T)     |
| [x]         | disk_free                                    | free disk space in GB                       | mountpoint path (/)       |
| [x]         | disk_perc                                    | disk usage in percent                       | mountpoint path (/)       |
| [x]         | disk_total                                   | total disk space in GB                      | mountpoint path (/)       |
| [x]         | disk_used                                    | used disk space in GB                       | mountpoint path (/)       |
| [x]         | entropy                                      | available entropy                           | NULL                      |
| [x]         | gid                                          | GID of current user                         | NULL                      |
| [x]         | hostname                                     | hostname                                    | NULL                      |
| [x]         | ipv4                                         | IPv4 address                                | interface name (eth0)     |
| [x]         | ipv6                                         | IPv6 address                                | interface name (eth0)     |
| [x]         | kernel_release                               | `uname -r`                                  | NULL                      |
| [ ]^1       | keyboard_indicators caps/num lock indicators | format string (c?n?)                        | see keyboard_indicators.c |
| [ ]^1       | keymap                                       | layout (variant) of current                 | NULL                      |
| [x]         | load_avg                                     | load average                                | NULL                      |
| [x]         | netspeed_rx                                  | receive network speed                       | interface name (wlan0)    |
| [x]         | netspeed_tx                                  | transfer network speed                      | interface name (wlan0)    |
| [x]         | num_files                                    | number of files in a directory  path        | (/home/foo/Inbox/cur)     |
| [x]         | ram_free                                     | free memory in GB                           | NULL                      |
| [x]         | ram_perc                                     | memory usage in percent                     | NULL                      |
| [x]         | ram_total                                    | total memory size in GB                     | NULL                      |
| [x]         | ram_used                                     | used memory in GB                           | NULL                      |
| [x]         | run_command                                  | custom shell command                        | command (echo foo)        |
| [x]         | swap_free                                    | free swap in GB                             | NULL                      |
| [x]         | swap_perc                                    | swap usage in percent                       | NULL                      |
| [x]         | swap_total                                   | total swap size in GB                       | NULL                      |
| [x]         | swap_used                                    | used swap in GB                             | NULL                      |
| [x]         | temp                                         | temperature in degree celsius   sensor file | (/sys/class/thermal/...)  |
|             |                                              |                                             | NULL on OpenBSD           |
|             |                                              |                                             | thermal zone on FreeBSD   |
|             |                                              |                                             | (tz0, tz1, etc.)          |
| [x]         | uid                                          | UID of current user                         | NULL                      |
| [x]         | up                                           | interface is running                        | interface name (eth0)     |
| [x]         | uptime                                       | system uptime                               | NULL                      |
| [x]         | username                                     | username of current user                    | NULL                      |
| [ ]         | vol_perc                                     | OSS/ALSA volume in percent                  | mixer file (/dev/mixer)   |
|             |                                              |                                             | NULL on OpenBSD/FreeBSD   |
| [x]         | wifi_essid                                   | WiFi ESSID                                  | interface name (wlan0)    |
| [x]         | wifi_perc                                    | WiFi signal in percent                      | interface name (wlan0)    |

> ^1: is not yet implemented due to lack of portability