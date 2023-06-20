from django.contrib.auth import authenticate
from django.contrib.auth.hashers import make_password, check_password
from rest_framework import viewsets
from rest_framework.decorators import action
from rest_framework.permissions import AllowAny, IsAuthenticated
from rest_framework.authtoken.models import Token
from rest_framework.response import Response

from app.models import AdminConfig
from app.serializers import AminConfigSerializer
from ocserv.modules.decorators import recaptcha
from ocserv.modules.handlers import OcservUserHandler, OcctlHandler
from ocserv.throttles import CustomThrottle


# TODO: add throttle for admin


class AdminViewSet(viewsets.ViewSet):
    permission_classes = [AllowAny]

    @action(detail=False, methods=["GET"], throttle_classes=[CustomThrottle(rate="10/min")])
    def config(self, request):
        admin_config = AdminConfig.objects.first()
        data = {
            "config": True if admin_config else False,
            "captcha_site_key": admin_config.captcha_site_key if admin_config else None,
        }
        return Response(data)

    @action(detail=False, methods=["POST"], url_path="create", throttle_classes=[CustomThrottle(rate="3/min")])
    def create_admin_configs(self, request):
        if AdminConfig.objects.all().exists():
            return Response({"error": ["Admin config exists!"]}, status=400)
        data = request.data
        serializer = AminConfigSerializer(data=data)
        serializer.is_valid(raise_exception=True)
        admin_config = serializer.save()
        token = Token.objects.create(user=admin_config)
        return Response({"token": token.key, "captcha_site_key": admin_config.captcha_site_key}, status=201)

    @recaptcha
    @action(detail=False, methods=["POST"], throttle_classes=[CustomThrottle(rate="20/hour")])
    def login(self, request):
        data = request.data
        data.get("username")
        if user := authenticate(request, username=data.get("username"), password=data.get("password")):
            token, _ = Token.objects.get_or_create(user=user)
            return Response({"token": token.key})
        return Response({"error": ["Invalid username or password"]}, status=400)

    @action(detail=False, methods=["DELETE"], permission_classes=[IsAuthenticated])
    def logout(self, request):
        Token.objects.get(user=request.user).delete()
        return Response(status=204)

    @action(detail=False, methods=["GET", "PATCH"], permission_classes=[IsAuthenticated])
    def configuration(self, request):
        data = request.data.copy()
        admin_config = AdminConfig.objects.first()
        if request.method == "GET":
            serializer = AminConfigSerializer(admin_config)
            return Response(serializer.data, status=200)
        if new_password := data.pop("new_password", None):
            if check_password((data.pop("password", None)), request.user.username):
                data["password"] = make_password(new_password)
            else:
                return Response({"error": ["Invalid old password"]}, status=400)
        serializer = AminConfigSerializer(data=data, instance=admin_config, partial=True)
        serializer.is_valid(raise_exception=True)
        serializer.save()
        return Response(status=202)

    @action(detail=False, methods=["GET"], permission_classes=[IsAuthenticated])
    def dashboard(self, request):
        online = OcservUserHandler()
        occtl = OcctlHandler()
        actions = [
            {"action": "show_ip_bans"},
            {"action": "show_status"},
            {"action": "show_iroutes"},
        ]
        server_stats = occtl.show(action=actions)
        result = {"online_users": online.online() or [], **server_stats}
        return Response(result)
