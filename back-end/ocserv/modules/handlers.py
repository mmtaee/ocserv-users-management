import json
import os
import subprocess

from django.conf import settings

from ocserv.modules.logger import Logger
from ocserv.modules.methods import user_key_creator, ip_bans_creator

logger = Logger()


class OcservServiceHandler:
    @staticmethod
    def subprocess_handler(mode="status", **kwargs):
        lines = kwargs.get("lines", 20)
        if settings.DOCKERIZED:
            command = ["tail", settings.OCSERV_LOG_FILE, f"-n {lines}"]
        else:
            if mode == "journal":
                command = ["sudo", "journalctl", "ocserv.service", f"-n {lines}"]
            else:
                command = ["sudo", "systemctl", mode, "ocserv.service", "--output=json-pretty"]
        p = subprocess.Popen(command, stdout=subprocess.PIPE)
        (output, err) = p.communicate()
        output = output.decode("utf-8")
        if err:
            logger.log(level="critical", message=f"subprocess handler in OcservServiceHandler class({err})")
            return False
        return output.splitlines() if output else None

    def status(self):
        output = self.subprocess_handler()
        return output

    def restart(self):
        self.subprocess_handler(mode="restart")
        output = self.subprocess_handler()
        return output

    def journalctl(self, lines):
        return self.subprocess_handler(mode="journal", **{"lines": lines})


class OcservGroupHandler:
    GROUP_DIR = "/etc/ocserv/groups"

    @staticmethod
    def reload():
        try:
            command = f"sudo /usr/bin/occtl reload"
            os.system(command)
        except Exception as e:
            logger.log(level="critical", message=f"occtl reload service error ({e})")
            return False
        return True

    def add_or_update(self, name, configs=None):
        path = f"{self.GROUP_DIR}/{name}"
        try:
            if not os.path.exists(path):
                touch_command = f"sudo touch {path}"
                subprocess.run(touch_command, shell=True)
            if configs:
                config_str = ""
                for key, val in configs.items():
                    if key.startswith("dns"):
                        key = key[:-1]
                    config_str += f"{key}={val}\n"
            else:
                config_str = "# remove configs by admin \n"
            echo_command = f"echo '{config_str}' | sudo tee {path}"
            subprocess.run(echo_command, shell=True)
        except Exception as e:
            logger.log(level="critical", message=f"add or update ocserv group error ({e})")
            return False
        return self.reload()

    def destroy(self, name):
        path = f"{self.GROUP_DIR}/{name}"
        try:
            if os.path.exists(path):
                rm_command = f"sudo rm {path}"
                subprocess.run(rm_command, shell=True)
            else:
                logger.log(level="warning", message=f"delete group config error (FileNotFoundError)")
        except Exception as e:
            logger.log(level="critical", message=f"destroy ocserv group error ({e})")
        self.reload()

    def update_defaults(self, configs=None):
        dir_path = "/etc/ocserv/defaults"
        path = dir_path + "/group.conf"
        try:
            if not os.path.isdir(dir_path):
                dir_command = f"sudo mkdir {dir_path}"
                subprocess.run(dir_command, shell=True)
            if not os.path.exists(path):
                touch_command = f"sudo touch {path}"
                subprocess.run(touch_command, shell=True)
            if configs:
                config_str = ""
                for key, val in configs.items():
                    if key.startswith("dns"):
                        key = key[:-1]
                    config_str += f"{key}={val}\n"
            else:
                config_str = "# remove configs by admin \n"
            echo_command = f"echo '{config_str}' | sudo tee {path}"
            subprocess.run(echo_command, shell=True)
        except Exception as e:
            logger.log(level="critical", message=f"update defaults ocserv group error ({e})")
        self.reload()

    # def sync_db(self):
    #     defaults_path = f"{self.GROUP_DIR}/defaults/group.conf"
    #     group_path = f"{self.GROUP_DIR}/groups"
    #     data = {
    #         "defaults": {},
    #         "groups": [],
    #     }
    #     try:
    #         groups = [
    #             path
    #             for path in os.listdir(group_path)
    #             if not path.startswith(".") and not path.endswith(".conf") and os.path.isfile(path)
    #         ]
    #     except Exception as e:
    #         logger.log(level="critical", message=f"sync_db ocserv group error ({e})")
    #         return data
    #     for group in groups:
    #         path = f"{group_path}/group"
    #         with open(path, "r") as f:
    #             configs = dict(line.strip().replace(" ", "").split("=") for line in f.readlines() if not line.startswith("#"))
    #             data["groups"].append({"name": group, "configs": configs})
    #             f.close()
    #     with open(defaults_path, "r") as default_f:
    #         configs = dict(line.strip().replace(" ", "").split("=") for line in default_f.readlines() if not line.startswith("#"))
    #         data["defaults"] = configs
    #         default_f.close()
    #     return data


class OcservUserHandler:
    def __init__(self, username=None):
        self.username = username

    def change_group(self, password, group):
        try:
            command = f"sudo /usr/bin/ocpasswd"
            if group:
                command += f" -g {group}"
            command += f" -c /etc/ocserv/ocpasswd {self.username}"
            os.system(command)
        except Exception as e:
            logger.log(level="critical", message=f"change user group error ({e})")
            return False
        return True

    def status_handler(self, active=True):
        """
        ocserv lock and unlock method
        """
        try:
            command = f'sudo /usr/bin/ocpasswd {"-l" if not active else "-u"} -c /etc/ocserv/ocpasswd {self.username}'
            os.system(command)
        except Exception as e:
            logger.log(level="critical", message=f"change user active error ({e})")
            return False
        return True

    def add_or_update(self, password, group=None, active=True):
        try:
            command = f'/usr/bin/echo -e "{password}\n{password}\n" | sudo /usr/bin/ocpasswd'
            if group:
                command += f" -g {group}"
            command += f" -c /etc/ocserv/ocpasswd {self.username}"
            os.system(command)
        except Exception as e:
            logger.log(level="critical", message=f"add user error ({e})")
            return False
        self.status_handler(active)
        return True

    def delete(self):
        try:
            command = f"sudo /usr/bin/ocpasswd  -c /etc/ocserv/ocpasswd -d {self.username}"
            os.system(command)
        except Exception as e:
            logger.log(level="critical", message=f"delete user error ({e})")
            return False
        return True

    def disconnect(self):
        p = subprocess.Popen(
            ["sudo", "/usr/bin/occtl", "disconnect", "user", f"{self.username}"],
            stdout=subprocess.PIPE,
        )
        output, err = p.communicate()
        if output:
            output = output.decode("utf-8")
            if output.strip() != f"user '{self.username}' was disconnected":
                return False
        return True

    @staticmethod
    def online():
        users = []
        try:
            p = subprocess.Popen(
                ["sudo", "/usr/bin/occtl", "-j", "show", "users", "--output=json-pretty"],
                stdout=subprocess.PIPE,
            )
            _users, err = p.communicate()
            if _users and len(_users.decode("utf-8")) > 0:
                _users = json.loads(_users)
                users = user_key_creator(_users)
        except FileNotFoundError as e:
            logger.log(level="critical", message=f"online users error ({e})")
        return users


class OcctlHandler:
    @staticmethod
    def get_command(cmd_name):
        cmd = {
            "show_ip_bans": ["-j", "show", "ip", "bans"],
            "show_ip_ban_points": ["-j", "show", "ip", "bans", "points"],
            "unban_ip": ["unban", "ip"],
            "reload_configs": ["reload"],
            "show_status": ["show", "status"],
            "show_user": ["-j", "show", "user"],
            "show_users": ["-j", "show", "users"],
            "show_iroutes": ["-j", "show", "iroutes"],
            "disconnect_user": ["disconnect", "user"],
            "disconnect_id": ["disconnect", "id"],
            # "show_sessions_all": ["-j" "show", "sessions", "all"],
            # "show_sessions_valid": ["show", "sessions", "valid"],
            # "show_events": ["show", "events"],
        }
        return cmd.get(cmd_name, [])

    @staticmethod
    def subprocess_handler(command):
        exc = ["sudo", "/usr/bin/occtl"] + command
        try:
            p = subprocess.Popen(exc, stdout=subprocess.PIPE)
            (output, err) = p.communicate()
            output = output.decode("utf-8")
            if err:
                logger.log(level="critical", message=f"subprocess handler in OcctlHandler class({err})")
                return False
            return output
        except Exception as e:
            logger.log(level="critical", message=f"occtl command error ({e}), command: {' '.join(exc)}")
        raise ValueError("subprocess_handler not done, see log error")

    def output(self, action, extra_commands):
        command = self.get_command(action)
        command += extra_commands
        command = list(filter(None, command))
        output = self.subprocess_handler(command)
        return output

    def show(self, action):
        result = {}
        if isinstance(action, list):
            for act in action:
                key = act.get("action")
                args = act.get("args", [])
                result[key.replace(" ", "_")] = self.output(key, args)
        elif isinstance(action, dict):
            key = action.get("action")
            result = {key.replace(" ", "_"): self.output(key, action.get("args", []))}
        else:
            result = {action.replace(" ", "_"): self.output(action)}
        if "show_users" in result or "show_user" in result:
            result["show_users"] = user_key_creator(
                result["show_users"] if "show_users" in result else result["show_user"]
            )
            result.pop("show_user", None)
        if "show_ip_bans" in result or "show_ip_ban_points" in result:
            result["show_ip_bans"] = ip_bans_creator(
                result["show_ip_bans"] if "show_ip_bans" in result else result["show_ip_ban_points"]
            )
            result.pop("show_ip_ban_points", None)
        return result

    def reload(self):
        command = ["reload"]
        return self.subprocess_handler(command)
