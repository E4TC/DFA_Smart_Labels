# dfa_smart_labels

## Development Setup

### Prerequisite

- `nox` is available via command line (e.g. installed with `pipx`).
- Python 3.12 is installed (see `pyenv`)

### Backend

To set up the django project execute the following commands. In the last step this will start the django
development server

```shell
nox -s create_dev_env
# activate venv (e.g. .venv/bin/activate)

git init
git lfs install
pre-commit install

python manage.py migrate
python manage.py createsuperuser  # Enter username and password
python manage.py runserver
```

### Troubleshooting MacOS

There may be a problem installing the GraphVIZ dependency required by `django-extensions` on MacOS with 
Apple Silicon CPUs. This can be resolved by setting some environment variable for the compiler:

```shell
export CFLAGS="-I $(brew --prefix graphviz)/include"
export LDFLAGS="-L $(brew --prefix graphviz)/lib"
```

Executing these command before running `pip install requirements.txt` or `pre-commit` should solve the problem.

## Requirement files

 - `requirement.txt`: includes all direct project related Python dependencies
 - `requirement.mypy`: includes additional Python dependencies to run mypy
 - `requirement.prod`: includes additional Python dependencies to run in a production environment (e.g. PostgreSQL DB)


## Environment Variablen

| Variable               | Default         | Description                                                  |
|------------------------|-----------------|--------------------------------------------------------------|
| `DJANGO_SECRET_KEY`    | An insecure key | Secret key for the django project                            |
| `DJANGO_DEBUG`         | false           | Set to `true` run the Django project as root                 |
| `DJANGO_ALLOWED_HOSTS` | localhost       | Set the allowed host (separated by `,`)                      |
| `DJANGO_STATIC_URL`    | /static/        | URL to static files                                          |
| `DJANGO_STATIC_ROOT`   |                 | Path to static files                                         |
| `DJANGO_MEDIA_URL`     | /media/         | URL to media files                                           |
| `DJANGO_MEDIA_ROOT`    |                 | Path to media files                                          |
| `DJANGO_LOG_FILE`      |                 | Path to log file (no logging if not specified)               |
| `DJANGO_LOG_LEVEL`     | WARNING         | Log Level for Django (DEBUG, INFO, WARNING, ERROR, CRITICAL) |
| `DB_USE_SQLITE`        | false           | Set to `true` to use SQLite instead of MySQL / MariaDB       |
| `DB_HOST`              |                 | Host of the database                                         |
| `DB_PORT`              | 3306            | Port of the database                                         |
| `DB_NAME`              |                 | Name of the database                                         |
| `DB_USERNAME`          |                 | Username for the database                                    |
| `DB_PASSWORD`          |                 | Password for the database                                    |

## Deployment

When deploying the app the following steps and commands should be executed

```shell
# 1. Set up and activate virtual environment
# 2. Copy .env.example to .env and set the at least the following environment variables
#    - DJANGO_SECRET_KEY=some_secret_key
#    - DJANGO_DEBUG=False
#    - DJANGO_ALLOWED_HOSTS=your.domain.com
#    - DJANGO_STATIC_ROOT=/path/to/static/files
#    - DJANGO_MEDIA_ROOT=/path/to/media/files
#    - DJANGO_LOG_FILE/path/to/log/file
#    - DB_HOST=host_of_database
#    - DB_PORT=port_of_database
#    - DB_NAME=name_of_database
#    - DB_USERNAME=user_of_database
#    - DB_PASSWORD=password_for_database_user

pip install -r requirements.txt -r requirements.prod
python manage.py migrate
python manage.py compilemessages
python manage.py collectstatic
```

## Commands

### Generating Translations
If you are using the JavaScript translation functions from django `gettext`, `ngettext`, ... there is a simple way to generate and compile translations for them as well. This script will also work if you are not using the JavaScript functions.

```shell
./scripts/makemessages.sh
```

### Generating ER Diagramm
Execute the following command to generate an ER-Diagram of the models. This command will only work if executed
with `DEBUG = True`

```shell
nox -s graph_models
```
