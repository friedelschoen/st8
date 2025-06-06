# st8 - a fast, minimal status bar

`st8` is a clone of [`slstatus`](https://tools.suckless.org/slstatus/) rewritten in Go.

## Outputs

To be complainant with multiple systems, `st8` understands multiple outputs:

- **`stdout`**: prints every block delimited by `|` to stdout.
- **`xsetname`**: sets roots (the window managers) name to blocks delimited by `|`.
- **`sway`**: use `swaybar-protocol(7)` via stdout/stdin. Should work on i3 but is untested yet.

## Format Syntax

`st8` uses blocks (like in `i3blocks` or `dwmblocks`), one block should have one function. The configuration,
usually located at `$XDG_CONFIG_HOME/st8/{status,notify}.conf` or `~/.config/st8/{status,notify}.conf` follows the [INI-format](https://en.wikipedia.org/wiki/INI_file).

Keys and values are stripped from leading and trailing whitespaces, except a value can be quoted (using `"`) and the whitespaces between the quotes are conserved.

### Blocks

The section-name corresponds to the [component](docs/components.md). Some components take arguments which can be specified in the section. For example:

```ini
[battery_perc]
battery = BAT0
```

### Special Keys

Some keys may be overwritten by the component or ignored by the output-driver.

- **text-color**: Color of the text: `#RRGGBB` or `#RRGGBBAA`
- **background-color**: Color of the block-background: `#RRGGBB` or `#RRGGBBAA`
- **border-color**: Color of the block-border: `#RRGGBB` or `#RRGGBBAA`
- **border-top**: Width in pixel of the top-border (default `1px`)
- **border-bottom**: Width in pixel of the bottom-border (default `1px`)
- **border-left**: Width in pixel of the left-border (default `1px`)
- **border-right**: Width in pixel of the right-border (default `1px`)
- **width**: Minimum width of this block, either in pixels (`12px`), in space-widths (`6wh`) or the width of a string `ab-cd`
- **align**: Alignment of the text if block is wider, either `left`, `center` or `right` (default `left`)
- **separate**: Put a separator right to this block (default `yes`)
- **separator-width**: With separator right to this block (default `yes`)
- **markup**: Use markup: either `none` or [`pango`](https://docs.gtk.org/Pango/pango_markup.html) (default `none`)
- **format**: Format the resulting text, first occurrence of `{}` is replaced with the actual text and can contain optional padding: (default `{}`)
  Padding is done like in C's `printf`:
  - `-010` pads with `0`'s until a length of 10 chars, aligns the text left
  - `5` pads with spaces until a length of 5 chars, aligns the text right
  - `_7` pads with `_`'s until a length of 7 chars, aligns the text right

## Notes:

- If a function fails or is unavailable, it returns an `<error>`.

## Components

| function                  | description                          |
| ------------------------- | ------------------------------------ |
| `battery_perc`            | battery percentage                   |
| `battery_remaining`       | battery name (BAT0)                  |
| `battery_state`           | battery charging state               |
| `cat`                     | read arbitrary file                  |
| `cpu_freq`                | cpu frequency in MHz                 |
| `cpu_perc`<sup>1</sup>    | cpu usage in percent                 |
| `datetime`                | date and time                        |
| `disk_free`               | free disk space in GB                |
| `disk_perc`               | disk usage in percent                |
| `disk_total`              | total disk space in GB               |
| `disk_used`               | used disk space in GB                |
| `entropy`                 | available entropy                    |
| `gid`                     | GID of current user                  |
| `hostname`                | hostname                             |
| `ipv4`                    | IPv4 address                         |
| `ipv6`                    | IPv6 address                         |
| `kernel_release`          | current kernel                       |
| `load_avg`                | load average                         |
|                           |                                      |
| `netspeed_rx`<sup>1</sup> | receive network speed                |
| `netspeed_tx`<sup>1</sup> | transfer network speed               |
| `num_files`               | number of files in a directory  path |
| `ram_free`                | free memory in GB                    |
| `ram_perc`                | memory usage in percent              |
| `ram_total`               | total memory size in GB              |
| `ram_used`                | used memory in GB                    |
| `run_command`             | runs shell command                   |
| `swap_free`               | free swap in GB                      |
| `swap_perc`               | swap usage in percent                |
| `swap_total`              | total swap size in GB                |
| `swap_used`               | used swap in GB                      |
| `temp`                    | temperature in degree celsius        |
|                           |                                      |
| `uid`                     | UID of current user                  |
| `up`                      | interface is running                 |
| `uptime`                  | system uptime                        |
| `username`                | username of current user             |
| `wifi_essid`              | WiFi ESSID                           |
| `wifi_perc`               | WiFi signal in percent               |

New in `st8`:

| function                     | description                           |
| ---------------------------- | ------------------------------------- |
| `notify_appname`<sup>2</sup> | senders name                          |
| `notify_appicon`<sup>2</sup> | senders icon                          |
| `notify_summary`<sup>2</sup> | notification summary                  |
| `notify_body`<sup>2</sup>    | notification content                  |
| `notify_actions`<sup>2</sup> | notification actions, comma-seperated |
| `period_command`             | runs command at an interval           |
| `counter`                    | counter which increases by clicking   |


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