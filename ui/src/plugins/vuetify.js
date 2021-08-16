import Vue from 'vue';
import Vuetify from 'vuetify/lib';
import colors from 'vuetify/lib/util/colors';

Vue.use(Vuetify);

export default new Vuetify({
  theme: {
    themes: {
      dark: {
        primary: '#21CFF3',
        accent: '#FF4081',
        secondary: '#ffe18d',
        success: colors.green.accent2,
        info: '#2196F3',
        warning: '#FB8C00',
        error: '#FF5252',
        running: '#21CFF3',
        tag: '#00227b',
      },
      light: {
        primary: '#1976D2',
        accent: '#e91e63',
        secondary: '#30b1dc',
        success: colors.green.darken1,
        info: '#2196F3',
        warning: '#FB8C00',
        error: '#FF5252',
        running: colors.blue.darken1,
        tag: '#6f74dd',
      },
    },
  },
});
