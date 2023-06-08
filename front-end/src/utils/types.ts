interface URLS {
    admin: Object,
    users: Object
}

interface UrlsParams {
    urlName: keyof URLS,
    urlPath: string,
    pk?: string | number | null,
    params?: object | null
}

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

export {
    URLS,
    UrlsParams,
    AminLogin,
    AdminConfig,
    OcservConfigItems,
}
