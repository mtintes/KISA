package main

import (
	"fmt"
	"log"

	"fyne.io/fyne/theme"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/google/uuid"
)

func configurePage(myApp fyne.App) *fyne.Container {
	border := configureBorder(myApp)

	configureContainer := container.New(layout.NewGridLayout(1), border)
	return configureContainer
}

func configureBorder(myApp fyne.App) *fyne.Container {
	var border *fyne.Container

	inputPicos := binding.NewUntypedList()
	for _, i := range mainSetup.InputPicos {
		inputPicos.Append(i)
	}

	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.ContentAddIcon(), func() {
			id := uuid.New().String()

			newPico := InputPico{Name: fmt.Sprintf("pico-%s", id[:4]), Inputs: []InputPin{
				InputPin{Topic: fmt.Sprintf("topic-%s", id[:4])},
			}}
			inputPicos.Append(newPico)
			mainSetup.InputPicos = append(mainSetup.InputPicos, newPico)
			for _, pico := range mainSetup.InputPicos {
				log.Println(pico.Name)
				for _, pin := range pico.Inputs {
					log.Println("topic: ", pin.Topic)
					log.Println("pin number: ", pin.PinNumber)
				}
			}
		}),
	)

	border = container.NewBorder(toolbar, nil, nil, nil, widget.NewListWithData(inputPicos,
		func() fyne.CanvasObject {
			return container.New(layout.NewAdaptiveGridLayout(7), widget.NewLabel(""))
		},
		func(di binding.DataItem, co fyne.CanvasObject) {
			v, _ := di.(binding.Untyped).Get()

			input := v.(InputPico)
			inputPins := binding.NewUntypedList()
			inputPins.AddListener(binding.NewDataListener(func() {
				newUpdate, _ := inputPins.Get()
				input.Inputs = interfaceToInputPin(newUpdate)
				//log.Println(input.Inputs)
			}))

			for _, pin := range input.Inputs {
				inputPins.Append(pin)
			}

			ctr, _ := co.(*fyne.Container)

			l := ctr.Objects[0].(*widget.Label)
			l.SetText(input.Name)
			// // a := ctr.Objects[2].(*widget.Label)
			if len(ctr.Objects) == 1 {
				ctr.Add(layout.NewSpacer())
				ctr.Add(layout.NewSpacer())
				ctr.Add(layout.NewSpacer())
				ctr.Add(layout.NewSpacer())
				ctr.Add(widget.NewButtonWithIcon("", theme.ContentRemoveIcon(), func() {
				}))

				ctr.Add(widget.NewButtonWithIcon("", theme.DocumentCreateIcon(), func() {
					configurePinsWindow(inputPins, myApp)
					log.Println(input.Name)
				}))
			}

			// if len(cont2.Objects) == 1 {
			// 	cont2.RemoveAll()
			// 	cont2.Add(widget.NewListWithData(inputPins, func() fyne.CanvasObject {
			// 		return widget.NewLabel("New")
			// 	},
			// 		func(di binding.DataItem, co fyne.CanvasObject) {
			// 			// v, _ := di.(binding.Untyped).Get()
			// 			// input := v.(inputPin)
			// 			// ctr, _ := co.(*fyne.Container)
			// 			// l := ctr.Objects[0].(*widget.Label)
			// 			// l.SetText(input.topic)
			// 			// log.Println(input)
			// 		}))
			// }

			// a.SetText(input.inputs[0].topic)
		}))

	return border
}
