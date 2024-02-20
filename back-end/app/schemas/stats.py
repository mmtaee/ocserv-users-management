from drf_yasg import openapi
from drf_yasg.utils import swagger_auto_schema

schemas = {
    "list": {
        "responses": {
            200: openapi.Response(
                description="Successful Response",
                examples={
                    "application/json": {
                        "total": {"total_rx": 0, "total_tx": 0},
                        "result": {"january": 0, "february": 0},
                        "months": ["january", "february"],
                        "current_month": {"total_rx": 0, "total_tx": 0, "month": "january"},
                    }
                },
            ),
        },
    }
}


def get_ocserv_stats_schema(schema_name: str, pk=None, method=None):
    swagger_data: dict = schemas.get(schema_name)
    name = schema_name.replace("_", " ").title()
    if swagger_data:
        swagger_data.update({"operation_id": f"Ocserv Stats - {name}"})
    if method:
        swagger_data["method"] = method.lower()
    return swagger_auto_schema(**swagger_data)
