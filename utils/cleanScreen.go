package utils

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

// create a map for storing clear funcs
var clear map[string]func()

func init() {
	clear = make(map[string]func())
	clear["linux"] = func() {
		// Linux clear
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["darwin"] = func() {
		// Mac clear
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["windows"] = func() {
		// Windows clear
		cmd := exec.Command("cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

// Limpia la pantalla de terminal
func ClearScreen() {
	// runtime.GOOS -> linux, windows, darwin
	value, ok := clear[runtime.GOOS]

	// if we defined a clear func for that platform:
	if ok {
		// we execute it
		value()
	} else {
		// unsupported platform
		fmt.Println("-----------------------------------------------------")
	}
}
