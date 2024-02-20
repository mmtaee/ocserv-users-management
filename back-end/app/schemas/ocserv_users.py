from drf_yasg import openapi
from drf_yasg.utils import swagger_auto_schema

from app.serializers import OcservUserSerializer

ocserv_user_example = {
    "id": 0,
    "username": "string",
    "password": "string",
    "active": True,
    "create": "2019-08-24",
    "expire_date": "2019-08-24",
    "deactivate_date": "2019-08-24",
    "desc": "string",
    "traffic": 1,
    "default_traffic": 0,
    "tx": "string",
    "rx": "string",
    "group": 0,
    "online": False,
    "group_name": "string",
}

schemas = {
    "list": {
        "manual_parameters": [
            openapi.Parameter(
                name="ascending",
                in_=openapi.IN_QUERY,
                type=openapi.TYPE_STRING,
                description="sort ascending by id",
                required=False,
            ),
            openapi.Parameter(
                name="username",
                in_=openapi.IN_QUERY,
                type=openapi.TYPE_STRING,
                description="Ocserv username filter include",
                required=False,
            ),
        ],
        "responses": {
            200: openapi.Response(
                description="Successful Response",
                examples={
                    "application/json": {
                        "result": [ocserv_user_example],
                        "pages": 1,
                        "page": 1,
                        "total_count": 0,
                    }
                },
            ),
        },
    },
    "create": {
        "request_body": OcservUserSerializer,
        "responses": {
            201: openapi.Response(
                description="Successful Response",
                examples={"application/json": ocserv_user_example},
            ),
            400: openapi.Response(
                description="Bad Request",
                examples={
                    "application/json": {
                        "error": ["Ocserv User not created", "Ocserv User exists"]
                    }
                },
            ),
        },
    },
    "retrieve": {
        "responses": {
            200: openapi.Response(
                description="Successful Response",
                examples={"application/json": ocserv_user_example},
            ),
            404: openapi.Response(
                description="Not Found",
                examples={
                    "application/json": {"error": ["Ocserv user does not exist"]},
                },
            ),
        }
    },
    "partial_update": {
        "request_body": OcservUserSerializer,
        "responses": {
            201: openapi.Response(
                description="Successful Response",
                examples={"application/json": ocserv_user_example},
            ),
            404: openapi.Response(
                description="Not Found",
                examples={
                    "application/json": {"error": ["Ocserv user does not exist"]},
                },
            ),
            400: openapi.Response(
                description="Bad Request",
                examples={
                    "application/json": {"error": ["Ocserv user does not updated"]},
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
                    "application/json": {"error": ["Ocserv user does not exist"]},
                },
            ),
            400: openapi.Response(
                description="Bad Request",
                examples={
                    "application/json": {"error": ["Ocserv user does not deleted"]},
                },
            ),
        }
    },
    "disconnect": {
        "responses": {
            "202": openapi.Response(
                description="Successful Response",
                examples={
                    "application/json": {"message": ["Ocserv user disconnected successfully"]},
                },
            ),
            404: openapi.Response(
                description="Not Found",
                examples={
                    "application/json": {"error": ["Ocserv user does not exist"]},
                },
            ),
            400: openapi.Response(
                description="Bad Request",
                examples={
                    "application/json": {"error": ["Ocserv User Disconnect Failed"]},
                },
            ),
        }
    },
    "user_status": {
        "request_body": openapi.Schema(
            required=["status"],
            type=openapi.TYPE_OBJECT,
            properties={
                "status": openapi.Schema(
                    type=openapi.TYPE_STRING, description="User status Boolean"
                )
            },
        ),
        "responses": {
            "202": openapi.Response(
                description="Successful Response",
                examples={
                    "application/json": {"message": ["Ocserv user status changed"]},
                },
            ),
            404: openapi.Response(
                description="Not Found",
                examples={
                    "application/json": {"error": ["Ocserv user does not exist"]},
                },
            ),
            400: openapi.Response(
                description="Bad Request",
                examples={
                    "application/json": {
                        "error": [
                            "Status of user not found in body",
                            "Ocserv User change status Failed",
                        ]
                    },
                },
            ),
        },
    },
    "sync_ocpasswd": {
        "responses": {
            200: openapi.Response(
                description="Successful Response",
                examples={
                    "application/json": {
                        "result": [ocserv_user_example],
                        "new_users": ["username"],
                        "pages": 1,
                        "page": 1,
                        "total_count": 0,
                    }
                },
            ),
        },
    },
}


def get_ocserv_user_schema(schema_name: str, pk=None, method=None):
    swagger_data: dict = schemas.get(schema_name)
    name = schema_name.replace("_", " ").title()
    if swagger_data:
        swagger_data.update({"operation_id": f"Ocserv Users - {name}"})
    if method:
        swagger_data["method"] = method.lower()
    pk_url = openapi.Parameter(
        name="id", in_=openapi.IN_PATH, type=openapi.TYPE_STRING, description="User ID"
    )
    if pk:
        if "manual_parameters" in swagger_data:
            swagger_data["manual_parameters"].insert(0, pk_url)
        else:
            swagger_data["manual_parameters"] = [pk_url]
    return swagger_auto_schema(**swagger_data)
