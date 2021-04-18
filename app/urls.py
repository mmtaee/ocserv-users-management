from django.urls import path

from .views import *

app_name = 'app' 


urlpatterns = [
    path('', Home.as_view(), name='home'),
    path('add_user/', AddUser.as_view(), name='add_user'),
    path('del_user/', DelUser.as_view(), name='del_user'),
    path('handler_user/', HandlerUser.as_view(), name='handler_user'),
    path('login/', Login.as_view(), name='login'),
    path('logout/', Logout.as_view(), name='logout'),
    path('change_password/', ChangePassword.as_view(), name='change_password'),
    path('sync_db/', SyncDb.as_view(), name='sync_db'),   
]
