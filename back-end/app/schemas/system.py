from drf_yasg import openapi
from drf_yasg.utils import swagger_auto_schema

schemas = {
    "action_log_list": {
        "responses": {
            200: openapi.Response(
                description="Successful Response",
                examples={"application/json": {"logs": []}},
            ),
        },
    },
    "clear_action_log": {
        "responses": {
            200: openapi.Response(
                description="Successful Response",
                examples={"application/json": {"message": ["Backend logs successfully cleared"]}},
            ),
        }
    },
    "ocserv_service_status": {
        "responses": {
            200: openapi.Response(
                description="Successful Response",
                examples={"application/json": {"status": [], "dockerized": False}},
            ),
        }
    },
    "ocserv_service_restart": {
        "responses": {
            202: openapi.Response(
                description="Successful Response",
                examples={"application/json": {"status": [], "dockerized": False}},
            ),
            429: openapi.Response(
                description="Too Many Requests",
                schema=openapi.Schema(
                    type=openapi.TYPE_OBJECT,
                    properties={
                        "message": openapi.Schema(type=openapi.TYPE_STRING),
                    },
                ),
            ),
        }
    },
    "ocserv_service_journal": {
        "responses": {
            200: openapi.Response(
                description="Successful Response",
                examples={"application/json": {"logs": []}},
            ),
        },
    },
}


def get_ocserv_system_schema(schema_name: str, pk=None, method=None):
    swagger_data: dict = schemas.get(schema_name)
    name = schema_name.replace("_", " ").title()
    if swagger_data:
        swagger_data.update({"operation_id": f"Ocserv System - {name}"})
    if method:
        swagger_data["method"] = method.lower()
    return swagger_auto_schema(**swagger_data)
