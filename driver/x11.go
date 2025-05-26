//go:build x11
// +build x11

package driver

// #cgo pkg-config: x11
// #include <stdlib.h>
// #include <X11/Xlib.h>
import "C"
import (
	"fmt"
	"unsafe"
)

type Display struct{ ptr *C.Display }

func init() {
	Drivers["xsetname"] = &Display{}
}

func (dpy *Display) Init() error {
	dpy.ptr = (*C.Display)(C.XOpenDisplay(nil))
	if dpy.ptr == nil {
		return fmt.Errorf("unable to open display")
	}
	return nil
}

func (dpy *Display) Close() error {
	if dpy.ptr != nil {
		C.XCloseDisplay(dpy.ptr)
		dpy.ptr = nil
	}
	return nil
}

func (dpy *Display) SetText(text string) error {
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	C.XStoreName(dpy.ptr, C.XDefaultRootWindow(dpy.ptr), ctext)
	C.XFlush(dpy.ptr)
	return nil
}
