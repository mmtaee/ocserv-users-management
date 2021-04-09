from django import forms

from .models import * 


class AddUserForm(forms.ModelForm):

    class Meta:
        model = OcservUser
        exclude = ('oc_active',)

    def __init__(self, *args, **kwargs):
        super().__init__(*args, **kwargs)
        for field in self.fields:
            self.fields[field].widget.attrs.update({'class': 'form-control'})

    def clean_oc_username(self):
        oc_username = self.cleaned_data['oc_username']
        obj = OcservUser.objects.filter(oc_username=oc_username)
        if obj.exists():
            raise forms.ValidationError("Username is Already Exist")
        return oc_username