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

func interfaceToInputPin(in []interface{}) []InputPin {
	var response []InputPin

	for _, object := range in {
		response = append(response, object.(InputPin))
	}

	return response
}
