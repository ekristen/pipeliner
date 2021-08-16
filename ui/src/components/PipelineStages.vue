<template>
  <v-breadcrumbs :items="stages" class="pipeline-stages pa-0">
    <template v-slot:item="{ item }">
      <v-breadcrumbs-item>
        <v-menu offset-y auto>
          <template v-slot:activator="{ on, attrs }">
            <v-btn
              x-small
              outlined
              fab
              :color="color(item.state)"
              v-bind="attrs"
              v-on="on"
              @click="getJobs(item.pipeline_id, item.id)"
            >
              <v-icon>{{ icon(item.state) }}</v-icon>
            </v-btn>
          </template>
          <PipelineStageBuilds :pipeline_id="item.pipeline_id" :stage_idx="item.index" />
        </v-menu>
      </v-breadcrumbs-item>
    </template>
    <template v-slot:divider>
      <v-icon>mdi-minus</v-icon>
    </template>
  </v-breadcrumbs>
</template>

<style>
li.v-breadcrumbs__divider {
  padding: 0px !important;
}
</style>

<script>
import { mapState } from "vuex";

import PipelineStageBuilds from "@/components/PipelineStageBuilds";

import State from "@/mixins/state";

export default {
  name: "PipelineStages",

  props: ["id"],

  mixins: [State],

  created() {
    const pipelineId = this.id;
    this.$store.dispatch("getPipelineStages", { id: pipelineId });
  },

  components: {
    PipelineStageBuilds,
  },

  computed: mapState({
    stages(state) {
      const pipelineId = this.id;
      return state.stages
        .filter((s) => s.pipeline_id == pipelineId)
        .sort((a, b) => {
          if (a.index < b.index) {
            return -1;
          }
          if (a.index > b.index) {
            return 1;
          }
          return 0;
        });
    },
    builds(state) {
      const pipelineId = this.id;
      return state.builds
        .filter((b) => b.pipeline_id == pipelineId)
        .sort((a, b) => {
          if (a.id < b.id) {
            return -1;
          }
          if (a.id > b.id) {
            return 1;
          }
          return 0;
        });
    },
  }),

  data() {
    return {
      stagesLoaded: false,
      buildsLoaded: {},
    };
  },

  methods: {
    getJobs(pipeline_id, stage_id) {
      if (
        typeof this.buildsLoaded[`${pipeline_id}_${stage_id}`] === "undefined"
      ) {
        this.buildsLoaded[`${pipeline_id}_${stage_id}`] = true;
        this.$store.dispatch("getPipelineStageBuilds", {
          pipeline_id,
          stage_id,
        });
      }
    },
  },

  getters: {
    getStageBuilds: (state) => (index) => {
      return state.builds.filter((b) => b.stage_index == index);
    },
  },
};
</script>