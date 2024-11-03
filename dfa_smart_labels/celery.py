from __future__ import absolute_import, unicode_literals
import os
from celery import Celery

os.environ.setdefault("DJANGO_SETTINGS_MODULE", "dfa_smart_labels.settings")
app = Celery("dfa_smart_label")
app.config_from_object("django.conf:settings", namespace="CELERY")
app.autodiscover_tasks()