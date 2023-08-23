from datetime import datetime

from django.utils import timezone
from rest_framework import viewsets
from rest_framework.decorators import action
from rest_framework.permissions import IsAuthenticated
from rest_framework.response import Response

from app.models import OcservUser, OcservGroup
from app.serializers import OcservUserSerializer
from ocserv.modules.handlers import OcservUserHandler
from ocserv.modules.methods import pagination


user_handler = OcservUserHandler()


class OcservUsersViewSet(viewsets.ViewSet):
    permission_classes = [IsAuthenticated]

    def list(self, request):
        online_users = user_handler.online()
        online_users = [user.get("username") for user in online_users if user.get("username")]
        users = OcservUser.objects.all().select_related("group")
        if username := request.GET.get("username"):
            users = users.filter(username__icontains=username)
        data = pagination(request, users.order_by("-id"), OcservUserSerializer, context={"online_users": online_users})
        return Response(data)

    def create(self, request):
        data = request.data
        username = data.get("username")
        try:
            group = OcservGroup.objects.get(pk=data.get("group"))
        except OcservGroup.DoesNotExist:
            return Response({"error": ["Ocserv group does not exists"]}, status=404)
        if OcservUser.objects.filter(username=username).exists():
            return Response({"error": ["Ocserv User exists"]}, status=404)
        user_handler.username = username
        result = user_handler.add_or_update(password=data.get("password"), group=group.name, active=data.get("active"))
        if result:
            serializer = OcservUserSerializer(data=data)
            serializer.is_valid(raise_exception=True)
            serializer.save()
            return Response(serializer.data, status=201)
        return Response({"error": ["Ocserv User not created"]}, status=400)

    def retrieve(self, request, pk=None):
        try:
            user = OcservUser.objects.select_related("group").get(pk=pk)
        except OcservUser.DoesNotExist:
            return Response({"error": ["Ocserv user does not exist"]}, status=404)
        serializer = OcservUserSerializer(user)
        return Response(serializer.data)

    def partial_update(self, request, pk=None):
        update = True
        result = True
        try:
            user = OcservUser.objects.select_related("group").get(pk=pk)
        except OcservUser.DoesNotExist:
            return Response({"error": ["Ocserv user does not exist"]}, status=404)
        data = request.data.copy()
        data.pop("username", None)
        old_password = user.password
        expire_date = data.get("expire_date")
        if data.get("password") == old_password:
            update = False
        if expire_date and datetime.strptime(expire_date, "%Y-%m-%d").date() <= timezone.now().date():
            data["active"] = False
            update = True
        if update:
            user_handler.username = user.username
            result = user_handler.add_or_update(
                password=data.get("password"), group=user.group.name, active=data.get("active")
            )
        if result:
            serializer = OcservUserSerializer(instance=user, data=data, partial=True)
            serializer.is_valid(raise_exception=True)
            serializer.save()
            return Response(serializer.data, status=202)
        return Response({"error": ["Ocserv user does not updated"]}, status=400)

    def destroy(self, request, pk=None):
        try:
            user = OcservUser.objects.get(pk=pk)
        except OcservUser.DoesNotExist:
            return Response({"error": ["Ocserv user does not exist"]}, status=404)
        user_handler.username = user.username
        result = user_handler.delete()
        if result:
            user.delete()
            return Response(status=204)
        return Response({"error": ["Ocserv user does not deleted"]}, status=400)

    @action(detail=True, methods=["POST"])
    def disconnect(self, request, pk=None):
        try:
            user = OcservUser.objects.select_related("group").get(pk=pk)
        except OcservUser.DoesNotExist:
            return Response({"error": ["Ocserv user does not exist"]}, status=404)
        user_handler.username = user.username
        result = user_handler.disconnect()
        return Response(status=202) if result else Response({"error": ["Ocserv User Disconnect Failed"]}, status=400)

    @action(detail=True, methods=["POST"], url_path="status")
    def user_status_handler(self, request, pk=None):
        try:
            if status := request.data.get("status", None) is None:
                raise NotImplemented
            user = OcservUser.objects.select_related("group").get(pk=pk)
        except OcservUser.DoesNotExist:
            return Response({"error": ["Ocserv user does not exist"]}, status=404)
        except NotImplemented:
            return Response({"error": ["Status of user not found in body"]}, status=400)
        last_status = user.active
        if last_status != status:
            user_handler.username = user.username
            result = user_handler.status_handler(active=status)
            if not result:
                return Response({"error": ["Ocserv User change status Failed"]}, status=400)
        user.active = status
        user.save()
        return Response(status=202)

    # @action(detail=True, methods=["POST"], url_path="group")
    # def change_group(self, request, pk=None):
    #     try:
    #         if group := request.data.get("group", None) is None:
    #             raise NotImplemented
    #         user = OcservUser.objects.select_related("group").get(pk=pk)
    #     except OcservUser.DoesNotExist:
    #         return Response({"error": ["Ocserv user does not exist"]}, status=404)
    #     except NotImplemented:
    #         return Response({"error": ["group of user not found in body"]}, status=400)
    #     last_group = user.group.id
    #     if last_group != group:
    #         try:
    #             group = OcservGroup.objects.get(pk=group)
    #         except OcservGroup.DoesNotExist:
    #             return Response({"error": ["invalid group name"]}, status=400)
    #         result = user_handler.change_group(password=user.password, group=group.name)
    #         if not result:
    #             return Response({"error": ["Ocserv User change group Failed"]}, status=400)
    #     user.group = group
    #     user.save()
    #     return Response(status=202)
