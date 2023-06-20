package main

import (
	"bytes"
	"log"
	"net"
	"text/template"
)

func BuildPythonCode(controller *Controller) string {
	if controller.Direction == 0 {
		return InputController(controller)
	} else {
		return OutputController(controller)
	}
}

func InputController(controller *Controller) string {

	log.Println(GetLocalIP())

	publisherTemplate := `import network
	import socket
	import time
	import ujson
	import urequests
	from picozero import Pot # Pot is short for Potentiometer
	from machine import Pin, I2C
	from secret import ssid,password
	from umqtt.simple2 import MQTTClient
	import random
	wlan = network.WLAN(network.STA_IF)
	wlan.active(True)
	wlan.connect(ssid, password)
	
	print(ssid)
	
	# Wait for connect or fail
	max_wait = 10
	while max_wait > 0:
		if wlan.status() < 0 or wlan.status() >= 3:
			break
		max_wait -= 1
		print('waiting for connection...')
		time.sleep(1)
	
	# Handle connection error
	if wlan.status() != 3:
		raise RuntimeError('network connection failed')
	else:
		print('connected')
		status = wlan.ifconfig()
		print( 'ip = ' + status[0] )
	
	#r = urequests.get("http://192.168.0.98:8080/")
	#print(r.content)
	#r.close()
	#post_data = ujson.dumps({'some': 'garbage'})
	dial = Pot(2)
	mqtt_server = '192.168.0.98'
	client_id = 'pi1'
	topic_pub = b'TomsHardware'
	topic_msg = b'Movement Detected'
	
	def mqtt_connect():
		client = MQTTClient(client_id, mqtt_server, keepalive=3600)
		client.connect()
		print('Connected to %s MQTT Broker'%(mqtt_server))
		return client
	
	def reconnect():
		print('Failed to connect to the MQTT Broker. Reconnecting...')
		time.sleep(5)
		machine.reset()
	
	try:
		client = mqtt_connect()
	except OSError as e:
		reconnect()
	while True:
			print(dial.value)
			topic_msg = str(dial.value)
			client.publish(topic_pub, topic_msg)
			time.sleep(.1)`

	tmpl, err := template.New("publisher").Parse(publisherTemplate)
	if err != nil {
		panic(err)
	}
	buf := new(bytes.Buffer)
	tmpl.Execute(buf, controller)

	return buf.String()
}

func OutputController(controller *Controller) string {
	subscriberTemplate := `import network
import socket
import time
import ujson
import urequests
from machine import Pin, I2C
from secret import ssid,password
from umqtt.simple2 import MQTTClient
import random
wlan = network.WLAN(network.STA_IF)
wlan.active(True)
wlan.connect(ssid, password)

print(ssid)
# Wait for connect or fail
max_wait = 10
while max_wait > 0:
    if wlan.status() < 0 or wlan.status() >= 3:
        break
    max_wait -= 1
    print('waiting for connection...')
    time.sleep(1)

# Handle connection error
if wlan.status() != 3:
    raise RuntimeError('network connection failed')
else:
    print('connected')
    status = wlan.ifconfig()
    print( 'ip = ' + status[0] )

mqtt_server = '{{ .MqttServer }}'
client_id = '{{ .Name }}'
topic_pub = b'{{ .Topic }}'
topic_msg = b'Movement Detected'

def mqtt_connect():
    client = MQTTClient(client_id, mqtt_server, keepalive=3600)
    client.connect()
    print('Connected to %s MQTT Broker'%(mqtt_server))
    return client

def reconnect():
    print('Failed to connect to the MQTT Broker. Reconnecting...')
    time.sleep(5)
    machine.reset()

try:
    client = mqtt_connect()
except OSError as e:
    reconnect()
while True:
        client.publish(topic_pub, topic_msg)
        time.sleep(.1)`

	tmpl, err := template.New("subscriber").Parse(subscriberTemplate)
	if err != nil {
		panic(err)
	}
	buf := new(bytes.Buffer)
	tmpl.Execute(buf, controller)

	return buf.String()
}

func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}
