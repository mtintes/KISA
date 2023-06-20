package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("Controller List")
	w.Resize(fyne.Size{700, 700})

	data := dummyData()
	Controllers := &ControllerApp{data: data, visible: data.all()}

	content := getTabs(Controllers, w)

	menu := fyne.NewMainMenu(
		fyne.NewMenu("File",
			fyne.NewMenuItem("Load", func() {
				data.load(w, Controllers)
			}),
			fyne.NewMenuItem("Save...", func() {
				data.save(w)
			}),
			fyne.NewMenuItem("Save As...", func() {
				data.saveAs(w)
			}),
			fyne.NewMenuItem("Quit", func() {
				a.Quit()
			}),
		),
	)

	w.SetContent(content)
	w.SetMainMenu(menu)

	if len(data.all()) > 0 {
		Controllers.setController(data.all()[0])
	}
	w.ShowAndRun()
}

func getTabs(Controllers *ControllerApp, w fyne.Window) *container.AppTabs {
	return container.NewAppTabs(
		container.NewTabItem("Controllers", Controllers.makeControllerUI(w)),
		container.NewTabItem("Broker", widget.NewLabel("World")),
	)
}
