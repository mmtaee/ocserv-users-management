from rest_framework import viewsets
from rest_framework.decorators import action
from rest_framework.permissions import IsAuthenticated
from rest_framework.response import Response

from app.schemas.occtl import get_ocserv_occtl_schema
from ocserv.modules.handlers import OcctlHandler
from ocserv.throttles import custom_throttle

occtl_handler = OcctlHandler()


class OcctlViewSet(viewsets.ViewSet):
    permission_classes = [IsAuthenticated]

    @get_ocserv_occtl_schema("occtl_show_result")
    @action(detail=False, methods=["GET"], url_path="command/(?P<action_command>[^/.]+)")
    def occtl_show_result(self, request, action_command=None):
        """
        unban ip >> Unban the specified IP <br />
        reload >> Reloads the server configuration <br />
        show status >> Prints the status and statistics of the server <br />
        show ip bans >> Prints the banned IP addresses <br />
        show ip ban points >> Prints all the known IP addresses which have points <br />
        show user >> Prints information on the specified user <br />
        show users >> Prints the connected users <br />
        show iroutes >> Prints the routes provided by users of the server <br />
        show events	>> Provides information about connecting users <br />
        """
        result = occtl_handler.show(
            action={"action": action_command, "args": [request.GET.get("args")]}
        )
        if not result:
            return Response({"error": [f"Occtl command ({action_command}) not done"]}, status=400)
        return Response(result, status=200)

    @get_ocserv_occtl_schema("reload_server")
    @custom_throttle(rate="1/minute")
    @action(detail=False, methods=["GET"], url_path="reload")
    def reload_server(self, request):
        result = occtl_handler.reload()
        if not result:
            return Response({"error": ["Occtl command (reload) not done"]}, status=400)
        return Response({"message": "Ocserv service successfully reloaded"}, status=202)
