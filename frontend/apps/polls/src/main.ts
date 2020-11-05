import 'primeicons/primeicons.css'
import Button from 'primevue/button'
import DataView from 'primevue/dataview'
import InputText from 'primevue/inputtext'
import 'primevue/resources/primevue.min.css'
import 'primevue/resources/themes/saga-green/theme.css'
import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import store from './store'

const app = createApp(App)

app.component('InputText', InputText)
app.component('Button', Button)
app.component('DataView', DataView)

app.use(store)
    .use(router)
    .mount('#app')
