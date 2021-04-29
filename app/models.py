from django.db import models
from django.utils import timezone

class OcservUser(models.Model):
    oc_username = models.CharField(max_length=16)
    oc_password = models.CharField(max_length=32)
    oc_active = models.BooleanField(default=True)
    create = models.DateTimeField(auto_now_add=True)
    expire_date = models.DateField(null=True, blank=True)
    desc = models.TextField(null=True, blank=True)

    def __str__(self):
        return self.oc_username

    @property
    def expire(self):
        if self.expire_date and self.expire_date <= timezone.now().date():
            return True
        return False

    class Meta:
        verbose_name = "Ocserv User"
        verbose_name_plural = "Ocserv Users"
        ordering = ["-create",]


