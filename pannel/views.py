from django.shortcuts import render, get_object_or_404
from django.views import generic, View
from django.utils.decorators import method_decorator
from django.utils.translation import ugettext as _
from django.http import Http404, JsonResponse
from django.contrib.auth.hashers import make_password
from django.contrib import messages
from django.conf import settings
from django.utils import translation

from .forms import AddAccountsForm
from .decorators import *
from .models import *
from .task import status_ocserv, restart_ocserv

import os
from datetime import datetime, timedelta
from dateutil.relativedelta import *

class MainPageView(generic.ListView):
    template_name = 'account/main.html'
    paginate_by = 10
    queryset = Users.objects.filter(order_expire__gte=datetime.today(), order_expire__lte=datetime.today()+timedelta(days=15), lock=False).order_by("-name")


class AddAccountView(View):
    template_name = 'account/add.html'
    form_class = AddAccountsForm

    def get(self, request, *args, **kwargs):
        form = self.form_class
        context = {
            'form' : form,
        }
        return render(request, self.template_name, context)

    def post(self, request, *args, **kwargs):
        form = self.form_class(request.POST)
        print(form.is_valid())
        if form.is_valid():
            name = form.cleaned_data.get('name')
            password = form.cleaned_data.get('password')
            order_date = form.cleaned_data.get("order_date")
            month = form.cleaned_data.get('month_expire')
            new = form.save(commit=False)
            new.password = make_password(password)
            new.user = request.user
            new.order_expire = order_date + relativedelta(days=int(month)*31)
            new.save()

            command = f'/usr/bin/echo -e "{password}\n{password}\n"|sudo /usr/bin/ocpasswd -c /etc/ocserv/ocpasswd {name}'
            os.system(command)

            msg = f"{name} Added to Ocserv Successfully"
            messages.success(request, msg)
            context = {
                'form' : self.form_class,
            }
        else :
            context = {
                'form' : form,
            }
        return render(self.request, self.template_name, context)


class ListAccountView(generic.ListView):
    template_name = "account/list.html"
    queryset = Users.objects.filter(lock=False).order_by("-name")
    paginate_by = 10


class EditPageView(generic.ListView):
    template_name = "account/edit.html"
    queryset = Users.objects.all().order_by('-name')
    paginate_by = 10


class EditAccountView(generic.RedirectView):

    def dispatch(self, request, *args, **kwargs):
        mode = self.kwargs.get('name', None)
        id = self.kwargs.get('id')
        user = get_object_or_404(Users, id=id)

        if not mode :
            return super().dispatch(request, *args, **kwargs)

        if mode == "unlock" :
            user.lock = False
            user.order_date = datetime.now().date()
            user.order_expire = datetime.now().date()
            user.save()

            command = f'ocpasswd -c /etc/ocserv/ocpasswd -u {user.name}'
            os.system(command)

        elif mode == "lock" :
            user.lock = True
            user.save()

            command = f'ocpasswd -c /etc/ocserv/ocpasswd -l {user.name}'
            os.system(command)

        elif mode == "delete" :
            user.delete()

            command = f'ocpasswd -c /etc/ocserv/ocpasswd -d {user.name}'
            os.system(command)

        return super().dispatch(request, *args, **kwargs)

    def post(self, request, *args, **kwargs):
        id = self.kwargs.get('id')
        user = get_object_or_404(Users, id=id)
        password = request.POST.get('password', None)
        month = request.POST.get('month', None)
        next_url = request.META.get('HTTP_REFERER', None)

        if password :
            new_password = make_password(password)
            user.password = new_password
            user.save()

            command = f"echo -e '{new_password}\n{new_password}\n'|ocpasswd -c /etc/ocserv/ocpasswd {user.name}"
            os.system(command)

        if month :
            from datetime import datetime
            old_expiry = datetime.strptime(str(user.order_expire), '%Y-%m-%d').date()
            new_expiry = old_expiry + relativedelta(days=int(month)*31)
            user.order_expire = new_expiry
            user.save()

        if next_url :
            return redirect(next_url)
        
        return super().post(request, *args, **kwargs)

    def get_redirect_url(self, *args, **kwargs):
        lang = self.request.LANGUAGE_CODE
        self.url = f"/{lang}/pannel/edit/"
        return super().get_redirect_url(*args, **kwargs)

     
class ServiceView(generic.TemplateView):
    template_name = "service.html"

    def dispatch(self, request, *args, **kwargs):
        name = self.kwargs.get('name', None)

        if name is None :
            return super().dispatch(request, *args, **kwargs)

        if name == "status" :
            service_status = status_ocserv.delay()
            result = service_status.get()

            if 'status' not in result:
                messages.error(request, f"Service : {result['service']} , Disabled")

            elif "running" in result['status']:
                messages.success(request, f"Service : {result['service']}")
                messages.success(request, f"Status : {result['status']}")
                messages.success(request, f"Since : {result['since']}")
                messages.success(request, f"Uptime : {result['uptime']}")

            else :
                messages.warning(request, f"Service : {result['service']}")
                messages.warning(request, f"Status : {result['status']}")
                messages.warning(request, f"Since : {result['since']}")
                messages.warning(request, f"Downtime : {result['uptime']}")
        
        if name == "restart":
            service_restart = restart_ocserv.delay()
            result = service_restart.get()

            if 'status' not in result:
                messages.warning(request, f"Restart : {result['restart']}")
                messages.error(request, f"Service : {result['service']} , Disabled")

            elif "running" in result['status']:
                messages.success(request, f"Restart : {result['restart']}")
                messages.success(request, f"Service : {result['service']}")
                messages.success(request, f"Status : {result['status']}")
                messages.success(request, f"Since : {result['since']}")
                messages.success(request, f"Uptime : {result['uptime']}")

            else :
                messages.success(request, f"Restart : {result['restart']}")
                messages.warning(request, f"Service : {result['service']}")
                messages.warning(request, f"Status : {result['status']}")
                messages.warning(request, f"Since : {result['since']}")
                messages.warning(request, f"Downtime : {result['uptime']}")

        return redirect("pannel:service_view")
        
    
# Ajax

class GetAccountsAjaxView(View):

    def get(self, request, *args, **kwargs):
        if request.is_ajax:
            queryset = Users.objects.all().order_by("-name")
            query = None
            for item in queryset :
                if query == None :
                    query = ""
                    query += str(item.name)
                else :
                    query = query +  "," + str(item.name)
            return JsonResponse(data={'results': query.rstrip()}, status=200)
        return JsonResponse({}, status=400)


class GetEditFormAjaxView(View):
    template_name = "ajax/edit.html"

    def get(self, request, *args, **kwargs):
        if request.is_ajax:
            name = (request.GET.get("name", None)).strip()
            lang = (request.GET.get("lang", None)).strip()
            user = Users.objects.filter(name__iexact=name).first()

            if lang == "en" :
                translation.activate('en')

            context = {
                'user_edit' : user,
            }
            return render(request, self.template_name, context)
            
        return JsonResponse({}, status=400)
