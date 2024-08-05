package main

import (
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"time"
)

func GetClient(id_optional ...string) mqtt.Client {
	opts := mqtt.NewClientOptions()
	opts.AddBroker("10.101.206.6:1883")
	opts.SetUsername("e4tc")
	opts.SetPassword("6gD$kQ2o9^Fa956f")

	id := "e4tc_client"
	if len(id_optional) > 0 {
		id = id_optional[0]
	}
	opts.SetClientID(id)

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		fmt.Sprintf("Error connecting to MQTT broker:", token.Error())
	}
	return client
}

func AnnounceService(service string) {
	client := GetClient("e4tc_dfa_announce_client")
	event := &Event{Timestamp: time.Now().UnixNano() / 1e6, TimestampHr: time.Now().Format("2006-01-02 15:04:05"), Event: "announcement"}
	jsonObj, _ := json.Marshal(event)
	token := client.Publish("meta/svc/service_"+service+"/announced", 0, true, jsonObj)
	if token.Wait() && token.Error() != nil {
		fmt.Sprintf("Error connecting to MQTT broker:", token.Error())
	}
	client.Disconnect(100)
}

func StartHeartbeat(service string) {
	client := GetClient("e4tc_dfa_heartbeat_client")
	for range time.Tick(time.Second * 5) {
		event := &Event{Timestamp: time.Now().UnixNano() / 1e6, TimestampHr: time.Now().Format("2006-01-02 15:04:05"), Event: "heartbeat"}
		jsonObj, _ := json.Marshal(event)
		token := client.Publish("meta/svc/service_"+service+"/heartbeat", 0, false, jsonObj)
		if token.Wait() && token.Error() != nil {
			fmt.Sprintf("Error connecting to MQTT broker:", token.Error())
		}
	}
}

func on_message_receive() {

}

func Publish(topic string, jsonObj []byte) {
	client := GetClient("e4tc_dfa_pub_client")
	token := client.Publish(topic, 0, false, jsonObj)
	if token.Wait() && token.Error() != nil {
		fmt.Sprintf("Error connecting to MQTT broker:", token.Error())
	}
	client.Disconnect(100)
}

// dfa/#
type Event struct {
	Timestamp   int64  `json:"timestamp"`
	TimestampHr string `json:"timestamp_hr"`
	Event       string `json:"event"`
}
