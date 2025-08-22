import {defineStore} from "pinia";
import {type ModelsUser, SystemUsersApi} from "@/api";
import {getAuthorization} from "@/utils/request.ts";


export const useUserStore = defineStore('user', {
    state: (): ModelsUser => ({
        is_admin: false,
        uid: "",
        username: "",
        created_at: undefined,
        updated_at: undefined,
        last_login: "",
    }),

    actions: {
        async getProfile() {
            const api = new SystemUsersApi()
            try {
                const res = await api.systemUsersProfileGet(getAuthorization())
                if (res.data) {
                    this.setUser(res.data)
                }
            } catch (error) {
                console.error("Failed to fetch user profile", error)
                this.clearUser()
            }
        },

        setUser(user: ModelsUser) {
            Object.assign(this, user)
        },

        clearUser() {
            Object.assign(this, {
                _: 0,
                is_admin: false,
                uid: "",
                username: "",
                created_at: undefined,
                updated_at: undefined,
                last_login: ""
            })
        },
    },
    getters: {
        user(state): ModelsUser | null {
            return state;
        },
        isAdmin(state: ModelsUser) {
            return state.is_admin
        }
    }
})