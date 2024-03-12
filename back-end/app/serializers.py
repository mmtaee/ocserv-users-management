from django.contrib.auth.models import User
from rest_framework import serializers

from .models import AdminPanelConfiguration, OcservGroup, OcservUser, MonthlyTrafficStat


class UserSerializer(serializers.ModelSerializer):
    is_admin = serializers.SerializerMethodField(default=False)

    class Meta:
        model = User
        fields = ("id", "username", "is_admin")

    @staticmethod
    def get_is_admin(instance):
        return instance.is_superuser


class AminConfigSerializer(serializers.ModelSerializer):
    class Meta:
        model = AdminPanelConfiguration
        fields = (
            "captcha_site_key",
            "captcha_secret_key",
            "default_traffic",
            "default_configs",
        )
        read_only_fields = ("id",)


class OcservGroupSerializer(serializers.ModelSerializer):
    class Meta:
        model = OcservGroup
        fields = "__all__"
        read_only_fields = ("id",)
        extra_kwargs = {"name": {"trim_whitespace": True}}


class OcservUserSerializer(serializers.ModelSerializer):
    class Meta:
        model = OcservUser
        fields = "__all__"
        read_only_fields = ("id",)
        extra_kwargs = {
            "username": {"trim_whitespace": True},
            "password": {"trim_whitespace": True},
        }

    def to_representation(self, instance):
        rep = super().to_representation(instance)
        online_users = self.context.get("online_users", [])
        rep["password"] = (
            "Hashed by ocserv" if not instance.password else instance.password
        )
        rep["online"] = False if instance.username not in online_users else True
        rep["group"] = instance.group.id
        rep["group_name"] = instance.group.name
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
