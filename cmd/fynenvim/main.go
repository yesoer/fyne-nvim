package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	nvim "github.com/yesoer/fyne-nvim"
)

func main() {
	a := app.New()
	w := a.NewWindow("Fyne NeoVim Example")
	w.Resize(fyne.NewSize(900, 600))

	nvim := nvim.New("./")
	nvim.Resize(fyne.NewSize(900, 600))
	w.SetContent(nvim)
	w.Canvas().Focus(nvim)

	fmt.Println("show and run")
	w.ShowAndRun()
}
