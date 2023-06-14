package main

import (
	"fmt"
	"reflect"
	"strconv"

	"fyne.io/fyne/theme"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/google/uuid"
)

func configurePinsWindow(inputPins binding.UntypedList, app fyne.App) {

	configurePinsWindow := app.NewWindow("Input Pins")
	configurePinsWindow.SetContent(inputsContainer(inputPins))
	configurePinsWindow.Resize(fyne.NewSize(500, 500))
	configurePinsWindow.Show()
}

func inputsContainer(inputPins binding.UntypedList) *fyne.Container {

	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.ContentAddIcon(), func() {
			id := uuid.New().String()

			inputPins.Append(InputPin{Topic: fmt.Sprintf("topic-%s", id[:4])})

		}),
	)
	labelbar := container.New(layout.NewGridLayout(2), toolbar, widget.NewLabel(""), widget.NewLabel("topic"), widget.NewLabel("pinNumber"))

	border := container.NewBorder(labelbar, nil, nil, nil, widget.NewListWithData(inputPins,
		func() fyne.CanvasObject {
			return container.New(layout.NewGridLayout(2), widget.NewLabel("a"))
		},
		func(di binding.DataItem, co fyne.CanvasObject) {
			v, _ := di.(binding.Untyped).Get()

			input := v.(InputPin)
			ctr, _ := co.(*fyne.Container)
			if len(ctr.Objects) == 1 {

				ctr.RemoveAll()
				topicEntry := widget.NewEntry()
				topicEntry.SetText(input.Topic)
				topicEntry.OnChanged = func(s string) {
					val := reflect.ValueOf(di.(interface{})).Elem().FieldByName("index").Int()
					input.Topic = s
					inputPins.SetValue(int(val), input)
				}

				pinNumber := widget.NewEntry()
				pinNumber.SetText(strconv.FormatInt(int64(input.PinNumber), 10))
				pinNumber.OnChanged = func(s string) {
					val := reflect.ValueOf(di.(interface{})).Elem().FieldByName("index").Int()
					intConv, _ := strconv.Atoi(s)
					input.PinNumber = intConv
					inputPins.SetValue(int(val), input)
				}

				ctr.Add(topicEntry)
				ctr.Add(pinNumber)

				//ctr.Add(entryLabel(input.Topic, "Topic", &input, &inputPins, di))

			}

		}))

	return border
}

// func entryLabel(s string, field string, input *InputPin, boundList *binding.UntypedList, di binding.DataItem) *widget.Entry {
// 	entry := widget.NewEntry()
// 	entry.SetText(s)
// 	entry.OnChanged = func(updateString string) {
// 		index := reflect.ValueOf(di.(interface{})).Elem().FieldByName("index").Int()
// 		updatedInput := reflect.ValueOf(input).Elem().FieldByName(field)

// 		//test.SetString(updateString) //.SetString(updateString)
// 		fmt.Println("test", updatedInput.CanSet())
// 		updatedInput.SetString(s)
// 		boundList.SetValue(int(index), *input)
// 	}
// 	return entry
// }
