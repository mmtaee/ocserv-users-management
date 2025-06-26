import {defineStore} from "pinia";
import {SystemUsersApi} from "@/api";
import {getAuthorization} from "@/utils/request.ts";
import {formatDateTime} from "@/utils/convertors.ts";

interface UserState {
    uid: string;
    username: string;
    isAdmin: boolean;
    lastLogin: string;
}

export const useUserStore = defineStore('user', {
    state: (): UserState => ({
        uid: "",
        username: "",
        isAdmin: false,
        lastLogin: "",
    }),

    actions: {
        getProfile() {
            const api = new SystemUsersApi()
            api.systemUsersProfileGet(getAuthorization()).then((res) => {
                if (res.data) {
                    this.uid = res.data.uid;
                    this.username = res.data.username;
                    this.isAdmin = res.data.is_admin;
                    this.lastLogin = formatDateTime(res.data.last_login, "");
                }
            })
        },
        setUser(user: UserState) {
            this.uid = user.uid;
            this.username = user.username;
            this.isAdmin = user.isAdmin;
            this.lastLogin = user.lastLogin;
        }
    }
})