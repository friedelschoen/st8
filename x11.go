package main

// #cgo pkg-config: x11
// #include <stdlib.h>
// #include <X11/Xlib.h>
import "C"
import "unsafe"

type Display C.Display

func OpenDisplay() *Display {
	return (*Display)(C.XOpenDisplay(nil))
}

func (dpy *Display) cptr() *C.Display {
	return (*C.Display)(dpy)
}

func (dpy *Display) Close() {
	C.XCloseDisplay(dpy.cptr())
}

func (dpy *Display) StoreName(text string) {
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	C.XStoreName(dpy.cptr(), C.XDefaultRootWindow(dpy.cptr()), ctext)
	C.XFlush(dpy.cptr())
}
