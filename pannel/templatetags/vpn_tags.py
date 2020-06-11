from django import template
from django.shortcuts import reverse

import re

register = template.Library()


@register.simple_tag
def change_current_url_lang(url, lang):
    if lang == "fa":
        return f"/{lang}{url}"
    path = re.findall(r'/\w./(.*)', url)[0]
    return f'/{path}'

