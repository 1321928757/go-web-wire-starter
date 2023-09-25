import Vue from 'vue'
import App from './App.vue'

// Element
import ElementUI from 'element-ui'
import 'element-ui/lib/theme-chalk/index.css'
import axios from 'axios'
import lodash from 'lodash'

Vue.config.productionTip = false

Vue.prototype.$axios= axios
Vue.prototype.$lodash= lodash

// Element
Vue.use(ElementUI)

new Vue({
  render: h => h(App),
}).$mount('#app')
