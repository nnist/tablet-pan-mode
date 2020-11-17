package devices

import (
	"fmt"
	evdev "github.com/gvalkov/golang-evdev"
	"time"
)

// eventCodeInEvents checks whether a given code is in a list of input events.
func eventCodeInEvents(code uint16, events []evdev.InputEvent) bool {
	for _, e := range events {
		if e.Code == code {
			return true
		}
	}
	return false
}

// watchDeviceForEventCode watches a device for a given event code and returns
// true on a channel if the event is triggered.
func WatchDeviceForEventCode(c chan bool, dev *evdev.InputDevice, code uint16) {
	ticker := time.NewTicker(time.Millisecond * 50)
	defer ticker.Stop()

	for range ticker.C {
		events, err := dev.Read()
		if err != nil {
			fmt.Println(err)
			panic("Could not read device events.")
		}
		if eventCodeInEvents(code, events) {
			c <- true
		} else {
			c <- false
		}
	}
}
