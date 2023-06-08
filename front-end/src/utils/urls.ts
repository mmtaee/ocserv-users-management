import { URLS } from "@/utils/types"

const urlsPath: URLS = {
    admin: {
        config: "/config/",
        createConfig: "/create/",
        login: "/login/",
        configuration: "/configuration/",
        logout: "/logout/",
        dashboard: "/dashboard/",
    },       
    users: {
        users: "/pk/",
        disconnect: "/pk/disconnect/",
        syncOcpasswd: "/sync/",
    },
};

export {
    urlsPath,
}