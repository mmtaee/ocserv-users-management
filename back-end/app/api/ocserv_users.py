from django.core.cache import cache
from rest_framework import viewsets
from rest_framework.decorators import action
from rest_framework.permissions import IsAuthenticated
from rest_framework.response import Response

from ocserv.modules.handlers import OcservUserHandler


class OcservUsersViewSet(viewsets.ViewSet):
    permission_classes = [IsAuthenticated]
