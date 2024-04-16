from unittest.mock import patch

from app.api.ocserv_groups import OcservGroupsViewSet
from app.models import OcservGroup
from app.tests import OcservTestAbstract


class OcservGroupApiTest(OcservTestAbstract):

    @patch("ocserv.modules.handlers.OcservGroupHandler.add_or_update")
    def setUp(self, *args, **kwargs) -> None:
        OcservGroup.objects.get_or_create(name="test_group", desc="test group")

    def test_group_list(self):
        self.token = None
        request = self.factory.get("/groups/", headers=self.get_header)
        response = OcservGroupsViewSet.as_view({"get": "list"})(request)
        self.assertTrue(
            len(list(filter(lambda x: x.get("name") == "test_group", response.data.get("result"))))
            > 0,
            "len of groups must be greater than zero",
        )
        self.assertEqual(response.data.get("total_count"), 1)
        self.assertEqual(response.data.get("page"), 1)
        self.assertEqual(response.data.get("pages"), 1)

    def test_group_detail(self):
        pass

    def test_group_create(self):
        pass

    def test_group_update(self):
        pass

    def test_group_delete(self):
        pass
