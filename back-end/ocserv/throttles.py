from django.conf import settings
from django.utils.decorators import method_decorator

from rest_framework import exceptions
from rest_framework.throttling import SimpleRateThrottle

from functools import wraps


class CustomThrottle(SimpleRateThrottle):
    THROTTLE_RATES = {}

    def __init__(self, rate="1/minute"):
        if not getattr(self, "rate", None):
            self.rate = rate
        super().__init__()

    def get_cache_key(self, request, view):
        ident = self.get_ident(request)
        return self.cache_format % {"scope": self.scope, "ident": ident}


def custom_throttle(rate, check_docker=False):
    def throttle_decorator(view_func):
        @wraps(view_func)
        def _wrap(request, *args, **kwargs):
            if check_docker and settings.DOCKERIZED:
                return view_func(request, *args, **kwargs)
            throttle_obj = CustomThrottle(rate)
            throttle_durations = []
            if not throttle_obj.allow_request(request, view_func):
                throttle_durations.append(throttle_obj.wait())
            if throttle_durations:
                durations = [duration for duration in throttle_durations if duration is not None]
                duration = max(durations, default=None)
                raise exceptions.Throttled(duration)
            return view_func(request, *args, **kwargs)

        return _wrap

    return method_decorator(throttle_decorator)
