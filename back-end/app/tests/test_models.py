import os
import random
from datetime import datetime
from unittest import TestCase
from unittest.mock import patch
from django.core.management import call_command
from rest_framework.exceptions import ValidationError
from app.models import (
    AdminPanelConfiguration,
    OcservGroup,
    OcservUser,
    MonthlyTrafficStat,
)
from decimal import Decimal

from ocserv.settings import BASE_DIR

default_configs = {
    "routes": ["192.168.1.6", "192.168.2.6"],
    "dns1": ["8.8.8.8"],
    "ipv4-network": "172.16.12.1/22",
    "rx-data-per-sec": 500000,
}

update_default_configs = default_configs
update_default_configs.update({"mtu": 1500})


class ModelsTestCase(TestCase):
    def __init__(self, *args, **kwargs):
        super().__init__(*args, **kwargs)
        self.counter = int(os.environ.get("COUNTER", default=100))

    @classmethod
    def setUpClass(cls):
        test_db_path = BASE_DIR / "db/db_test.sqlite3"
        if os.path.exists(test_db_path):
            os.remove(test_db_path)
        call_command("migrate")
        super().setUpClass()

    def create_bulk_ocserv_groups(self):
        bulk_obj = []
        for i in range(self.counter):
            bulk_obj.append(
                OcservGroup(
                    name=f"test{i}",
                    desc=f"test{i} group",
                    configs=update_default_configs,
                )
            )
        OcservGroup.objects.bulk_create(bulk_obj)

    def create_bulk_ocserv_users(self, prefix=None):
        bulk_users = []
        groups = OcservGroup.objects.all()
        for i in range(self.counter):
            random_index = random.randint(0, groups.count() - 1)
            if prefix:
                username = f"{prefix}-test-{i+1}"
            else:
                username = f"test{i + 1}"
            bulk_users.append(
                OcservUser(
                    group=groups[random_index],
                    username=username,
                    password="1234",
                    active=True,
                    rx=Decimal(str(random.randint(0, 100))),
                    tx=Decimal(str(random.randint(0, 100))),
                )
            )
        users = OcservUser.objects.bulk_create(bulk_users)
        return users

    @patch("ocserv.modules.handlers.OcservGroupHandler.update_defaults")
    def setUp(self, mock_data) -> None:
        if not OcservGroup.objects.filter(name="defaults").exists():
            self.group = OcservGroup.objects.create(
                name="defaults", desc="defaults group"
            )
        else:
            self.group = OcservGroup.objects.get(name="defaults")
        if (admin := AdminPanelConfiguration.objects.last()) is None:
            admin = AdminPanelConfiguration.objects.create(
                captcha_site_key=os.environ.get("CAPTCHA_SITE_KEY"),
                captcha_secret_key=os.environ.get("CAPTCHA_SECRET_KEY"),
                default_traffic=10,
                default_configs=default_configs,
            )
        self.admin = admin

    @patch("ocserv.modules.handlers.OcservGroupHandler.update_defaults")
    def test_admin_panel_configuration_create_extra(self, *args):
        with self.assertRaises(ValidationError) as context:
            AdminPanelConfiguration.objects.create(
                captcha_site_key=os.environ.get("CAPTCHA_SITE_KEY"),
                captcha_secret_key=os.environ.get("CAPTCHA_SECRET_KEY"),
                default_traffic=10,
                default_configs=default_configs,
            )
        self.assertEqual(
            str(context.exception.detail["error"][0]),
            "Amin Configuration already exists",
        )

    @patch("ocserv.modules.handlers.OcservGroupHandler.update_defaults")
    def test_admin_panel_configuration_update(self, *args):
        self.admin.default_configs = update_default_configs
        self.admin.save()

    def test_admin_panel_configuration_get(self):
        instance = AdminPanelConfiguration.objects.last()
        self.assertEquals(instance.default_traffic, 10)
        self.assertEqual(instance.default_configs, update_default_configs)

    def test_ocserv_create_default_group(self):
        with self.assertRaises(ValidationError) as context:
            OcservGroup.objects.create(name="defaults", desc="defaults group")
        self.assertEqual(
            str(context.exception.detail["error"][0]),
            "Invalid name (defaults) for group",
        )

    def test_ocserv_group_create(self):
        self.create_bulk_ocserv_groups()
        self.assertEqual(OcservGroup.objects.count(), self.counter + 1)

    def test_ocserv_group_get(self):
        queryset = OcservGroup.objects.filter(name__startswith="test")
        self.assertEqual(queryset.count(), self.counter)

    def test_ocserv_group_edit(self):
        queryset = OcservGroup.objects.filter(name__startswith="test")
        for num, obj in enumerate(queryset):
            obj.name = f"test-x{num+1000}"
            obj.desc = f"test-x{num+1000} group"
        OcservGroup.objects.bulk_update(queryset, fields=["name", "desc"])
        self.assertEqual(
            OcservGroup.objects.filter(desc__startswith="test-x").count(), self.counter
        )

    @patch("ocserv.modules.handlers.OcservGroupHandler.destroy")
    def test_ocserv_group_remove(self, *args):
        OcservGroup.objects.filter(name__startswith="test").delete()
        self.assertEqual(OcservGroup.objects.count(), 1)

    @patch("ocserv.modules.handlers.OcservUserHandler.add_or_update")
    def test_ocserv_user_create(self, *args):
        self.create_bulk_ocserv_users()
        self.assertEqual(
            OcservUser.objects.filter(username__startswith="test").count(), self.counter
        )

    @patch("ocserv.modules.handlers.OcservUserHandler.add_or_update")
    def test_ocserv_user_edit(self, *args):
        queryset = OcservUser.objects.all()
        for num, obj in enumerate(queryset):
            n = f"update-{num + 1000}"
            obj.username = f"test-{n}"
            obj.password = f"test-{n}-1234"
        OcservUser.objects.bulk_update(queryset, fields=["username", "password"])
        self.assertEqual(
            OcservUser.objects.get(username="test-update-1001").password,
            f"test-update-1001-1234",
        )

    @patch("ocserv.modules.handlers.OcservUserHandler.delete")
    def test_ocserv_user_remove(self, *args):
        OcservUser.objects.filter(username__startswith="test").delete()
        self.assertEqual(OcservUser.objects.count(), 0)

    def test_monthly_traffic_stat(self):
        bulk_stats = []
        today = datetime.now().today()
        if OcservUser.objects.count() == 0:
            self.create_bulk_ocserv_users(prefix="monthly")
        for user in (queryset := OcservUser.objects.all()):
            bulk_stats.append(
                MonthlyTrafficStat(
                    user=user,
                    year=today.year,
                    month=today.month,
                    rx=user.rx,
                    tx=user.tx,
                )
            )
        MonthlyTrafficStat.objects.bulk_create(bulk_stats)
        self.assertEqual(MonthlyTrafficStat.objects.count(), len(queryset))


# DEBUG=False COUNTER=1000 ./manage.py test --settings=ocserv.settings_test
