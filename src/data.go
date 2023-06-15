package main

const (
	input  = 0
	output = 1
)

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

func dummyData() *ControllerList {
	return &ControllerList{
		Controllers: []*Controller{
			{Name: "New Pico", Id: "1", Pins: []*Pin{
				{PinNumber: 2},
				{PinNumber: 3},
			}},
			{Name: "New Pico", Id: "1", Pins: []*Pin{
				{PinNumber: 1, Topic: "a"},
				{PinNumber: 2, Topic: "b"},
				{PinNumber: 3, Topic: "c"},
				{PinNumber: 4, Topic: "d"},
				{PinNumber: 5, Topic: "e"},
				{PinNumber: 6, Topic: "f"},
				{PinNumber: 7, Topic: "g"},
				{PinNumber: 8, Topic: "h"},
			}},
		}}
}
