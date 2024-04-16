import os
import random
from datetime import datetime
from decimal import Decimal
from unittest import TestCase
from unittest.mock import patch

from rest_framework.exceptions import ValidationError

from app.models import OcservGroup, OcservUser, AdminPanelConfiguration, MonthlyTrafficStat
from app.tests import default_configs, update_default_configs


class ModelsTestCase(TestCase):
    def __init__(self, *args, **kwargs):
        super().__init__(*args, **kwargs)
        self.counter = int(os.environ.get("COUNTER", default=10))

    def setUp(self) -> None:
        self.admin_configs, _ = AdminPanelConfiguration.objects.get_or_create(
            id=1,
            defaults={
                "captcha_site_key": os.environ.get("CAPTCHA_SITE_KEY"),
                "captcha_secret_key": os.environ.get("CAPTCHA_SECRET_KEY"),
                "default_traffic": 10,
                "default_configs": default_configs,
            },
        )

    def tearDown(self) -> None:
        OcservGroup.objects.exclude(name="defaults").delete()

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

    def create_bulk_ocserv_users(self):
        bulk_users = []
        groups = OcservGroup.objects.all()
        for i in range(self.counter):
            random_index = random.randint(0, groups.count() - 1)
            bulk_users.append(
                OcservUser(
                    group=groups[random_index],
                    username=f"test{i + 1}",
                    password="1234",
                    active=True,
                    rx=Decimal(str(random.randint(0, 100))),
                    tx=Decimal(str(random.randint(0, 100))),
                )
            )
        OcservUser.objects.bulk_create(bulk_users)

    @patch("ocserv.modules.handlers.OcservGroupHandler.update_defaults")
    def test_model_admin_panel_configuration_create_extra(self, *args):
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
    def test_model_admin_panel_configuration_update(self, *args):
        self.admin_configs.default_configs = update_default_configs
        self.admin_configs.save()

    def test_model_admin_panel_configuration_get(self):
        instance = AdminPanelConfiguration.objects.get(id=1)
        self.assertEquals(instance.default_traffic, 10)
        self.assertEqual(instance.default_configs, update_default_configs)

    def test_model_ocserv_create_default_group(self):
        with self.assertRaises(ValidationError) as context:
            OcservGroup.objects.create(name="defaults", desc="defaults group")
        self.assertEqual(
            str(context.exception.detail["error"][0]),
            "Invalid name (defaults) for group",
        )

    @patch("ocserv.modules.handlers.OcservGroupHandler.add_or_update")
    def test_model_ocserv_group_create(self, *args, **kwargs):
        self.create_bulk_ocserv_groups()
        self.assertEqual(OcservGroup.objects.count(), self.counter + 1)
        for i in range(self.counter):
            OcservGroup.objects.create(
                name=f"test-non-bulk-{i}",
                desc=f"test{i} non bulk group",
                configs=update_default_configs,
            )
        self.assertEqual(OcservGroup.objects.count(), (self.counter * 2) + 1)

    @patch("ocserv.modules.handlers.OcservGroupHandler.add_or_update")
    def test_model_ocserv_group_edit(self, *args, **kwargs):
        obj = OcservGroup.objects.create(
            name="test",
            desc="test group",
            configs=update_default_configs,
        )
        with self.assertRaises(ValidationError) as context:
            obj.name = "defaults"
            obj.save()
        self.assertEqual(
            str(context.exception.detail["error"][0]),
            "Invalid name (defaults) for group",
        )
        obj.name = "test_edit"
        obj.desc = "test group edit"
        obj.save()

    @patch("ocserv.modules.handlers.OcservGroupHandler.destroy")
    def test_model_ocserv_group_remove(self, *args):
        OcservGroup.objects.filter(name__startswith="test").delete()
        self.assertEqual(OcservGroup.objects.count(), 1)

    @patch("ocserv.modules.handlers.OcservUserHandler.add_or_update")
    def test_model_ocserv_user_create(self, *args):
        self.create_bulk_ocserv_users()
        self.assertEqual(
            OcservUser.objects.filter(username__startswith="test").count(), self.counter
        )

        groups = OcservGroup.objects.all()
        for i in range(self.counter):
            random_index = random.randint(0, groups.count() - 1)
            OcservUser.objects.create(
                group=groups[random_index],
                username=f"test-non-bulk-{i + 1}",
                password="1234",
                active=True,
                rx=Decimal(str(random.randint(0, 100))),
                tx=Decimal(str(random.randint(0, 100))),
            )
        self.assertEqual(
            OcservUser.objects.filter(username__startswith="test").count(), self.counter * 2
        )

    @patch("ocserv.modules.handlers.OcservUserHandler.add_or_update")
    @patch("ocserv.modules.handlers.OcservUserHandler.delete")
    def test_model_ocserv_user_edit(self, *args):
        user = OcservUser.objects.get(username="init_user")
        self.assertEqual(
            user.password,
            "1234",
        )
        user.group = OcservGroup.objects.last()
        user.username = "edited_user"
        user.password = "edited_password"
        user.save()
        self.assertEqual(
            OcservUser.objects.get(username="edited_user").password,
            "edited_password",
        )
        user.delete()

    @patch("ocserv.modules.handlers.OcservUserHandler.delete")
    def test_model_ocserv_user_remove(self, *args):
        OcservUser.objects.filter(username__startswith="test").delete()
        self.assertEqual(OcservUser.objects.count(), 0)

    def test_model_monthly_traffic_stat(self):
        bulk_stats = []
        today = datetime.now().today()
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
