import os
from unittest import TestCase
from unittest.mock import patch

from django.contrib.auth.models import User
from django.core.management import call_command

from app.models import OcservGroup, AdminPanelConfiguration
from ocserv.settings import BASE_DIR

test_db_path = BASE_DIR / "db/db_test.sqlite3"
if os.path.exists(test_db_path):
    os.remove(test_db_path)

call_command("migrate")

default_configs = {
    "routes": ["192.168.1.6", "192.168.2.6"],
    "dns1": ["8.8.8.8"],
    "ipv4-network": "172.16.12.1/22",
    "rx-data-per-sec": 500000,
}

update_default_configs = default_configs
update_default_configs.update({"mtu": 1500})


class OcservTestAbstract(TestCase):

    def __init__(self, *args, **kwargs):
        super().__init__(*args, **kwargs)
        self.token = None

    @patch("ocserv.modules.handlers.OcservGroupHandler.update_defaults")
    def setUp(self, *args, **kwargs) -> None:
        if not OcservGroup.objects.filter(name="defaults").exists():
            self.group = OcservGroup.objects.create(name="defaults", desc="defaults group")
        else:
            self.group = OcservGroup.objects.get(name="defaults")
        if (admin := AdminPanelConfiguration.objects.last()) is None:
            admin = AdminPanelConfiguration.objects.create(
                captcha_site_key=os.environ.get("CAPTCHA_SITE_KEY"),
                captcha_secret_key=os.environ.get("CAPTCHA_SECRET_KEY"),
                default_traffic=10,
                default_configs=default_configs,
            )
            User.objects.create_superuser(
                username="test_admin",
                password="test_admin_passwd",
                is_superuser=True,
            )
        self.admin = admin
        if (staff := User.objects.filter(username="setup_test_staff").last()) is None:
            staff = User(
                username="setup_test_staff",
                password="setup_test_staff_passwd",
                is_superuser=False,
                is_staff=False,
            )
            User.objects.bulk_create(
                [
                    staff,
                    User(
                        username="setup_test_staff2",
                        password="setup_test_staff_passwd",
                        is_superuser=False,
                        is_staff=False,
                    ),
                    User(
                        username="setup_test_staff3",
                        password="setup_test_staff_passwd",
                        is_superuser=False,
                        is_staff=False,
                    ),
                ]
            )
        self.staff = staff

    def check_status_and_errors(self, response, status, error_msg=None):
        self.assertEqual(response.status_code, int(status))
        if response.status_code in [400, 403, 404] and error_msg:
            self.assertEqual(response.data["error"][0], error_msg)

    @property
    def auth_header(self):
        return {"Authorization": f"Token {self.token}"}
