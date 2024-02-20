from drf_yasg import openapi
from drf_yasg.utils import swagger_auto_schema

schemas = {
    "occtl_show_result": {
        "manual_parameters": [
            openapi.Parameter(
                name="action_command",
                in_=openapi.IN_PATH,
                type=openapi.TYPE_STRING,
                description="Occtl action command name",
                required=True,
            ),
            openapi.Parameter(
                name="args",
                in_=openapi.IN_QUERY,
                type=openapi.TYPE_STRING,
                description="Occtl action command args",
                required=False,
            ),
        ],
        "responses": {
            200: openapi.Response(
                description="Successful Response (dynamic result for each command and occtl action)",
                examples={"application/json": {"dynamic_keys_for_each_command": "value"}},
            ),
            400: openapi.Response(
                description="Bad Request",
                examples={
                    "application/json": {"error": ["Occtl command ({action_command}) not done"]},
                },
            ),
        },
    },
    "reload_server": {
        "responses": {
            200: openapi.Response(
                description="Successful Response (dynamic result for each command and occtl action)",
                examples={"application/json": {"message": "Ocserv service successfully reloaded"}},
            ),
            400: openapi.Response(
                description="Bad Request",
                examples={
                    "application/json": {"error": ["Occtl command (reload) not done"]},
                },
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
        },
    },
    "ocserv_service_journal": {},
}


def get_ocserv_occtl_schema(schema_name: str, pk=None, method=None):
    swagger_data: dict = schemas.get(schema_name)
    name = schema_name.replace("_", " ").title()
    if swagger_data:
        swagger_data.update({"operation_id": f"Ocserv Occtl - {name}"})
    if method:
        swagger_data["method"] = method.lower()
    return swagger_auto_schema(**swagger_data)
