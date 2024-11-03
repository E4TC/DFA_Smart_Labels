from django.core.management.base import BaseCommand
from backend.tasks import start_mqtt_listener

class Command(BaseCommand):
    help = "Start the MQTT listener"

    def handle(self, *args, **kwargs):
        start_mqtt_listener()