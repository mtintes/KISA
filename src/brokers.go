package main

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/mochi-co/mqtt/v2"
	"github.com/mochi-co/mqtt/v2/hooks/auth"
	"github.com/mochi-co/mqtt/v2/listeners"
)

func makeBrokerUI(w fyne.Window) fyne.CanvasObject {
	var server *mqtt.Server
	var startButton *widget.Button
	var stopButton *widget.Button

	startButton = widget.NewButtonWithIcon("Start Broker", theme.MediaPlayIcon(), func() {
		server = StartServer()
		stopButton.Hidden = false
		startButton.Hidden = true
	})
	stopButton = widget.NewButtonWithIcon("Stop Broker", theme.MediaStopIcon(), func() {
		err := server.Close()
		if err != nil {
			log.Fatal(err)
		}
		startButton.Hidden = false
		stopButton.Hidden = true
	})
	stopButton.Hidden = true

	return container.NewCenter(startButton, stopButton)
}

func StartServer() *mqtt.Server {
	server := mqtt.New(nil)
	_ = server.AddHook(new(auth.AllowHook), nil)

	tcp := listeners.NewTCP("t1", ":1883", nil)

	err := server.AddListener(tcp)
	if err != nil {
		log.Fatal(err)
	}

	err = server.Serve()
	if err != nil {
		log.Fatal(err)
	}
	return server
}
