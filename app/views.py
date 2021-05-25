from django.shortcuts import render, redirect, get_object_or_404
from django.http import HttpResponse, JsonResponse
from django.contrib.auth import login as django_login , logout as django_logout, authenticate
from django.contrib.auth.models import User
from django.conf import settings
from django.views.generic import *
from django.contrib.auth.forms import PasswordChangeForm
from django.utils.decorators import method_decorator
from django.contrib.auth.decorators import login_required
from django.contrib.auth import update_session_auth_hash
from django.db.models import Q
from django.core.serializers import serialize

import os, subprocess, re, json, subprocess
from ratelimit.decorators import ratelimit

from.models import *
from .forms import *
from .decorators import *

@method_decorator(ratelimit(key='ip', rate='10/d'), name='dispatch')
@method_decorator(reCAPTCHA, name='post')
class Login(View):
    template_name = 'login.html'

    def get(self, request, *args, **kwargs):
        if request.user.is_authenticated:
            return redirect("app:home")
        context = {
            'GOOGLE_CAPTCHA_SITE_KEY' : settings.GOOGLE_CAPTCHA_SITE_KEY
        }
        return render(request, self.template_name, context=context)

    def post(self, request, *args, **kwargs):
        context = {
                'GOOGLE_CAPTCHA_SITE_KEY' : settings.GOOGLE_CAPTCHA_SITE_KEY
        }
        if request.recaptcha_is_valid:
            username = request.POST.get('username')
            password = request.POST.get('password')
            auth = authenticate(username=username, password=password)
            if auth:
                django_login(request, auth)
                if password == 'admin':
                    return redirect('app:change_password')
                return redirect("app:home")
            context ['error'] = 'Invalid username or password'
        return render(request, self.template_name, context=context)


@method_decorator(ratelimit(key='ip', rate='50/d'), name='dispatch')
class Logout(RedirectView):
    url = "/login/"

    def get(self, request, *args, **kwargs):
        django_logout(request)
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
    queryset = OcservUser.objects.all().order_by("create")
    paginate_by = 10

    def get_context_data(self, *args, **kwargs):
        context = super().get_context_data(*args, **kwargs)
        context["search"] = True
        return context

@method_decorator(login_required, name='dispatch')
class SearchView(View):
    template_name = "home.html"

    def get(self, request, *args, **kwargs):
        id = kwargs.get("id")
        if id:
            oc_user = OcservUser.objects.filter(id=id)
            return render(request, self.template_name, context={'object_list' : oc_user, 'search' : False})
        code = "<script>window.close();</script>"
        return HttpResponse(code)



@method_decorator(login_required, name='dispatch')
class AddUser(View):
    template_name = "add_user.html"

    def get(self, request, *args, **kwargs):
        context = {
            'form' : AddUserForm,
            'last_users' : OcservUser.objects.all().order_by("-create")[:7],
        }
        return render(request, self.template_name, context)

    def post(self, request, *args, **kwargs):
        form = AddUserForm(request.POST)
        if form.is_valid():
            username = form.cleaned_data.get("oc_username")
            password = form.cleaned_data.get("oc_password")
            form.save()
            command = f'/usr/bin/echo -e "{password}\n{password}\n"|sudo /usr/bin/ocpasswd -c /etc/ocserv/ocpasswd {username}'
            os.system(command)
            context = {
                'form' : AddUserForm,
                'success' : username,
            }
        else:
            context = {
                'form' : form,
                'error' : True,
            }
        context['last_users'] = OcservUser.objects.all().order_by("-create")[:7]
        return render(request, self.template_name, context)


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


@method_decorator(login_required, name='dispatch')
class SyncDb(View):

    def get(self, request, *args, **kwargs):
        if request.is_ajax :
            ocpasswd_users = {}
            with open("/etc/ocserv/ocpasswd", "r") as f:
                lines = f.readlines()
                for line in lines:
                    user_items = line.strip().split(":")
                    ocpasswd_users[user_items[0]] = False if  user_items[2].startswith("!") else True
            oc_queryset = OcservUser.objects.all()
            for user in oc_queryset:
                if user.oc_username not in ocpasswd_users:
                    user.delete()
                else:
                    user.oc_active = ocpasswd_users[user.oc_username]
                    user.save()
            for user in ocpasswd_users:
                filtered = oc_queryset.filter(oc_username=user)
                if not filtered.exists():              
                    OcservUser.objects.create(
                        oc_username=user,
                        oc_password="OCSERV_HASH_PASSWORD",
                        oc_active=ocpasswd_users[user]
                    )
            return JsonResponse({}, status=200)
        return JsonResponse({}, status=400)


@method_decorator(login_required, name='dispatch')
class Service(TemplateView):
    template_name = "service.html"


@method_decorator(login_required, name='dispatch')
class ServiceHandler(View):

    def get(self, request, *args, **kwargs):
        if request.is_ajax :
            service_status = {}
            mode = (request.GET.get("mode", None)).strip() 
            password = (request.GET.get("password", None)).strip() 
            service = (request.GET.get("service", None)).strip() 
            if not request.user.check_password(password):
                return JsonResponse({'error' : 'Invalid password'}, status=400)

            def subprocess_handler(_mode):
                p =  subprocess.Popen(["systemctl", _mode, service, "--output=json-pretty"], stdout=subprocess.PIPE)
                (output, err) = p.communicate()
                output = output.decode('utf-8')
                if err:
                    return False
                elif output:
                    return output
                # for restart mode result
                return None

            def status_handler(output):
                status_regx= r"Active:(.*) since (.*);(.*)"
                active = False
                for line in output.splitlines():
                    status_search = re.search(status_regx, line)
                    if status_search:
                        service_status['status'] = status_search.group(1).strip()
                        service_status['since'] = status_search.group(2).strip()
                        service_status['uptime'] = status_search.group(3).strip()
                        active = True
                if not active:
                    service_status['status'] = "Deactive: (Not Running)"
                service_status['service'] = service
                return service_status

            if mode == "restart":
                output = subprocess_handler(mode)
                if output is None:
                    output = subprocess_handler("status")
                    if output and output is not None:
                        service_status = status_handler(output)
                else:
                    service_status['status'] = "Restarting Error: (Not Running)"
            
            elif mode == "status":
                output = subprocess_handler(mode)
                if output and output is not None:
                    service_status = status_handler(output)

            return JsonResponse(service_status, status=200)
        return JsonResponse({}, status=400) 



@method_decorator(login_required, name='dispatch')
class SerchUserHandler(View):
    
    def get(self, request, *args, **kwargs):
        if request.is_ajax :
            search_param = (request.GET.get("search_param", None)).strip() 
            oc_users = OcservUser.objects.filter(
                Q(oc_username__icontains=search_param)|
                Q(oc_password__icontains=search_param)|
                Q(desc__icontains=search_param)
            )
            response = list(map( lambda user:{
                "id" : user.id,
                "value" : user.oc_username ,
                "oc_password": user.oc_password ,
                "oc_active": user.oc_active ,
                "expire_date": user.expire_date.strftime("%Y-%m-%d") if user.expire_date else None,
                "desc": user.desc ,
            }, oc_users ))
            return JsonResponse(response, status=200, safe=False)
        return JsonResponse({}, status=400) 


@method_decorator(login_required, name='dispatch')
class OnlineUsers(View):
    
    def get(self, request, *args, **kwargs):
        if request.is_ajax :
            p =  subprocess.Popen(["sudo", "occtl", "-j",  "show", "users", "--output=json-pretty"], stdout=subprocess.PIPE)
            (output, err) = p.communicate()
            if output:
                output = output.decode('utf-8')
                if len(output) > 0:
                    output = json.loads(output)
                    output = [i['Username'] for i in output]
                else:
                    output = []
                return JsonResponse(output, status=200, safe=False)
        return JsonResponse({}, status=400) 