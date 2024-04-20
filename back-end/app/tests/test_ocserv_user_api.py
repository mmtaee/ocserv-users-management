import json
from datetime import datetime
from unittest.mock import patch

from app.api.ocserv_users import OcservUsersViewSet
from app.tests import OcservTestAbstract, staff_username, staff_password


class OcservGroupApiTest(OcservTestAbstract):

    def __init__(self, *args, **kwargs):
        super().__init__(*args, **kwargs)
        self.staff_token = None

    def get_header(self, staff=False) -> dict:
        if staff and not self.staff_token:
            self.staff_token = self.login(staff_username, password=staff_password)
            return {"Authorization": f"Token {self.staff_token}"}
        return super().get_header

    def user_list(self, staff=False):
        request = self.factory.get("/users/", headers=self.get_header(staff=staff))
        response = OcservUsersViewSet.as_view({"get": "list"})(request)
        self.check_status_and_errors(response, 200)
        self.assertTrue(
            len(
                list(
                    filter(
                        lambda x: x.get("username") == "init_user",
                        response.data.get("result"),
                    )
                )
            )
            > 0,
            "len of users must be greater than zero",
        )
        self.assertIn(response.data.get("total_count"), range(1, 10))
        self.assertEqual(response.data.get("page"), 1)
        self.assertEqual(response.data.get("pages"), 1)

    @patch("ocserv.modules.handlers.OcservUserHandler.online")
    def test_user_list(self, *args, **kwargs):
        self.user_list()
        self.user_list(staff=True)

    def user_detail(self, staff=False):
        request = self.factory.get("/users/1/", headers=self.get_header(staff=staff))
        response = OcservUsersViewSet.as_view({"get": "retrieve"})(request, pk=1)
        self.check_status_and_errors(response, 200)

    def test_user_detail(self):
        self.user_detail()
        self.user_detail(staff=True)

    def user_create(self, group_id=1, username="test1", staff=False):
        data = {"group": group_id, "username": username, "password": "test1234"}
        request = self.factory.post(
            "/users/", headers=self.get_header(staff=staff), data=data
        )
        response = OcservUsersViewSet.as_view({"post": "create"})(request)
        if group_id != 1:
            self.check_status_and_errors(response, 404, "Ocserv group does not exists")
            return
        if response.status_code != 201:
            self.check_status_and_errors(response, 400, "Ocserv User exists")
            return
        self.check_status_and_errors(response, 201)
        self.assertIsInstance(response.data, dict)
        self.assertEqual(response.data.get("group"), 1)
        return response.data

    @patch("ocserv.modules.handlers.OcservUserHandler.add_or_update")
    def test_user_create(self, *args, **kwargs):
        self.user_create()
        self.user_create(staff=True, username="test2")

        self.user_create()  # 400 error
        self.user_create(staff=True, username="test2")  # 400 error

        self.user_create(group_id=0)  # 400 invalid group id
        self.user_create(group_id=0, staff=True)  # 400 invalid group id

    def user_update(self, group_id=None, user_id=None, staff=False, **kwargs):
        if group_id is None:
            group_id = 1
        if user_id is None:
            user_id = 1
        username = kwargs.pop("username", None)
        password = kwargs.pop("password", "test12345")
        data = {
            "group": group_id,
            "password": password,
            "expire_date": str(datetime.now()),
        }
        if username:
            data["username"] = username
        request = self.factory.patch(
            f"/users/{user_id}/",
            content_type="application/json",
            headers=self.get_header(staff=staff),
            data=json.dumps(data),
        )
        response = OcservUsersViewSet.as_view({"patch": "partial_update"})(
            request, pk=user_id
        )
        if user_id != 1:
            self.check_status_and_errors(response, 404, "Ocserv user does not exist")
            return
        if group_id and response.status_code == 400:
            self.assertEqual(
                response.data["group"][0],
                f'Invalid pk "{group_id}" - object does not exist.',
            )
            return
        if username:
            # send new username. username should not be allowed to change
            self.assertNotEqual(response.data["username"], username)
        self.check_status_and_errors(response, 202)
        self.assertFalse(response.data["active"], False)

    @patch("ocserv.modules.handlers.OcservUserHandler.add_or_update")
    @patch("ocserv.modules.handlers.OcservUserHandler.status_handler")
    def test_user_update(self, *args, **kwargs):
        self.user_update(user_id=1)
        self.user_update(user_id=10)
        self.user_update(user_id=1, group_id=10)  # 400 invalid group id
        self.user_update(
            username="masoud"
        )  # send new username. username should not be allowed to change

    def user_delete(self, user_id=None, *args, **kwargs):
        request = self.factory.delete(
            f"/users/{user_id}/",
            headers=self.get_header(),
        )
        response = OcservUsersViewSet.as_view({"delete": "destroy"})(
            request, pk=user_id
        )
        if user_id is None or response.status_code == 404:
            self.check_status_and_errors(response, 404, "Ocserv user does not exist")
            return
        self.check_status_and_errors(response, 204)

    @patch("ocserv.modules.handlers.OcservUserHandler.delete")
    @patch("ocserv.modules.handlers.OcservUserHandler.add_or_update")
    def test_user_delete(self, *args, **kwargs):
        result = self.user_create()
        user_id = result["id"]
        self.user_delete(user_id=user_id)
        self.user_delete(user_id=50)  # 400 error: invalid user id
        self.user_delete()  # 400 error: None user id

    @patch("ocserv.modules.handlers.OcservUserHandler.add_or_update")
    @patch("ocserv.modules.handlers.OcservUserHandler.disconnect")
    def test_user_disconnect(self, mock_disconnect, mock_add):
        result = self.user_create()
        user_id = result["id"]
        mock_disconnect.return_value = True
        request = self.factory.post(
            f"/users/{user_id}/disconnect/",
            headers=self.get_header(),
        )
        response = OcservUsersViewSet.as_view({"post": "disconnect"})(
            request, pk=user_id
        )
        self.check_status_and_errors(response, 202)

        mock_disconnect.return_value = False
        request = self.factory.post(
            f"/users/{user_id}/disconnect/",
            headers=self.get_header(),
        )
        response = OcservUsersViewSet.as_view({"post": "disconnect"})(
            request, pk=user_id
        )
        self.check_status_and_errors(response, 400, "Ocserv User Disconnect Failed")

    @patch("ocserv.modules.handlers.OcservUserHandler.add_or_update")
    @patch("ocserv.modules.handlers.OcservUserHandler.status_handler")
    def test_user_status_handler(self, *args, **kwargs):
        result = self.user_create(username="test_user_status")
        user_id = result["id"]
        data = {"status": 0}

        # send None status get error 400
        request = self.factory.post(
            f"/users/{user_id}/status/",
            headers=self.get_header(),
        )
        response = OcservUsersViewSet.as_view({"post": "user_status_handler"})(
            request, pk=user_id
        )
        self.check_status_and_errors(response, 400, "status of user not found in data")

        # send invalid user 404

        request = self.factory.post(
            "/users/100/status/", headers=self.get_header(), data=data
        )
        response = OcservUsersViewSet.as_view({"post": "user_status_handler"})(
            request, pk=100
        )
        self.check_status_and_errors(response, 404, "ocserv user does not exist")

        # success status
        request = self.factory.post(
            f"/users/{user_id}/status/", headers=self.get_header(), data=data
        )
        response = OcservUsersViewSet.as_view({"post": "user_status_handler"})(
            request, pk=user_id
        )
        self.check_status_and_errors(response, 202)

    def test_user_sync(self):
        pass
