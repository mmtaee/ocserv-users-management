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
        meta: {title: "Home"}
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
    {
        path: '/ocserv-groups',
        name: 'OcservGroupsPage',
        component: () => import('../views/OcservGroupsViews.vue'),
        meta: {title: "Ocserv Groups"}
    },
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
        localStorage.removeItem("token")
        const configStore = useConfigStore()
        if (configStore.config.setup) {
            next("/")
            return;
        }
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