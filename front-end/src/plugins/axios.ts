import axios, { AxiosInstance, AxiosRequestConfig, AxiosResponse } from 'axios';
import router from './router';


let config: AxiosRequestConfig = {
    baseURL: process.env.NODE_ENV == "production" ? `${location.origin}/api` : "http://127.0.0.1:8000/api",
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
            if (error.response.status == 400) {
                return error.response
            }
            // if (error.response.status == 500) {
            // }
            // else {
            //     if (error.response.status == 401) {
            //         localStorage.removeItem('token');
            //         router.push({ name: "Login" });
            //     }
            //     if (error.response.status == 405) {
            //         console.log("method not allowed")

            //     }
            //     if (error.response.status == 400) {
            //         console.log("error 400: ", error.response.data)
            //     }
            //     return error.response
            // }
            return Promise.reject(error);
        }

    }
);

export default _axios
