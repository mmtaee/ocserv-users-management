from django.shortcuts import redirect
from django.http import Http404
from django.conf import settings
from django.contrib import messages
from django.core.exceptions import PermissionDenied
from django.utils.translation import ugettext as _

# from functools import wraps
import requests

from .models import BlockIP


def superuser_required(view_func):

    def wrap(request, *args, **kwargs):

        if not request.user.is_superuser :
            raise PermissionDenied()
        
        return view_func(request, *args, **kwargs)

    return wrap


def user_access(view_func):

    def wrap(request, *args, **kwargs):
        ip = request.META.get('REMOTE_ADDR')
        user_ip, create = BlockIP.objects.get_or_create(ip=ip)
        if user_ip.block == True :
            raise Http404(_("Your access to this website is blocked"))

        return view_func(request, *args, **kwargs)

    return wrap


def check_recaptcha (view_func):

    def wrap(request, *args, **kwargs):
        request.recaptcha_is_valid = None
        if request.method == 'POST':
            recaptcha_response = request.POST.get('g-recaptcha-response')
            data = {
                'secret': settings.GOOGLE_RECAPTCHA_SECRET_KEY,
                'response': recaptcha_response
            }
            r = requests.post('https://www.google.com/recaptcha/api/siteverify', data=data)
            result = r.json()
            if result['success']:
                request.recaptcha_is_valid = True
            else:
                request.recaptcha_is_valid = False
                messages.error(request, _('Invalid reCAPTCHA. Please try again.'))
        return view_func(request, *args, **kwargs)

    return wrap


# TODO : change return path !Important
def anonymous_required(view_func, redirect_to=None):

    def wrap(request, *args, **kwargs):

        if request.user is not None and request.user.is_authenticated :
            return redirect("pannel:main")

        return view_func(request, *args, **kwargs)

    return wrap
