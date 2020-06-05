from django.urls import path
from .views import *

app_name = "home"

urlpatterns = [
    path('', HomePageView.as_view(), name="home"),
    path('home/', HomePageView.as_view(),),
    path('login/', LoginPageView.as_view(), name="login"),
    path('logout/', LogoutView.as_view(), name="logout"),

    # Ajax
    # path('home/', HomePageView.as_view(),),
]

