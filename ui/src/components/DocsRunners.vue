<template>
  <v-expansion-panels v-model="panel" class="mt-5">
    <v-expansion-panel v-bind:key="0">
      <v-expansion-panel-header>Quick Start</v-expansion-panel-header>
      <v-expansion-panel-content>
        <ol>
          <li>
            Download the GitLab Runner for your platform
            <a
              href="https://gitlab.com/gitlab-org/gitlab-runner/-/releases/v13.7.0"
              target="_blank"
            >Releases</a>
          </li>
          <li>
            Register your GitLab Runner with Pipeliner, per sure to get the IP of the host your Pipeliner Server is running on.
            <br />
            <code>gitlab-runner register -n --run-untagged -r pipeliner --executor docker --docker-image ubuntu:20.04 -c config.toml -u {{ origin }}</code>
          </li>
          <li>
            Start your GitLab Runner
            <code>gitlab-runner run -c config.toml</code>
          </li>
          <li>Verify your GitLab Runner has come online on this page.</li>
          <li>(Optional) Change your runner's concurrency by editing the config.toml and changing the concurrency value (example: from 1 to 5)</li>
        </ol>
      </v-expansion-panel-content>
    </v-expansion-panel>
  </v-expansion-panels>
</template>

<script>
export default {
  name: "DocsRunner",
  data: () => ({
    panel: null,
  }),
  props: ["runners"],
  mounted() {
    if (this.runners.length > 0) {
      this.panel = null;
      return;
    }

    this.panel = 0;
  },
  computed: {
    origin() {
      return window.location.origin;
    },
  },
};
</script>