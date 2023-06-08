from django.contrib.auth.hashers import make_password
from rest_framework import serializers

from .models import AdminConfig


class AminConfigSerializer(serializers.ModelSerializer):
    class Meta:
        model = AdminConfig
        fields = ("username", "password", "captcha_site_key", "captcha_secret_key", "default_traffic")
        read_only_fields = ("id",)

    def create(self, validated_data):
        validated_data["password"] = make_password(validated_data["password"])
        return AdminConfig.objects.create(**validated_data)

    def to_representation(self, instance):
        rep = super().to_representation(instance)
        rep.pop("password")
        return rep


from rest_framework import serializers

from .models import OcservGroup


class OcservGroupSerializer(serializers.ModelSerializer):
    class Meta:
        model = OcservGroup
        fields = "__all__"
        read_only_fields = ("id",)
        extra_kwargs = {"name": {"trim_whitespace": True}}


from rest_framework import serializers

from .models import OcservUser, MonthlyTrafficStat


class OcservUserSerializer(serializers.ModelSerializer):
    class Meta:
        model = OcservUser
        fields = "__all__"
        read_only_fields = ("id",)
        extra_kwargs = {"username": {"trim_whitespace": True}}

    def to_representation(self, instance):
        rep = super().to_representation(instance)
        online_users = self.context.get("online_users", [])
        rep["password"] = "Hashed by ocserv" if not instance.password else instance.password
        rep["online"] = False if instance.username not in online_users else True
        rep["group"] = instance.group.name
        return rep


class MonthlyTrafficStatSerializer(serializers.ModelSerializer):
    class Meta:
        model = MonthlyTrafficStat
        fields = "__all__"
        read_only_fields = ("id",)

    def to_representation(self, instance):
        rep = super().to_representation(instance)
        rep["user"] = instance.user.username
        return rep