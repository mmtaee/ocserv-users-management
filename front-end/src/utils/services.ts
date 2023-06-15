import _axios from "@/plugins/axios";
import { AxiosResponse } from "axios";
import { AdminConfig, AminLogin, Config, Dashboard, UserPagination, GroupPagination, OcservUser, OcservGroup, Occtl } from "./types";

class Services {
    private status_code: number = 500

    private axiosMethods: { [K: string]: Function } = {
        get: _axios.get,
        post: _axios.post,
        patch: _axios.patch,
        put: _axios.put,
        delete: _axios.delete
    }
    public async request(method: string, url: string, data?: object): Promise<AxiosResponse> {
        let response: AxiosResponse = await this.axiosMethods[method](url = url, data = data);
        this.status_code = response.status
        return response
    }
    public status(): number {
        return this.status_code
    }
}

class AdminServiceApi extends Services {
    public async config(): Promise<Config> {
        let method: string = "get"
        let url = "/admin/config/"
        let response: AxiosResponse = await this.request(method, url)
        return response.data
    }
    public async login(data: AminLogin): Promise<{ token: string }> {
        let method: string = "post"
        let url = "/admin/login/"
        let response: AxiosResponse = await this.request(method, url, data)
        return response.data
    }
    public async logout(): Promise<void> {
        let method: string = "delete"
        let url = "/admin/logout/"
        await this.request(method, url)
    }
    public async create_configs(data: object): Promise<Config> {
        let method: string = "post"
        let url = "/admin/create/"
        let response: AxiosResponse = await this.request(method, url, data)
        return response.data
    }
    public async patch_configuration(data: object): Promise<null> {
        let method: string = "patch"
        let url = "/admin/configuration/"
        await this.request(method, url, data)
        return null
    }
    public async get_configuration(): Promise<AdminConfig> {
        let method: string = "get"
        let url = "/admin/configuration/"
        let response: AxiosResponse = await this.request(method, url)
        return response.data
    }
    public async dashboard(): Promise<Dashboard> {
        let method: string = "get"
        let url = "/admin/dashboard/"
        let response: AxiosResponse = await this.request(method, url)
        return response.data
    }
}

class OcservUserApi extends Services {
    public async users(): Promise<UserPagination> {
        let method: string = "get"
        let url = "/users/"
        let response: AxiosResponse = await this.request(method, url)
        return response.data
    }
    public async create_user(data: OcservUser): Promise<OcservUser> {
        let method: string = "post"
        let url = "/users/"
        let response: AxiosResponse = await this.request(method, url, data = data)
        return response.data
    }
    public async update_user(pk: number, data: OcservUser): Promise<OcservUser> {
        let method: string = "patch"
        let url = `/users/${pk}/`
        let response: AxiosResponse = await this.request(method, url, data = data)
        return response.data
    }
    public async delete_user(pk: number): Promise<OcservUser> {
        let method: string = "delete"
        let url = `/users/${pk}/`
        let response: AxiosResponse = await this.request(method, url)
        return response.data
    }
}

class OcservGroupApi extends Services {
    public async groups(): Promise<GroupPagination> {
        let method: string = "get"
        let url = "/groups/"
        let response: AxiosResponse = await this.request(method, url)
        return response.data
    }
    public async create_group(data: OcservGroup): Promise<OcservGroup> {
        let method: string = "post"
        let url = "/groups/"
        let response: AxiosResponse = await this.request(method, url, data = data)
        return response.data
    }
    public async update_group(pk: number, data: OcservGroup): Promise<OcservGroup> {
        let method: string = "patch"
        let url = `/groups/${pk}/`
        let response: AxiosResponse = await this.request(method, url, data = data)
        return response.data
    }
    public async delete_group(pk: number): Promise<OcservGroup> {
        let method: string = "delete"
        let url = `/groups/${pk}/`
        let response: AxiosResponse = await this.request(method, url)
        return response.data
    }
}

class OcctlServiceApi extends Services {
    public async config(config: string, args: string): Promise<Occtl> {
        let method: string = "get"
        let url = `/occtl/command/${config}/?args=${args}`
        let response: AxiosResponse = await this.request(method, url)
        return response.data
    }
    public async reload(): Promise<null> {
        let method: string = "get"
        let url = "/occtl/reload/"
        let response: AxiosResponse = await this.request(method, url)
        return response.data
    }
}

const adminServiceApi = new AdminServiceApi()
const ocservUserApi = new OcservUserApi()
const ocservGroupApi = new OcservGroupApi()
const occtlServiceApi = new OcctlServiceApi()
export {
    adminServiceApi,
    ocservUserApi,
    ocservGroupApi,
    occtlServiceApi,

}