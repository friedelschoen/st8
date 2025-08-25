module github.com/friedelschoen/st8

go 1.24.5

require (
	github.com/friedelschoen/ctxmenu v0.0.0
	github.com/godbus/dbus/v5 v5.1.0
	github.com/ncruces/go-strftime v0.1.9
	github.com/spf13/pflag v1.0.7
)

require (
	github.com/friedelschoen/wayland v0.0.0 // indirect
	golang.org/x/image v0.30.0 // indirect
	golang.org/x/text v0.28.0 // indirect
)

replace github.com/friedelschoen/ctxmenu v0.0.0 => ../ctxmenu

replace github.com/friedelschoen/wayland v0.0.0 => ../wayland
