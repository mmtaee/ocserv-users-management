declare interface AminLogin {
    username: string | null;
    password: string | null;
}

declare interface AdminConfig {
    username: string | null;
    password: string | null;
    new_password: string | null;
    captcha_site_key: string | null;
    captcha_secret_key: string | null;
    default_configs: object | null;
    default_traffic: number | null
}

declare type OcservConfigItems = {
    label: string;
    model: string;
};

declare interface Config {
    captcha_site_key: string | null,
    config: boolean,
    token?: string
}

declare interface Dashboard {
    online_users: Array<string>,
    server_stats: object,
}



export {
    AminLogin,
    AdminConfig,
    OcservConfigItems,
    Config,
    Dashboard,
}
