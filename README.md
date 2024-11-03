# dfa_smart_labels
# activate venv (e.g. .venv/bin/activate)

git init
git lfs install
pre-commit install

python manage.py migrate
python manage.py createsuperuser  # Enter username and password
python manage.py runserver

# start MQTT Listener
python manage.py listen_mqtt

# start Celery worker:
celery -A myproject worker --loglevel=info