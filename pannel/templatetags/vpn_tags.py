from django import template
from django.shortcuts import reverse

import re

register = template.Library()


@register.simple_tag
def change_current_url_lang(url, lang):
    path = re.findall(r'/(fa|en)/(.+)', url)[0][1]
    return f'/{lang}/{path}'
