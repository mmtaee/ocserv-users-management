import subprocess

from django.conf import settings
from django.contrib.auth import authenticate
from django.contrib.auth.hashers import make_password, check_password
from rest_framework import viewsets
from rest_framework.decorators import action
from rest_framework.permissions import AllowAny, IsAuthenticated
from rest_framework.authtoken.models import Token
from rest_framework.response import Response

from app.models import AdminConfig
from app.schemas.admin import get_admin_schema
from app.serializers import AminConfigSerializer
from ocserv.modules.decorators import recaptcha
from ocserv.modules.handlers import OcservUserHandler, OcctlHandler
from ocserv.throttles import custom_throttle


def grep_key(key):
    command = f"grep -r {key} {settings.SOCKET_PASSWD_FILE}"
    result = subprocess.run(command.split(" "), capture_output=True, text=True)
    if result.stderr:
        return False
    if result.stdout:
        return result.stdout
    return


def socket_passwd(key, val=None, delete=False):
    if not val and not delete:
        return grep_key(key)
    elif delete:
        command = f"sed -i /{key}/d {settings.SOCKET_PASSWD_FILE}"
        result = subprocess.run(command.split(" "), capture_output=True, text=True)
        if result.stderr:
            return False
        if result.stdout:
            return result.stdout
    else:
        if not grep_key(key):
            with open(settings.SOCKET_PASSWD_FILE, "a") as f:
                f.write(f"{key}:{val}\n")


class AdminViewSet(viewsets.ViewSet):
    permission_classes = [AllowAny]

    @get_admin_schema("config", security=False)
    @custom_throttle(rate="10/minutes")
    @action(detail=False, methods=["GET"])
    def config(self, request):
        admin_config = AdminConfig.objects.first()
        data = {
            "config": True if admin_config else False,
            "captcha_site_key": admin_config.captcha_site_key if admin_config else None,
        }
        return Response(data)

    @get_admin_schema("create_admin_configs", security=False)
    @custom_throttle(rate="3/minutes")
    @action(detail=False, methods=["POST"], url_path="create")
    def create_admin_configs(self, request):
        if AdminConfig.objects.all().exists():
            return Response({"error": ["Admin config exists!"]}, status=400)
        data = request.data
        serializer = AminConfigSerializer(data=data)
        serializer.is_valid(raise_exception=True)
        admin_config = serializer.save()
        token = Token.objects.create(user=admin_config)
        # socket_passwd(key=admin_config.uu_id, val=token.key)
        return Response(
            {
                "token": token.key,
                "captcha_site_key": admin_config.captcha_site_key,
                "user": admin_config.uu_id,
            },
            status=201,
        )

    @get_admin_schema("login")
    @custom_throttle(rate="10/hours")
    @recaptcha
    @action(detail=False, methods=["POST"])
    def login(self, request):
        data = request.data
        data.get("username")
        if user := authenticate(
            request, username=data.get("username"), password=data.get("password")
        ):
            token, _ = Token.objects.get_or_create(user=user)
            token = token.key
            # socket_passwd(key=user.adminconfig.uu_id, val=token)
            return Response({"token": token, "user": user.adminconfig.uu_id})
        return Response({"error": ["Invalid username or password"]}, status=400)

    @get_admin_schema("logout")
    @action(detail=False, methods=["DELETE"], permission_classes=[IsAuthenticated])
    def logout(self, request):
        token = Token.objects.get(user=request.user)
        # socket_passwd(key=token.user.adminconfig.uu_id, delete=True)
        token.delete()
        return Response(status=204)

    @get_admin_schema("configuration_get", method="GET")
    @get_admin_schema("configuration_patch", method="PATCH")
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
        return Response(serializer.data, status=202)

    @get_admin_schema("dashboard")
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
