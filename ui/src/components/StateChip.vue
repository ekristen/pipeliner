<template>
  <v-chip v-bind="$attrs" :color="color(state)" label outlined :to="to">
    <v-icon left>{{ icon(state) }}</v-icon>
    {{ state }}
  </v-chip>
</template>

<script>
import State from "@/mixins/state";

export default {
  name: "StateChip",

  props: ["state", "to"],

  mixins: [State],
};
</script>