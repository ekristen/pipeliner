<template>
  <v-row class="mt-4">
    <v-col>
      <v-card>
        <v-card-text>
          <div class="title font-weight-bold">Workflows</div>
          <v-divider class="my-2"></v-divider>
          <span class="text-h2">{{ stats.workflows }}</span>
        </v-card-text>
      </v-card>
    </v-col>
    <v-col>
      <v-card>
        <v-card-text>
          <div class="title font-weight-bold">Pipelines</div>
          <v-divider class="my-2"></v-divider>
          <span class="text-h2">{{ stats.pipelines }}</span>
        </v-card-text>
      </v-card>
    </v-col>
    <v-col>
      <v-card>
        <v-card-text>
          <div class="title font-weight-bold">Jobs</div>
          <v-divider class="my-2"></v-divider>
          <span class="text-h2">{{ stats.jobs }}</span>
        </v-card-text>
      </v-card>
    </v-col>
    <v-col>
      <v-card>
        <v-card-text>
          <div class="title font-weight-bold">Artifacts</div>
          <v-divider class="my-2"></v-divider>
          <span class="text-h2">{{ stats.artifacts }}</span>
        </v-card-text>
      </v-card>
    </v-col>
    <v-col>
      <v-card>
        <v-card-text>
          <div class="title font-weight-bold">Runners</div>
          <v-divider class="my-2"></v-divider>
          <span class="text-h2">{{ stats.runners }}</span>
        </v-card-text>
      </v-card>
    </v-col>
  </v-row>
</template>

<script>
import { mapState } from "vuex";

export default {
  name: "Dashboard",

  created() {
    this.$store.dispatch("getStats");
  },

  data: () => ({
    alert: true,
  }),

  computed: {
    ...mapState({
      stats: (state) => state.stats,
    }),
  },
};
</script>
