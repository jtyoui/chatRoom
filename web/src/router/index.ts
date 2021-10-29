import {createRouter, createWebHistory} from "vue-router";
import Login from "../views/Login.vue";


const routes = [
    {
        path: '/login',
        name: 'Login',
        component: Login
    },
    {
        path: '/',
        name: 'Home',
        component: () => import('../views/Home.vue')
    }
]


const router = createRouter({
        history: createWebHistory(),
        routes: routes
    }
)
export default router
