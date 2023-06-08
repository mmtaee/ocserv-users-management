import os
import django
import subprocess
import re
from datetime import datetime


def check_stats():
    logger = Logger()
    cmd = "journalctl -fu ocserv"
    if logfile := os.environ.get("LOG_FILE"):
        cmd = f"tail -f {logfile}"
    process = subprocess.Popen(cmd.split(" "), stdout=subprocess.PIPE)
    last_log_entry = "start script"
    while True:
        line = process.stdout.readline().decode().strip()
        if line.startswith(last_log_entry):
            continue
        last_log_entry = line
        regex_disconnect = "\w{3}\s\d{2}\s\d{2}\:\d{2}\:\d{2}\s\w.+ocserv\[\d+\]\:\socserv\[\d+\]\:\smain\[(\w+)\]"
        regex_disconnect += "\W.+rx\:\W(\d+)\,\Wtx\:\W(\d+)"
        if re.match(regex_disconnect, line):
            line_data = re.findall(regex_disconnect, line)
            if line_data:
                line_data = line_data[0]
                username = line_data[0].strip()
                rx = line_data[1]
                tx = line_data[2]
                try:
                    ocserv_user = OcservUser.objects.get(username=username)
                    ocserv_user.rx += rx
                    ocserv_user.tx += tx
                    stat, _ = MonthlyTrafficStat.objects.get_or_create(
                        user=ocserv_user, year=datetime.now().year, month=datetime.now().month
                    )
                    stat.rx += rx
                    stat.tx += tx
                    stat.save()
                    if (ocserv_user.traffic == OcservUser.MONTHLY and stat.tx >= ocserv_user.default_traffic) or (
                        ocserv_user.traffic == OcservUser.TOTALLY and ocserv_user.rx >= ocserv_user.default_traffic
                    ):
                        result = user_handler(action="deactivate", username=ocserv_user.username)
                        if result and result.get("action"):
                            ocserv_user.active = False
                        else:
                            logger.log(
                                level="critical", message=f"deactivate for user with that username ({username}) failed"
                            )
                    ocserv_user.save()
                except OcservUser.DoesNotExist:
                    logger.log(level="warning", message=f"user with that username ({username}) does not exists in db")
            else:
                logger.log(level="critical", message="unprocessable ocserv log to calculate user-rx and user-tx")
                logger.log(level="info", message=line)


if __name__ == "__main__":
    os.environ.setdefault("DJANGO_SETTINGS_MODULE", "ocserv.settings")
    django.setup()
    check_stats()
    # from .handlers import user_handler
    # from .models import MonthlyTrafficStat
    # from .models import OcservUser
    from ocserv.modules.logger import Logger
