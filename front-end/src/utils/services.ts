import _axios from "@/plugins/axios";
import { AxiosResponse } from "axios";
import { AdminConfig, AminLogin, Config, Dashboard, UserPagination, GroupPagination, OcservUser, OcservGroup, Occtl } from "./types";

class Services {
    private status_code: number = 500
    private responseData: any
    public method: string = "get"
    public baseUrl: string = ""
    public path: string = ""

    private axiosMethods: { [K: string]: Function } = {
        get: _axios.get,
        post: _axios.post,
        patch: _axios.patch,
        put: _axios.put,
        delete: _axios.delete
    }
    public async request(data?: object): Promise<any> {
        let url: string = this.baseUrl + this.path
        let response: AxiosResponse = await this.axiosMethods[this.method](url = url, data = data);
        this.status_code = response.status
        this.responseData = response.data
        return this.data()
    }
    public data(): any {
        return this.responseData
    }
    public status(): number {
        return this.status_code
    }
}

class AdminServiceApi extends Services {
    constructor() {
        super()
        this.baseUrl = "/admin/"
        this.path = ""
    }
    public async config(): Promise<Config> {
        this.method = "get"
        this.path = "config/"
        return this.request()
    }
    public async login(data: AminLogin): Promise<{ token: string }> {
        this.method = "post"
        this.path = "login/"
        return this.request(data)
    }
    public async logout(): Promise<void> {
        this.method = "delete"
        this.path = "logout/"
        await this.request()
    }
    public async create_configs(data: object): Promise<Config> {
        this.method = "post"
        this.path = "create/"
        return await this.request(data)
    }
    public async patch_configuration(data: object): Promise<null> {
        this.method = "patch"
        this.path = "configuration/"
        await this.request(data)
        return null
    }
    public async get_configuration(): Promise<AdminConfig> {
        this.method = "get"
        this.path = "configuration/"
        return await this.request()
    }

    public async dashboard(): Promise<Dashboard> {
        this.method = "get"
        this.path = "dashboard/"
        return this.request()
    }
}

class OcservUserApi extends Services {
    constructor() {
        super()
        this.baseUrl = "/users/"
        this.path = ""
    }
    public async users(): Promise<UserPagination> {
        this.method = "get"
        return this.request()
    }
    public async create_user(data: OcservUser): Promise<OcservUser> {
        this.method = "post"
        return this.request(data = data)
    }
    public async update_user(pk: number, data: OcservUser): Promise<OcservUser> {
        this.method = "patch"
        this.path = `${pk}/`
        return this.request(data = data)
    }
    public async delete_user(pk: number): Promise<OcservUser> {
        this.method = "delete"
        this.path = `${pk}/`
        return this.request()
    }
    public async disconnect_user(pk: number): Promise<OcservUser> {
        this.method = "post"
        this.path = `${pk}/disconnect/`
        return this.request()
    }
}

class OcservGroupApi extends Services {
    constructor() {
        super()
        this.baseUrl = "/groups/"
        this.path = ""
    }
    public async groups(args?: string | null): Promise<GroupPagination> {
        this.method = "get"
        this.path = args ? `?args=${args}` : ""
        return this.request()
    }
    public async create_group(data: OcservGroup): Promise<OcservGroup> {
        this.method = "post"
        return this.request(data = data)
    }
    public async update_group(pk: number, data: OcservGroup): Promise<OcservGroup> {
        this.method = "patch"
        this.path = `${pk}/`
        return this.request(data = data)
    }
    public async delete_group(pk: number): Promise<OcservGroup> {
        this.method = "delete"
        this.path = `${pk}/`
        return this.request()
    }
}

class OcctlServiceApi extends Services {
    constructor() {
        super()
        this.baseUrl = "/occtl/"
        this.method = ""
        this.path = ""
    }
    public async config(config: string, args: string): Promise<Occtl> {
        this.method = "get"
        this.path = `command/${config}/?args=${args}`
        return this.request()
    }
    public async reload(): Promise<null> {
        this.method = "get"
        this.path = "reload/"
        return this.request()
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