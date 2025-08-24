import {defineStore} from "pinia";
import {OCCTLApi} from "@/api";


interface ServerState {
    Version: string
    OcctlVersion: string
}


export const useServerStore = defineStore('server', {
    state: (): ServerState => ({
        Version: "",
        OcctlVersion: ""
    }),
    actions: {
        async getServerInfo() {
            const api = new OCCTLApi()
            await api.occtlServerInfoGet().then((res) => {
                if (res.data) {
                    this.Version = res.data.version || ""
                    this.OcctlVersion = (res.data.occtl_version || "").replace(/\n/g, '<br />')
                }
            })
        }
    },
    getters: {
        versionInfo: (state) => state.Version,
        occtlVersionInfo: (state) => state.OcctlVersion
    }
})