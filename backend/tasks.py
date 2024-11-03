import json
import os
import paho.mqtt.client as mqtt
from celery import shared_task
from django.utils.dateparse import parse_datetime
from .models import Order, OrderPosition


MQTT_BROKER=os.environ.get("MQTT_BROKER")
MQTT_PORT=int(os.environ.get("MQTT_PORT"))
MQTT_USER=os.environ.get("MQTT_USER")
MQTT_PASSWORD=os.environ.get("MQTT_PASSWORD")
MQTT_TOPIC = "dfa/order/#"

@shared_task
def process_mqtt_message(data):
    order_data = data.get("new")
    order_id = order_data.get("order_id")

    if Order.objects.filter(order_id=order_id).exists():
        print(f"Order {order_id} already exists. Skipping.")
        return

    order = Order.objects.create(
        order_id=order_id,
        description=order_data.get("description"),
        customer=order_data.get("customer"),
        created=order_data.get("created"),
        created_hr=parse_datetime(order_data.get("created_hr")),
    )

    for position_data in order_data.get("positions", []):
        OrderPosition.objects.create(
            guid=position_data.get("guid"),
            order=order,
            position=position_data.get("position"),
            description=position_data.get("article"),
            quantity=position_data.get("quantity"),
            state=position_data.get("state"),
        )

    print(f"Order {order_id} created successfully.")

def on_connect(client, userdata, flags, rc):
    print("Connected to MQTT broker with result code " + str(rc))
    client.subscribe(MQTT_TOPIC)

def on_message(client, userdata, msg):
    data = json.loads(msg.payload.decode())
    process_mqtt_message.delay(data) 

def start_mqtt_listener():
    client = mqtt.Client()
    client.on_connect = on_connect
    client.on_message = on_message
    client.username_pw_set(MQTT_USER, MQTT_PASSWORD)
    client.connect(MQTT_BROKER, MQTT_PORT, 60)
    client.loop_forever()