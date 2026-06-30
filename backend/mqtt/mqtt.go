package mqttclient

import (
	"crypto/tls"
	"encoding/json"
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	"smart-door-backend/constant"
	"smart-door-backend/handlers"
	data "smart-door-backend/Types"
)

var Client mqtt.Client

// =====================================================
// MESSAGE HANDLERS - dipanggil saat ada pesan masuk dari ESP32
// =====================================================

var sensorHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("[MQTT] sensor: %s\n", msg.Payload())
	var sensor data.SensorData
	if err := json.Unmarshal(msg.Payload(), &sensor); err != nil {
		fmt.Println("[MQTT] Gagal parse sensor JSON:", err)
		return
	}
	if err := handlers.SaveSensorData(sensor); err != nil {
		fmt.Println("[MongoDB] Gagal simpan sensor:", err)
	}
}

var accessHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("[MQTT] access: %s\n", msg.Payload())
	var access data.AccessLog
	if err := json.Unmarshal(msg.Payload(), &access); err != nil {
		fmt.Println("[MQTT] Gagal parse access JSON:", err)
		return
	}
	if err := handlers.SaveAccessLog(access); err != nil {
		fmt.Println("[MongoDB] Gagal simpan access log:", err)
	}
}

var doorHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("[MQTT] door: %s\n", msg.Payload())
	var door data.DoorStatus
	if err := json.Unmarshal(msg.Payload(), &door); err != nil {
		fmt.Println("[MQTT] Gagal parse door JSON:", err)
		return
	}
	if err := handlers.SaveDoorStatus(door); err != nil {
		fmt.Println("[MongoDB] Gagal simpan door status:", err)
	}
}

var statusHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("[MQTT] status: %s\n", msg.Payload())
	var status data.DeviceStatus
	if err := json.Unmarshal(msg.Payload(), &status); err != nil {
		fmt.Println("[MQTT] Gagal parse status JSON:", err)
		return
	}
	if err := handlers.SaveDeviceStatus(status); err != nil {
		fmt.Println("[MongoDB] Gagal simpan device status:", err)
	}
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("[MQTT] Connected to HiveMQ Cloud")
	subscribeAll(client)
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("[MQTT] Connection lost: %v\n", err)
}

// =====================================================
// INIT & CONNECT
// =====================================================

func InitMQTT() error {
	broker := constant.MQTT_URL
	port := constant.MQTT_PORT

	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tls://%s:%s", broker, port))
	opts.SetClientID("go-backend-smart-door")
	opts.SetUsername(constant.MQTT_USERNAME)
	opts.SetPassword(constant.MQTT_PASSWORD)
	opts.SetTLSConfig(&tls.Config{InsecureSkipVerify: true})
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	opts.SetAutoReconnect(true)
	opts.SetConnectRetry(true)

	Client = mqtt.NewClient(opts)
	if token := Client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func subscribeAll(client mqtt.Client) {
	topics := map[string]mqtt.MessageHandler{
		constant.TOPIC_SENSOR: sensorHandler,
		constant.TOPIC_ACCESS: accessHandler,
		constant.TOPIC_DOOR:   doorHandler,
		constant.TOPIC_STATUS: statusHandler,
	}

	for topic, handler := range topics {
		token := client.Subscribe(topic, 1, handler)
		token.Wait()
		if token.Error() != nil {
			fmt.Printf("[MQTT] Gagal subscribe %s: %v\n", topic, token.Error())
		} else {
			fmt.Printf("[MQTT] Subscribed: %s\n", topic)
		}
	}
}

// =====================================================
// PUBLISH COMMAND - dipanggil dari HTTP endpoint Frontend
// =====================================================

func PublishDoorCommand(command string) error {
	cmd := data.DoorCommand{Command: command}
	payload, err := json.Marshal(cmd)
	if err != nil {
		return err
	}

	token := Client.Publish(constant.TOPIC_CMD_DOOR, 1, false, payload)
	token.Wait()
	return token.Error()
}

func PublishReset() error {
	token := Client.Publish(constant.TOPIC_CMD_RESET, 1, false, "{}")
	token.Wait()
	return token.Error()
}
