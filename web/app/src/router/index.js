import {createRouter, createWebHistory} from 'vue-router'
import Home from '@/views/Home'
import EndpointDetails from "@/views/EndpointDetails";
import SuiteDetails from '@/views/SuiteDetails';
import AdminPanel from '@/views/AdminPanel';
import CertificateMonitor from '@/views/CertificateMonitor';

const routes = [
    {
        path: '/',
        name: 'Home',
        component: Home
    },
    {
        path: '/endpoints/:key',
        name: 'EndpointDetails',
        component: EndpointDetails,
    },
    {
        path: '/suites/:key',
        name: 'SuiteDetails',
        component: SuiteDetails
    },
    {
        path: '/admin',
        name: 'Admin',
        component: AdminPanel
    },
    {
        path: '/certificates',
        name: 'Certificates',
        component: CertificateMonitor
    }
];

const router = createRouter({
    history: createWebHistory(process.env.BASE_URL),
    routes
});

export default router;
