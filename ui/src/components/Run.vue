<template>
  <v-btn v-bind="$attrs" @click="run">Run</v-btn>
</template>

<script>
export default {
  name: "Run",

  props: ["workflow"],

  methods: {
    run() {
      this.$store
        .dispatch("createPipeline", {
          id: this.$props.workflow.id,
        })
        .then((pipeline) => {
          this.$router.push({
            name: "PipelinesView",
            params: { id: pipeline.id },
          });
        });
    },
  },
};
</script>