package main

import (
	"fmt"
	"github.com/bendahl/uinput"
)

func main() {
	fmt.Println("test")
	keyboard, err := uinput.CreateKeyboard("/dev/uinput", []byte("testkeyboard"))
	if err != nil {
		return
	}
	defer keyboard.Close()
	keyboard.KeyPress(uinput.KeyA)
}
