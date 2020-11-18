package devices

import (
	"fmt"
	evdev "github.com/gvalkov/golang-evdev"
	"time"
)

type Keyboard struct {
	Active bool
}

// eventCodeInEvents checks whether a given code is in a list of input events.
func eventCodeInEvents(code uint16, events []evdev.InputEvent) bool {
	for _, e := range events {
		if e.Code == code {
			return true
		}
	}
	return false
}

// WatchKeyboard watches a keyboard for a given event code and returns a struct
// containing its status.
func WatchKeyboard(keyboard *Keyboard, dev *evdev.InputDevice, code uint16) {
	ticker := time.NewTicker(time.Millisecond * 50)
	defer ticker.Stop()

	for range ticker.C {
		events, err := dev.Read()
		if err != nil {
			fmt.Println(err)
			panic("Could not read keyboard events.")
		}
		if eventCodeInEvents(code, events) {
			keyboard.Active = true
		} else {
			keyboard.Active = false
		}
	}
}
