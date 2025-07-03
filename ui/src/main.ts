import {createApp} from 'vue'
import App from './App.vue'
import vuetify from "@/plugins/vuetify.ts";
import i18n from "@/plugins/i18n.ts";
import router from "@/plugins/router.ts";
import {useConfigStore} from "@/stores/config.ts";
import {createPinia} from "pinia";
import {useUserStore} from "@/stores/user.ts";


const app = createApp(App)

app.use(createPinia())

;(async () => {
    const configStore = useConfigStore()
    const setup = await configStore.getConfig()

    app.use(vuetify)
    app.use(i18n)
    app.use(router)

    if (!setup) {
        await router.push({name: 'SetupPage'})
    } else {
        if (localStorage.getItem("token")) {
            const userStore = useUserStore()
            userStore.getProfile()
        }
    }

    app.mount('#app')

    const preloader = document.getElementById('preloader')
    if (preloader) {
        preloader.style.opacity = '0'
        setTimeout(() => preloader.remove(), 500)
    }
})()

