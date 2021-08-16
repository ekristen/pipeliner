<template>
  <v-btn class="mx-2" v-bind="$attrs" dark small color="grey" @click="browse">Browse</v-btn>
</template>

<script>
import { mapState } from "vuex";

export default {
  name: "DownloadArtifact",

  props: ["build", "type"],

  created() {
    this.$store.dispatch("getBuildArtifact", {
      id: this.build.id,
      type: this.type,
    });
  },

  computed: {
    ...mapState({
      artifact(state) {
        return state.artifacts.find((a) => {
          return a.build_id == this.build.id && a.type == this.type;
        });
      },
    }),
  },

  methods: {
    browse() {
      this.$router.push({
        name: "ArtifactsView",
        params: { id: this.artifact.build_id, artifact_id: this.artifact.id },
      });
    },
  },
};
</script>
