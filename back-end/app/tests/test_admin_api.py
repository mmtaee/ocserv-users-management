import json
from unittest import mock
from unittest.mock import patch

from app.api.admin import AdminViewSet
from app.tests import SetUpTestAbstract, default_configs, update_default_configs

from rest_framework.test import APIRequestFactory


class AdminApiTest(SetUpTestAbstract):

    def __init__(self, *args, **kwargs):
        super().__init__(*args, **kwargs)
        self.factory = APIRequestFactory()
        self.token = None
        self.admin_password = "test_admin_passwd"

    def check_status(self, response, status):
        self.assertEqual(response.status_code, int(status))

    @property
    def auth_header(self):
        if not self.token:
            self.login()
        return {"Authorization": f"Token {self.token}"}

    def login(self):
        data = {"username": "test_admin", "password": self.admin_password}
        request = self.factory.post("/admin/login/", data=data)
        response = AdminViewSet.as_view({"post": "login"})(request)
        if response.status_code == 400:
            self.admin_password = "new_test_admin_passwd"
            data = {"username": "test_admin", "password": self.admin_password}
            request = self.factory.post("/admin/login/", data=data)
            response = AdminViewSet.as_view({"post": "login"})(request)
        self.check_status(response, 200)
        self.assertEqual(response.data["user"]["username"], "test_admin")
        self.assertIn("token", response.data)
        self.token = response.data["token"]

    def test_admin_config(self):
        request = self.factory.get("/admin/config/")
        response = AdminViewSet.as_view({"get": "config"})(request)
        self.check_status(response, 200)
        self.assertEqual(response.data["config"], True)

    def test_create_admin_config(self):
        data = {
            "username": "test_admin_fake",
            "password": "test_admin_passwd_fake",
            **default_configs,
        }
        request = self.factory.post("/admin/create/", data=data)
        response = AdminViewSet.as_view({"post": "create_admin_configs"})(request)
        if response.status_code != 400:
            self.check_status(response, 201)
            self.assertEqual(response.data["user"]["username"], "test_admin")
        else:
            self.assertEqual(response.data["error"][0], "Admin config exists!")

    def test_login(self):
        data = {"username": "test_admin_fake", "password": "test_admin_passwd_fake"}
        request = self.factory.post("/admin/login/", data=data)
        response = AdminViewSet.as_view({"post": "login"})(request)
        if response.status_code == 400:
            self.assertEqual(response.data["error"][0], "Invalid username or password")
        else:
            self.login()

    def test_logout(self):
        request = self.factory.delete("/admin/logout/", headers=self.auth_header)
        response = AdminViewSet.as_view({"delete": "logout"})(request)
        self.check_status(response, 204)

    def test_configuration_get(self):
        request = self.factory.get("/admin/configuration/", headers=self.auth_header)
        response = AdminViewSet.as_view({"get": "configuration"})(request)
        self.check_status(response, 200)
        self.assertIn("default_configs", response.data)
        self.assertIn("captcha_secret_key", response.data)
        self.assertEqual(response.data["default_configs"]["ipv4-network"], "172.16.12.1/22")
        self.assertEqual(response.data["default_traffic"], 10)

    @patch("ocserv.modules.handlers.OcservGroupHandler.update_defaults")
    def test_configuration_update(self, *args, **kwargs):
        request = self.factory.patch(
            "/admin/configuration/", headers=self.auth_header, data=update_default_configs
        )
        response = AdminViewSet.as_view({"patch": "configuration"})(request)
        self.check_status(response, 202)
        self.assertEqual(response.data["default_configs"]["mtu"], 1500)

    @patch("ocserv.modules.handlers.OcctlHandler.show")
    @patch("ocserv.modules.handlers.OcservUserHandler.online")
    def test_dashboard(self, mock_data_online, mock_data_show):
        mock_data_online.return_value = []
        mock_data_show.return_value = {
            "show_ip_bans": [],
            "show_status": "Note: the printed statistics are not real-time; session time\n"
            "as well as RX and TX data are updated on user disconnect\n\nGeneral info:"
            "\n\tStatus: online\n\tServer PID: 22\n\tSec-mod PID: 26\n\t"
            "Up since: 2024-04-14 12:02 (   28s)\n\tActive sessions: 0\n\t"
            "Total sessions: 0\n\tTotal authentication failures: 0\n\t"
            "IPs in ban list: 0\n\nCurrent stats period:\n\t"
            "Last stats reset: 2024-04-14 12:02 (   28s)\n\t"
            "Sessions handled: 0\n\tTimed out sessions: 0\n\t"
            "Timed out (idle) sessions: 0\n\tClosed due to error sessions: 0\n\t"
            "Authentication failures: 0\n\tAverage auth time:     0s\n\t"
            "Max auth time:     0s\n\tAverage session time:     0s\n\t"
            "Max session time:     0s\n\tRX: 0 bytes\n\tTX: 0 bytes\n",
            "show_iroutes": {},
        }
        request = self.factory.get("/admin/dashboard/", headers=self.auth_header)
        response = AdminViewSet.as_view({"get": "dashboard"})(request)
        self.check_status(response, 200)
        self.assertIn("Note", response.data["show_status"], "Note is not present in show_status")
        self.assertEqual(response.data["show_ip_bans"], [])
        self.assertEqual(response.data["show_iroutes"], {})

    def test_change_password(self):
        data = {
            "old_password": self.admin_password,
            "password": "new_test_admin_passwd",
        }
        request = self.factory.post("/admin/change_password/", headers=self.auth_header, data=data)
        response = AdminViewSet.as_view({"post": "change_password"})(request)
        self.check_status(response, 202)
        self.admin_password = "new_test_admin_passwd"

    def test_staff_list(self):
        pass

    def test_staff_create(self):
        # TODO: try create with staff user
        # TODO: try create with admin user
        # TODO: try create exists user
        pass

    def test_staff_delete(self):
        # TODO: try delete with staff user
        # TODO: try delete with admin user
        # TODO: try delete non exists user
        pass
