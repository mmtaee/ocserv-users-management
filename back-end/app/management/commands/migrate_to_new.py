from django.core.management.base import BaseCommand, CommandError
from django.db.utils import IntegrityError
import sqlite3

from app.models import OcservUser, OcservGroup


class Command(BaseCommand):
    help = "Migrate data from old database to new database"

    def __init__(self):
        super().__init__()
        self.path = None

    def fetch_old_users(self):
        try:
            conn = sqlite3.connect(self.path)
            cursor = conn.cursor()
        except sqlite3.Error as e:
            raise CommandError(f"Database error: {str(e)}")
        query = 'SELECT "app_ocservuser"."oc_username","app_ocservuser"."oc_password","app_ocservuser"."oc_active",'
        query += '"app_ocservuser"."expire_date","app_ocservuser"."desc" FROM "app_ocservuser"'
        try:
            cursor.execute(query)
            old_users = cursor.fetchall()
        except sqlite3.OperationalError as e:
            raise CommandError(f"Error executing SQL query: {str(e)}")
        finally:
            conn.close()
        return old_users

    def add_arguments(self, parser):
        parser.add_argument("--old-path", type=str, required=True, help="Path to the old SQLite database")

    def handle(self, *args, **options):
        self.path = options["old_path"]
        old_users = self.fetch_old_users()
        group = OcservGroup.objects.get(name="defaults")
        for user in old_users:
            try:
                OcservUser.objects.create(
                    group=group, username=user[0], password=user[1], active=user[2], expire_date=user[3], desc=user[4]
                )
            except IntegrityError:
                self.stdout.write(self.style.ERROR(f"User with username({user[0]}) already exists"))
                continue


#   /var/www/site/back-end/venv/bin/python3 manage.py migrate_to_new --old-path /OLD_PATH/db.sqlite3
#   python3 /app/manage.py migrate_to_new --old-path /OLD_PATH/db.sqlite3
