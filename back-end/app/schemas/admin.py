from drf_yasg import openapi
from drf_yasg.utils import swagger_auto_schema


ocserv_configs_sample = {
    "rx-data-per-sec": "string",
    "tx-data-per-sec": "string",
    "max-same-clients": "string",
    "ipv4-network": "string",
    "dns1": "string",
    "dns2": "string",
    "no-udp": "string",
    "keepalive": "string",
    "dpd": "string",
    "mobile-dpd": "string",
    "tunnel-all-dns": "string",
    "restrict-user-to-routes": "string",
    "stats-report-time": "string",
    "mtu": "string",
    "idle-timeout": "string",
    "mobile-idle-timeout": "string",
    "session-timeout": "string",
    "no_routes": "string",
    "routes": "string",
}

ocserv_configs_openapi_properties = {
    "rx-data-per-sec": openapi.Schema(type=openapi.TYPE_STRING),
    "tx-data-per-sec": openapi.Schema(type=openapi.TYPE_STRING),
    "max-same-clients": openapi.Schema(type=openapi.TYPE_STRING),
    "ipv4-network": openapi.Schema(type=openapi.TYPE_STRING),
    "dns1": openapi.Schema(type=openapi.TYPE_STRING),
    "dns2": openapi.Schema(type=openapi.TYPE_STRING),
    "no-udp": openapi.Schema(type=openapi.TYPE_STRING),
    "keepalive": openapi.Schema(type=openapi.TYPE_STRING),
    "dpd": openapi.Schema(type=openapi.TYPE_STRING),
    "mobile-dpd": openapi.Schema(type=openapi.TYPE_STRING),
    "tunnel-all-dns": openapi.Schema(type=openapi.TYPE_STRING),
    "restrict-user-to-routes": openapi.Schema(type=openapi.TYPE_STRING),
    "stats-report-time": openapi.Schema(type=openapi.TYPE_STRING),
    "mtu": openapi.Schema(type=openapi.TYPE_STRING),
    "idle-timeout": openapi.Schema(type=openapi.TYPE_STRING),
    "mobile-idle-timeout": openapi.Schema(type=openapi.TYPE_STRING),
    "session-timeout": openapi.Schema(type=openapi.TYPE_STRING),
    "no_routes": openapi.Schema(type=openapi.TYPE_STRING),
    "routes": openapi.Schema(type=openapi.TYPE_STRING),
}


ocserv_configs_openapi = openapi.Schema(
    type=openapi.TYPE_OBJECT,
    description="Ocserv Defaults Group Configs",
    properties=ocserv_configs_openapi_properties,
)


amin_config_serializer_schema = openapi.Response(
    description="Successful Response",
    examples={
        "application/json": {
            "username": "string",
            "captcha_site_key": "string",
            "captcha_secret_key": "string",
            "default_traffic": 0,
            "default_configs": ocserv_configs_sample,
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
    "default_configs": ocserv_configs_openapi,
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
                        "captcha_site_key": "string",
                    },
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
        }
    },
    "create_admin_configs": {
        "request_body": admin_config_request_body_create,
        "responses": {
            201: openapi.Response(
                description="Successful Response",
                examples={
                    "application/json": {
                        "token": "string",
                        "captcha_site_key": "string",
                        "user": "string",
                    }
                },
            ),
            400: openapi.Response(
                description="Bad Request",
                examples={
                    "application/json": {"error": ["Admin config exists!"]},
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
    "login": {
        "responses": {
            201: openapi.Response(
                description="Successful Response",
                examples={"application/json": {"token": "string", "user": "string"}},
            ),
            400: openapi.Response(
                description="Bad Request",
                examples={
                    "application/json": {"error": ["Invalid username or password"]},
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
    "logout": {"responses": {204: "string"}},
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


def get_admin_schema(schema_name: str, method=None, security=True):
    swagger_data: dict = schemas.get(schema_name)
    name = schema_name.replace("_", " ").title()
    if swagger_data:
        swagger_data.update({"operation_id": f"Admin - {name}"})
    if method:
        swagger_data["method"] = method.lower()
    if not security:
        swagger_data["security"] = []
    return swagger_auto_schema(**swagger_data)
