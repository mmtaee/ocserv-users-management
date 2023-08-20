from rest_framework import viewsets
from rest_framework.permissions import IsAuthenticated
from rest_framework.response import Response

from app.models import OcservGroup
from app.serializers import OcservGroupSerializer
from ocserv.modules.handlers import OcservGroupHandler
from ocserv.modules.methods import pagination

group_handler = OcservGroupHandler()


class OcservGroupsViewSet(viewsets.ViewSet):
    permission_classes = [IsAuthenticated]

    def list(self, request):
        groups = OcservGroup.objects.all()
        if request.GET.get("args") != "defaults":
            groups = groups.exclude(name="defaults").order_by("-id")
        data = pagination(request, groups, OcservGroupSerializer)
        return Response(data)

    def create(self, request):
        data = request.data
        if data.get("name") == "defaults":
            return Response({"error": ["Name 'defaults' is not a valid name for group"]}, status=400)
        result = group_handler.add_or_update(name=data.get("name"), configs=data.get("configs"))
        if not result:
            return Response({"error": ["Ocserv group does not created"]}, status=400)
        serializer = OcservGroupSerializer(data=data)
        serializer.is_valid(raise_exception=True)
        serializer.save()
        return Response(serializer.data, status=201)

    def retrieve(self, request, pk=None):
        try:
            group = OcservGroup.objects.get(pk=pk)
        except OcservGroup.DoesNotExist:
            return Response({"error": ["Ocserv group does not exists"]}, status=404)
        serializer = OcservGroupSerializer(group)
        return Response(serializer.data)

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

    def destroy(self, request, pk=None):
        try:
            group = OcservGroup.objects.get(pk=pk)
        except OcservGroup.DoesNotExist:
            return Response({"error": ["Ocserv group does not exists"]}, status=404)
        group_handler.destroy(name=group.name)
        group.delete()
        return Response(status=204)
