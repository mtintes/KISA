package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func makeConnectionUI() *fyne.Container {
	return container.NewVBox(
		widget.NewLabel("Connection"),
		layout.NewSpacer(),
	)

}
