from rest_framework import viewsets
from rest_framework.decorators import action
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

    @action(detail=False, methods=["GET"], url_path="(?P<action_command>[^/.]+)")
    def occtl_show_result(self, request, action_command=None):
        extra = []
        if action_command == "show_user":
            if user := request.GET.get("user") is None:
                return Response({"error": ["Action 'show_user' need username in query params"]}, status=400)
            extra.append(user)
        result = occtl_handler.show(action=action_command, extra=extra)
        if not result:
            return Response({"error": [f"Occtl command ({action_command}) not done"]}, status=400)
        return Response(status=202)

    @action(detail=False, methods=["GET"], throttle_classes=[CustomThrottle])
    def reload(self, request):
        result = occtl_handler.reload()
        if not result:
            return Response({"error": ["Occtl command (reload) not done"]}, status=400)
        return Response(status=202)

    @action(detail=False, methods=["GET"])
    def unban_ip(self, request):
        if ip := request.GET.get("ip") is None:
            return Response({"error": ["Action 'unban_ip' need ip in query params"]}, status=400)
        result = occtl_handler.unban_ip(ip)
        if not result:
            return Response({"error": [f"Occtl command (unban_ip) not done for ip({ip})"]}, status=400)
        return Response(status=202)