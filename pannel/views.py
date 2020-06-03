from django.shortcuts import render, get_object_or_404
from django.views import generic, View
from django.utils.decorators import method_decorator
from django.contrib.auth import login, authenticate, logout
from django.utils.translation import ugettext as _
from django.http import Http404 
from django.contrib.auth.hashers import make_password
from django.contrib import messages
from django.conf import settings

from jalali_date import date2jalali

from .forms import UserLoginForm, AddAccountsForm
from .decorators import *
from .models import *

import os
from dateutil.relativedelta import *


# TODO : main page template
class HomePageView(generic.TemplateView):
    template_name = "home.html"


class LoginPageView(View):
    template_name = "login.html"
    form_class = UserLoginForm
    model = BlockIP

    @method_decorator(anonymous_required)
    def dispatch(self, request, *args, **kwargs):
        return super().dispatch(request, *args, **kwargs)

    def get(self, request, *args, **kwargs):
        form = self.form_class
        context = {
            'form' : form,
            'site_key' : settings.GOOGLE_RECAPTCHA_SITE_KEY,
        }
        return render(request, self.template_name, context)

    @method_decorator(check_recaptcha)
    def post(self, request, *args, **kwargs):
        form = self.form_class(request=request, data=request.POST)
        _ip = request.META.get('REMOTE_ADDR')
        user_ip, create = self.model.objects.get_or_create(ip=_ip)

        if form.is_valid() and request.recaptcha_is_valid:
            username = form.cleaned_data.get('username')
            password = form.cleaned_data.get('password')
            user = authenticate(username=username, password=password)
            if user is not None:
                login(request, user)
                remember_me = form.cleaned_data['remember_me']  
                if remember_me:
                    self.request.session.set_expiry(2592000)
                user_ip.faild_try = 0
                user_ip.save()
                return redirect('pannel:main')

        user_ip.faild_try += 1
        user_ip.save()
        if user_ip.faild_try > 3 :
            user_ip.block = True
            user_ip.save()
            raise Http404()
 
        context = {
            'form' : form,
            'site_key' : settings.GOOGLE_RECAPTCHA_SITE_KEY,
        }
        return render(self.request, self.template_name, context)


class LogoutView(generic.RedirectView):
    url = "/"

    def get(self, request, *args, **kwargs):
        logout(request)
        return super().get(request, *args, **kwargs)


class MainPageView(generic.TemplateView):
    template_name = 'account/main.html'


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

        if form.is_valid():
            name = form.cleaned_data.get('name')
            password = form.cleaned_data.get('password')
            new = form.save(commit=False)
            new.password = make_password(password)
            new.user = request.user
            new.save()

            # command = f'/usr/bin/echo -e "{password}\n{password}\n"|sudo /usr/bin/ocpasswd -c /etc/ocserv/ocpasswd {name}'
            # os.system(command)

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
    queryset = Users.objects.filter(lock=False)
    paginate_by = 10


class EditPageView(generic.ListView):
    template_name = "account/edit.html"
    paginate_by = 10
    queryset = Users.objects.all().order_by('-name')


class EditAccountView(generic.RedirectView):

    def post(self, request, *args, **kwargs):
        id = self.kwargs.get('id')
        user = get_object_or_404(Users, id=id)
        password = request.POST.get('password', None)
        month = request.POST.get('month', None)
        if password :
            new_password = make_password(password)
            user.password = new_password
            user.save()
        if month :
            from datetime import datetime
            old_expiry = datetime.strptime(str(user.order_expire), '%Y-%m-%d').date()
            new_expiry = old_expiry + relativedelta(months=int(month))
            user.order_expire = new_expiry
            user.save()
        
        return super().post(request, *args, **kwargs)

    def get_redirect_url(self, *args, **kwargs):
        lang = self.request.LANGUAGE_CODE
        self.url = f"/{lang}/edit/"
        return super().get_redirect_url(*args, **kwargs)

 
class Edit2AccountView(generic.RedirectView):
    
    def dispatch(self, request, *args, **kwargs):
        mode = self.kwargs.get('name')
        id = self.kwargs.get('id')
        user = get_object_or_404(Users, id=id)
        if mode == "unlock" :
            user.lock = False
            user.save()
            
        elif mode == "lock" :
            user.lock = True
            user.save()

        elif mode == "delete" :
            user.delete()

        return super().post(request, *args, **kwargs)

    def get_redirect_url(self, *args, **kwargs):
        lang = self.request.LANGUAGE_CODE
        self.url = f"/{lang}/edit/"
        return super().get_redirect_url(*args, **kwargs)

    

    