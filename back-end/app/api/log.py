from rest_framework import viewsets
from rest_framework.decorators import action
from rest_framework.permissions import IsAuthenticated
from rest_framework.response import Response

from ocserv.modules.handlers import OcservServiceHandler
from ocserv.modules.logger import Logger
from ocserv.throttles import CustomThrottle

service_handler = OcservServiceHandler()
logger = Logger()


class SystemViewSet(viewsets.ViewSet):
    permission_classes = [IsAuthenticated]
    @action(detail=False, methods=["GET"], url_path="action_log/list")
    def action_log_list(self, request):
        logs = logger.read()
        return Response({"logs": logs})

    @action(detail=False, methods=["DELETE"], url_path="action_log/clear")
    def clear_action_log(self, request):
        logger.clear()
        return Response(status=204)

    @action(detail=False, methods=["GET"], url_path="ocserv/service/status")
    def ocserv_service_status(self):
        status = service_handler.status()
        return Response({"status": status})

    @action(detail=False, methods=["GET"], url_path="ocserv/service/logs")
    def ocserv_service_log(self, request):
        lines = request.GET.get("lines", 20)
        logs = service_handler.journalctl_log(lines)
        return Response({"logs": logs})

    @action(detail=False, methods=["GET"], url_path="ocserv/service/restart", throttle_classes=[CustomThrottle])
    def ocserv_service_restart(self, request):
        service_handler.restart()
        return Response(status=202)
