from django.contrib.auth.hashers import make_password
from django.contrib.auth.models import User

from app.api.admin import AdminViewSet
from app.tests import OcservTestAbstract


class StaffApiTest(OcservTestAbstract):

    def __init__(self, *args, **kwargs):
        super().__init__(*args, **kwargs)
        self.token = None

    def setUp(self) -> None:
        staff_username = "test_staff"
        staff_password = "test_staff_password"
        User.objects.get_or_create(
            username=staff_username,
            defaults={
                "password": make_password(staff_password),
                "is_staff": False,
                "is_superuser": False,
            },
        )
        self.token = self.login(staff_username, staff_password)

    def test_staff_list(self):
        request = self.factory.get("/admin/staffs/", headers=self.get_header)
        response = AdminViewSet.as_view({"get": "staffs"})(request)
        self.check_status_and_errors(response, 403, "you have not access to this route")

    def test_staff_create(self):
        data = {"username": "setup_test_staff", "password": "setup_test_staff_passwd"}
        request = self.factory.post("/admin/staffs/", headers=self.get_header, data=data)
        response = AdminViewSet.as_view({"post": "staffs"})(request)
        self.check_status_and_errors(response, 403, "you have not access to this route")

    def test_staff_delete(self):
        request = self.factory.delete("/admin/staffs/3/", headers=self.get_header)
        response = AdminViewSet.as_view({"delete": "delete_staff"})(request, pk=3)
        self.check_status_and_errors(response, 403, "you have not access to this route")
