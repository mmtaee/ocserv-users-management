import axios, { AxiosInstance, AxiosRequestConfig, AxiosResponse } from 'axios';
import router from './router';
import { urlsPath } from '../utils/urls';
import { UrlsParams } from '@/utils/types';


let config: AxiosRequestConfig = {
    baseURL: process.env.NODE_ENV == "production" ? location.origin : "http://127.0.0.1:8000/api",
};

const _axios: AxiosInstance = axios.create(config);

_axios.interceptors.request.use(
    (config: AxiosRequestConfig) => {
        let token = window.localStorage.getItem('token');
        if (token) {
            config.headers['Authorization'] = `Token ${token}`;
        }
        return config;
    },
    (error) => {
        return Promise.reject(error);
    }
);

_axios.interceptors.response.use(
    (response: AxiosResponse) => {
        return response;
    },
    (error) => {
        if (error.response) {
            if (error.response.status == 500) {
                return Promise.reject(error);
            } else {
                if (error.response.status == 401) {
                    localStorage.removeItem('token');
                    router.push({ name: "Login" });
                }
                if (error.response.status == 405) {
                    console.log("method not allowed")
                    
                }
                if (error.response.status == 400) {
                    console.log("error 400: ", error.response.data)
                } 
                return error.response
            }

        }

    }
);


async function httpRequest(method: string, urlParams: UrlsParams, data?: object | null): Promise<AxiosResponse> {
    const AxiosMethods: { [K: string]: Function } = {
        get: _axios.get,
        post: _axios.post,
        patch: _axios.patch,
        put: _axios.put,
        delete: _axios.delete
    }
    const result: any = Object.entries(urlsPath[urlParams.urlName]).find(item => {
        return item[0] == urlParams.urlPath
    });
    let url: string = `${urlParams.urlName}${result[1]}`
    url = urlParams.pk ? url.replace(/pk/g, urlParams.pk.toString()) : url.replace("/pk/", "/")
    return await AxiosMethods[method](url = url, data = data, { params: urlParams.params });
}

export default httpRequest
