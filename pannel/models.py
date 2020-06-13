from django.db import models
from django.utils.translation import ugettext as _
from django.contrib.auth.models import User


class Users(models.Model):
    user = models.ForeignKey(User, on_delete=models.SET_NULL, null=True)
    name = models.CharField(_('name'), max_length=200,)
    family = models.CharField(_('family'), max_length=20, null=True, blank=True)
    password = models.CharField(_('password'), max_length=200)
    tell_number = models.CharField(max_length=12, null=True, blank=True)
    order_date = models.DateField()
    order_expire = models.DateField(null=True, blank=True)
    lock = models.BooleanField(_('lock'), default=False)

    def __str__(self):
        return self.name

    class Meta:
        verbose_name = _("User Account")
        verbose_name_plural = "User Accounts"


class BlockIP(models.Model):
    ip = models.CharField(max_length=30)
    block = models.BooleanField(default=False)
    faild_try = models.PositiveSmallIntegerField(default=0)

    def __str__(self):
        return self.ip