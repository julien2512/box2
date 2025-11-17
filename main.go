// Package main launches the calculator app
//
//go:generate fyne bundle -o data.go Icon.png
package main
import (
	"fyne.io/fyne/v2/app"
	"syscall"
)

func main() {
	// switches the Windows system clock from 15.6ms to 1ms.
	// from https://github.com/golang/go/issues/61042
	winmmDLL := syscall.NewLazyDLL("winmm.dll")
	procTimeBeginPeriod := winmmDLL.NewProc("timeBeginPeriod")
	procTimeBeginPeriod.Call(uintptr(1))

	app := app.New()
	app.SetIcon(resourceIconPng)

	c := newBox()
	c.loadUI(app)
	app.Run()
}