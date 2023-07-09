package main

import (
	"encoding/json"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
)

const (
	input  = 0
	output = 1
)

type ConnectionList struct {
	Connections []*Connection
}

type Connection struct {
	Input  *SimpleController
	Output *SimpleController
}

type SimpleController struct {
	Name      string
	PinNumber int
	Topic     string
}

type Pin struct {
	Topic     string
	PinNumber int
	MinValue  float64
	MaxValue  float64
	Direction string
}
type Controller struct {
	Id, Name string

	Direction int
	Pins      []*Pin
}

type ControllerList struct {
	filepath    string
	Controllers []*Controller
}

func (l *ControllerList) add(t *Controller) {
	l.Controllers = append([]*Controller{t}, l.Controllers...)
}

func (l *ControllerList) all() []*Controller {
	var items []*Controller
	items = append(items, l.Controllers...)

	return items
}

func (l *ControllerList) inputs() *ConnectionList {
	var items []*Connection

	for _, c := range l.Controllers {
		if c.Direction == input {
			items = append(items, &Connection{
				Input: &SimpleController{
					Name: c.Name,
				},
				Output: &SimpleController{
					Name: "",
				},
			})
		}
	}

	return &ConnectionList{Connections: items}
}

func (l *ControllerList) save(w fyne.Window) {
	if l.filepath == "" {
		dialog.NewFileSave(func(u fyne.URIWriteCloser, err error) {
			if err == nil && u != nil {
				l.filepath = u.URI().Path()
				d1, err := json.Marshal(l.Controllers)
				if err != nil {
					panic(err)
				}

				u.Write([]byte(d1))
				u.Close()
			}
		}, w).Show()
		return
	} else {
		d1, err := json.Marshal(l.Controllers)
		if err != nil {
			panic(err)
		}
		err = os.WriteFile(l.filepath, d1, 0644)
		if err != nil {
			panic(err)
		}
	}
}

func (l *ControllerList) saveAs(w fyne.Window) {
	dialog.NewFileSave(func(u fyne.URIWriteCloser, err error) {
		if err == nil && u != nil {
			l.filepath = u.URI().Path()
			d1, err := json.Marshal(l.Controllers)
			if err != nil {
				panic(err)
			}
			u.Write([]byte(d1))
			u.Close()
		}
	}, w).Show()
}

func (l *ControllerList) load(w fyne.Window, apps *Apps, app fyne.App) {
	dialog.NewFileOpen(func(u fyne.URIReadCloser, err error) {
		l.filepath = u.URI().Path()
		data, err := os.ReadFile(l.filepath)
		if err != nil {
			panic(err)
		}
		err = json.Unmarshal(data, &l.Controllers)
		if err != nil {
			panic(err)
		}

		tabs := getTabs(apps, w, app)

		w.SetContent(tabs)
		megaRefresh(apps)
	}, w).Show()

}

func dummyData() *ControllerList {
	return &ControllerList{
		Controllers: []*Controller{
			{
				Id:        "1",
				Name:      "Controller 1",
				Direction: input,
				Pins: []*Pin{
					{
						Topic:     "a",
						PinNumber: 1,
					},
				},
			},
		}}
}

func megaRefresh(apps *Apps) {
	apps.Controllers.refreshData()
	apps.Calibrations.refreshData()
	apps.Connections.refreshData()
}
