import os
from pathlib import Path
from decouple import config

BASE_DIR = Path(__file__).resolve().parent.parent

SECRET_KEY = os.environ.get("SECRET_KEY", config("SECRET_KEY", default="DJANGO_SECRET_KEY"))

DEBUG = bool(os.environ.get("DEBUG", config("DEBUG", "False")).title()) == True

ALLOWED_HOSTS = ["*"]

INSTALLED_APPS = [
    "django.contrib.admin",
    "django.contrib.auth",
    "django.contrib.contenttypes",
    "django.contrib.sessions",
    "django.contrib.messages",
    "django.contrib.staticfiles",
    "rest_framework",
    "rest_framework.authtoken",
    "corsheaders",
    "app.apps.AppConfig",
]

REST_FRAMEWORK = {
    "DEFAULT_AUTHENTICATION_CLASSES": [
        "rest_framework.authentication.TokenAuthentication",
    ],
}

MIDDLEWARE = [
    "django.middleware.security.SecurityMiddleware",
    "django.contrib.sessions.middleware.SessionMiddleware",
    "corsheaders.middleware.CorsMiddleware",
    "django.middleware.common.CommonMiddleware",
    "django.middleware.csrf.CsrfViewMiddleware",
    "django.contrib.auth.middleware.AuthenticationMiddleware",
    "django.contrib.messages.middleware.MessageMiddleware",
    "django.middleware.clickjacking.XFrameOptionsMiddleware",
]

ROOT_URLCONF = "ocserv.urls"

TEMPLATES = [
    {
        "BACKEND": "django.template.backends.django.DjangoTemplates",
        "DIRS": [],
        "APP_DIRS": True,
        "OPTIONS": {
            "context_processors": [
                "django.template.context_processors.debug",
                "django.template.context_processors.request",
                "django.contrib.auth.context_processors.auth",
                "django.contrib.messages.context_processors.messages",
            ],
        },
    },
]

WSGI_APPLICATION = "ocserv.wsgi.application"

if not os.path.exists("./db"):
    Path("./db").mkdir(parents=True, exist_ok=True)

DATABASES = {
    "default": {
        "ENGINE": "django.db.backends.sqlite3",
        "NAME": BASE_DIR / "db/db.sqlite3",
    }
}

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

LANGUAGE_CODE = "en-us"

TIME_ZONE = "UTC"

USE_I18N = True

USE_TZ = True

STATIC_URL = "static/"

DEFAULT_AUTO_FIELD = "django.db.models.BigAutoField"

if DEBUG:
    if os.system("pip freeze | grep drf-yasg") != 0:
        os.system("python3 -m pip install drf-yasg")
    CORS_ALLOW_ALL_ORIGINS = True
    LOG_PATH = Path.joinpath(BASE_DIR, "log.txt")
    INSTALLED_APPS += ["drf_yasg"]
    SWAGGER_SETTINGS = {
        "USE_SESSION_AUTH": False,
        "SECURITY_DEFINITIONS": {
            "Token": {
                "type": "apiKey",
                "name": "Authorization",
                "in": "header",
                "description": "Authorization: Token TOKEN",
            }
        },
        "SHOW_REQUEST_HEADERS": True,
    }
else:
    CORS_ALLOW_ALL_ORIGINS = False
    CORS_ALLOWED_ORIGINS = list(
        filter(None, os.environ.get("CORS_ALLOWED", config("CORS_ALLOWED")).split(","))
    )
    LOG_PATH = "/var/log/backend.log"

SOCKET_PASSWD_FILE = "/var/log/socket_passwd"

OSCERV_CONFIG_KEYS = [
    "rx-data-per-sec",
    "tx-data-per-sec",
    "max-same-clients",
    "ipv4-network",
    "dns1",
    "dns2",
    "no-udp",
    "keepalive",
    "dpd",
    "mobile-dpd",
    "tunnel-all-dns",
    "restrict-user-to-routes",
    "stats-report-time",
    "mtu",
    "idle-timeout",
    "mobile-idle-timeout",
    "session-timeout",
    "no_routes",
    "routes",
]

DOCKERIZED = os.environ.get("DOCKERIZED", config("DOCKERIZED", "False")).title() == "True"

OCSERV_LOG_FILE = os.environ.get("OCSERV_LOG_FILE", config("OCSERV_LOG_FILE", None))
