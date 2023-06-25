from django.core.management.base import BaseCommand
from django.utils import timezone

from app.models import OcservUser
from ocserv.modules.handlers import OcservUserHandler
from ocserv.modules.logger import Logger


class Command(BaseCommand):
    help = "Deactivate Expired Accounts or activated monthly account"

    def handle(self, *args, **kwargs):
        logger = Logger()
        users = OcservUser.objects.all()

        # deactivate expired accounts:
        expired_users = users.filter(active=True, expire_date__lte=timezone.now())
        user_handler = OcservUserHandler()
        for user in expired_users:
            user_handler.username = user.username
            result = user_handler.status_handler(active=False)
            if result:
                user.active = False
                user.deactivate_date = None
                user.save()
            else:
                logger.log(level="critical", message=f"deactivate for user with that username ({user.username}) failed")

        # activate monthly accounts
        deactivate_accounts = users.filter(
            active=False,
            traffic=OcservUser.MONTHLY,
            expire_date__gt=timezone.now(),
            deactivate_date__isnull=False,
            deactivate_date__month__lt=timezone.now().month,
        )
        for user in deactivate_accounts:
            user_handler.username = user.username
            result = user_handler.status_handler(active=True)
            if result:
                user.active = True
                user.save()
            else:
                logger.log(level="critical", message=f"activate for user with that username ({user.username}) failed")
