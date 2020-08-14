from django import template
from django.shortcuts import reverse

import re

register = template.Library()


@register.simple_tag
def change_current_url_lang(url, current_lang):
    if current_lang == "fa":
        return f"/en{url}"
    path = re.findall(r'/\w./(.*)', url)[0]
    return f'/{path}'

