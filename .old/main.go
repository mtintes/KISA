package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var MyApp fyne.App
var mainSetup setup

func main() {
	mainSetup = setup{
		InputPicos: []InputPico{},
	}

	MyApp = app.New()
	myWindow := MyApp.NewWindow("Grid Layout")
	myWindow.Resize(fyne.NewSize(500, 500))
	myWindow.SetMaster()
	myWindow.SetContent(tabMenu())
	myWindow.ShowAndRun()
}

func tabMenu() fyne.CanvasObject {
	tabs := container.NewAppTabs(
		container.NewTabItem("Inputs", configurePage(MyApp)),
		container.NewTabItem("Setup", widget.NewLabel("World")),
		container.NewTabItem("Record", widget.NewLabel("World")),
		container.NewTabItem("Play", widget.NewLabel("World")),
	)

	tabs.SetTabLocation(container.TabLocationTop)
	return tabs
}
