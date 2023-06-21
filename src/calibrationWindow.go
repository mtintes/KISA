package main

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/rs/zerolog/log"
)

func calibratePinUI(pin *Pin, a fyne.App) {
	window := a.NewWindow("Calibrate Pin")
	window.Resize(fyne.Size{300, 300})
	quit := make(chan bool)
	startPositionButton := widget.NewButton("Start Position", func() {

		go listen(pin.Topic, quit)

	})

	middlePositionButton := widget.NewButton("Middle Position", func() {
		quit <- true
		log.Print("Middle Position")
	})

	endPositionButton := widget.NewButton("End Position", func() {
		log.Print("End Position")
	})

	window.SetContent(container.NewVBox(
		startPositionButton, middlePositionButton, endPositionButton,
	))
	window.Show()
}

func connect(clientId string) mqtt.Client {
	opts := createClientOptions(clientId)
	client := mqtt.NewClient(opts)
	token := client.Connect()
	for !token.WaitTimeout(3 * time.Second) {
	}
	if err := token.Error(); err != nil {
		log.Print(err)
	}
	return client
}

func createClientOptions(clientId string) *mqtt.ClientOptions {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s", "localhost:1883"))

	opts.SetClientID(clientId)
	return opts
}

func listen(topic string, quit chan bool) {
	client := connect("KISA Calibrator")
	client.Subscribe(topic, 0, func(client mqtt.Client, msg mqtt.Message) {
		fmt.Printf("* [%s] %s\n", msg.Topic(), string(msg.Payload()))

	OUT:
		for {
			select {
			case <-quit:
				client.Unsubscribe("a")
				client.Disconnect(10)
			default:
				break OUT
			}
		}

	})

}
