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

    @patch("ocserv.modules.handlers.OcservGroupHandler.add_or_update")
    def setUp(self, *args, **kwargs) -> None:
        self.staff_username = "test_staff"
        self.staff_password = "test_staff_password"
        OcservGroup.objects.get_or_create(name="test_group", desc="test group")
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
            self.staff_token = self.login(self.staff_username, password=self.staff_password)
            return {"Authorization": f"Token {self.staff_token}"}
        return super().get_header

    def group_list(self, staff=False):
        request = self.factory.get("/groups/", headers=self.get_header(staff=staff))
        response = OcservGroupsViewSet.as_view({"get": "list"})(request)
        self.check_status_and_errors(response, 200)
        self.assertTrue(
            len(list(filter(lambda x: x.get("name") == "test_group", response.data.get("result"))))
            > 0,
            "len of groups must be greater than zero",
        )
        self.assertIn(response.data.get("total_count"), [1, 3])
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
        request = self.factory.post("/groups/", headers=self.get_header(staff=staff), data=data)
        response = OcservGroupsViewSet.as_view({"post": "create"})(request)

    @patch("ocserv.modules.handlers.OcservGroupHandler.add_or_update")
    def test_group_create(self, *args, **kwargs):
        self.group_create()
        self.group_create(staff=True, name="test group3")

    def test_group_update(self):
        pass

    def test_group_delete(self):
        pass
