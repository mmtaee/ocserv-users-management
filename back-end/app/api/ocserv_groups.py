from rest_framework import viewsets
from rest_framework.permissions import IsAuthenticated
from rest_framework.response import Response

from app.models import OcservGroup
from app.schemas.ocserv_groups import get_ocserv_group_schema
from app.serializers import OcservGroupSerializer
from ocserv.modules.handlers import OcservGroupHandler
from ocserv.modules.methods import pagination

group_handler = OcservGroupHandler()


class OcservGroupsViewSet(viewsets.ViewSet):
    permission_classes = [IsAuthenticated]

    @get_ocserv_group_schema("list")
    def list(self, request):
        groups = OcservGroup.objects.all()
        if request.GET.get("args") != "defaults":
            groups = groups.exclude(name="defaults").order_by("id")
        if name := request.GET.get("name"):
            groups = groups.filter(name__icontains=name).order_by(
                "id"
                if (ascending := request.GET.get("ascending")) and eval(ascending.title())
                else "-id"
            )
        data = pagination(request, groups, OcservGroupSerializer)
        return Response(data)

    @get_ocserv_group_schema("create")
    def create(self, request):
        data = request.data
        if data.get("name") == "defaults":
            return Response(
                {"error": ["Name 'defaults' is not a valid name for group"]}, status=400
            )
        result = group_handler.add_or_update(name=data.get("name"), configs=data.get("configs"))
        if not result:
            return Response({"error": ["Ocserv group does not created"]}, status=400)
        serializer = OcservGroupSerializer(data=data)
        serializer.is_valid(raise_exception=True)
        serializer.save()
        return Response(serializer.data, status=201)

    @get_ocserv_group_schema("retrieve", pk=True)
    def retrieve(self, request, pk=None):
        try:
            group = OcservGroup.objects.get(pk=pk)
        except OcservGroup.DoesNotExist:
            return Response({"error": ["Ocserv group does not exists"]}, status=404)
        serializer = OcservGroupSerializer(group)
        return Response(serializer.data)

    @get_ocserv_group_schema("partial_update", pk=True)
    def partial_update(self, request, pk=None):
        data = request.data
        try:
            group = OcservGroup.objects.get(pk=pk)
        except OcservGroup.DoesNotExist:
            return Response({"error": ["Ocserv group does not exists"]}, status=404)
        result = group_handler.add_or_update(name=data.get("name"), configs=data.get("configs"))
        if not result:
            return Response({"error": ["Ocserv group does not updated"]}, status=400)
        serializer = OcservGroupSerializer(instance=group, data=data, partial=True)
        serializer.is_valid(raise_exception=True)
        serializer.save()
        return Response(serializer.data, status=202)

    @get_ocserv_group_schema("destroy", pk=True)
    def destroy(self, request, pk=None):
        try:
            group = OcservGroup.objects.get(pk=pk)
        except OcservGroup.DoesNotExist:
            return Response({"error": ["Ocserv group does not exists"]}, status=404)
        if group.name == "defaults":
            return Response({"error": ["You can not delete defaults Ocserv group"]}, status=400)
        group_handler.destroy(name=group.name)
        group.delete()
        return Response(status=204)
