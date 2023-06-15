package main

import (
	"sort"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/google/uuid"
)

type ControllerApp struct {
	data    *ControllerList
	visible []*Controller
	current *Controller

	Controllers *widget.List
	Id          *widget.Label
	Name        *widget.Entry
	Direction   *widget.RadioGroup
	Pins        *fyne.Container
}

func (a *ControllerApp) makeUI() fyne.CanvasObject {
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

	a.Name = widget.NewEntry()
	a.Name.OnChanged = func(text string) {
		if a.current == nil {
			return
		}

		a.current.Name = text
		a.Controllers.Refresh()
	}

	a.Direction = widget.NewRadioGroup([]string{"In", "Out"}, func(direction string) {
		if a.current == nil {
			return
		}

		if direction == "In" {
			a.current.Direction = input
		} else {
			a.current.Direction = output
		}
	})

	a.Pins = container.New(layout.NewBorderLayout(nil, nil, nil, nil), widget.NewList(
		func() int {
			if a.current != nil {
				return len(a.current.Pins)
			}
			return 1
		},
		func() fyne.CanvasObject {
			pins := container.NewGridWithColumns(2, widget.NewEntry(), widget.NewEntry())
			return pins
		}, func(lii widget.ListItemID, co fyne.CanvasObject) {
			if a.current != nil {
				Pin := a.current.Pins[lii]
				box := co.(*fyne.Container)
				pinNumber := box.Objects[0].(*widget.Entry)
				bindPin := binding.BindInt(&Pin.PinNumber)
				bindPinToString := binding.IntToString(bindPin)
				pinNumber.Bind(bindPinToString)
				pinNumber.PlaceHolder = "pin #"

				topic := box.Objects[1].(*widget.Entry)
				topic.PlaceHolder = "topic_name"

				bindTopic := binding.BindString(&Pin.Topic)
				topic.Bind(bindTopic)

			}

		}))

	form := widget.NewForm(
		widget.NewFormItem("ID", a.Id),
		widget.NewFormItem("Name", a.Name),
		widget.NewFormItem("Priority", a.Direction),
	)

	details := container.NewBorder(form, nil, nil, nil, a.Pins)
	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.ContentAddIcon(), func() {
			newId := uuid.New()
			Controller := &Controller{
				Id:   newId.String(),
				Name: "New Pico",
				Pins: []*Pin{
					{PinNumber: 4},
				}}
			a.data.add(Controller)
			a.setController(Controller)
			a.refreshData()
		}),
	)
	return container.NewBorder(toolbar, nil, a.Controllers, nil, details)
}

func (a *ControllerApp) refreshData() {
	// hide done
	a.visible = a.data.all()
	a.Controllers.Refresh()
}

func (a *ControllerApp) setController(t *Controller) {
	a.current = t

	sort.Slice(a.current.Pins, func(i, j int) bool {
		return a.current.Pins[i].PinNumber < a.current.Pins[j].PinNumber
	})

	a.Name.SetText(t.Name)
	a.Id.SetText(t.Id)
	if t.Direction == input {
		a.Direction.SetSelected("In")
	} else {
		a.Direction.SetSelected("Out")
	}

	a.Pins.Refresh()

}

func main() {
	a := app.New()
	w := a.NewWindow("Controller List")
	w.Resize(fyne.Size{500, 500})

	data := dummyData()
	Controllers := &ControllerApp{data: data, visible: data.all()}
	w.SetContent(Controllers.makeUI())
	if len(data.all()) > 0 {
		Controllers.setController(data.all()[0])
	}
	w.ShowAndRun()
}
