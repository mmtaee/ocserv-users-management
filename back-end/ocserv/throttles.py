from rest_framework.throttling import SimpleRateThrottle


class CustomThrottle(SimpleRateThrottle):
    THROTTLE_RATES = {}

    def __init__(self, rate="1/minute"):
        if not getattr(self, "rate", None):
            self.rate = rate
        super().__init__()

    def get_cache_key(self, request, view):
        ident = self.get_ident(request)
        return self.cache_format % {"scope": self.scope, "ident": ident}
