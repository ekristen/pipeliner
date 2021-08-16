<template>
  <v-navigation-drawer v-bind="$attrs" app clipped permanent :mini-variant="mini" class="pt-2">
    <v-list dense nav>
      <v-list-item v-for="item in items" :key="item.title" :to="item.to" link>
        <v-list-item-icon>
          <v-icon>{{ item.icon }}</v-icon>
        </v-list-item-icon>

        <v-list-item-content>
          <v-list-item-title>{{ item.title }}</v-list-item-title>
        </v-list-item-content>
      </v-list-item>
    </v-list>

    <template v-slot:append>
      <v-list-item @click="alwaysMini = !alwaysMini" link>
        <v-list-item-icon>
          <v-icon>{{collapseIcon}}</v-icon>
        </v-list-item-icon>

        <v-list-item-content>
          <v-list-item-title>Collapse</v-list-item-title>
        </v-list-item-content>
      </v-list-item>
    </template>
  </v-navigation-drawer>
</template>

<script>
export default {
  name: "Navigation",

  data() {
    return {
      alwaysMini: false,
      items: [
        {
          title: "Dashboard",
          icon: "mdi-view-dashboard-outline",
          to: "/",
        },
        {
          title: "Pipelines",
          icon: "mdi-rocket-outline",
          to: "/pipelines",
        },
        { title: "Jobs", icon: "mdi-rocket-launch-outline", to: "/builds" },
        { title: "Workflows", icon: "mdi-puzzle-outline", to: "/workflows" },
        { title: "Runners", icon: "mdi-network-outline", to: "/runners" },
        { title: "Variables", icon: "mdi-variable", to: "/variables" },
        { title: "Register Tokens", icon: "mdi-signal", to: "/tokens" },
      ],
      right: null,
    };
  },

  computed: {
    mini() {
      return this.alwaysMini || this.$vuetify.breakpoint.mdAndDown;
    },
    collapseIcon() {
      if (this.mini) {
        return "mdi-arrow-expand-right";
      } else {
        return "mdi-arrow-collapse-left";
      }
    },
  },
};
</script>
