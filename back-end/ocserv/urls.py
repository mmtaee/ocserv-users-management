from django.conf import settings
from django.contrib import admin
from django.urls import path, include
from rest_framework.permissions import AllowAny


from .routers import router

urlpatterns = [
    path("api/", include(router.urls)),
]

if settings.DEBUG:
    from drf_yasg.views import get_schema_view
    from drf_yasg import openapi

    schema_view_v1 = get_schema_view(
        openapi.Info(
            title="Ocserv Api Backend Service Api Documents",
            default_version="v1",
            description="""
            Ocserv Api Backend Service Api Documents.
            """,
        ),
        public=True,
        permission_classes=[AllowAny],
        patterns=urlpatterns,
    )

    urlpatterns += [
        path("admin/", admin.site.urls),
        path("doc/", schema_view_v1.with_ui("redoc", cache_timeout=0), name="redoc"),
    ]
