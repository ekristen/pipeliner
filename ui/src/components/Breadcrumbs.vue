<template>
  <v-breadcrumbs v-if="breadcrumbs" :items="breadcrumbs"></v-breadcrumbs>
</template>
<script>
import { mapState } from "vuex";

export default {
  name: "Breadcrumbs",

  computed: {
    ...mapState({
      breadcrumbs: (state) => state.breadcumbs,
    }),
  },
};
</script>