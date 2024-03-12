from django.core.cache import cache
from django.utils.decorators import method_decorator
from rest_framework.response import Response

from functools import wraps
import requests

from app.models import AdminPanelConfiguration


def check_recaptcha(view_func):
    @wraps(view_func)
    def _wrapper(request, *args, **kwargs):
        cache.clear()
        config = cache.get("admin_config")
        if not config:
            config = AdminPanelConfiguration.objects.last()
            if not config:
                return Response(status=400)
            cache.set("admin_config", config)
        if config.captcha_secret_key:
            data = {"secret": config.captcha_secret_key, "response": request.data.get("token")}
            response = requests.post(
                url="https://www.google.com/recaptcha/api/siteverify",
                data=data,
            )
            result = response.json()
            if not result["success"]:
                return Response(
                    {"error": ["Captcha challenge failed"]},
                    status=400,
                )
        return view_func(request, *args, **kwargs)

    return _wrapper


recaptcha = method_decorator(check_recaptcha)

