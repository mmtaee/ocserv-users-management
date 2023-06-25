import Vue from 'vue';
import VueRouter, { RouteConfig } from 'vue-router';
import Dashboard from "@/views/Dashboard.vue"

const originalPush = VueRouter.prototype.push;
VueRouter.prototype.push = function push(location: any) {
  return (originalPush as Function).call(this, location).catch((err: any) => err);
};

Vue.use(VueRouter);

const routes: RouteConfig[] = [
  {
    path: '/',
    name: 'Dashboard',
    component: Dashboard,
    meta: {
      title: "Dashboard",
      requireAuth: true
    }
  },
  {
    path: '/config',
    name: "Config",
    component: () => import("@/views/Configs.vue"),
    meta: {
      title: "Config",
      requireAuth: false
    }
  },
  {
    path: '/login',
    name: "Login",
    component: () => import("@/views/Login.vue"),
    meta: {
      title: "Login",
      requireAuth: false
    }
  },
  {
    path: '*',
    name: 'NotFound',
    component: () => import("@/views/NotFound.vue"),
    meta: {
      title: "404  Not Found)",
      requireAuth: false
    },
  },
  {
    path: '/configuration',
    name: "Configuration",
    component: () => import("@/views/Configuration.vue"),
    meta: {
      title: "Configuration",
      requireAuth: true
    }
  },
  {
    path: '/groups',
    name: "Groups",
    component: () => import("@/views/Groups.vue"),
    meta: {
      title: "Groups",
      requireAuth: true
    }
  },
  {
    path: '/users',
    name: "Users",
    component: () => import("@/views/Users.vue"),
    meta: {
      title: "Users",
      requireAuth: true
    }
  },
  {
    path: '/stats',
    name: "Stats",
    component: () => import("@/views/Stats.vue"),
    meta: {
      title: "Statistics",
      requireAuth: true
    }
  },
  {
    path: '/logs',
    name: "Logs",
    component: () => import("@/views/Logs.vue"),
    meta: {
      title: "Logs",
      requireAuth: true
    }
  },
  {
    path: '/occtl',
    name: "Occtl",
    component: () => import("@/views/Occtl.vue"),
    meta: {
      title: "Occtl",
      requireAuth: true
    }
  },
]

const router = new VueRouter({
  mode: 'history',
  base: process.env.BASE_URL,
  routes
})

router.beforeEach((to, from, next) => {
  document.title = to.meta?.title;
  if (!to.meta?.requireAuth) {
    next()
  }
  if (to.meta?.requireAuth) {
    if (!localStorage.getItem("token")) {
      next("/login")
    }
  }
  if (to.name == "Login" && localStorage.getItem("token")) {
    next("/")
  }
  next()
});

export default router
