

# def pagination(request, queryset, serializer, context=None):
#     if context is None:
#         context = {}
#     query_count = queryset.count()
#     page = int(request.GET.get("page", 1))
#     if query_count == 0:
#         return {"result": [], "pages": 1, "page": int(page if str(page).isnumeric() else 1)}
#     item_per_page = request.GET.get("item_per_page")
#     item_per_page = 100 if not str(item_per_page).isnumeric() else item_per_page
#     if query_count < int(item_per_page):
#         item_per_page = query_count
#     paginator = Paginator(queryset, item_per_page)
#     try:
#         obj = paginator.page(page).object_list
#     except ZeroDivisionError:
#         return queryset.none(), 0
#     except (EmptyPage, PageNotAnInteger, InvalidPage):
#         obj = paginator.page(1).object_list
#     pages = paginator.num_pages
#     if pages in [0, 1]:
#         pages = 1
#     serializer = serializer(instance=obj, context=context, many=True)
#     return {
#         "result": serializer.data,
#         "pages": pages,
#         "page": int(page if str(page).isnumeric() else 1),
#         "total_count": query_count,
#     }