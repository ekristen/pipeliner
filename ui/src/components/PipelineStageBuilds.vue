<template>
  <v-list dense>
    <v-list-item
      v-for="build in builds"
      :key="build.id"
      :to="{ name: 'BuildsView', params: { id: build.id } }"
    >
      <v-list-item-icon>
        <StateIcon :state="build.state" />
      </v-list-item-icon>
      <v-list-item-title>
        {{
        build.name
        }}
      </v-list-item-title>
      <v-list-item-icon v-if="build.retried">
        <v-icon>mdi-sync</v-icon>
      </v-list-item-icon>
    </v-list-item>
  </v-list>
</template>
<script>
import { mapState } from "vuex";

import StateIcon from "@/components/StateIcon";

export default {
  name: "PipelineStageBuilds",

  props: ["pipeline_id", "stage_idx"],

  components: {
    StateIcon,
  },

  computed: {
    ...mapState({
      builds: function (state) {
        const collator = new Intl.Collator("en", {
          numeric: true,
          sensitivity: "base",
        });
        return state.builds
          .filter(
            (b) =>
              b.pipeline_id === this.pipeline_id &&
              b.stage_idx === this.stage_idx
          )
          .sort((a, b) => {
            if (a.stage_idx < b.stage_idx) return -1;
            if (a.stage_idx > b.stage_idx) return 1;
            if (collator.compare(a.name, b.name) < 0) return -1;
            if (collator.compare(a.name, b.name) > 0) return 1;
            if (a.id < b.id) return 1;
            if (a.id > b.id) return -1;
            return 0;
          });
      },
    }),
  },
};
</script>