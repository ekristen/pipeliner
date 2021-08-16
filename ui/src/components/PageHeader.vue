<template>
  <v-row v-if="show">
    <v-col class="pa-0">
      <v-card elevation="0" flat tile :height="height">
        <v-toolbar height="48px">
          <v-toolbar-title class="pt-3" v-if="title">
            <span>{{ title }}</span>
          </v-toolbar-title>
          <template v-slot:extension height="40px" v-if="crumbs.length > 0">
            <v-breadcrumbs :items="crumbs" class="pl-1"></v-breadcrumbs>
          </template>
        </v-toolbar>
      </v-card>
    </v-col>
  </v-row>
</template>
<script>
export default {
  name: "PageHeader",

  props: ["title"],

  computed: {
    crumbs() {
      return [];
    },
    height() {
      let total = 0;
      if (this.title) {
        total += 40;
      }
      if (this.crumbs.length > 0) {
        total += 48;
      }

      return `${total}px`;
    },
    show() {
      if (this.title || this.crumbs.length > 0) {
        return true;
      }
      return false;
    },
  },
};
</script>