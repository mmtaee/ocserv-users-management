import logging
import os

from django.conf import settings
from django.utils import timezone


log_level_to_nember = {
    "CRITICAL": 50,
    "FATAL": 50,
    "ERROR": 40,
    "WARNING": 30,
    "WARN": 30,
    "INFO": 20,
    "DEBUG": 10,
    "NOTSET": 0,
}


class Logger:
    LOG_PATH = settings.LOG_PATH
    __INSTANCE = None

    def __init__(self, stdout=False):
        self.stdout = stdout

    def __new__(cls, *args, **kwargs):
        if not hasattr(cls, "__instance"):
            cls.__INSTANCE = super().__new__(cls)
        return cls.__INSTANCE

    def to_stdout(self, msg):
        if self.stdout:
            logging.log(level=50, msg=msg)
        else:
            print(msg)

    def open_file(self, mode="a"):
        try:
            with open(self.LOG_PATH, mode) as f:
                return f
        except PermissionError as e:
            self.to_stdout(str(e))
            raise e

    def log(self, level, message):
        message = f"[{level.title()}] {timezone.now()} - {message.lower()}\n"
        f = self.open_file()
        f.write(message)
        f.close()
        if self.stdout:
            logging.log(level=log_level_to_nember.get(level.lower(), 0), msg=message)

    def clear(self):
        f = self.open_file("w")
        f.write("## Logs cleared by admin")
        f.close()

    def read(self):
        if not os.path.exists(self.LOG_PATH):
            os.mknod(self.LOG_PATH)
        f = self.open_file("r")
        lines = f.readlines()
        f.close()
        return lines
