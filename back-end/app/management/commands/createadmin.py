from django.contrib.auth.models import User
from django.core.management import BaseCommand


class Command(BaseCommand):
    help = "Create Staff User for panel"

    def add_arguments(self, parser):
        parser.add_argument(
            "-u",
            "--username",
            type=str,
            required=True,
            help="Amin Username",
        )
        parser.add_argument(
            "-p",
            "--password",
            type=str,
            required=True,
            help="Amin Password",
        )

    def handle(self, *args, **options):
        username = options["username"]
        if not User.objects.filter(username=username).exists():
            user = User.objects.create_user(
                username=username,
                password=options["password"],
                is_staff=False,
                is_superuser=False,
            )
            self.stdout.write(f"User with username ({user.username}) created.")
        else:
            self.stdout.write(
                self.style.ERROR(f"User with username ({username}) already exists.")
            )
