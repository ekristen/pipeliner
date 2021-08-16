import Vue from 'vue'
import App from './App.vue'
import router from './router'
import store from './store'
import vuetify from './plugins/vuetify'

import * as JSONBig from 'json-bigint'
import VueNativeSock from 'vue-native-websocket'
import VueMoment from 'vue-moment'
import VueHighlightJS from 'vue-highlightjs'

import 'highlight.js/styles/solarized-dark.css'

const jsonBig = JSONBig({ storeAsString: true })

Vue.config.productionTip = false

const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
Vue.use(VueNativeSock, `${protocol}${window.location.host}/v1/ws`, {
  store,
  format: 'json',
  reconnection: true, // (Boolean) whether to reconnect automatically (false)
  passToStoreHandler: function (eventName, event) {
    if (!eventName.startsWith('SOCKET_')) { return }
    let method = 'commit'
    let target = eventName.toUpperCase()
    let msg = event
    if (this.format === 'json' && event.data) {
      msg = jsonBig.parse(event.data)
      if (msg.mutation) {
        target = [msg.namespace || '', msg.mutation].filter((e) => !!e).join('/')
      } else if (msg.action) {
        method = 'dispatch'
        target = [msg.namespace || '', msg.action].filter((e) => !!e).join('/')
      }
    }
    this.store[method](target, msg)
  }
})

Vue.use(VueHighlightJS);
Vue.use(VueMoment);

new Vue({
  router,
  store,
  vuetify,
  render: h => h(App)
}).$mount('#app')
