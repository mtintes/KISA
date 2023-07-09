package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type ConnectionApp struct {
	data    *ControllerList
	inputs  []*Connection
	current *Connection

	Connections *widget.List
}

func (a *ConnectionApp) makeConnectionUI() *fyne.Container {
	a.Connections = widget.NewList(
		func() int {
			return len(a.inputs)
		},
		func() fyne.CanvasObject {
			item := container.NewWithoutLayout(widget.NewLabel("Item x"))
			return item
		},
		func(i int, c fyne.CanvasObject) {
			Connection := a.inputs[i]
			box := c.(*fyne.Container)
			label := box.Objects[0].(*widget.Label)
			label.SetText(Connection.Input.Name)
		})

	a.Connections.OnSelected = func(id int) {
		a.setConnection(a.inputs[id])
	}

	return container.NewBorder(nil, nil, nil, nil, a.Connections)
}

func (a *ConnectionApp) setConnection(Connection *Connection) {
	a.current = Connection
	a.Connections.Refresh()
}

func (a *ConnectionApp) refreshData() {
	// hide done
	a.inputs = a.data.inputs().Connections
	a.Connections.Refresh()
}
