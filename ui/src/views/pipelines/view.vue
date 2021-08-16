<template>
  <v-sheet elevation="2" class="mt-4" v-if="pipeline">
    <div>
      <StateChip :state="pipeline.state" class="ma-2 ml-3" />
      <span class="text-h6">
        Pipeline #{{
        pipeline.id
        }}
        triggered {{ pipeline.created_at }}
      </span>
      <div class="float-right mt-2">
        <v-btn
          color="warning"
          small
          class="mr-2"
          v-if="isCancelable(pipeline.state)"
          @click="cancel(pipeline)"
        >Cancel running</v-btn>
        <v-btn
          outlined
          color="error"
          small
          class="mr-2"
          v-show="!isCancelable(pipeline.state)"
        >Delete</v-btn>
      </div>
    </div>

    <v-divider></v-divider>
    <v-data-table
      :headers="headers"
      :items="builds"
      group-by="stage"
      disable-sort
      disable-pagination
      hide-default-footer
      :loading="loading"
    >
      <template v-slot:item.state="{item}">
        <StateChip :state="item.state" :to="{ name: 'BuildsView', params: { id: item.id } }" />
      </template>
      <template v-slot:item.id="{item}">
        <router-link
          :to="{ name: 'BuildsView', params: { id: item.id } }"
          class="mr-3"
        >#{{ item.id }}</router-link>
        <JobTags :job="item" x-small />
        <v-tooltip right v-if="item.retried" small>
          <template v-slot:activator="{ on, attrs }">
            <v-icon v-bind="attrs" v-on="on">mdi-sync</v-icon>
          </template>
          <span>This job was retried</span>
        </v-tooltip>
      </template>
      <template v-slot:item.duration="{item}">{{ duration(item) }}</template>
      <template v-slot:item.actions="{item}">
        <DownloadArtifact :build="item" type="archive" v-if="hasArtifacts(item, artifacts)" />

        <v-btn
          class="mx-2"
          fab
          x-small
          color="info darken-2"
          @click="retryJob(item)"
          v-if="retryableState(item)"
        >
          <v-icon>mdi-replay</v-icon>
        </v-btn>
        <v-btn class="mx-2" fab x-small color="error darken-2" v-if="isCancelable(item.state)">
          <v-icon>mdi-stop</v-icon>
        </v-btn>
        <v-btn
          class="mx-2"
          fab
          x-small
          color="success darken-2"
          v-if="isManualJob(item)"
          v-on:click="runJob(item)"
        >
          <v-icon>mdi-play</v-icon>
        </v-btn>
      </template>
      <template v-slot:group.header="{group}">
        <td colspan="5">
          <PipelineBuildsGroup :name="group" :pipeline="pipeline" />
        </td>
      </template>
    </v-data-table>
  </v-sheet>
</template>

<style lang="scss" scoped>
tbody {
  tr:hover {
    background-color: transparent !important;
  }
}
</style>

<script>
//import Trace from "@/components/Trace";

import PipelineBuildsGroup from "@/components/PipelineBuildsGroup";
import DownloadArtifact from "@/components/DownloadArtifact";
import StateChip from "@/components/StateChip";
import JobTags from "@/components/JobTags";

import { mapState } from "vuex";

export default {
  name: "BuildsView",

  data: () => ({
    loading: true,
    headers: [
      {
        text: "Status",
        value: "state",
        width: "1%",
      },
      {
        text: "Job ID",
        value: "id",
      },
      {
        text: "Name",
        value: "name",
      },
      {
        text: "Duration",
        value: "duration",
      },
      {
        text: "",
        value: "actions",
      },
    ],
    tag: "",
    current_stage: "",
  }),

  created() {
    this.$store.dispatch("getPipeline", { id: this.$route.params.id });
    this.$store.dispatch("getPipelineBuilds", { id: this.$route.params.id });
    this.$store.dispatch("getPipelineStages", { id: this.$route.params.id });
    this.$store.dispatch("getPipelineArtifacts", { id: this.$route.params.id });
  },

  components: {
    DownloadArtifact,
    StateChip,
    JobTags,
    PipelineBuildsGroup,
  },

  watch: {
    builds() {
      this.loading = false;
    },
  },

  computed: {
    ...mapState({
      artifacts(state) {
        return state.artifacts.filter(
          (a) => a.pipeline_id == this.$route.params.id
        );
      },
      pipeline(state) {
        return state.pipelines.find((d) => d.id === this.$route.params.id);
      },
      builds(state) {
        const collator = new Intl.Collator("en", {
          numeric: true,
          sensitivity: "base",
        });

        return state.builds
          .filter((d) => d.pipeline_id === this.$route.params.id)
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
      stages(state) {
        return state.stages
          .filter((d) => d.pipeline_id === this.$route.params.id)
          .reduce((stgs, s) => {
            const found = stgs.find((e) => e.id === s.id);
            if (found) {
              return stgs;
            } else {
              return [...stgs, s];
            }
          }, [])
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
    }),
  },

  methods: {
    retryJob(job) {
      this.$store.dispatch("retryJob", { id: job.id });
    },
    runJob(job) {
      this.$store.dispatch("runJob", { id: job.id });
    },
    duration(job) {
      const start = Date.parse(job.started_at);
      const finish = Date.parse(job.finished_at);

      const duration = Math.floor((finish - start) / 1000);

      if (isNaN(duration)) {
        return "";
      }

      return `${duration} seconds`;
    },

    getStageState(name) {
      const stage = this.stages.find((s) => s.name == name);
      if (stage) {
        return stage.state;
      }
      return "unknown";
    },

    isCurrentStage(stage) {
      let ret = false;
      if (stage === this.current_stage) {
        ret = true;
      }

      this.current_stage = stage;

      return ret;
    },

    retryableState(build) {
      const retryable = ["success", "failed"];

      if (build.retried) {
        return false;
      }

      if (retryable.indexOf(build.state) !== -1) {
        return true;
      }

      return false;
    },

    hasArtifacts: (build, artifacts) => {
      if (typeof build.artifacts !== "undefined") {
        if (build.artifacts.length > 0) return true;
      }

      if (typeof artifacts !== "undefined") {
        let found = artifacts.find(
          (a) => a.build_id == build.id && a.type == "archive"
        );
        if (found) return true;
      }

      return false;
    },

    isManualJob(build) {
      if (build.state === "manual") {
        return true;
      }
      return false;
    },

    isCancelable(state) {
      if (["created", "pending", "running"].indexOf(state) !== -1) {
        return true;
      }
      return false;
    },

    delete(pipeline) {
      this.$store.dispatch("deletePipeline", { id: pipeline.id });
      this.$router.push({
        name: "PipelinesView",
      });
    },

    cancel(pipeline) {
      this.$store.dispatch("cancelPipeline", { id: pipeline.id });
    },
  },
};
</script>
