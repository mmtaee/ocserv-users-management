import type {RouteRecordRaw} from 'vue-router'
import {createRouter, createWebHistory} from 'vue-router'
import HomeView from "@/views/HomeView.vue";
import {useConfigStore} from "@/stores/config.ts";
import {useIsSmallDisplay} from "@/stores/display.ts";


const routes: Array<RouteRecordRaw> = [
    {
        path: '/mobile-not-allowed',
        name: 'MobileNotAllowed',
        component: () => import('@/views/MobileNotAllowedPage.vue'),
    },
    {
        path: '/',
        name: 'HomePage',
        component: HomeView,
        meta: {
            title: "Home",
        }
    },
    {
        path: '/setup',
        name: 'SetupPage',
        component: () => import('../views/SetupView.vue'),
        meta: {
            title: "Setup",
            desktopOnly: true
        }
    },
    {
        path: '/login',
        name: 'LoginPage',
        component: () => import('../views/LoginView.vue'),
        meta: {
            title: "Login",
        }
    },
    {
        path: '/account',
        name: 'AccountPage',
        component: () => import('../views/AccountView.vue'),
        meta: {
            title: "Account",
            desktopOnly: true
        }
    },
    // {
    //     path: '/config',
    //     name: 'ConfigPage',
    //     component: () => import('../views/ConfigView.vue'),
    // },
    // {
    //     path: '/change_password',
    //     name: 'ChangePasswordPage',
    //     component: () => import('../views/ChangePasswordView.vue'),
    // },
    // {
    //     path: '/staffs',
    //     name: 'StaffsPage',
    //     component: () => import('../views/StaffView.vue'),
    // },
    // {
    //     path: '/oc_user',
    //     name: 'OcservUserPage',
    //     component: () => import('../views/OcservUserView.vue'),
    // },
    // {
    //     path: '/error',
    //     name: 'ErrorPage',
    //     component: () => import('../views/ErrorView.vue'),
    // },
]

const router = createRouter({
    history: createWebHistory(import.meta.env.BASE_URL),
    routes: routes,
})


router.beforeEach((to, _from, next) => {
    let token = localStorage.getItem('token') || null

    if (to.meta?.title) {
        document.title = to.meta.title as string
    } else {
        document.title = to.name as string
    }


    if (to.path === '/setup') {
        if (token !== null) {
            next("/")
            return;
        }
        const configStore = useConfigStore()
        if (configStore.config.setup) {
            next("/")
            return;
        }
        next()
        return;
    }

    if (to.path === '/login') {
        if (token !== null) {
            next("/")
            return;
        }
    }

    if (token === null && to.path !== '/login') {
        next("/login")
        return;
    }

    const smallDisplay = useIsSmallDisplay()
    if (to.meta.desktopOnly && smallDisplay.isSmallDisplay) {
        next("/mobile-not-allowed");
        return;
    }

    next()
})


export default router