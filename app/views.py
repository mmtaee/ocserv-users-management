from django.shortcuts import render, redirect, get_object_or_404
from django.http import HttpResponse
from django.contrib.auth import login , logout, authenticate
from django.contrib.auth.models import User
from django.conf import settings
from django.views.generic import *
from django.contrib.auth.forms import PasswordChangeForm
from django.utils.decorators import method_decorator
from django.contrib.auth.decorators import login_required
from django.contrib.auth import update_session_auth_hash
from django.http import JsonResponse

import os

from.models import *
from .forms import *

class Login(View):
    template_name = 'login.html'

    def get(self, request, *args, **kwargs):
        if request.user.is_authenticated:
            return redirect("app:home")
        return render(request, self.template_name)

    def post(self, request, *args, **kwargs):
        username = request.POST.get('username')
        password = request.POST.get('password')
        auth = authenticate(username=username, password=password)
        if auth:
            login(request, auth)
            if password == 'admin':
                return redirect('app:change_password')
            return redirect("app:home")
        context = {
            'error' : 'Invalid username or password',
        }
        return render(request, self.template_name, context=context)


class Logout(RedirectView):
    url = "/login/"

    def get(self, request, *args, **kwargs):
        logout(request)
        return super().get(request, *args, **kwargs)


@method_decorator(login_required, name='dispatch')
class ChangePassword(View):
    template_name = "change_password.html"

    def get(self, request, *args, **kwargs):
        return render(request, self.template_name)

    def post(self, request, *args, **kwargs):
        username = request.user.username
        old_password = request.POST.get('old_password')
        new_password = request.POST.get('new_password')
        if len(new_password) < 8:
            context = {
                'error' : 'Password must 8 chars',
            }
            return render(request, self.template_name, context)
        auth = authenticate(username=username, password=old_password)
        if auth:
            user = request.user
            user.set_password(new_password)
            user.save()
            update_session_auth_hash(request, user)           
            return redirect("app:home")
        context = {
            'error' : 'Invalid old password',
        }
        return render(request, self.template_name)


@method_decorator(login_required, name='dispatch')
class Home(ListView):
    template_name = "home.html"
    queryset = OcservUser.objects.all()


@method_decorator(login_required, name='dispatch')
class AddUser(View):
    template_name = 'add_user.html'

    def get(self, request, *args, **kwargs):
        context = {
            'form' : AddUserForm,
        }
        return render(request, self.template_name, context)

    def post(self, request, *args, **kwargs):
        form = AddUserForm(request.POST)
        if form.is_valid():
            username = form.cleaned_data.get("oc_username", None)
            form.save()
            command = f'/usr/bin/echo -e "{password}\n{password}\n"|sudo /usr/bin/ocpasswd -c /etc/ocserv/ocpasswd {username}'
            os.system(command)
            code = "<script>window.close();</script>"
            return HttpResponse(code)

@method_decorator(login_required, name='dispatch')
class DelUser(View):

    def get(self, request, *args, **kwargs):
        if request.is_ajax :
            user_id = (request.GET.get("user_id", None)).strip()
            obj = OcservUser.objects.get(id=user_id)
            command = f'sudo /usr/bin/ocpasswd  -c /etc/ocserv/ocpasswd -d {obj.oc_username}'
            os.system(command)
            obj.delete()
            return JsonResponse({}, status=200)
        return JsonResponse({}, status=400)



@method_decorator(login_required, name='dispatch')
class HandlerUser(View):
    def get(self, request, *args, **kwargs):
        if request.is_ajax :
            user_id = (request.GET.get("user_id", None)).strip()
            action = (request.GET.get("action", None)).strip()
            obj = OcservUser.objects.get(id=user_id)
            if action == 'active':
                obj.oc_active = True
                command = f'sudo /usr/bin/ocpasswd  -c /etc/ocserv/ocpasswd -u {obj.oc_username}'
                os.system(command)
            else:
                obj.oc_active = False
                command = f'sudo /usr/bin/ocpasswd  -c /etc/ocserv/ocpasswd -l {obj.oc_username}'
                os.system(command)
            obj.save()
            return JsonResponse({}, status=200)
        return JsonResponse({}, status=400)

