package main

import (
	"fmt"
	"strconv"
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
	phase := make(chan string)
	currentPhaseToSend := 0

	startCalibrationButtonOn = widget.NewButton("Start Calibration", func() {

		switch currentPhaseToSend {
		case 0:
			go listen(client, pin, quit, phase, startCalibrationButtonOn, instructionsLabel)
			phase <- "start"
			currentPhaseToSend++

		case 1:
			phase <- "middle"
			currentPhaseToSend++
		case 2:
			phase <- "end"
		}
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

func listen(client mqtt.Client, pin *Pin, quit chan bool, next chan string, startCalibrationButtonOn *widget.Button, instructionsLabel *widget.Label) {
	biggest := 0.0
	smallest := 0.0
	count := 0
	currentPhase := ""

	client.Subscribe(pin.Topic, 0, func(client mqtt.Client, msg mqtt.Message) {
		fmt.Printf("* [%s] %s\n", msg.Topic(), string(msg.Payload()))
		value, _ := strconv.ParseFloat(string(msg.Payload()), 64)
		if currentPhase == "start" {
			if count < 10 {
				fmt.Println("starting phase")
				pin.MinValue = calculateMinAverageInput(biggest, smallest, value)
				count++
				startCalibrationButtonOn.Disable()
			} else {
				startCalibrationButtonOn.Enable()
				fmt.Println("waiting for next phase")
				instructionsLabel.SetText("Put your mechanism in a middle position then click 'Start Calibration' button")
			}
		} else if currentPhase == "middle" {
			if count < 10 {
				fmt.Println("middle phase")
				pin.Direction = calculateDirectionInput(biggest, smallest, value)
				count++
				startCalibrationButtonOn.Disable()
			} else {
				instructionsLabel.SetText("Put your mechanism in a fully open position (max value) then click 'Start Calibration' button")
				fmt.Println("waiting for next phase")
				startCalibrationButtonOn.Enable()
			}
		} else if currentPhase == "end" {
			if count < 10 {
				fmt.Println("end phase")
				pin.MaxValue = calculateMaxAverageInput(biggest, smallest, value)
				count++
				startCalibrationButtonOn.Disable()
			} else {
				fmt.Println("waiting for next phase")
				startCalibrationButtonOn.Enable()
			}
		}

		select {
		case newPhase := <-next:
			currentPhase = newPhase
			count = 0
		case <-quit:
			client.Unsubscribe("a")
			client.Disconnect(250)
			return
		default:
		}
	})

}

func calculateMinAverageInput(biggest, smallest, value float64) float64 {

	if value > biggest {
		biggest = value
	} else if value < smallest {
		smallest = value
	}

	return 0.0
}

func calculateMaxAverageInput(biggest, smallest, value float64) float64 {

	if value > biggest {
		biggest = value
	} else if value < smallest {
		smallest = value
	}

	return 0.0
}

func calculateDirectionInput(biggest, smallest, value float64) string {

	return "clockwise"
}
