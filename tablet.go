package main

import (
	"fmt"
	"github.com/bendahl/uinput"
	evdev "github.com/gvalkov/golang-evdev"
	"strings"
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
func watchDeviceForEventCode(c chan bool, dev *evdev.InputDevice, code uint16) {
	for {
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

// watchPenInRange returns true when the pen is in range of the tablet.
func watchPenInRange(c chan bool, dev *evdev.InputDevice) {
	for {
		events, err := dev.Read()
		if err != nil {
			fmt.Println(err)
			panic("Could not read device events.")
		}
		for _, e := range events {
			if e.Code == evdev.ABS_DISTANCE {
				if e.Value != 0 {
					c <- true
				} else {
					c <- false
				}
			}
		}
	}
}

func main() {
	keyboard, err := uinput.CreateKeyboard("/dev/uinput", []byte("testkeyboard"))
	if err != nil {
		return
	}
	defer keyboard.Close()

	mouse, err := uinput.CreateMouse("/dev/uinput", []byte("testmouse"))
	if err != nil {
		return
	}
	defer mouse.Close()

	device_glob := "/dev/input/event*"
	devices, _ := evdev.ListInputDevices(device_glob)

	var pen *evdev.InputDevice = nil
	var kbd *evdev.InputDevice = nil

	for _, dev := range devices {
		if strings.Contains(dev.Name, "Wacom") && strings.Contains(dev.Name, "Pen") {
			pen = dev
		}
		if strings.Contains(dev.Name, "Kinesis") &&
			strings.Contains(dev.Name, "Freestyle") &&
			strings.Contains(dev.Phys, "input0") {
			kbd = dev
		}
	}

	if pen == nil {
		panic("Could not open pen device.")
	}
	if kbd == nil {
		panic("Could not open keyboard device.")
	}

	var pen_active bool
	var key_active bool

	penChan := make(chan bool)
	kbdChan := make(chan bool)

	go watchPenInRange(penChan, pen)
	go watchDeviceForEventCode(kbdChan, kbd, evdev.KEY_CAPSLOCK)

	for {
		key_active = <-kbdChan
		pen_active = <-penChan
		if pen_active && key_active {
			fmt.Println("Button pressed and pen is in range!")
		}
	}
}
