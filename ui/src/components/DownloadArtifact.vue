<template>
  <v-btn class="mx-2" v-bind="$attrs" dark small color="grey" @click="download">
    <v-icon dark>mdi-cloud-download</v-icon>
  </v-btn>
</template>

<script>
//import axios from "axios";
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
    // Source: https://codepen.io/nigamshirish/pen/ZMpvRa
    forceFileDownload(response, name) {
      console.log(name);
      const url = window.URL.createObjectURL(new Blob([response.data]));
      const link = document.createElement("a");
      link.href = url;
      link.setAttribute("download", name); //or any other extension
      document.body.appendChild(link);
      link.click();
    },

    download() {
      window.open(
        `${window.location.origin}/v1/artifacts/${this.artifact.id}/download`,
        "_blank"
      );
      /*
      const artifact = this.artifact;
      console.log(artifact);
      axios({
        method: "get",
        url: `${window.location.origin}/v1/artifacts/${artifact.id}/download`,
        responseType: "arraybuffer",
      })
        .then((response) => {
          this.forceFileDownload(response, artifact.name);
        })
        .catch(() => console.log("error occured"));
      */
    },
  },
};
</script>
