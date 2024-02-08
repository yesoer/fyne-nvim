# Fyne-Nvim
 
This go package contains a custom neovim widget for the fyne framework.
 
## Install

To install the application :
```sh
go install github.com/yesoer/fyne-nvim/cmd/fynenvim@latest
```

To add it to your own project :
```sh
go get github.com/yesoer/fyne-nvim
```

## Usage as a widget

Using the fyne neovim widget in your project is pretty straight forward,
as can be seen from the cmd/fynenvim/main.go :

```go
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

	nvim := nvim.New()
	nvim.Resize(fyne.NewSize(900, 600))
	w.SetContent(nvim)
	w.Canvas().Focus(nvim)

	fmt.Println("show and run")
	w.ShowAndRun()
}
```
