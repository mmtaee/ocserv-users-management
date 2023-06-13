from django.contrib import admin

from app.models import OcservUser, OcservGroup

admin.site.register(OcservGroup)
admin.site.register(OcservUser)
