<template>
  <v-icon :color="color(state)" left>{{ icon(state) }}</v-icon>
</template>

<script>
import State from "@/mixins/state";

export default {
  name: "StateIcon",

  props: ["state"],

  mixins: [State],
};
</script>