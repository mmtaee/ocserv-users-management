from rest_framework.routers import DefaultRouter

from app.api.admin import AdminViewSet
from app.api.ocserv_users import OcservUsersViewSet

router = DefaultRouter()

router.register("admin", AdminViewSet, basename="admin")
router.register("users", OcservUsersViewSet, basename="users")
