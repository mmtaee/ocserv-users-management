from django.core.management.base import BaseCommand
from django.utils import timezone

from .handlers import user_handler
from .models import OcservUser
from ocserv.logger import Logger


class Command(BaseCommand):
    help = "Deactivate Expired Accounts or activated monthly account"

    def handle(self, *args, **kwargs):
        logger = Logger()
        users = OcservUser.objects.all()
        # deactivate expired accounts:
        expired_users = users.filter(active=True, expire_date__lte=timezone.now())
        for user in expired_users:
            result = user_handler(action="deactivate", username=user.username)
            if result:
                user.active = False
                user.save()
            else:
                logger.log(level="critical", message=f"deactivate for user with that username ({user.username}) failed")
        # activate monthly accounts
        deactivate_accounts = users.filter(active=False, traffic=OcservUser.MONTHLY)
        for user in deactivate_accounts:
            result = user_handler(action="active", username=user.username)
            if result:
                user.active = True
                user.save()
            else:
                logger.log(level="critical", message=f"activate for user with that username ({user.username}) failed")
