# Fyne-Nvim
 
This go package contains a custom neovim widget for the fyne framework.

## Table of Contents
- [Get Started](#get-started)
  - [Usage as a Widget](#usage-as-a-widget)
- [Developer Notes](#developer-notes)
  - [Contributions](#contributions)
  - [Project Structure](#project-structure)
  - [Resources](#resources)
 
## Get Started

To install the application :
```sh
go install github.com/yesoer/fyne-nvim/cmd/fynenvim@latest
```

To add it to your own project :
```sh
go get github.com/yesoer/fyne-nvim
```

### Usage as a Widget

Using the fyne neovim widget in your project is pretty straight forward,
as can be seen from the `cmd/fynenvim/main.go` :

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

## Developer Notes

### Contributions

When confused, no worries, just publish what you have. 
A not perfectly structured contribution is still far better than nothing.

#### Branch Naming 

Branch names should look like this
`type>/<name>`

`<type>` is one of the following (extend if needed) :

| type | when to use      |
|------|------------------|
| feat | any new features |
| maintenance | any work on docs, git workflows, tests etc. |
| refactor | when refactoring existing parts of the application |
| fix  | bug fixes        |
| test | testing environments/throwaway branches |

`<name>` is a short description of what you are doing, words should be seperated using '-'.

#### Commit Messages

More specific distinction happens in **commit messages** which should be structured
as follows :

```
<type>(<scope>): <subject>
```

- **type**
Must be one of the following:

  * **feat**: A new feature
  * **fix**: A bug fix
  * **docs**: Documentation only changes
  * **style**: Changes that do not affect the meaning of the code (white-space, formatting, missing
    semi-colons, etc)
  * **refactor**: A code change that neither fixes a bug nor adds a feature
  * **perf**: A code change that improves performance
  * **test**: Adding missing or correcting existing tests
  * **chore**: Changes to the build process or auxiliary tools and libraries such as documentation
  generation

- **scope** refers to the part of the software, which usually will be best identified by the package name.

- **subject** gives a short idea of what was done/what the intend of the commit is.

As for the **commit body** there is no mandatory structure as of now.

#### Other Tips/Notes on Contributing

**Issues and Pull Requests** for now will not have any set guidelines.

Check out [Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments) for commmon code review topics around golang.

### Project Structure

Maybe this description helps you getting started and understanding the project.
Don't be afraid to ask for more information.

| location    | Description |
|-------------|-------------|
| cmd/        | Contains the fynenvim executable code |
| nvim.go     | Implements the widget interface i.e. is the center of this project |
| render.go   | Implements the renderer for our widget as required for custom widgets |
| input.go    | Using the mappings from keymap.go this forwards inputs from Fyne to Neovim |
| output.go   | Provides functions to write runes etc. to the textgrid which visualizes Neovim |
| events.go   | Process the events received from Neovim (uses output.go to forward visual changes to Fyne) |

### Resources

[Neovim UI Event Docs](https://neovim.io/doc/user/ui.html)
[Neovim UI API Docs](https://neovim.io/doc/user/api.html#api-ui)
[Fyne Custom Widget Docs](https://docs.fyne.io/extend/custom-widget) with a more detailed example [on github](https://github.com/stuartdd2/developer.fyne.io/blob/master/extend/custom-widget.md)

To me the [fyne terminal project](https://github.com/fyne-io/terminal) was very helpful as well.
