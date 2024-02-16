import os.path

from django.core.management.base import BaseCommand, CommandError
from django.db.utils import IntegrityError
import sqlite3

from app.models import OcservUser, OcservGroup, AdminConfig


class Command(BaseCommand):
    help = "Migrate data from old database to new database"

    def __init__(self):
        super().__init__()
        self.traffic = None
        self.group = None
        self.path = None

    def add_arguments(self, parser):
        parser.add_argument(
            "--old-path",
            type=str,
            required=True,
            help="Path to the old SQLite database",
        )
        parser.add_argument(
            "--free-traffic",
            action="store_true",
            default=False,
            help="migrate users with free usage traffic",
        )

    def handle(self, *args, **options):
        self.path = options["old_path"]
        if not os.path.isfile(options["old_path"]):
            raise FileNotFoundError(self.path)
        try:
            self.group = OcservGroup.objects.get(name="defaults")
        except OcservGroup.DoesNotExist:
            raise CommandError(
                "Default admin configs not found. Try to create it from panel first!"
            )
        if not AdminConfig.objects.exists():
            raise CommandError(
                "Default admin configs not found. Try to create it from panel first!"
            )
        self.traffic = (
            OcservUser.FREE if options["free_traffic"] else OcservUser.MONTHLY
        )
        user_exists = list(OcservUser.objects.values_list("username", flat=True))
        conn = sqlite3.connect(self.path)
        if conn is None:
            raise ConnectionError("Database Connection Error")
        cursor = conn.cursor()
        try:
            cursor.execute(
                "SELECT username,password,active,expire_date,desc FROM app_ocservuser;"
            )
            while True:
                batch = cursor.fetchmany(100)
                if not batch:
                    break
                for user in batch:
                    username = user[0]
                    if username not in user_exists:
                        try:
                            OcservUser.objects.create(
                                group=self.group,
                                username=username,
                                password=user[1],
                                active=user[2],
                                expire_date=user[3],
                                desc=user[4],
                                traffic=self.traffic,
                            )
                            self.stdout.write(
                                self.style.SUCCESS(
                                    f"User with username ({user[0]}) added."
                                )
                            )
                        except IntegrityError:
                            self.stdout.write(
                                self.style.ERROR(
                                    f"User with username ({user[0]}) already exists."
                                )
                            )
                            continue
                    else:
                        self.stdout.write(
                            self.style.ERROR(
                                f"User with username ({user[0]}) already exists."
                            )
                        )
        except sqlite3.OperationalError as e:
            raise CommandError(f"Error executing SQL query: {str(e)}")
        finally:
            conn.close()


#   /var/www/site/back-end/venv/bin/python3 manage.py migrate_to_new --old-path /OLD_PATH/db.sqlite3 --free-traffic
#   python3 /app/manage.py migrate_to_new --old-path /OLD_PATH/db.sqlite3
