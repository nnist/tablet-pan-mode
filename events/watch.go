package events

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

// PenDevice contains the position and range status of a pen device.
type PenDevice struct {
	X      int32
	Y      int32
	Active bool
}

// watchPen returns the position and range of a pen device.
func WatchPen(c chan PenDevice, dev *evdev.InputDevice) {
	var x int32
	var y int32
	var active bool

	ticker := time.NewTicker(time.Millisecond * 50)
	defer ticker.Stop()

	for range ticker.C {
		events, err := dev.Read()
		if err != nil {
			fmt.Println(err)
			panic("Could not read device events.")
		}

		for _, e := range events {
			if e.Type == evdev.EV_ABS {
				if e.Code == evdev.ABS_DISTANCE {
					if e.Value != 0 {
						active = true
					} else {
						active = false
					}
				} else if e.Code == evdev.ABS_X {
					x = e.Value
				} else if e.Code == evdev.ABS_Y {
					y = e.Value
				}
			}
		}
		c <- PenDevice{x, y, active}
	}
}
