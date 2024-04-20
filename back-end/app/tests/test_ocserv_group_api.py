from unittest.mock import patch

from django.contrib.auth.hashers import make_password
from django.contrib.auth.models import User

from app.api.ocserv_groups import OcservGroupsViewSet
from app.models import OcservGroup
from app.tests import OcservTestAbstract


class OcservGroupApiTest(OcservTestAbstract):

    def __init__(self, *args, **kwargs):
        super().__init__(*args, **kwargs)
        self.staff_token = None
        self.group_list_id = []
        self.staff_username = "test_staff"
        self.staff_password = "test_staff_password"

    @patch("ocserv.modules.handlers.OcservGroupHandler.add_or_update")
    def setUp(self, *args, **kwargs) -> None:
        OcservGroup.objects.get_or_create(name="test_group", desc="test group")

        obj1, _ = OcservGroup.objects.get_or_create(
            name="test_group_delete1", desc="test group"
        )
        self.group_list_id.append(obj1.id)

        obj2, _ = OcservGroup.objects.get_or_create(
            name="test_group_delete2", desc="test group"
        )
        self.group_list_id.append(obj2.id)

        obj3, _ = OcservGroup.objects.get_or_create(
            name="test_group_delete3", desc="test group"
        )
        self.group_list_id.append(obj3.id)

        obj4, _ = OcservGroup.objects.get_or_create(
            name="test_group_delete4", desc="test group"
        )
        self.group_list_id.append(obj4.id)

        User.objects.get_or_create(
            username=self.staff_username,
            defaults={
                "password": make_password(self.staff_password),
                "is_staff": False,
                "is_superuser": False,
            },
        )

    def get_header(self, staff=False) -> dict:
        if staff and not self.staff_token:
            self.staff_token = self.login(
                self.staff_username, password=self.staff_password
            )
            return {"Authorization": f"Token {self.staff_token}"}
        return super().get_header

    def group_list(self, staff=False):
        request = self.factory.get("/groups/", headers=self.get_header(staff=staff))
        response = OcservGroupsViewSet.as_view({"get": "list"})(request)
        self.check_status_and_errors(response, 200)
        self.assertTrue(
            len(
                list(
                    filter(
                        lambda x: x.get("name") == "test_group",
                        response.data.get("result"),
                    )
                )
            )
            > 0,
            "len of groups must be greater than zero",
        )
        self.assertIn(response.data.get("total_count"), range(1, 10))
        self.assertEqual(response.data.get("page"), 1)
        self.assertEqual(response.data.get("pages"), 1)

    def test_group_list(self):
        self.group_list()
        self.group_list(staff=True)

    def group_detail(self, staff=False):
        request = self.factory.get("/groups/1/", headers=self.get_header(staff=staff))
        response = OcservGroupsViewSet.as_view({"get": "retrieve"})(request, pk=1)
        self.check_status_and_errors(response, 200)

    def test_group_detail(self):
        self.group_detail()
        self.group_detail(staff=True)

    def group_create(self, staff=False, name="test group2"):
        data = {"name": name, "decs": "test group desc"}
        request = self.factory.post(
            "/groups/", headers=self.get_header(staff=staff), data=data
        )
        response = OcservGroupsViewSet.as_view({"post": "create"})(request)
        self.check_status_and_errors(response, 201)
        self.assertIsInstance(response.data, dict)
        self.assertEquals(response.data["name"], name.replace(" ", "_"))

    @patch("ocserv.modules.handlers.OcservGroupHandler.add_or_update")
    def test_group_create(self, *args, **kwargs):
        self.group_create()
        self.group_create(staff=True, name="test group3")

    def group_update(self, group_id=None, staff=False, new_name=None):
        _group_id = group_id if group_id else 100  # override id with non exists id
        data = {
            "desc": "test group desc",
        }
        if new_name:
            data["name"] = new_name
        request = self.factory.patch(
            f"/groups/{_group_id}/", headers=self.get_header(staff=staff), data=data
        )
        response = OcservGroupsViewSet.as_view({"patch": "partial_update"})(
            request, pk=_group_id
        )
        if group_id:
            if group_id == 1:
                self.check_status_and_errors(
                    response, 400, "Invalid name (defaults) for group"
                )
            else:
                self.check_status_and_errors(response, 202)
                self.assertEqual(response.data.get("desc"), data["desc"])
                if new_name:
                    self.assertEqual(response.data.get("name"), new_name)
        else:
            self.check_status_and_errors(response, 404, "Ocserv group does not exists")

    @patch("ocserv.modules.handlers.OcservGroupHandler.add_or_update")
    @patch("ocserv.modules.handlers.OcservGroupHandler.delete")
    def test_group_update(self, *args, **kwargs):
        self.group_update(group_id=self.group_list_id[0], new_name="changed_name_admin")
        self.group_update(
            group_id=self.group_list_id[1], staff=True, new_name="changed_name_staff"
        )

        self.group_update()  # 404 test
        self.group_update(staff=True)  # 404 test

        self.group_update(group_id=1)  # 400 error
        self.group_update(group_id=1, staff=True)  # 400 error

        with self.assertRaises(IndexError):
            self.group_update(
                group_id=self.group_list_id[100], new_name="changed_name_admin"
            )
            self.group_update(
                group_id=self.group_list_id[100],
                staff=True,
                new_name="changed_name_staff",
            )

    def group_delete(self, group_id=None, staff=False):
        _group_id = group_id if group_id else 100  # override id with non exists id
        request = self.factory.delete(
            f"/groups/{group_id}/", headers=self.get_header(staff=staff)
        )
        response = OcservGroupsViewSet.as_view({"delete": "destroy"})(
            request, pk=group_id
        )
        if group_id:
            if group_id == 1:
                self.check_status_and_errors(
                    response, 400, "You can not delete defaults Ocserv group"
                )
            else:
                self.check_status_and_errors(response, 204)
        else:
            self.check_status_and_errors(response, 404, "Ocserv group does not exists")

    @patch("ocserv.modules.handlers.OcservGroupHandler.delete")
    def test_group_delete(self, *args, **kwargs):
        self.group_delete(group_id=self.group_list_id.pop())
        self.group_delete(group_id=self.group_list_id.pop(), staff=True)

        self.group_delete()  # 404 test
        self.group_delete(staff=True)  # 404 test

        self.group_delete(group_id=1)  # 400 error
        self.group_delete(group_id=1, staff=True)  # 400 error

        with self.assertRaises(IndexError):
            self.group_delete(group_id=self.group_list_id[100])
            self.group_delete(group_id=self.group_list_id[100], staff=True)
