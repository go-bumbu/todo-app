import { createApp } from 'vue'
import App from './App.vue'
const app = createApp(App)

import CustomTheme from '@/theme.js'
import 'primeflex/primeflex.css'
import 'primeicons/primeicons.css'

import '@/assets/style.scss'

import PrimeVue from 'primevue/config'

app.use(PrimeVue, {
    // Default theme configuration
    theme: {
        preset: CustomTheme,
        options: {
            prefix: 'c',
            darkModeSelector: 'system',
            cssLayer: false
        }
    }
})

// pinia store
import { createPinia } from 'pinia'
app.use(createPinia())

// // initialize toast service
// import ToastService from 'primevue/toastservice';
// app.use(ToastService);

// add the app router
import router from './router'
app.use(router)

// focus trap
import FocusTrap from 'primevue/focustrap'
app.directive('focustrap', FocusTrap)

app.mount('#app')
