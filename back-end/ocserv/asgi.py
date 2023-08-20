import os

from django.core.asgi import get_asgi_application

os.environ.setdefault('DJANGO_SETTINGS_MODULE', 'ocserv.settings')

application = get_asgi_application()


# import subprocess
# import json
# import os

# from django.core.asgi import get_asgi_application
# from django.conf import settings

# os.environ.setdefault("DJANGO_SETTINGS_MODULE", "ocserv.settings")

# django_application = get_asgi_application()
# receivers = set()


# def socket_passwd(key):
#     command = f"grep -r {key} {settings.SOCKET_PASSWD_FILE}"
#     result = subprocess.run(command.split(" "), capture_output=True, text=True)
#     if result.stderr:
#         return False
#     if result.stdout:
#         return result.stdout


# async def websocket_application(scope, receive, send):
#     while True:
#         event = await receive()

#         if event["type"] == "websocket.connect":
#             query_string = scope["query_string"].decode()
#             if "ws" in query_string and query_string.split("=")[1] == settings.WS_TOKEN:
#                 await send({"type": "websocket.accept"})
#             if "token" in query_string:
#                 items = dict((i.split("=")[0], i.split("=")[1]) for i in query_string.split("&"))
#                 user_token = socket_passwd(items.get("user"))
#                 if user_token and user_token.split(":")[1].strip() == items.get("token"):
#                     receivers.add(send)
#                     await send({"type": "websocket.accept"})
#                     await send(
#                         {
#                             "type": "websocket.send",
#                             "text": "successful added to socket list",
#                         }
#                     )

#         if event["type"] == "websocket.disconnect":
#             if send in receivers:
#                 receivers.remove(send)

#         if event["type"] == "websocket.receive":
#             if send not in receivers and len(receivers):
#                 msg = event["text"]
#                 try:
#                     msg = json.loads(msg)
#                 except Exception as e:
#                     print("json convert error: ", e)
#                     return
#                 text = msg["Text"]
#                 print("****************")
#                 print(text)
#                 print(receivers)
#                 print("****************")

#                 for res in receivers:
#                     await res(
#                         {
#                             "type": "websocket.send",
#                             "text": text,
#                         }
#                     )


# async def application(scope, receive, send):
#     if scope["type"] == "http":
#         await django_application(scope, receive, send)
#     elif scope["type"] == "websocket":
#         await websocket_application(scope, receive, send)
#     else:
#         raise NotImplementedError(f"Unknown scope type {scope['type']}")


# # gunicorn ocserv.asgi:application -w 4 -k uvicorn.workers.UvicornWorker -b 0.0.0.0:8000
