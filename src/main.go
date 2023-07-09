package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
)

type Apps struct {
	Controllers  *ControllerApp
	Calibrations *CalibrationApp
	Connections  *ConnectionApp
}

func main() {
	a := app.New()
	w := a.NewWindow("Controller List")
	w.Resize(fyne.Size{1000, 700})

	data := dummyData()

	apps := &Apps{
		Controllers:  &ControllerApp{data: data, visible: data.all()},
		Calibrations: &CalibrationApp{data: data, visible: data.all()},
		Connections:  &ConnectionApp{data: data, inputs: data.inputs().Connections},
	}

	content := getTabs(apps, w, a)

	menu := fyne.NewMainMenu(
		fyne.NewMenu("File",
			fyne.NewMenuItem("Load", func() {
				data.load(w, apps, a)
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
		apps.Controllers.setController(data.all()[0])
		apps.Calibrations.setController(data.all()[0])
	}
	w.ShowAndRun()
}

func getTabs(apps *Apps, w fyne.Window, a fyne.App) *container.AppTabs {
	return container.NewAppTabs(
		container.NewTabItem("Controllers", apps.Controllers.makeControllerUI(w, apps)),
		container.NewTabItem("Broker", makeBrokerUI(w)),
		container.NewTabItem("Calibration", apps.Calibrations.makeCalibrationUI(w, a)),
		container.NewTabItem("Connection", apps.Connections.makeConnectionUI()),
	)
}
