package swaybarproto

import "syscall"

type Header struct {
	Version        int            `json:"version"`                // The protocol version to use. Currently, this must be 1
	ClickEvents    bool           `json:"click_events,omitempty"` // Whether to receive click event information to standard input
	ContinueSignal syscall.Signal `json:"cont_signal,omitempty"`  // The signal that swaybar should send to continue processing
	StopSignal     syscall.Signal `json:"stop_signal,omitempty"`  // The signal that swaybar should send to stop processing
}
