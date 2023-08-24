from rest_framework.routers import DefaultRouter

from app.api.admin import AdminViewSet
from app.api.log import SystemViewSet
from app.api.occtl import OcctlViewSet
from app.api.ocserv_groups import OcservGroupsViewSet
from app.api.ocserv_users import OcservUsersViewSet
from app.api.stats import StatViewSet

router = DefaultRouter()

router.register("admin", AdminViewSet, basename="admin")
router.register("users", OcservUsersViewSet, basename="users")
router.register("groups", OcservGroupsViewSet, basename="groups")
router.register("occtl", OcctlViewSet, basename="occtl")
router.register("stats", StatViewSet, basename="stats")
router.register("system", SystemViewSet, basename="system")
