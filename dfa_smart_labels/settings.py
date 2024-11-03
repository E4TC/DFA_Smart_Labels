import os
from pathlib import Path
from warnings import filterwarnings

import django
import django_stubs_ext
from dotenv import load_dotenv

env_loaded_successfully = load_dotenv()

if not env_loaded_successfully:
    raise SystemExit(
        "Environment file not found. You must add a file called .env to the root directory of the project.\n"
        "An example can be found in .env.example",
    )


APP_VERSION = "0.0.1"

BASE_DIR = Path(__file__).resolve().parent.parent

# SECURITY WARNING: keep the secret key used in production secret!
SECRET_KEY = os.environ.get(
    "DJANGO_SECRET_KEY",
    "django-insecure-(o+xx3h!%72w3eaf#csx6onmnm)bl%3jmhrv7adjsj9m@+)*90",
)

# SECURITY WARNING: don't run with debug turned on in production!
DEBUG = os.environ.get("DJANGO_DEBUG", "false").lower() == "true"
USE_HTTPS = os.environ.get("DJANGO_USE_HTTPS", "true").lower() == "true"
DJANGO_ALLOWED_HOST = os.environ.get("DJANGO_ALLOWED_HOST", "localhost")

ALLOWED_HOSTS = [DJANGO_ALLOWED_HOST, "localhost", "127.0.0.1"]

__trusted_origin = f'{"https" if USE_HTTPS else "http"}://{DJANGO_ALLOWED_HOST}'
CSRF_TRUSTED_ORIGINS = [__trusted_origin, "http://localhost:8000", "http://localhost"]

if "FORCE_SCRIPT_NAME" in os.environ.keys():
    FORCE_SCRIPT_NAME = os.environ.get("FORCE_SCRIPT_NAME")

INSTALLED_APPS = [
    "django.contrib.admin",
    "django.contrib.auth",
    "django.contrib.contenttypes",
    "django.contrib.sessions",
    "django.contrib.messages",
    "django.contrib.staticfiles",
    "user_management",
    "rest_framework",
    "drf_spectacular",
    
]

if DEBUG:
    INSTALLED_APPS += ["django_extensions"]

MIDDLEWARE = [
    "django.middleware.security.SecurityMiddleware",
    "django.contrib.sessions.middleware.SessionMiddleware",
    "django.middleware.common.CommonMiddleware",
    "django.middleware.csrf.CsrfViewMiddleware",
    "django.contrib.auth.middleware.AuthenticationMiddleware",
    "django.contrib.messages.middleware.MessageMiddleware",
    "django.middleware.clickjacking.XFrameOptionsMiddleware",
]

ROOT_URLCONF = "dfa_smart_labels.urls"
AUTH_USER_MODEL = "user_management.User"

TEMPLATES = [
    {
        "BACKEND": "django.template.backends.django.DjangoTemplates",
        "DIRS": [BASE_DIR / "templates"],
        "APP_DIRS": True,
        "OPTIONS": {
            "context_processors": [
                "django.template.context_processors.debug",
                "django.template.context_processors.request",
                "django.contrib.auth.context_processors.auth",
                "django.contrib.messages.context_processors.messages",
                "dfa_smart_labels.context_processors.default",
            ],
        },
    },
]

WSGI_APPLICATION = "dfa_smart_labels.wsgi.application"

# Database
# https://docs.djangoproject.com/en/4.0/ref/settings/#databases

DATABASES = {
    "default": {
        "ENGINE": "django.db.backends.sqlite3",
        "NAME": BASE_DIR / "db.sqlite3",
    },
}

if os.environ.get("DB_USE_SQLITE", "false").lower() != "true":
    DATABASES["default"] = {
        "ENGINE": "django.db.backends.mysql",
        "HOST": os.environ.get("DB_HOST"),
        "PORT": os.environ.get("DB_PORT", 5432),
        "NAME": os.environ.get("DB_NAME"),
        "USER": os.environ.get("DB_USERNAME"),
        "PASSWORD": os.environ.get("DB_PASSWORD"),
    }

# Password validation
# https://docs.djangoproject.com/en/4.0/ref/settings/#auth-password-validators

AUTH_PASSWORD_VALIDATORS = [
    {
        "NAME": "django.contrib.auth.password_validation.UserAttributeSimilarityValidator",
    },
    {
        "NAME": "django.contrib.auth.password_validation.MinimumLengthValidator",
    },
    {
        "NAME": "django.contrib.auth.password_validation.CommonPasswordValidator",
    },
    {
        "NAME": "django.contrib.auth.password_validation.NumericPasswordValidator",
    },
]

# Internationalization
# https://docs.djangoproject.com/en/4.0/topics/i18n/

LANGUAGE_CODE = "de"
LOCALE_PATHS = [os.path.join(BASE_DIR, "locale")]

TIME_ZONE = "Europe/Berlin"

USE_I18N = True

USE_TZ = True

# Static files (CSS, JavaScript, Images)
# https://docs.djangoproject.com/en/4.0/howto/static-files/

STATIC_URL = os.environ.get("DJANGO_STATIC_URL", "/static/")
STATIC_ROOT = os.environ.get("DJANGO_STATIC_ROOT")
STATICFILES_DIRS = [
    BASE_DIR / "static",
]
MEDIA_URL = os.environ.get("DJANGO_MEDIA_URL", "/media/")
MEDIA_ROOT = os.environ.get("DJANGO_MEDIA_ROOT")

# Default primary key field type
# https://docs.djangoproject.com/en/4.0/ref/settings/#default-auto-field

DEFAULT_AUTO_FIELD = "django.db.models.BigAutoField"

DATA_UPLOAD_MAX_NUMBER_FIELDS = 10000

if os.environ.get("DJANGO_LOG_FILE", "") != "":
    LOGGING = {
        "version": 1,
        "disable_existing_loggers": False,
        "handlers": {
            "file": {
                "level": os.environ.get("DJANGO_LOG_LEVEL", "WARNING"),
                "class": "logging.FileHandler",
                "filename": os.environ.get("DJANGO_LOG_FILE"),
            },
        },
        "loggers": {
            "django": {
                "handlers": ["file"],
                "level": os.environ.get("DJANGO_LOG_LEVEL", "WARNING"),
                "propagate": True,
            },
        },
    }
REST_FRAMEWORK = {
    "DEFAULT_PAGINATION_CLASS": "rest_framework.pagination.LimitOffsetPagination",
    "PAGE_SIZE": 20,
    "DEFAULT_FILTER_BACKENDS": ("django_filters.rest_framework.DjangoFilterBackend",),
    "DEFAULT_SCHEMA_CLASS": "drf_spectacular.openapi.AutoSchema",
}

SPECTACULAR_SETTINGS = {
    "TITLE": "dfa_smart_labels",
    "DESCRIPTION": "dfa_smart_labels",
    "VERSION": APP_VERSION,
}

django_stubs_ext.monkeypatch()


if django.VERSION[0] < 6:
    # only important if URLField is used
    # https://adamj.eu/tech/2023/12/07/django-fix-urlfield-assume-scheme-warnings/
    filterwarnings("ignore", "The FORMS_URLFIELD_ASSUME_HTTPS transitional setting is deprecated.")
    FORMS_URLFIELD_ASSUME_HTTPS = True
