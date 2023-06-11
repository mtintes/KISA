package main

type setup struct {
	InputPicos []InputPico
}

type InputPico struct {
	Name   string
	Inputs []InputPin
}

type InputPin struct {
	Topic     string
	PinNumber int
	MinValue  float64
	MaxValue  float64
	Direction string
}

func NewInputPins(topic string, pin int) InputPin {
	return InputPin{
		Topic:     topic,
		PinNumber: pin,
	}
}
