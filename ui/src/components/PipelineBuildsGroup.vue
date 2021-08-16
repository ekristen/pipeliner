<template>
  <span v-if="stage">
    <StateIcon :state="stage.state" />
    <strong>{{ stage.name }}</strong>
  </span>
</template>
<script>
import StateIcon from "@/components/StateIcon";

import { mapState } from "vuex";

export default {
  name: "PipelineBuildsGroup",

  props: ["name", "pipeline"],

  components: {
    StateIcon,
  },

  computed: {
    ...mapState({
      stage(state) {
        return state.stages.find(
          (s) => s.pipeline_id === this.pipeline.id && s.name === this.name
        );
      },
    }),
  },
};
</script>