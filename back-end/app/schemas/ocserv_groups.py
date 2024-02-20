from drf_yasg import openapi
from drf_yasg.utils import swagger_auto_schema

from app.schemas.admin import ocserv_configs_openapi_properties, ocserv_configs_sample


create_ocserv_configs_openapi = openapi.Schema(
    type=openapi.TYPE_OBJECT,
    required=["name", "desc"],
    description="Ocserv Defaults Group Configs",
    properties={
        "configs": openapi.Schema(
            type=openapi.TYPE_OBJECT, properties=ocserv_configs_openapi_properties
        ),
        "name": openapi.Schema(type=openapi.TYPE_STRING),
        "desc": openapi.Schema(type=openapi.TYPE_STRING),
    },
)

update_ocserv_configs_openapi = openapi.Schema(
    type=openapi.TYPE_OBJECT,
    required=[],
    description="Ocserv Defaults Group Configs",
    properties={
        "configs": openapi.Schema(
            type=openapi.TYPE_OBJECT, properties=ocserv_configs_openapi_properties
        ),
        "name": openapi.Schema(type=openapi.TYPE_STRING),
        "desc": openapi.Schema(type=openapi.TYPE_STRING),
    },
)


schemas = {
    "list": {
        "manual_parameters": [
            openapi.Parameter(
                name="args",
                in_=openapi.IN_QUERY,
                type=openapi.TYPE_STRING,
                description="Group filter args(defaults=only default group). if not exists will exclude group=defaults",
                required=False,
            ),
            openapi.Parameter(
                name="name",
                in_=openapi.IN_QUERY,
                type=openapi.TYPE_STRING,
                description="Group Name filter include",
                required=False,
            ),
        ],
        "responses": {
            200: openapi.Response(
                description="Successful Response",
                examples={
                    "application/json": {
                        "result": [
                            {
                                "id": 0,
                                "name": "string",
                                "desc": "string",
                                "configs": ocserv_configs_sample,
                            }
                        ],
                        "pages": 1,
                        "page": 1,
                        "total_count": 0,
                    }
                },
            ),
        },
    },
    "create": {
        "request_body": create_ocserv_configs_openapi,
        "responses": {
            201: openapi.Response(
                description="Successful Response",
                examples={
                    "application/json": {
                        "id": 0,
                        "name": "string",
                        "desc": "string",
                        "configs": ocserv_configs_sample,
                    }
                },
            ),
            400: openapi.Response(
                description="Bad Request",
                examples={
                    "application/json": {
                        "error": [
                            "Name 'defaults' is not a valid name for group",
                            "Ocserv group does not created",
                        ]
                    },
                },
            ),
        },
    },
    "retrieve": {
        "responses": {
            200: openapi.Response(
                description="Successful Response",
                examples={
                    "application/json": {
                        "id": 0,
                        "name": "string",
                        "desc": "string",
                        "configs": ocserv_configs_sample,
                    }
                },
            ),
            404: openapi.Response(
                description="Not Found",
                examples={
                    "application/json": {"error": ["Ocserv group does not exists"]},
                },
            ),
        }
    },
    "partial_update": {
        "request_body": update_ocserv_configs_openapi,
        "responses": {
            202: openapi.Response(
                description="Successful Response",
                examples={
                    "application/json": {
                        "id": 0,
                        "name": "string",
                        "desc": "string",
                        "configs": ocserv_configs_sample,
                    }
                },
            ),
            404: openapi.Response(
                description="Not Found",
                examples={
                    "application/json": {"error": ["Ocserv group does not exists"]},
                },
            ),
            400: openapi.Response(
                description="Bad Request",
                examples={
                    "application/json": {"error": ["Ocserv group does not updated"]},
                },
            ),
        },
    },
    "destroy": {
        "responses": {
            204: "",
            404: openapi.Response(
                description="Not Found",
                examples={
                    "application/json": {"error": ["Ocserv group does not exists"]},
                },
            ),
            400: openapi.Response(
                description="Bad Request",
                examples={
                    "application/json": {"error": ["You can not delete defaults Ocserv group"]},
                },
            ),
        }
    },
}


def get_ocserv_group_schema(schema_name: str, pk=None, method=None):
    swagger_data: dict = schemas.get(schema_name)
    name = schema_name.replace("_", " ").title()
    if swagger_data:
        swagger_data.update({"operation_id": f"Ocserv Groups - {name}"})
    if method:
        swagger_data["method"] = method.lower()
    pk_url = openapi.Parameter(
        name="id", in_=openapi.IN_PATH, type=openapi.TYPE_STRING, description="Group ID"
    )
    if pk:
        if "manual_parameters" in swagger_data:
            swagger_data["manual_parameters"].insert(0, pk_url)
        else:
            swagger_data["manual_parameters"] = [pk_url]
    return swagger_auto_schema(**swagger_data)
