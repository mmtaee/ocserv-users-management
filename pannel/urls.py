from django.urls import path
from .views import *

app_name = "pannel"


urlpatterns = [

    path('main/', MainPageView.as_view(), name="main"),
    path('add/', AddAccountView.as_view(), name="add"),
    path('list/', ListAccountView.as_view(), name="list"),
    path('edit/', EditPageView.as_view(), name="edit_view"),
    path('edit/<int:id>/<str:name>/', EditAccountView.as_view(), name="edit_account"),
    
    path('service/', ServiceView.as_view(), name="service_view"),
    path('service/<str:name>/', ServiceView.as_view(), name="service_action"),

    # Ajax
    path('ajax/get_accounts/', GetAccountsAjaxView.as_view()),
    path('ajax/get_editform/', GetEditFormAjaxView.as_view()),
    path('ajax/account_result/', GetEditFormAjaxView.as_view()),

    
]
