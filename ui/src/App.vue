<template>
  <v-app>
    <Navigation />

    <v-app-bar app clipped-left clipped-right>
      <v-app-bar-nav-icon class="mr-1"></v-app-bar-nav-icon>Pipeliner
    </v-app-bar>

    <!-- Sizes your content based upon application components -->
    <v-main>
      <!-- Provides the application the proper gutter -->
      <v-container fluid class="mt-0 pt-0">
        <!-- If using vue-router -->
        <router-view></router-view>
      </v-container>
    </v-main>

    <v-footer app padless>
      <v-col class="text-right" cols="12">
        v
        <strong>{{ version }}</strong>
      </v-col>
    </v-footer>

    <Disconnected />
  </v-app>
</template>

<script>
//import Breadcrumbs from "@/components/Breadcrumbs";
import Disconnected from "@/components/Disconnected";
import Navigation from "@/components/Navigation";

import { mapState } from "vuex";

export default {
  name: "App",
  data: () => ({}),
  components: {
    //Breadcrumbs,
    Disconnected,
    Navigation,
  },
  created() {
    this.$store.dispatch("getVersion");
  },
  computed: {
    ...mapState({
      version: (state) => state.version,
    }),
  },
};
</script>
