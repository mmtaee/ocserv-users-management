import json
import os
import subprocess

from ocserv.modules.logger import Logger

logger = Logger()


class OcservGroupHandler:
    GROUP_DIR = "/etc/ocserv/groups"

    @staticmethod
    def reload():
        try:
            command = f"/usr/bin/occtl reload"
            os.system(command)
        except Exception as e:
            logger.log(level="critical", message=f"occtl reload service error ({e})")

    def add_or_update(self, name, configs=None):
        path = f"{self.GROUP_DIR}/{name}"
        try:
            if not os.path.exists(path):
                os.mknod(path)
            with open(path, "w") as f:
                if configs:
                    for key, val in configs.items():
                        if key.startswith("dns"):
                            key = key[:-1]
                        config = f"{key}={val}\n"
                        f.write(config)
                else:
                    f.write("# remove configs by admin \n")
            f.close()
        except Exception as e:
            logger.log(level="critical", message=f"add or update ocserv group error ({e})")
        self.reload()

    def destroy(self, name):
        path = f"{self.GROUP_DIR}/{name}"
        try:
            if os.path.exists(path):
                os.remove(path)
            else:
                logger.log(level="warning", message=f"delete group config error (FileNotFoundError)")
        except Exception as e:
            logger.log(level="critical", message=f"destroy ocserv group error ({e})")
        self.reload()

    def update_defaults(self, configs=None):
        path = f"{self.GROUP_DIR}/defaults/group.conf"
        try:
            if not os.path.exists(path):
                os.mknod(path)
            with open(path, "w") as f:
                if configs:
                    for key, val in configs.items():
                        if key.startswith("dns"):
                            key = key[:-1]
                        config = f"{key}={val}\n"
                        f.write(config)
                else:
                    f.write("# remove configs by admin \n")
            f.close()
        except Exception as e:
            logger.log(level="critical", message=f"update defaults ocserv group error ({e})")
        self.reload()

    def sync_db(self):
        defaults_path = f"{self.GROUP_DIR}/defaults/group.conf"
        group_path = f"{self.GROUP_DIR}/groups"
        data = {
            "defaults": {},
            "groups": [],
        }
        try:
            groups = [
                path
                for path in os.listdir(group_path)
                if not path.startswith(".") and not path.endswith(".conf") and os.path.isfile(path)
            ]
        except Exception as e:
            logger.log(level="critical", message=f"sync_db ocserv group error ({e})")
            return data
        for group in groups:
            path = f"{group_path}/group"
            with open(path, "r") as f:
                configs = dict(
                    line.strip().replace(" ", "").split("=") for line in f.readlines() if not line.startswith("#")
                )
                data["groups"].append({"name": group, "configs": configs})
                f.close()
        with open(defaults_path, "r") as default_f:
            configs = dict(
                line.strip().replace(" ", "").split("=") for line in default_f.readlines() if not line.startswith("#")
            )
            data["defaults"] = configs
            default_f.close()
        return data


class OcservUserHandler:
    def __init__(self, username=None):
        self.username = username

    @staticmethod
    def online():
        users = []
        try:
            p = subprocess.Popen(
                ["/usr/bin/occtl", "-j", "show", "users", "--output=json-pretty"],
                stdout=subprocess.PIPE,
            )
            _users, err = p.communicate()
            if _users:
                if len(_users.decode("utf-8")) > 0:
                    _users = json.loads(_users)
                    users = [i["Username"] for i in _users]
        except FileNotFoundError as e:
            logger.log(level="critical", message=f"online users error ({e})")
        return {"users": users}


class OcctlHandler:
    @staticmethod
    def get_command(cmd_name):
        cmd = {
            "show_ip_bans": ["show", "ip", "bans"],
            "show_ip_ban_points": ["show", "ip", "bans", "points"],
            "unban_ip": ["unban", "ip"],
            "reload_configs": ["reload"],
            "show_status": ["show", "status"],
            "show_user": ["show", "user"],
            "show_users": ["show", "users"],
            "show_iroutes": ["show", "iroutes"],
            "show_events": ["show", "events"],
        }
        return cmd.get(cmd_name, [])

    @staticmethod
    def subprocess_handler(command):
        exc = ["/usr/bin/occtl"] + command
        try:
            p = subprocess.Popen(exc, stdout=subprocess.PIPE)
            (output, err) = p.communicate()
            output = output.decode("utf-8")
            if err:
                logger.log(level="critical", message=f"subprocess handler in OcctlHandler class({err})")
                return False
            return output.splitlines() if output else None
        except Exception as e:
            logger.log(level="critical", message=f"occtl command error ({e}), command: {' '.join(exc)}")
        return ""

    def output(self, action, extra_commands):
        command = self.get_command(action)
        command += extra_commands
        command = list(filter(None, command))
        output = self.subprocess_handler(command)
        return {action: output}

    def __call__(self, *args, **kwargs):
        result = {}
        action = kwargs.get("action")
        if isinstance(action, list):
            for act in action:
                result.update(**self.output(act.get("action"), act.get("extra_commands", [])))
        if isinstance(action, dict):
            result = self.output(action.get("action"), action.get("extra_commands", []))
        return result
