package main

import (
	"fmt"
	"github.com/bendahl/uinput"
	evdev "github.com/gvalkov/golang-evdev"
	"github.com/nnist/tablet-pan-mode/devices"
	"strings"
	"time"
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
	evdevDevices, _ := evdev.ListInputDevices(device_glob)

	var pen *evdev.InputDevice = nil
	var kbd *evdev.InputDevice = nil

	for _, dev := range evdevDevices {
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

	var penDev devices.PenDevice
	var key_active bool

	kbdChan := make(chan bool)

	go devices.WatchPen(&penDev, pen)
	go devices.WatchDeviceForEventCode(kbdChan, kbd, evdev.KEY_CAPSLOCK)

	ticker := time.NewTicker(time.Millisecond * 25)
	defer ticker.Stop()

	for range ticker.C {
		select {
		case key_active = <-kbdChan:
			fmt.Println("key: ", key_active)
		default:
		}

		if penDev.Active && key_active {
			fmt.Println("Button pressed and pen is in range!")
			fmt.Println("pen:", penDev.X, penDev.Y)
		}

	}
}
