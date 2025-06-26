import axios, {
    type AxiosInstance,
    type AxiosRequestConfig,
    type AxiosResponse,
    type InternalAxiosRequestConfig
} from 'axios'
import {type SnackbarItem, useSnackbarStore} from "@/stores/snackbar.ts";


export const BaseUrl = import.meta.env.VITE_API_URL ?? 'http://localhost:8080'
export const BasePath = import.meta.env.VITE_API_PATH ?? '/api';

const config: AxiosRequestConfig = {baseURL: BaseUrl + BasePath}

const api: AxiosInstance = axios.create(config)

api.interceptors.request.use(
    (config: InternalAxiosRequestConfig): InternalAxiosRequestConfig => {
        // config.headers = config.headers ?? {} as AxiosRequestHeaders
        // const token = localStorage.getItem('token')
        // if (token) {
        //     config.headers["Authorization"] = `Bearer ${token}`
        // }
        return config
    },
    (error) => {
        return Promise.reject(error)
    }
)

// Response interceptor: Handle 401/400
api.interceptors.response.use(
    (response: AxiosResponse): AxiosResponse => {
        return response
    },
    (error) => {
        const {response} = error
        if (response) {
            console.log("response", response)
            if (response.status === 401) {
                localStorage.removeItem('token')
                return Promise.resolve(response)
            }
            if (response.status === 400) {
                const snackbar = useSnackbarStore()
                let errorList = response.data.error
                if (Array.isArray(errorList)) {
                    const messages: SnackbarItem[] = errorList.map((message: string, index: number) => ({
                        id: index + 1,
                        message,
                        color: 'error',
                        timeout: (index + 1) * 1000 + 4000,
                    }))
                    snackbar.show(messages)
                } else {
                    snackbar.show({
                        id: 1,
                        message: typeof errorList === 'string' ? errorList : 'Unknown error',
                        color: 'error',
                        timeout: 4000,
                    })
                }
            }
        }
        return Promise.reject(error)
    }
)

export default api
