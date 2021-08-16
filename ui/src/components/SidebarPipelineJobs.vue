<template>
  <div>
    <v-divider></v-divider>
    <v-list>
      <v-subheader class="font-weight-black">
        Pipeline
        <router-link
          class="ml-2"
          :to="{ name: 'PipelinesView', params: { id: this.pipeline_id } }"
        >#{{ this.pipeline_id }}</router-link>
      </v-subheader>
      <v-list-item dense>
        <v-list-item-content>
          <v-select
            label="Stage"
            :items="stages"
            item-value="index"
            item-text="name"
            :value="this.selected_stage"
            @change="stageChanged"
          ></v-select>
        </v-list-item-content>
      </v-list-item>
    </v-list>

    <v-divider></v-divider>
    <v-list>
      <v-list-item dense v-for="build in builds" :key="build.id">
        <v-list-item-icon class="mr-1">
          <StateIcon :state="build.state" xsmall />
        </v-list-item-icon>
        <v-list-item-content>
          <router-link :to="{ name: 'BuildsView', params: { id: build.id } }">{{ build.name }}</router-link>
        </v-list-item-content>
        <v-list-item-icon>
          <v-tooltip left>
            <template v-slot:activator="{ on, attrs }">
              <v-icon v-bind="attrs" v-on="on">{{ jobIcon(build) }}</v-icon>
            </template>
            <span>{{ jobIconText(build) }}</span>
          </v-tooltip>
        </v-list-item-icon>
      </v-list-item>
    </v-list>
  </div>
</template>
<script>
import StateIcon from "@/components/StateIcon";

import { mapState } from "vuex";

export default {
  name: "SidebarPipelineJobs",

  props: ["build"],

  components: {
    StateIcon,
  },

  data() {
    return {
      build_id: 0,
      selected_stage: 0,
      pipeline_id: 0,
    };
  },

  computed: {
    ...mapState({
      stages(state) {
        return state.stages.filter((s) => s.pipeline_id === this.pipeline_id);
      },
      builds(state) {
        return state.builds
          .filter((b) => b.pipeline_id === this.pipeline_id)
          .filter((b) => b.stage_idx === this.selected_stage)
          .sort((a, b) => {
            if (a.stage_idx < b.stage_idx) return -1;
            if (a.stage_idx > b.stage_idx) return 1;
            if (a.name < b.name) return -1;
            if (a.name > b.name) return 1;
            if (a.id < b.id) return 1;
            if (a.id > b.id) return -1;
            return 0;
          });
      },
    }),
  },

  watch: {
    build(build) {
      this.build_id = build.id;
    },
  },

  methods: {
    stageChanged(value) {
      this.selected_stage = value;
      this.$store.dispatch("getPipelineBuilds", {
        id: this.pipeline_id,
      });
    },
    isCurrentBuild(build) {
      return build.id === this.build_id;
    },
    jobIcon(build) {
      if (build.id === this.build_id) {
        return "mdi-arrow-left";
      } else if (build.retried) {
        return "mdi-sync";
      } else {
        return "";
      }
    },
    jobIconText(build) {
      if (build.id === this.build_id) {
        return "Current Viewed Job";
      } else if (build.retried) {
        return "Retried Job";
      } else {
        return "";
      }
    },
  },

  created() {
    this.selected_stage = this.build.stage_idx;
    this.pipeline_id = this.build.pipeline_id;
    this.build_id = this.build.id;

    this.$store.dispatch("getPipelineStages", { id: this.build.pipeline_id });
    this.$store.dispatch("getPipelineBuilds", {
      id: this.pipeline_id,
    });
  },
};
</script>