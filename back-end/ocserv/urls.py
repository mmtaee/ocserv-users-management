from django.conf import settings
from django.contrib import admin
from django.urls import path, include

from .routers import router

urlpatterns = [
    path("api/", include(router.urls)),
]

if settings.DEBUG:
    urlpatterns += [
        path("admin/", admin.site.urls),
    ]
