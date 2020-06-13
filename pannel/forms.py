from django import forms
from django.contrib.auth.forms import AuthenticationForm
from django.contrib.auth.models import User
from django.utils.translation import gettext as _
from django.utils.translation import get_language
from jalali_date.fields import JalaliDateField, SplitJalaliDateTimeField
from jalali_date.widgets import AdminJalaliDateWidget, AdminSplitJalaliDateTime

from .models import *

from functools import partial

DateInput = partial(forms.DateInput, {'class': 'datepicker'})

class UserLoginForm(AuthenticationForm):
    username = forms.CharField(max_length=150, label=_("Username"))
    password = forms.CharField(max_length=32, widget=forms.PasswordInput, label=_("Password"))
    remember_me = forms.BooleanField(required=False)
    
    class Meta:
        model = User
        fields = ('username','password')


class AddAccountsForm(forms.ModelForm):
    password = forms.CharField(widget=forms.PasswordInput)
    # order_expire = forms.DateField(widget=DateInput())
    order_date = forms.DateField(widget=DateInput())
    month_expire = forms.IntegerField(min_value=0)
    
    class Meta:
        model = Users
        fields = ['name', 'family', 'password', 'tell_number', 'order_date', 'month_expire']
        # exclude = ['user', 'lock', 'order_expire']

    def __init__(self, *args, **kwargs):
        lang = get_language()
        super().__init__(*args, **kwargs)
        if lang == 'fa' :
            self.fields['order_date'] = JalaliDateField(widget=AdminJalaliDateWidget)
            # self.fields['order_expire'] = JalaliDateField(widget=AdminJalaliDateWidget)

        self.fields['family'].required = False
        self.fields['tell_number'].required = False

    def clean_name(self):
        clean_data = super().clean()
        name = self.cleaned_data['name']
        if Users.objects.filter(name__iexact=name):
            raise forms.ValidationError(_("Name is Already Exist"))
        return name

    # def clean_order_expire(self):
    #     clean_data = super().clean()
    #     order_date = self.cleaned_data['order_date']
    #     order_expire = self.cleaned_data['order_expire']
    #     if order_date >= order_expire :
    #         raise forms.ValidationError(_("The order expiry date cannot be earlier or equal than the order date"))
    #     return order_expire

