import json

from django.core.paginator import Paginator, EmptyPage, PageNotAnInteger, InvalidPage


def pagination(request, queryset, serializer, context=None):
    if context is None:
        context = {}
    query_count = queryset.count()
    page = int(request.GET.get("page", 1))
    if query_count == 0:
        return {"result": [], "pages": 1, "page": int(page if str(page).isnumeric() else 1)}
    item_per_page = request.GET.get("item_per_page")
    item_per_page = 100 if not str(item_per_page).isnumeric() else item_per_page
    if query_count < int(item_per_page):
        item_per_page = query_count
    paginator = Paginator(queryset, item_per_page)
    try:
        obj = paginator.page(page).object_list
    except ZeroDivisionError:
        return queryset.none(), 0
    except (EmptyPage, PageNotAnInteger, InvalidPage):
        obj = paginator.page(1).object_list
    pages = paginator.num_pages
    if pages in [0, 1]:
        pages = 1
    serializer = serializer(instance=obj, context=context, many=True)
    return {
        "result": serializer.data,
        "pages": pages,
        "page": int(page if str(page).isnumeric() else 1),
        "total_count": query_count,
    }


def user_key_creator(users):
    if isinstance(users, str):
        users = json.loads(users)
    return [
        {
            "username": i.get("Username"),
            "hostname": i.get("Hostname"),
            "device": i.get("Device"),
            "remote_ip": i.get("Remote IP"),
            "user_agent": i.get("User-Agent"),
            "since": i.get("_Connected at"),
            "connected_at": i.get("Connected at"),
            "average_rx": i.get("Average RX"),
            "average_tx": i.get("Average TX"),
        }
        for i in users
    ]


def ip_bans_creator(bans):
    if isinstance(bans, str) and len(bans) > 2:
        bans = json.loads(bans)
        return [
            {
                "ip": i.get("IP"),
                "since": i.get("Since"),
                "score": i.get("Score"),
            }
            for i in bans
        ]
    return []
