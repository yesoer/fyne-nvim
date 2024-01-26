package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {
	a := app.New()
	w := a.NewWindow("Fyne NeoVim Example")
	w.Resize(fyne.NewSize(800, 600))

	nvim := New()
	nvim.Resize(fyne.NewSize(800, 600))
	nvim.SetText("Hello, World!")
	w.SetContent(nvim)
	w.Canvas().Focus(nvim)

	fmt.Println("show and run")
	w.ShowAndRun()
}
