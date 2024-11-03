import paho.mqtt.client as mqtt
import os

MQTT_BROKER=os.environ.get("MQTT_BROKER")
MQTT_PORT=os.environ.get("MQTT_PORT")
MQTT_USER=os.environ.get("MQTT_USER")
MQTT_PASSWORD=os.environ.get("MQTT_PASSWORD")


def on_connect(client, userdata, flags, rc):
    print("Connected with result code " + str(rc))

def on_message(client, userdata, msg):
    print(msg.topic + " " + str(msg.payload))

def start_mqtt_client():
    client = mqtt.Client()
    client.on_connect = on_connect
    client.on_message = on_message
    client.username_pw_set(MQTT_USER, MQTT_PASSWORD)
    client.connect(MQTT_BROKER, MQTT_PORT, 60)
    client.subscribe("dfa/order/#")
    client.loop_start()