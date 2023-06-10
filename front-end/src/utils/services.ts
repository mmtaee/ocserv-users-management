import _axios from "@/plugins/axios";
import { AxiosResponse } from "axios";
import { AdminConfig, AminLogin, Config, Dashboard } from "./types";

class Services {
    private status_code: number = 500

    private axiosMethods: { [K: string]: Function } = {
        get: _axios.get,
        post: _axios.post,
        patch: _axios.patch,
        put: _axios.put,
        delete: _axios.delete
    }
    public async request(method: string, url: string, data?: object, params?: object): Promise<AxiosResponse> {
        let response: AxiosResponse = await this.axiosMethods[method](url = url, data = data, { params: params });
        console.log("response : ", response)
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


class OcservUserApi extends Services { }

const adminServiceApi = new AdminServiceApi()
const ocservUserApi = new OcservUserApi()

export {
    adminServiceApi,
    ocservUserApi,
}