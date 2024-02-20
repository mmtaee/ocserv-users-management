from django.db.models import Sum
from rest_framework import viewsets
from rest_framework.permissions import IsAuthenticated
from rest_framework.response import Response

from datetime import datetime
import calendar

from app.models import MonthlyTrafficStat
from app.schemas.stats import get_ocserv_stats_schema


class StatViewSet(viewsets.ViewSet):
    permission_classes = [IsAuthenticated]

    @get_ocserv_stats_schema("list")
    def list(self, request):
        result = {}
        year = datetime.now().year
        stats = MonthlyTrafficStat.objects.filter(year=year).order_by("month")
        months = list(stats.values_list("month", flat=True).distinct())
        _months = [calendar.month_name[i] for i in months]
        for i in months:
            res = stats.filter(month=i).aggregate(
                **{
                    "total_rx": Sum("rx"),
                    "total_tx": Sum("tx"),
                }
            )
            result[calendar.month_name[i]] = res
        total = stats.aggregate(
            **{
                "total_rx": Sum("rx"),
                "total_tx": Sum("tx"),
            }
        )
        total.update({"year": year})
        current_month = stats.filter(month=datetime.now().month).aggregate(
            **{
                "total_rx": Sum("rx"),
                "total_tx": Sum("tx"),
            }
        )
        current_month.update({"month": datetime.now().strftime("%B")})
        return Response(
            {"total": total, "result": result, "months": _months, "current_month": current_month}
        )
