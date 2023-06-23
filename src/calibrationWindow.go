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

	var instructionsLabel *widget.Label
	var startCalibrationButtonOn *widget.Button
	var exitButton *widget.Button
	client := connect("KISA Calibrator")
	quit := make(chan bool)

	startCalibrationButtonOn = widget.NewButton("Start Calibration", func() {
		go listen(client, pin.Topic, quit)
		startCalibrationButtonOn.Hidden = true
		exitButton.Hidden = false
	})

	exitButton = widget.NewButton("End Position Done", func() {
		quit <- true
		client.Unsubscribe("a")
		client.Disconnect(250)
		window.Close()
	})

	exitButton.Hidden = true

	instructionsLabel = widget.NewLabel("Put your mechanism in a closed postion (minimum value) then click Start Calibration button to begin calibration")
	instructionsLabel.Wrapping = fyne.TextWrapWord
	window.SetContent(container.NewGridWithRows(2,
		instructionsLabel, startCalibrationButtonOn, exitButton,
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

func listen(client mqtt.Client, topic string, quit chan bool) {
	client.Subscribe(topic, 0, func(client mqtt.Client, msg mqtt.Message) {
		fmt.Printf("* [%s] %s\n", msg.Topic(), string(msg.Payload()))

		//This has some garbage process that requires
		//a transaction to come in before it will
		//close the connection
	OUT:
		for {
			select {
			case <-quit:
				//close(quit)
				client.Unsubscribe("a")
				client.Disconnect(250)
				break OUT
			default:
				break OUT
			}
		}

	})

}
