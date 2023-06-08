import os

from django.conf import settings
from django.utils import timezone


class Logger:
    LOG_PATH = settings.LOG_PATH
    __INSTANCE = None

    def __new__(cls, *args, **kwargs):
        if not hasattr(cls, "__instance"):
            cls.__INSTANCE = super().__new__(cls)
        return cls.__INSTANCE

    def log(self, level, message):
        message = f"[{level.title()}] {timezone.now()} - {message.lower()}\n"
        with open(self.LOG_PATH, "a") as f:
            f.write(message)
            f.close()

    def clear(self):
        with open(self.LOG_PATH, "w") as f:
            f.write("Log clear by admin")
            f.close()

    def read(self):
        if not os.path.exists(self.LOG_PATH):
            os.mknod(self.LOG_PATH)
        with open(self.LOG_PATH, "r") as f:
            lines = f.readlines()
            f.close()
        return lines