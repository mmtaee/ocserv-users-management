from functools import wraps
from django.conf import settings
from django.contrib import messages

import requests

def reCAPTCHA(view_func):
    @wraps(view_func)
    def wrap(request, *args, **kwargs):
        if settings.GOOGLE_CAPTCHA_SECRET_KEY is not None:
            request.recaptcha_is_valid = None
            if request.method == 'POST':
                recaptcha_response = request.POST.get('g-recaptcha-response')
                data = {
                    'secret': settings.GOOGLE_CAPTCHA_SECRET_KEY,
                    'response': recaptcha_response
                }
                r = requests.post('https://www.google.com/recaptcha/api/siteverify', data=data)
                result = r.json()
                if result['success']:
                    request.recaptcha_is_valid = True
                else:
                    request.recaptcha_is_valid = False
                    messages.error(request, 'Invalid reCAPTCHA. Please try again.')
        else:
            request.recaptcha_is_valid = True
        return view_func(request, *args, **kwargs)
    return wrap