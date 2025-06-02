//go:build x11
// +build x11

package driver

// #cgo pkg-config: x11
// #include <stdlib.h>
// #include <X11/Xlib.h>
import "C"
import (
	"fmt"
	"strings"
	"unsafe"

	"github.com/friedelschoen/st8/component"
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

func (dpy *Display) SetText(line []component.Block) error {
	var out strings.Builder
	for i, block := range line {
		if i > 0 {
			out.WriteString(" | ")
		}
		out.WriteString(block.Text)
	}
	ctext := C.CString(out.String())
	defer C.free(unsafe.Pointer(ctext))
	C.XStoreName(dpy.ptr, C.XDefaultRootWindow(dpy.ptr), ctext)
	C.XFlush(dpy.ptr)
	return nil
}
