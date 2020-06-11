from django.contrib import admin
from django.urls import path, include
from django.conf.urls.i18n import i18n_patterns
from django.conf import settings
from django.conf.urls.static import static

from decorator_include import decorator_include
from pannel.decorators import *

urlpatterns = [
    path('i18n/', include('django.conf.urls.i18n')), 
]


urlpatterns += i18n_patterns(
    # TODO : change secret  url in deployment
    path('secret/', admin.site.urls),
    path('admin/', include('admin_honeypot.urls', namespace='admin_honeypot')),
    path('pannel/', decorator_include([user_access, superuser_required], 'pannel.urls')),
    path('', decorator_include(user_access, 'home.urls')),
    path('api/', include('api.urls')),
    prefix_default_language=False,
)

if settings.DEBUG:
   urlpatterns += static(settings.STATIC_URL, document_root=settings.STATIC_ROOT)
