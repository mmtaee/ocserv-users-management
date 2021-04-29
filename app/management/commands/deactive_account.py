from django.contrib.auth.models import User
from django.core.management.base import BaseCommand
from django.utils.crypto import get_random_string

import os

from app.models import *

class Command(BaseCommand):
    help = 'Deactivate Expired Accounts'

    def handle(self, *args, **kwargs):
        oscerv_users = OcservUser.objects.filter(oc_active=True)
        for user in oscerv_users:
            self.stdout.write(self.style.WARNING(f'progressing ocserv user {user.oc_username}'))
            if user.expire:
                try :
                    user.oc_active = False
                    command = f'sudo /usr/bin/ocpasswd  -c /etc/ocserv/ocpasswd -l {user.oc_username}'
                    os.system(command)
                    user.save()
                    self.stdout.write(self.style.SUCCESS(f'ocserv user {user.oc_username} deactivate with success'))
                except:
                    self.stdout.write(self.style.ERROR(f'ocserv user {user.oc_username} not deactivate'))
    

