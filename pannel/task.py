from __future__ import absolute_import, unicode_literals

from django.contrib import messages

from celery import shared_task
import subprocess, re, os


# TODO : change mysql test to ocserv and add sudo to list command
service = "mysql"

@shared_task
def status_ocserv(restart=False):
    
    p =  subprocess.Popen(["systemctl", "status",  service, "--output=json-pretty"], stdout=subprocess.PIPE)
    (output, err) = p.communicate()
    output = output.decode('utf-8')

    status_regx= r"Active:(.*) since (.*);(.*)"
    service_status = {}

    for line in output.splitlines():
        status_search = re.search(status_regx, line)
        if status_search:
            service_status['status'] = status_search.group(1).strip()
            service_status['since'] = status_search.group(2).strip()
            service_status['uptime'] = status_search.group(3).strip()

    service_status['service'] = service
    service_status['restart'] = restart
    return service_status


@shared_task
def restart_ocserv():
    p =  subprocess.Popen(["systemctl", "restart",  service, "--output=json-pretty"], stdout=subprocess.PIPE)
    (output, err) = p.communicate()
    output = output.decode('utf-8')
    if err:
        status_ocserv(restart=False)
    return status_ocserv(restart=True)



