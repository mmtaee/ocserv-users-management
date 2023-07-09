from django.db.models import Sum
from rest_framework import viewsets
from rest_framework.permissions import IsAuthenticated
from rest_framework.response import Response

from datetime import datetime

from app.models import MonthlyTrafficStat


class StatViewSet(viewsets.ViewSet):
    permission_classes = [IsAuthenticated]

    def list(self, request):
        result = {}
        year = datetime.now().year
        stats = MonthlyTrafficStat.objects.filter(year=year)
        months = list(stats.values_list("month", flat=True).distinct())
        for i in months:
            res = stats.filter(month=i).aggregate(
                **{
                    "total_rx": Sum("rx"),
                    "total_tx": Sum("tx"),
                }
            )
            result[i] = res
        total = stats.aggregate(
            **{
                "total_rx": Sum("rx"),
                "total_tx": Sum("tx"),
            }
        )
        total.update({"year": year})
        return Response({"total": total, "result": result, "months": months})
