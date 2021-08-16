<template>
  <v-sheet elevation="2" class="mt-4" v-if="runner">
    <v-col>
      <h3>{{ runner.name }}</h3>
      <v-divider></v-divider>
    </v-col>
  </v-sheet>
</template>

<script>
import { mapState } from "vuex";

export default {
  name: "RunnersView",

  created() {
    this.$store.dispatch("getRunner", { id: this.$route.params.id });
  },

  computed: {
    ...mapState({
      runner(state) {
        return state.runners.find((d) => d.id === this.$route.params.id);
      },
    }),
  },

  methods: {
    id(runner) {
      if (runner.id === runner.id_str) {
        return runner.id;
      } else {
        return runner.id_str;
      }
    },

    state(runner) {
      if (!runner.active) {
        return "inactive";
      }

      const contactedAt = new Date(runner.contacted_at);

      const now = new Date();

      const diff = now - contactedAt;

      if (diff > 60 * 15) {
        return "offline";
      }

      return "active";
    },
    fgcolor(state) {
      if (state == "success") {
        return "white";
      }

      return "black";
    },
    bgcolor(state) {
      if (state == "active") {
        return "green";
      } else if (state == "offline") {
        return "grey";
      }

      return "pink";
    },
    icon(state) {
      if (state == "success") {
        return "mdi-checkbox-marked-circle";
      } else if (state == "skipped") {
        return "mdi-debug-step-over";
      } else if (state == "failed") {
        return "mdi-minus-circle-outline";
      }

      return "mdi-help-circle-outline";
    },
  },
};
</script>
