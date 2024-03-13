import logging
from sqlite3 import IntegrityError

from django.contrib.auth import authenticate
from django.contrib.auth.hashers import make_password, check_password
from django.contrib.auth.models import User
from rest_framework import viewsets
from rest_framework.decorators import action
from rest_framework.permissions import AllowAny, IsAuthenticated
from rest_framework.authtoken.models import Token
from rest_framework.response import Response

from app.models import AdminPanelConfiguration
from app.schemas.admin import get_admin_schema
from app.serializers import AminConfigSerializer, UserSerializer
from ocserv.modules.decorators import recaptcha
from ocserv.modules.handlers import OcservUserHandler, OcctlHandler
from ocserv.throttles import custom_throttle


class AdminViewSet(viewsets.ViewSet):
    permission_classes = [AllowAny]

    @get_admin_schema("config", security=False)
    @custom_throttle(rate="30/minutes")
    @action(detail=False, methods=["GET"])
    def config(self, request):
        admin_config = AdminPanelConfiguration.objects.last()
        data = {
            "config": True if admin_config else False,
            "captcha_site_key": admin_config.captcha_site_key if admin_config else None,
        }
        if request.user.is_authenticated:
            data["user"] = {
                "username": request.user.username,
                "is_admin": request.user.is_superuser,
            }
        return Response(data)

    @get_admin_schema("create_admin_configs", security=False)
    @custom_throttle(rate="3/minutes")
    @action(detail=False, methods=["POST"], url_path="create")
    def create_admin_configs(self, request):
        if AdminPanelConfiguration.objects.exists():
            return Response({"error": ["Admin config exists!"]}, status=400)
        data = request.data
        try:
            user = User.objects.create_superuser(
                username=data.pop("username"),
                password=data.pop("password"),
                is_superuser=True,
            )
        except IntegrityError as e:
            logging.warning(f"user already exists: {e}")
        serializer = AminConfigSerializer(data=data)
        serializer.is_valid(raise_exception=True)
        admin_config = serializer.save()
        token = Token.objects.create(user=user)
        return Response(
            {
                "token": token.key,
                "captcha_site_key": admin_config.captcha_site_key,
                "user": {
                    "username": user.username,
                    "is_admin": user.is_superuser,
                },
            },
            status=201,
        )

    @get_admin_schema("login")
    @custom_throttle(rate="10/hours")
    @recaptcha
    @action(detail=False, methods=["POST"])
    def login(self, request):
        data = request.data
        if user := authenticate(
            request, username=data.get("username"), password=data.get("password")
        ):
            token, _ = Token.objects.get_or_create(user=user)
            token = token.key
            return Response(
                {
                    "token": token,
                    "user": {
                        "username": user.username,
                        "is_admin": user.is_superuser,
                    },
                }
            )
        return Response({"error": ["Invalid username or password"]}, status=400)

    @get_admin_schema("logout")
    @action(detail=False, methods=["DELETE"], permission_classes=[IsAuthenticated])
    def logout(self, request):
        token = Token.objects.get(user=request.user)
        token.delete()
        return Response(status=204)

    @get_admin_schema("configuration_get", method="GET")
    @get_admin_schema("configuration_patch", method="PATCH")
    @action(
        detail=False, methods=["GET", "PATCH"], permission_classes=[IsAuthenticated]
    )
    def configuration(self, request):
        data = request.data.copy()
        admin_config = AdminPanelConfiguration.objects.last()
        if request.method == "GET":
            serializer = AminConfigSerializer(admin_config)
            return Response(serializer.data, status=200)
        serializer = AminConfigSerializer(
            data=data,
            instance=admin_config,
            partial=True,
        )
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

    @get_admin_schema("change_password")
    @action(detail=False, methods=["POST"], permission_classes=[IsAuthenticated])
    def change_password(self, request):
        data = request.data
        if (old_password := data.get("old_password")) is None or (
            password := data.get("password")
        ) is None:
            return Response(
                {"error": ["old_password or password is required!"]}, status=400
            )
        if not check_password(old_password, request.user.password):
            return Response({"error": ["invalid old password"]}, status=400)
        try:
            request.user.password = make_password(password)
            request.user.save()
        except Exception as e:
            return Response({"error": [f"error: {e}"]}, status=400)
        return Response(status=202)

    @get_admin_schema("staffs_list", method="get")
    @get_admin_schema("staffs_create", method="post")
    @action(detail=False, methods=["GET", "POST"], permission_classes=[IsAuthenticated])
    def staffs(self, request):
        if not request.user.is_superuser:
            return Response(
                {"error": ["you have not access to this route"]}, status=403
            )
        if request.method == "GET":
            users = User.objects.all().exclude(id=request.user.id)
            serializer = UserSerializer(users, many=True)
            status = 200
        else:
            data = request.data
            user = User.objects.create_user(
                username=data["username"],
                password=data["password"],
                is_staff=False,
                is_superuser=False,
            )
            serializer = UserSerializer(user)
            status = 202
        return Response(serializer.data, status=status)

    @get_admin_schema("delete_staff")
    @action(
        detail=False,
        methods=["DELETE"],
        permission_classes=[IsAuthenticated],
        url_path="staffs/(?P<pk>[^/.]+)",
    )
    def delete_staff(self, request, pk=None):
        if not request.user.is_superuser:
            return Response(status=403)
        try:
            staff = User.objects.get(id=pk)
        except User.DoesNotExist:
            return Response({"error": ["Staff not found!"]}, status=404)
        if staff.is_superuser:
            return Response(
                {"error": ["you have not access to delete admin role"]}, status=403
            )
        staff.delete()
        return Response(status=204)
