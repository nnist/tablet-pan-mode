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
func WatchPen(pen *PenDevice, dev *evdev.InputDevice) {
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
				switch code := e.Code; code {
				case evdev.ABS_DISTANCE:
					if e.Value != 0 {
						pen.Active = true
					} else {
						pen.Active = false
					}
				case evdev.ABS_X:
					pen.X = e.Value
				case evdev.ABS_Y:
					pen.Y = e.Value
				}
			}
		}

	}
}
