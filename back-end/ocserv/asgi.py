import json
import os

from django.core.cache import cache
from django.core.asgi import get_asgi_application
from django.conf import settings

os.environ.setdefault("DJANGO_SETTINGS_MODULE", "ocserv.settings")

django_application = get_asgi_application()


async def websocket_application(scope, receive, send):
    connections = []
    throttle_rate = 10

    while True:
        event = await receive()
        throttle = True
        if event["type"] == "websocket.connect":
            await send({"type": "websocket.accept"})
        if event["type"] == "websocket.disconnect":
            if send in connections:
                connections.remove(send)
            break
        if event["type"] == "websocket.receive":
            client_ip = scope.get("client")[0]
            cache_key = f"websocket_throttle:{client_ip}"
            if cache_val := cache.get(cache_key) is None:
                cache.set(cache_key, 0)
                cache_val = 0
            msg = event["text"]
            if int(cache_val) >= throttle_rate:
                return
            try:
                msg = json.loads(msg)
            except Exception as e:
                print("json convert error: ", e)
                return

            if "Authorization" in msg:
                _token = msg["Authorization"].split(" ")[1]
                if cache.get(_token) and send not in connections:
                    throttle = False
                    connections.append(send)
            elif "WSToken" in msg:
                _token = msg["WSToken"]
                text = msg["Text"]
                if msg["WSToken"] == settings.WS_TOKEN:
                    throttle = False
                    for con in connections:
                        await con(
                            {
                                "type": "websocket.send",
                                "text": text,
                            }
                        )
            if throttle:
                cache.incr(cache_key)
                return


async def application(scope, receive, send):
    if scope["type"] == "http":
        await django_application(scope, receive, send)
    elif scope["type"] == "websocket":
        await websocket_application(scope, receive, send)
    else:
        raise NotImplementedError(f"Unknown scope type {scope['type']}")


# uvicorn ocserv.asgi:application --host 127.0.0.1 --port 8000
