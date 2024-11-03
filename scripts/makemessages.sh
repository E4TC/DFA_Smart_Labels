#!/bin/bash

set -e

./manage.py makemessages --locale de -i "venv/*"
./manage.py makemessages_djangojs --locale de -i "node_modules/*" -i "venv/*" -d djangojs -e 'tsx' -lang Python;
./manage.py compilemessages -v 0
