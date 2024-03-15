from django.conf import settings
from django.contrib.auth.models import User
from django.db import models

from rest_framework.exceptions import ValidationError as RestValidationError

from ocserv.modules.handlers import OcservGroupHandler, OcservUserHandler


class AdminPanelConfiguration(models.Model):
    captcha_site_key = models.TextField(null=True, blank=True)
    captcha_secret_key = models.TextField(null=True, blank=True)
    default_traffic = models.PositiveIntegerField(default=10)
    default_configs = models.JSONField(default=dict, null=True, blank=True)

    class Meta:
        verbose_name = "Admin Panel Configuration"
        verbose_name_plural = "Admin Panel Configurations"

    def save(self, *args, **kwargs):
        if not self.pk:
            if AdminPanelConfiguration.objects.exists():
                raise RestValidationError(
                    {"error": ["Amin Configuration already exists"]}
                )
            if not OcservGroup.objects.filter(name="defaults").exists():
                OcservGroup.objects.create(name="defaults", desc="defaults group")

        if self.default_configs and type(self.default_configs) == dict:
            new_configs = {}
            for key, val in self.default_configs.items():
                if key in settings.OSCERV_CONFIG_KEYS and val:
                    new_configs[key] = val
            self.default_configs = new_configs
        else:
            self.default_configs = {}
        OcservGroupHandler().update_defaults(self.default_configs)
        super().save(*args, **kwargs)


class OcservGroup(models.Model):
    name = models.CharField(max_length=128, unique=True)
    desc = models.TextField(null=True, blank=True)
    configs = models.JSONField(default=dict, null=True, blank=True)

    class Meta:
        verbose_name = "Ocserv Group"
        verbose_name_plural = "Ocserv Groups"

    def __str__(self):
        return self.name

    def save(self, *args, **kwargs):
        if (
            self.name == "defaults"
            and OcservGroup.objects.filter(name="defaults").exists()
        ):
            raise RestValidationError({"error": ["Invalid name (defaults) for group"]})
        if self.configs and type(self.configs) == dict:
            new_configs = {}
            for key, val in self.configs.items():
                if key in settings.OSCERV_CONFIG_KEYS:
                    new_configs[key] = val
            self.configs = new_configs
        else:
            self.configs = {}
        if self.name != "defaults":
            OcservGroupHandler().add_or_update(self.name, self.configs)
        super().save(*args, **kwargs)

    def delete(self, *args, **kwargs):
        if self.name == "defaults":
            return False
        OcservGroupHandler().destroy(self.name)
        super().delete(*args, **kwargs)


class OcservUser(models.Model):
    FREE = 1
    MONTHLY = 2
    TOTALLY = 3
    TRAFFIC_MODE_CHOICES = (
        (FREE, "free"),
        (MONTHLY, "monthly"),
        (TOTALLY, "totally"),
    )
    group = models.ForeignKey(OcservGroup, on_delete=models.CASCADE)
    username = models.CharField(max_length=32, unique=True)
    password = models.CharField(max_length=32, null=True, blank=True)
    active = models.BooleanField(default=False)
    create = models.DateField(auto_now_add=True)
    expire_date = models.DateField(null=True, blank=True)
    deactivate_date = models.DateField(null=True, blank=True)
    desc = models.TextField(null=True, blank=True)
    traffic = models.PositiveSmallIntegerField(
        choices=TRAFFIC_MODE_CHOICES, default=MONTHLY
    )
    default_traffic = models.PositiveIntegerField(default=0)
    tx = models.DecimalField(max_digits=14, decimal_places=8, default=0)
    rx = models.DecimalField(max_digits=14, decimal_places=8, default=0)

    def __init__(self, *args, **kwargs):
        super().__init__(*args, **kwargs)
        self.__group = (
            self.group if hasattr(self, "group") and getattr(self, "group") else None
        )

    class Meta:
        verbose_name = "Ocserv User"
        verbose_name_plural = "Ocserv Users"

    def __str__(self):
        return self.username

    def save(self, *args, **kwargs):
        if (admin_config := AdminPanelConfiguration.objects.last()) is None:
            raise RestValidationError({"error": ["Default Admin Configs Not Found"]})
        user_handler = OcservUserHandler(username=self.username)
        if not self.pk:
            if not self.default_traffic and self.traffic != self.FREE:
                self.default_traffic = admin_config.default_traffic
        if self.traffic != self.FREE and self.default_traffic < self.tx:
            self.active = False
        if self.traffic == self.FREE:
            self.default_traffic = 0
        user_handler.add_or_update(
            password=self.password,
            group=self.group.name if self.group.name != "defaults" else None,
            active=self.active,
        )
        super().save(*args, **kwargs)

    def delete(self, *args, **kwargs):
        user_handler = OcservUserHandler(username=self.username)
        user_handler.delete()
        super().delete(*args, **kwargs)


class MonthlyTrafficStat(models.Model):
    user = models.ForeignKey(OcservUser, on_delete=models.CASCADE, related_name="user")
    year = models.PositiveSmallIntegerField(default=2023)
    month = models.PositiveSmallIntegerField(default=1)
    tx = models.DecimalField(max_digits=14, decimal_places=8, default=0)
    rx = models.DecimalField(max_digits=14, decimal_places=8, default=0)

    class Meta:
        verbose_name = "Monthly Traffic Stat"
        verbose_name_plural = "Monthly Traffic Stats"

    def __str__(self):
        return f"{self.month} >> {self.tx}"
