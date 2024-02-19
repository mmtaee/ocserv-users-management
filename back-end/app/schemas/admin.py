from drf_yasg import openapi
from drf_yasg.utils import swagger_auto_schema

from app.serializers import AminConfigSerializer

amin_config_serializer_schema = openapi.Response(
    description="Successful Response",
    examples={
        "application/json": {
            "username": "string",
            "captcha_site_key": "string",
            "captcha_secret_key": "string",
            "default_traffic": 0,
            "default_configs": {
                "rx-data-per-sec": "",
                "tx-data-per-sec": "",
                "max-same-clients": "",
                "ipv4-network": "",
                "dns1": "",
                "dns2": "",
                "no-udp": "",
                "keepalive": "",
                "dpd": "",
                "mobile-dpd": "",
                "tunnel-all-dns": "",
                "restrict-user-to-routes": "",
                "stats-report-time": "",
                "mtu": "",
                "idle-timeout": "",
                "mobile-idle-timeout": "",
                "session-timeout": "",
                "no_routes": "",
                "routes": "",
            },
        },
    },
)


admin_config_request_body_properties = {
    "username": openapi.Schema(
        type=openapi.TYPE_STRING,
        description="Admin Username",
    ),
    "password": openapi.Schema(
        type=openapi.TYPE_STRING,
        description="Admin Login Password",
    ),
    "captcha_site_key": openapi.Schema(
        type=openapi.TYPE_STRING,
        description="Login Captcha Site Key",
    ),
    "captcha_secret_key": openapi.Schema(
        type=openapi.TYPE_STRING,
        description="Login Captcha Secret Key",
    ),
    "default_traffic": openapi.Schema(
        type=openapi.TYPE_STRING,
        description="Ocserv User Default Traffic",
    ),
    "default_configs": openapi.Schema(
        type=openapi.TYPE_OBJECT,
        description="Ocserv Defaults Group Configs",
        properties={
            "rx-data-per-sec": "",
            "tx-data-per-sec": "",
            "max-same-clients": "",
            "ipv4-network": "",
            "dns1": "",
            "dns2": "",
            "no-udp": "",
            "keepalive": "",
            "dpd": "",
            "mobile-dpd": "",
            "tunnel-all-dns": "",
            "restrict-user-to-routes": "",
            "stats-report-time": "",
            "mtu": "",
            "idle-timeout": "",
            "mobile-idle-timeout": "",
            "session-timeout": "",
            "no_routes": "",
            "routes": "",
        },
    ),
}


admin_config_request_body_create = openapi.Schema(
    required=[
        "username",
        "password",
        "captcha_site_key",
        "captcha_secret_key",
        "default_traffic",
        "default_configs",
    ],
    type=openapi.TYPE_OBJECT,
    properties=admin_config_request_body_properties,
)


admin_config_request_body_update = openapi.Schema(
    required=[],
    type=openapi.TYPE_OBJECT,
    properties={
        **admin_config_request_body_properties,
        **{
            "new_password": openapi.Schema(
                type=openapi.TYPE_STRING,
                description="Admin New Login Password",
            ),
        },
    },
)


schemas = {
    "config": {
        "responses": {
            200: openapi.Response(
                description="Successful Response",
                examples={
                    "application/json": {
                        "config": True,
                        "captcha_site_key": "",
                    },
                },
            )
        }
    },
    "create_admin_configs": {
        "request_body": admin_config_request_body_create,
        "responses": {
            201: openapi.Response(
                description="Successful Response",
                examples={"application/json": {"token": "", "captcha_site_key": "", "user": ""}},
            ),
            400: openapi.Response(
                description="Bad Request",
                examples={
                    "application/json": {"error": ["Admin config exists!"]},
                },
            ),
        },
    },
    "login": {
        "responses": {
            201: openapi.Response(
                description="Successful Response",
                examples={"application/json": {"token": "", "user": ""}},
            ),
            400: openapi.Response(
                description="Bad Request",
                examples={
                    "application/json": {"error": ["Invalid username or password"]},
                },
            ),
        },
    },
    "logout": {"responses": {204: ""}},
    "configuration_get": {
        "responses": {
            200: amin_config_serializer_schema,
        }
    },
    "configuration_patch": {
        "request_body": admin_config_request_body_update,
        "responses": {
            202: amin_config_serializer_schema,
            400: openapi.Response(
                description="Bad Request",
                examples={
                    "application/json": {"error": ["Invalid old password"]},
                },
            ),
        },
    },
    "dashboard": {
        "responses": {
            200: openapi.Response(
                description="Successful Response",
                examples={
                    "application/json": {"online_users": [], "show_users": [], "show_ip_bans": []}
                },
            ),
        }
    },
}


def get_admin_schema(schema_name: str, method=None):
    swagger_data: dict = schemas.get(schema_name)
    print(schema_name)
    name = schema_name.replace("_", " ").title()
    if swagger_data:
        swagger_data.update({"operation_id": f"Admin - {name}"})
    if method:
        swagger_data["method"] = method.lower()
    return swagger_auto_schema(**swagger_data)
