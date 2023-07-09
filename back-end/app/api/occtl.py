from rest_framework import viewsets
from rest_framework.decorators import action
from rest_framework.permissions import IsAuthenticated
from rest_framework.response import Response

from ocserv.modules.handlers import OcctlHandler
from ocserv.throttles import CustomThrottle

occtl_handler = OcctlHandler()


class OcctlViewSet(viewsets.ViewSet):
    """
    unban ip >> Unban the specified IP
    reload >> Reloads the server configuration
    show status >> Prints the status and statistics of the server
    show ip bans >> Prints the banned IP addresses
    show ip ban points >> Prints all the known IP addresses which have points
    show user >> Prints information on the specified user
    show users >> Prints the connected users
    show iroutes >> Prints the routes provided by users of the server
    show events	>> Provides information about connecting users
    """

    permission_classes = [IsAuthenticated]

    @action(detail=False, methods=["GET"], url_path="command/(?P<action_command>[^/.]+)")
    def occtl_show_result(self, request, action_command=None):
        result = occtl_handler.show(action={"action": action_command, "args": [request.GET.get("args")]})
        if not result:
            return Response({"error": [f"Occtl command ({action_command}) not done"]}, status=400)
        return Response(result, status=200)

    @action(detail=False, methods=["GET"], throttle_classes=[CustomThrottle])
    def reload(self, request):
        result = occtl_handler.reload()
        if not result:
            return Response({"error": ["Occtl command (reload) not done"]}, status=400)
        return Response(status=202)
