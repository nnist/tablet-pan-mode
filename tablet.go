package main

import (
	"fmt"
	"github.com/bendahl/uinput"
	evdev "github.com/gvalkov/golang-evdev"
	"github.com/nnist/tablet-pan-mode/events"
	"strings"
)

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

	go events.WatchPenInRange(penChan, pen)
	go events.WatchDeviceForEventCode(kbdChan, kbd, evdev.KEY_CAPSLOCK)

	for {
		key_active = <-kbdChan
		pen_active = <-penChan
		if pen_active && key_active {
			fmt.Println("Button pressed and pen is in range!")
		}
	}
}
