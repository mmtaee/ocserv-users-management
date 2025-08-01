import axios, {
    type AxiosInstance,
    type AxiosRequestConfig,
    type AxiosResponse,
    type InternalAxiosRequestConfig
} from 'axios'
import {type SnackbarItem, useSnackbarStore} from "@/stores/snackbar.ts";
import {useLoadingStore} from "@/stores/loading.ts";


export const BaseUrl = import.meta.env.VITE_API_URL ?? 'http://localhost:8080'
export const BasePath = import.meta.env.VITE_API_PATH ?? '/api';

const config: AxiosRequestConfig = {baseURL: BaseUrl + BasePath}

const api: AxiosInstance = axios.create(config)


api.interceptors.request.use(
    (config: InternalAxiosRequestConfig): InternalAxiosRequestConfig => {
        const loadingStore = useLoadingStore()
        loadingStore.show()
        return config
    },
    (error) => {
        const loadingStore = useLoadingStore()
        loadingStore.hide()
        return Promise.reject(error)
    }
)

// Response interceptor: Handle 401/400
api.interceptors.response.use(
    (response: AxiosResponse): AxiosResponse => {
        const loadingStore = useLoadingStore()
        loadingStore.hide()
        return response
    },
    (error) => {
        const loadingStore = useLoadingStore()
        loadingStore.hide()
        const {response} = error
        if (response) {
            if (response.status === 401) {
                localStorage.removeItem('token')
                window.location.href = '/login'
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
