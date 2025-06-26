import {defineStore} from "pinia";
import {SystemApi} from "@/api";


interface ConfigState {
    setup: boolean
    googleCaptchaSiteKey: string
}


export const useConfigStore = defineStore('config', {
    state: (): ConfigState => ({
        setup: false,
        googleCaptchaSiteKey: "",
    }),

    actions: {
        async getConfig() {
            const api = new SystemApi()
            await api.systemInitGet().then((res) => {
                if (res.data) {
                    this.googleCaptchaSiteKey = res.data.google_captcha_site_key || ""
                    this.setup = true
                }
            })
            return this.setup
        },
        setConfig(googleCaptchaSiteKey: string | undefined) {
            if (googleCaptchaSiteKey) {
                this.googleCaptchaSiteKey = googleCaptchaSiteKey
            }
            this.setup = true
        }
    },
    getters: {
        config(state): ConfigState {
            return {
                setup: state.setup,
                googleCaptchaSiteKey: state.googleCaptchaSiteKey,
            }
        },
    },
})