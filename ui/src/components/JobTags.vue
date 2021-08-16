<template>
  <Tags :tags="tags" v-if="tags" v-bind="$attrs" />
</template>

<script>
import { mapState } from "vuex";

import Tags from "@/components/Tags";

export default {
  name: "JobTags",

  props: ["job", "jobTags"],

  components: {
    Tags,
  },

  computed: {
    ...mapState({
      buildTags(state) {
        return state.build_tags.filter((t) => t.build_id === this.job.id);
      },
    }),
    tags() {
      return this.buildTags;
    },
  },

  created() {
    this.$store.dispatch("getJobTags", { id: this.job.id });
  },
};
</script>