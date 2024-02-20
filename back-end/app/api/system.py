from django.conf import settings
from rest_framework import viewsets
from rest_framework.decorators import action
from rest_framework.permissions import IsAuthenticated
from rest_framework.response import Response

from app.schemas.system import get_ocserv_system_schema
from ocserv.modules.handlers import OcservServiceHandler
from ocserv.modules.logger import Logger
from ocserv.throttles import custom_throttle

service_handler = OcservServiceHandler()
logger = Logger()


class SystemViewSet(viewsets.ViewSet):
    permission_classes = [IsAuthenticated]

    @get_ocserv_system_schema("action_log_list")
    @action(detail=False, methods=["GET"], url_path="action/list")
    def action_log_list(self, request):
        logs = logger.read()
        return Response({"logs": logs})

    @get_ocserv_system_schema("clear_action_log")
    @action(detail=False, methods=["DELETE"], url_path="action/clear")
    def clear_action_log(self, request):
        logger.clear()
        return Response({"message": ["Backend logs successfully cleared"]}, status=202)

    @get_ocserv_system_schema("ocserv_service_status")
    @action(detail=False, methods=["GET"], url_path="ocserv/status")
    def ocserv_service_status(self, request):
        """
        Dockerized key say to user, that this result is coming from docker container
        """
        status = service_handler.status()
        return Response({"status": status, "dockerized": settings.DOCKERIZED})

    @get_ocserv_system_schema("ocserv_service_restart")
    @custom_throttle(rate="1/minute", check_docker=True)
    @action(detail=False, methods=["GET"], url_path="ocserv/restart")
    def ocserv_service_restart(self, request):
        """
        Dockerized key say to user, that this result is coming from docker container
        """
        service_handler.restart()
        status = service_handler.status()
        return Response({"status": status, "dockerized": settings.DOCKERIZED}, status=202)

    @get_ocserv_system_schema("ocserv_service_journal")
    @action(detail=False, methods=["GET"], url_path="ocserv/journal")
    def ocserv_service_journal(self, request):
        lines = request.GET.get("lines", 20)
        logs = service_handler.journalctl(lines)
        return Response({"logs": logs})
