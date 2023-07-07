package main

import (
	"sort"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type CalibrationApp struct {
	data    *ControllerList
	visible []*Controller
	current *Controller

	Controllers *widget.List
	Id          *widget.Label
	Name        *widget.Label
	Direction   *widget.Label
	Pins        *fyne.Container
}

func (a *CalibrationApp) makeCalibrationUI(w fyne.Window, app fyne.App) *fyne.Container {
	a.Controllers = widget.NewList(
		func() int {
			return len(a.visible)
		},
		func() fyne.CanvasObject {
			sidebar := container.NewWithoutLayout(widget.NewLabel("Item x"))
			return sidebar
		},
		func(i int, c fyne.CanvasObject) {
			Controller := a.visible[i]
			box := c.(*fyne.Container)
			label := box.Objects[0].(*widget.Label)
			label.SetText(Controller.Name)
		})

	a.Controllers.OnSelected = func(id int) {
		a.setController(a.visible[id])
	}

	a.Id = widget.NewLabel("")

	a.Name = widget.NewLabel("")

	a.Direction = widget.NewLabel("")

	a.Pins = container.NewBorder(container.NewGridWithColumns(6,
		widget.NewLabel("pin"),
		widget.NewLabel("topic"),
		widget.NewLabel("MinValue"),
		widget.NewLabel("MaxValue"),
		widget.NewLabel("Direction"),
		widget.NewLabel("")),
		nil,
		nil,
		nil,
		container.New(layout.NewBorderLayout(nil, nil, nil, nil), widget.NewList(
			func() int {
				if a.current != nil {
					return len(a.current.Pins)
				}
				return 1
			},
			func() fyne.CanvasObject {
				pins := container.NewGridWithColumns(6,
					widget.NewLabel(""),
					widget.NewLabel(""),
					widget.NewLabel(""),
					widget.NewLabel(""),
					widget.NewLabel(""),
					widget.NewButton("calibrate", func() {}))
				return pins
			}, func(lii widget.ListItemID, co fyne.CanvasObject) {
				if a.current != nil {
					Pin := a.current.Pins[lii]
					box := co.(*fyne.Container)
					pinNumber := box.Objects[0].(*widget.Label)
					bindPin := binding.BindInt(&Pin.PinNumber)
					bindPinToString := binding.IntToString(bindPin)
					pinNumber.Bind(bindPinToString)

					topic := box.Objects[1].(*widget.Label)

					bindTopic := binding.BindString(&Pin.Topic)
					topic.Bind(bindTopic)

					bindMinValue := binding.BindFloat(&Pin.MinValue)
					bindMinValueToString := binding.FloatToString(bindMinValue)
					minValue := box.Objects[2].(*widget.Label)
					minValue.Bind(bindMinValueToString)

					bindMaxValue := binding.BindFloat(&Pin.MaxValue)
					bindMaxValueToString := binding.FloatToString(bindMaxValue)
					maxValue := box.Objects[3].(*widget.Label)
					maxValue.Bind(bindMaxValueToString)

					direction := box.Objects[4].(*widget.Label)
					bindDirection := binding.BindString(&Pin.Direction)
					direction.Bind(bindDirection)

					button := box.Objects[5].(*widget.Button)
					button.OnTapped = func() {
						calibratePinUI(Pin, app)
					}

				}

			})))

	form := widget.NewForm(
		widget.NewFormItem("ID", a.Id),
		widget.NewFormItem("Name", a.Name),
		widget.NewFormItem("Priority", a.Direction),
	)

	details := container.NewBorder(form, nil, nil, nil, a.Pins)

	return container.NewBorder(nil, nil, a.Controllers, nil, details)
}

func (a *CalibrationApp) refreshData() {
	// hide done
	a.visible = a.data.all()
	a.Controllers.Refresh()
}

func (a *CalibrationApp) setController(t *Controller) {
	a.current = t

	sort.Slice(a.current.Pins, func(i, j int) bool {
		return a.current.Pins[i].PinNumber < a.current.Pins[j].PinNumber
	})

	a.Name.SetText(t.Name)
	a.Id.SetText(t.Id)
	if t.Direction == input {
		a.Direction.SetText("In")
	} else {
		a.Direction.SetText("Out")
	}

	a.Pins.Refresh()

}
