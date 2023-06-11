from rest_framework.throttling import SimpleRateThrottle


class CustomThrottle(SimpleRateThrottle):
    def __init__(self, rate="1/minute", scope=None):
        self.rate = rate
        self.scope = scope
        super().__init__()

    def allow_request(self, request, view):
        return False if self.throttle_success() else True
