import _axios from "@/plugins/axios";
import { AxiosError, AxiosResponse } from "axios";
import {
  User,
  AdminConfig,
  AminLogin,
  Config,
  Dashboard,
  UserPagination,
  GroupPagination,
  OcservUser,
  OcservGroup,
  Occtl,
  Stats,
  URLParams,
  SyncResponse,
} from "./types";
import store from "@/plugins/store";

class Services {
  private status_code: number = 500;
  public method: string = "get";
  public baseUrl: string = "";
  public path: string = "";
  public params: URLParams | null = null;
  public overlay: boolean = true;
  private axiosMethods: { [K: string]: Function } = {
    get: _axios.get,
    post: _axios.post,
    patch: _axios.patch,
    put: _axios.put,
    delete: _axios.delete,
  };

  private validatePath(): void {
    if (!this.path.endsWith("/")) {
      this.path = this.path + "/";
    }
  }

  public async request(data?: any): Promise<any> {
    if (this.overlay) {
      store.commit("setLoadingOverlay", {
        active: true,
        text: "Requesting ...",
      });
    }
    this.validatePath();
    var url: string = this.baseUrl + this.path;
    if (this.params) {
      url += "?";
      Object.keys(this.params).forEach((key, index) => {
        let val = this.params![key];
        if (index != 0) {
          url += "&";
        }
        url += `${key}=${val}`;
      });
    }
    return await this.axiosMethods[this.method]((url = url), (data = data))
      .then((response: AxiosResponse) => {
        if (response) {
          this.status_code = response.status;
          if (this.status_code == 400) {
            store.commit("setSnackBar", {
              text: response.data.error.join("<br/>"),
              color: "error",
            });
            return {};
          } else if (this.status_code == 401) {
            localStorage.removeItem("token");
            location.href = "/login";
          } else if (this.status_code == 403) {
            store.commit("setSnackBar", {
              text: "forbiden error",
              color: "warning",
            });
            return {};
          }
          return response.data;
        } else {
          throw new Error("Server error");
        }
      })
      .catch((error: AxiosError) => {
        this.status_code = error.response?.status!;
        if (error.response?.data) {
          let e = error.response.data;
          if (e.detail) {
            store.commit("setSnackBar", {
              text: e.detail,
              color: "orange",
            });
          } else if (e.error) {
            store.commit("setSnackBar", {
              text: e.error instanceof Array ? e.error.join(",") : e.error,
              color: "orange",
            });
          }
        } else {
          store.commit("setSnackBar", {
            text: "response failed from server",
            color: "error",
          });
        }
        return Promise.reject(error);
      })
      .finally((_: null) => {
        if (this.overlay) {
          store.commit("setLoadingOverlay", {
            active: false,
            text: null,
          });
        }
        this.params = null;
        this.path = "";
      });
  }

  public status(): number {
    return this.status_code;
  }
}

class AdminServiceApi extends Services {
  constructor() {
    super();
    this.baseUrl = "/admin/";
    this.path = "";
  }
  public async config(): Promise<Config> {
    this.method = "get";
    this.path = "config/";
    return this.request();
  }
  public async login(data: AminLogin): Promise<{ token: string; user: User }> {
    this.method = "post";
    this.path = "login/";
    return this.request(data);
  }
  public async logout(): Promise<void> {
    this.method = "delete";
    this.path = "logout/";
    await this.request();
  }
  public async create_configs(data: object): Promise<Config> {
    this.method = "post";
    this.path = "create/";
    return await this.request(data);
  }
  public async patch_configuration(data: object): Promise<null> {
    this.method = "patch";
    this.path = "configuration/";
    await this.request(data);
    return null;
  }
  public async get_configuration(): Promise<AdminConfig> {
    this.method = "get";
    this.path = "configuration/";
    return await this.request();
  }
  public async dashboard(): Promise<Dashboard> {
    this.method = "get";
    this.path = "dashboard/";
    return this.request();
  }
  public async change_password(data: object): Promise<void> {
    this.method = "post";
    this.path = "change_password/";
    return this.request(data);
  }
  public async get_staff(): Promise<User> {
    this.method = "get";
    this.path = "staffs/";
    return this.request();
  }

  public async create_staff(data: object): Promise<User> {
    this.method = "post";
    this.path = "staffs/";
    return this.request(data);
  }

  public async delete_staff(id: number) {
    this.method = "delete";
    this.path = `staffs/${id}/`;
    return this.request();
  }
}

class OcservUserApi extends Services {
  constructor() {
    super();
    this.baseUrl = "/users/";
    this.path = "";
  }
  public async users(params?: URLParams | null): Promise<UserPagination> {
    this.method = "get";
    this.path = "";
    if (params) {
      this.params = params;
    }
    return this.request();
  }
  public async create_user(data: OcservUser): Promise<OcservUser> {
    this.method = "post";
    this.path = "";
    return this.request((data = data));
  }
  public async update_user(pk: number, data: OcservUser): Promise<OcservUser> {
    this.method = "patch";
    this.path = `${pk}/`;
    return this.request((data = data));
  }
  public async delete_user(pk: number): Promise<OcservUser> {
    this.method = "delete";
    this.path = `${pk}/`;
    return this.request();
  }
  public async disconnect_user(pk: number): Promise<OcservUser> {
    this.method = "post";
    this.path = `${pk}/disconnect/`;
    return this.request();
  }
  public async sync_ocpasswd(params?: URLParams | null): Promise<SyncResponse> {
    this.method = "post";
    this.path = "sync/";
    if (params) {
      this.params = params;
    }
    return this.request();
  }
}

class OcservGroupApi extends Services {
  constructor() {
    super();
    this.baseUrl = "/groups/";
    this.path = "";
  }
  public async groups(params?: URLParams | null): Promise<GroupPagination> {
    this.method = "get";
    if (params) {
      this.params = params;
    }
    return this.request();
  }
  public async create_group(data: OcservGroup): Promise<OcservGroup> {
    this.path = "";
    this.method = "post";
    return this.request((data = data));
  }
  public async update_group(
    pk: number,
    data: OcservGroup
  ): Promise<OcservGroup> {
    this.method = "patch";
    this.path = `${pk}/`;
    return this.request((data = data));
  }
  public async delete_group(pk: number): Promise<OcservGroup> {
    this.method = "delete";
    this.path = `${pk}/`;
    return this.request();
  }
}

class OcctlServiceApi extends Services {
  constructor() {
    super();
    this.baseUrl = "/occtl/";
    this.method = "";
    this.path = "";
  }
  public async config(config: string, args: string): Promise<Occtl> {
    this.method = "get";
    this.path = `command/${config}/?args=${args}`;
    return this.request();
  }
  public async reload(): Promise<null> {
    this.method = "get";
    this.path = "reload/";
    return this.request();
  }
}

class StatsServiceApi extends Services {
  constructor() {
    super();
    this.baseUrl = "/stats/";
    this.method = "";
    this.path = "";
  }
  public async get_stats(): Promise<Stats> {
    this.method = "get";
    this.path = "";
    return this.request();
  }
}

class SystemServiceApi extends Services {
  constructor() {
    super();
    this.baseUrl = "/system/";
    this.method = "";
    this.path = "";
  }
  public async get_action_logs(): Promise<{ logs: string[] }> {
    this.method = "get";
    this.path = "action/list/";
    return this.request();
  }
  public async clear_action_logs(): Promise<null> {
    this.method = "delete";
    this.path = "action/clear/";
    return this.request();
  }
  public async ocserv_status(): Promise<{
    status: string[];
    dockerized?: boolean;
  }> {
    this.method = "get";
    this.path = "ocserv/status";
    return this.request();
  }
  public async ocserv_restart(): Promise<{
    status: string[];
    dockerized?: boolean;
  }> {
    this.method = "get";
    this.path = "ocserv/restart";
    return this.request();
  }
  public async journal(
    lines?: number,
    overlay: boolean = true
  ): Promise<{ logs: string[] }> {
    this.overlay = overlay;
    if (!Boolean(lines)) lines = 20;
    this.method = "get";
    this.path = `ocserv/journal/?lines=${lines}`;
    return this.request();
  }
}

const adminServiceApi = new AdminServiceApi();
const ocservUserApi = new OcservUserApi();
const ocservGroupApi = new OcservGroupApi();
const occtlServiceApi = new OcctlServiceApi();
const statsServiceApi = new StatsServiceApi();
const systemServiceApi = new SystemServiceApi();

export {
  adminServiceApi,
  ocservUserApi,
  ocservGroupApi,
  occtlServiceApi,
  statsServiceApi,
  systemServiceApi,
};
