<template>
  <div>
    Artifacts
    <v-data-table :headers="headers" :items="files" hide-default-footer></v-data-table>
  </div>
</template>
<script>
import { mapState } from "vuex";

export default {
  name: "ArtifactsView",

  data() {
    return {
      headers: [
        {
          text: "File",
          align: "start",
          sortable: false,
          value: "name",
        },
        { text: "Size", value: "size" },
      ],
    };
  },

  mounted() {
    this.$store.dispatch("getArtifactFiles", {
      id: this.$route.params.artifact_id,
    });
  },

  computed: {
    ...mapState({
      files(state) {
        let artifactId = this.$route.params.artifact_id;
        return state.artifact_files.filter((f) => {
          return f.artifact_id == artifactId;
        });
      },
    }),
  },
};
</script>