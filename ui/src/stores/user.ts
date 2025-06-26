import {defineStore} from "pinia";
import {SystemUsersApi} from "@/api";
import {getAuthorization} from "@/utils/request.ts";

interface UserState {
    uid: string;
    username: string;
    isAdmin: boolean;
}

export const useUserStore = defineStore('user', {
    state: (): UserState => ({
        uid: "",
        username: "",
        isAdmin: false,
    }),

    actions: {
        getProfile() {
            const api = new SystemUsersApi()
            api.systemUsersProfileGet(getAuthorization()).then((res) => {
                if (res.data) {
                    this.uid = res.data.uid;
                    this.username = res.data.username;
                    this.isAdmin = res.data.is_admin;
                }
            })
        },
        setUser(user: UserState) {
            this.uid = user.uid;
            this.username = user.username;
            this.isAdmin = user.isAdmin;
        },
        clearUser() {
            this.uid = "";
            this.username = "";
            this.isAdmin = false;
        },

    },
    getters: {
        user(state): UserState | null {
            return state;
        },
    }
})