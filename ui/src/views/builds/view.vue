<template>
  <div class="mt-4" v-if="build">
    <v-navigation-drawer app clipped right>
      <v-list-item>
        <v-list-item-content>
          <v-list-item-title class="title">{{ build.name }}</v-list-item-title>
        </v-list-item-content>
        <v-list-item-action>
          <v-btn
            outlined
            small
            class="float-right"
            v-if="isInteruptible(build)"
            @click="cancelJob(build)"
          >Cancel</v-btn>
          <v-btn
            outlined
            small
            class="float-right"
            v-if="isRetryable(build)"
            @click="retryJob(build)"
          >Retry</v-btn>
        </v-list-item-action>
      </v-list-item>

      <v-divider></v-divider>

      <v-list>
        <v-subheader class="font-weight-black">Metadata</v-subheader>
        <v-list-item dense>
          <v-list-item-content>
            <v-list-item-title>
              <strong>Duration:</strong>
              {{ build.duration }} seconds
            </v-list-item-title>
          </v-list-item-content>
        </v-list-item>
        <v-list-item dense>
          <v-list-item-content>
            <v-list-item-title>
              <strong>Timeout:</strong>
              {{ [build.timeout, "seconds"] | duration("asSeconds") }} seconds
            </v-list-item-title>
          </v-list-item-content>
        </v-list-item>
        <v-list-item dense>
          <v-list-item-content>
            <v-list-item-title>
              <strong>Runner:</strong>
              <router-link
                class="ml-1"
                :to="{name: 'RunnersView', params: { id: build.runner_id }}"
              >{{ build.runner_id }}</router-link>
            </v-list-item-title>
          </v-list-item-content>
        </v-list-item>
        <v-list-item dense v-if="build.tags">
          <v-list-item-content>
            <v-list-item-title>
              <strong>Tags:</strong>
              <JobTags :job="build" :jobTags="build.tags" class="ml-2" x-small />
            </v-list-item-title>
          </v-list-item-content>
        </v-list-item>
      </v-list>

      <div v-if="artifact">
        <v-divider></v-divider>
        <v-list>
          <v-subheader class="font-weight-black">Job Artifacts</v-subheader>
          <v-list-item dense>
            <v-list-item-content>
              <p class="caption">
                These artifacts are the latest. They will not be deleted (even
                if expired) until newer artifacts are available.
              </p>
              <v-btn small class="ma-0" v-if="artifact.expires_at">Keep</v-btn>
              <DownloadArtifact :build="build" type="archive" />
              <BrowseArtifact :build="build" type="archive" />
            </v-list-item-content>
          </v-list-item>
        </v-list>
      </div>

      <SidebarPipelineJobs :build="build" />
    </v-navigation-drawer>

    <h3>
      <StateChip :state="build.state" />
      Job #{{ build.id }} triggered {{ build.created_at | moment("from") }} by
      AUTHOR (queued for
      {{ [queuedFor(build), "seconds"] | duration("asSeconds") }} seconds)
    </h3>
    <v-divider></v-divider>
    <br />
    <Trace :build="build" />
  </div>
</template>

<script>
import BrowseArtifact from "@/components/BrowseArtifact";
import DownloadArtifact from "@/components/DownloadArtifact";
import StateChip from "@/components/StateChip";
import JobTags from "@/components/JobTags";
import Trace from "@/components/Trace";
import SidebarPipelineJobs from "@/components/SidebarPipelineJobs";

import { mapState } from "vuex";

export default {
  name: "BuildsView",

  data() {
    return {
      stage: null,
    };
  },

  created() {
    this.fetchData();
  },

  mounted() {
    this.$socket.sendObj({
      action: "subscribe",
    });
  },
  umounted() {
    this.$socket.sendObj({
      action: "unsubscribe",
    });
  },

  watch: {
    $route() {
      this.fetchData();
    },
  },

  components: {
    DownloadArtifact,
    BrowseArtifact,
    StateChip,
    JobTags,
    Trace,
    SidebarPipelineJobs,
  },

  computed: {
    ...mapState({
      builds: (state) => state.builds,
      build(state) {
        return state.builds.find((d) => d.id === this.$route.params.id);
      },
      artifact(state) {
        return state.artifacts.find(
          (a) => a.build_id == this.$route.params.id && a.type == "archive"
        );
      },
      artifacts(state) {
        return state.artifacts.filter(
          (a) => a.build_id == this.$route.params.id
        );
      },
      stages(state) {
        return state.stages.filter(
          (s) => s.pipeline_id === this.build.pipeline_id
        );
      },
    }),
  },

  methods: {
    cancelJob(job) {
      this.$store.dispatch("cancelJob", { id: job.id });
    },
    async retryJob(job) {
      const newJob = await this.$store.dispatch("retryJob", { id: job.id });
      this.$router.push({
        name: "BuildsView",
        params: { id: newJob.id },
      });
    },
    fetchData() {
      this.$store.dispatch("getBuild", { id: this.$route.params.id });
      this.$store.dispatch("getBuildArtifact", {
        id: this.$route.params.id,
        type: "archive",
      });
    },
    fgcolor(status) {
      if (status == "success") {
        return "white";
      }

      return "black";
    },
    bgcolor(status) {
      if (status == "success") {
        return "green";
      } else if (status == "skipped") {
        return "grey";
      }

      return "pink";
    },
    icon(status) {
      if (status == "success") {
        return "mdi-checkbox-marked-circle";
      } else if (status == "skipped") {
        return "mdi-debug-step-over";
      } else if (status == "failed") {
        return "mdi-minus-circle-outline";
      }

      return "mdi-help-circle-outline";
    },
    queuedFor(build) {
      const started = Date.parse(build.started_at);
      const created = Date.parse(build.created_at);

      const duration = Math.floor((started - created) / 1000);

      return duration;
    },
    hasArtifacts: (build, artifacts) => {
      if (typeof build.artifacts !== "undefined") {
        if (build.artifacts.length > 0) return true;
      }

      if (typeof artifacts !== "undefined") {
        let found = artifacts.find((a) => a.build_id == build.id);
        if (found) return true;
      }

      return false;
    },
    isInteruptible: (build) => {
      const states = ["running", "created", "pending"];
      if (states.indexOf(build.state) !== -1) {
        return true;
      }
      return false;
    },
    isRetryable: (build) => {
      const states = ["success", "failure", "canceled"];
      if (states.indexOf(build.state) !== -1 && build.retried == false) {
        return true;
      }

      return false;
    },
  },
};
</script>
