from django.urls import path
from .views import *

app_name = "pannel"


urlpatterns = [
    path('', HomePageView.as_view(), name="home"),
    path('home/', HomePageView.as_view(),),
    path('login/', LoginPageView.as_view(), name="login"),
    path('logout/', LogoutView.as_view(), name="logout"),

    path('main/', MainPageView.as_view(), name="main"),
    path('add/', AddAccountView.as_view(), name="add"),
    path('list/', ListAccountView.as_view(), name="list"),
    path('edit/', EditPageView.as_view(), name="edit_view"),
    path('edit/<int:id>/<str:name>/', EditAccountView.as_view(), name="edit_account"),

]
