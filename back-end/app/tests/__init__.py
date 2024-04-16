import os
from unittest import TestCase
from unittest.mock import patch

from django.conf import settings
from django.core.management import call_command
from rest_framework.test import APIRequestFactory

from app.api.admin import AdminViewSet
from app.models import OcservGroup, AdminPanelConfiguration, OcservUser

test_db_path = settings.BASE_DIR / "db/db_test.sqlite3"
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
admin_configs = None
admin_username = "main_test_admin"
admin_password = "main_test_admin_passwd"


@patch("ocserv.modules.handlers.OcservGroupHandler.update_defaults")
@patch("ocserv.modules.handlers.OcservUserHandler.add_or_update")
def init_db(*args, **kwargs):
    global admin_configs

    if (admin_configs := AdminPanelConfiguration.objects.last()) is None:
        admin_configs = AdminPanelConfiguration.objects.create(
            captcha_site_key=os.environ.get("CAPTCHA_SITE_KEY"),
            captcha_secret_key=os.environ.get("CAPTCHA_SECRET_KEY"),
            default_traffic=10,
            default_configs=default_configs,
        )
    OcservUser.objects.create(
        group=OcservGroup.objects.get(name="defaults"),
        username="init_user",
        password="1234",
        active=True,
        rx=0,
        tx=0,
    )


init_db()


class OcservTestAbstract(TestCase):
    def __init__(self, *args, **kwargs):
        super().__init__(*args, **kwargs)
        self.factory = APIRequestFactory()
        self.token = None

    def login(self, username=admin_username, password=admin_password, status=200) -> str | None:
        request = self.factory.post(
            "/admin/login/",
            data={
                "username": username,
                "password": password,
            },
        )
        response = AdminViewSet.as_view({"post": "login"})(request)
        self.check_status_and_errors(response, status=status)
        if status == 200:
            self.assertEqual(response.data["user"]["username"], username)
            self.assertIn("token", response.data)
            return response.data.get("token")
        return None

    @property
    def get_header(self) -> dict:
        token = self.token
        if not token:
            token = self.login()
        return {"Authorization": f"Token {token}"}

    def check_status_and_errors(self, response, status, error_msg=None):
        self.assertEqual(response.status_code, int(status))
        if response.status_code in [400, 403, 404] and error_msg:
            self.assertEqual(response.data["error"][0], error_msg)
