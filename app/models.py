from django.db import models


class OcservUser(models.Model):
    oc_username = models.CharField(max_length=16)
    oc_password = models.CharField(max_length=32)
    oc_active = models.BooleanField(default=True)
    create = models.DateTimeField(auto_now_add=True)
    
    def __str__(self):
        return self.oc_username

    class Meta:
        verbose_name = "Ocserv User"
        verbose_name_plural = "Ocserv Users"
        ordering = ["-create",]


