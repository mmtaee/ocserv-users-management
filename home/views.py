from django.shortcuts import render, get_object_or_404
from django.views import generic, View
from django.utils.decorators import method_decorator
from django.contrib.auth import login, authenticate, logout
from django.http import JsonResponse, Http404
from django.utils import translation

from pannel.forms import *
from pannel.decorators import *
from pannel.task import get_user_expiry

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

    def get(self, request, *args, **kwargs):
        logout(request)
        return super().get(request, *args, **kwargs)

    def get_redirect_url(self, *args, **kwargs):
        lang = self.request.LANGUAGE_CODE
        if lang == "fa" :
            self.url = f"/{lang}/"
        else:
            self.url = "/"
        return super().get_redirect_url(*args, **kwargs)

# Ajax
class GetAccountExpiryDateAjaxView(View):
    template_name = "ajax_home.html"

    def get(self, request, *args, **kwargs):
        if request.is_ajax:
            name = (request.GET.get("name", None)).strip()
            password = (request.GET.get("password", None)).strip()
            lang = (request.GET.get("lang", None)).strip()
            if lang == "fa" :
                translation.activate('fa')
            _ip = request.META.get('REMOTE_ADDR')
            user = get_user_expiry.delay(name, password, lang, _ip)
            user = user.get()
            if user :
                context = {
                    'name' : user['name'],
                    'order' : user['order'],
                    'expire' : user['expire'],
                }
                return render(request, self.template_name, context)
                
        return JsonResponse({}, status=400)
