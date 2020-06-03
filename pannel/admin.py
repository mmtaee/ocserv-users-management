from django.contrib import admin

from .models import BlockIP, Users

admin.site.register(Users)
admin.site.register(BlockIP)