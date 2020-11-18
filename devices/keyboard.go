package devices

import (
	"fmt"
	evdev "github.com/gvalkov/golang-evdev"
	"time"
)

type Keyboard struct {
	Active   bool
	lastTime int64
}

// WatchKeyboard watches a keyboard for a given event code and returns a struct
// containing its status.
func WatchKeyboard(keyboard *Keyboard, dev *evdev.InputDevice, code uint16) {
	// Create a timeout
	go func() {
		ticker := time.NewTicker(time.Millisecond * 50)
		defer ticker.Stop()
		for range ticker.C {
			keyboard.Active = time.Now().UnixNano()-keyboard.lastTime < 100000000
		}
	}()

	for {
		events, err := dev.Read()
		if err != nil {
			fmt.Println(err)
			panic("Could not read keyboard events.")
		}
		for _, e := range events {
			if e.Code == code {
				keyboard.Active = true
				keyboard.lastTime = e.Time.Nano()
			}
		}
	}
}
